package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type Document struct {
	ID        int             `json:"id"`
	PageUrl   string          `json:"page_url"`
	WebsiteID uuid.UUID       `json:"website_id"`
	Content   string          `json:"content" binding:"required"`
	Embedding pgvector.Vector `json:"embedding"`
	CreatedAt time.Time       `json:"created_at"`
}

type QueryDocument struct {
	ID         int     `json:"id"`
	Content    string  `json:"content" binding:"required"`
	Similarity float32 `json:"similarity"`
}

type QueryDocumentList []QueryDocument

type DocumentsList []Document

type Query struct {
	Query string `json:"query" binding:"required"`
}

type ChunkReq struct {
	Size    int    `json:"size"`
	Content string `json:"content"`
	Overlap int    `json:"overlap"`
}
