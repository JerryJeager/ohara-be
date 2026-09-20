package documents

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/JerryJeager/ohara-be/config"
	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/service/websites"
	"github.com/JerryJeager/ohara-be/internal/utils"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"google.golang.org/genai"
)

type DocumentSv interface {
	EmbedDocument(ctx context.Context, content *models.Document) error
	QueryDocument(ctx context.Context, query *models.Query) (string, error)
	ChunkDocument(ctx context.Context, chunkReq *models.ChunkReq) []string
	QueryWebsiteDocument(ctx context.Context, websiteID uuid.UUID, query *models.Query) (string, error)
}

type DocumentServ struct {
	repo DocumentStore
}

func NewDocumentService(repo DocumentStore) *DocumentServ {
	return &DocumentServ{repo: repo}
}

func ChunkAndEmbedDocument(content, title, pageUrl string, websiteID uuid.UUID) error {
	if title == "" {
		title = "none"
	}

	chunks := utils.ChunkByFixedSizeWithOverlap(content, 1000, 200)
	var documents []models.Document
	var contents []*genai.Content
	for _, chunk := range chunks {
		embeddingInput := fmt.Sprintf("title: %s | text: %s", title, chunk)
		contents = append(contents,
			genai.NewContentFromText(embeddingInput, genai.RoleUser),
		)
	}

	ctx := context.Background()

	result, err := config.AI.Models.EmbedContent(ctx,
		"gemini-embedding-2",
		contents,
		nil,
	)
	if err != nil {
		log.Printf("failed to embed chunk=> error: %s", err.Error())
		return err
	}

	if len(result.Embeddings) != len(chunks) {
		return fmt.Errorf(
			"embedding count mismatch: got %d embeddings for %d chunks",
			len(result.Embeddings),
			len(chunks),
		)
	}

	for i, emb := range result.Embeddings {
		var newDoc models.Document
		newDoc.Content = chunks[i]
		newDoc.Embedding = pgvector.NewVector(emb.Values)
		newDoc.WebsiteID = websiteID
		newDoc.PageUrl = pageUrl
		documents = append(documents, newDoc)
	}
	if len(documents) == 0 {
		_ = websites.UpdateWebsiteStatus(websiteID, "failed")
		log.Printf("all chunks failed to embed for website %s", websiteID)
		return errors.New("no chunks were successfully embedded")
	}

	return CreateDocuments(&documents)
}

func (s *DocumentServ) EmbedDocument(ctx context.Context, content *models.Document) error {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return err
	}

	contents := []*genai.Content{
		genai.NewContentFromText(content.Content, genai.RoleUser),
	}
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-2",
		contents,
		nil,
	)
	if err != nil {
		return err
	}

	content.Embedding = pgvector.NewVector(result.Embeddings[0].Values)

	return s.repo.CreateDocument(ctx, content)
}

func (s *DocumentServ) QueryDocument(ctx context.Context, query *models.Query) (string, error) {
	queryContext := "no available relevant data" //default context
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	contents := []*genai.Content{
		genai.NewContentFromText(query.Query, genai.RoleUser),
	}
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-2",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}

	queryEmbedding := pgvector.NewVector(result.Embeddings[0].Values)

	documentList, err := s.repo.GetDocumentChunks(ctx, queryEmbedding)
	if err != nil {
		return "", err
	}

	if documentList != nil {
		queryContext = ""
		for _, d := range *documentList {
			queryContext += fmt.Sprintf("%s \n", d.Content)
		}
	}

	prompt := fmt.Sprintf(`
		You are a helpful assistant.

		Answer the user's question using only the provided context.

		If the answer cannot be found in the context, say:
		"I don't have enough information to answer that."

		Context:
		%s

		Question:
		%s
	`, queryContext, query.Query)

	contents = []*genai.Content{
		genai.NewContentFromText(prompt, genai.RoleUser),
	}

	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.7-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}

	return response.Text(), nil
}

func (s *DocumentServ) QueryWebsiteDocument(ctx context.Context, websiteID uuid.UUID, query *models.Query) (string, error) {
	queryContext := "no available relevant data" //default context

	embeddingInput := fmt.Sprintf("task: question answering | query: %s", query.Query)
	contents := []*genai.Content{
		genai.NewContentFromText(embeddingInput, genai.RoleUser),
	}
	result, err := config.AI.Models.EmbedContent(ctx,
		"gemini-embedding-2",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}

	queryEmbedding := pgvector.NewVector(result.Embeddings[0].Values)

	documentList, err := s.repo.GetDocumentChunksForQuery(ctx, websiteID, queryEmbedding)
	if err != nil {
		return "", err
	}

	if documentList != nil {
		queryContext = ""
		for _, d := range *documentList {
			queryContext += fmt.Sprintf("%s \n", d.Content)
		}
	}

	prompt := fmt.Sprintf(`
		You are a helpful assistant.

		Answer the user's question using only the provided context.

		If the answer cannot be found in the context, say:
		"I don't have enough information to answer that."

		Context:
		%s

		Question:
		%s
	`, queryContext, query.Query)

	contents = []*genai.Content{
		genai.NewContentFromText(prompt, genai.RoleUser),
	}

	response, err := config.AI.Models.GenerateContent(
		ctx,
		"gemini-3.7-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}

	return response.Text(), nil
}

func (s *DocumentServ) ChunkDocument(ctx context.Context, chunkReq *models.ChunkReq) []string {
	return utils.ChunkBySentence(chunkReq.Content)
}
