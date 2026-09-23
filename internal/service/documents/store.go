package documents

import (
	"context"

	"github.com/JerryJeager/ohara-be/config"
	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type DocumentStore interface {
	CreateDocument(ctx context.Context, document *models.Document) error
	GetDocumentChunks(ctx context.Context, embedding pgvector.Vector) (*models.QueryDocumentList, error)
	GetDocumentChunksForQuery(ctx context.Context, websiteID uuid.UUID, embedding pgvector.Vector) (*models.QueryDocumentList, error)
	GetWebsite(ctx context.Context, websiteID uuid.UUID) (*models.Website, error)
}

type DocumentRepo struct {
	client *gorm.DB
}

func NewDocumentRepo(client *gorm.DB) *DocumentRepo {
	return &DocumentRepo{client: client}
}

func (r *DocumentRepo) CreateDocument(ctx context.Context, document *models.Document) error {
	return r.client.WithContext(ctx).Create(document).Error
}

func CreateDocuments(documents *[]models.Document) error {
	return config.Session.Create(documents).Error
}

func (r *DocumentRepo) GetDocumentChunks(ctx context.Context, embedding pgvector.Vector) (*models.QueryDocumentList, error) {
	var documentList models.QueryDocumentList

	// query := `
	// 	SELECT * FROM documents ORDER BY embedding <=> ? LIMIT 2
	// `

	// 	-- equivalent to similarity > 0.7 (1 - 0.7)
	query := `
		SELECT id, content, 1 - (embedding <=> ?) AS similarity
		FROM documents
		WHERE embedding <=> ? < 0.3
		ORDER BY embedding <=> ? ASC
		LIMIT 2;

	`

	if err := r.client.WithContext(ctx).Raw(query, embedding, embedding, embedding).Scan(&documentList).Error; err != nil {
		return nil, err
	}

	return &documentList, nil
}

func (r *DocumentRepo) GetDocumentChunksForQuery(ctx context.Context, websiteID uuid.UUID, embedding pgvector.Vector) (*models.QueryDocumentList, error) {
	var documentList models.QueryDocumentList
	// 	-- equivalent to similarity > 0.6 (1 - 0.6)
	query := `
		SELECT id, content, 1 - (embedding <=> ?) AS similarity
		FROM documents
		WHERE website_id = ? AND embedding <=> ? < 0.4
		ORDER BY embedding <=> ? ASC
		LIMIT 3;

	`

	if err := r.client.WithContext(ctx).Raw(query, embedding, websiteID, embedding, embedding).Scan(&documentList).Error; err != nil {
		return nil, err
	}

	return &documentList, nil
}

func (r *DocumentRepo) GetWebsite(ctx context.Context, websiteID uuid.UUID) (*models.Website, error) {
	var website models.Website
	if err := r.client.WithContext(ctx).Where("id = ?", websiteID).First(&website).Error; err != nil {
		return nil, err
	}
	return &website, nil
}
