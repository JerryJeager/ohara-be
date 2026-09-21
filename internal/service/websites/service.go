package websites

import (
	"context"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/google/uuid"
)

type WebsiteSv interface {
	CreateWebsite(ctx context.Context, website *models.Website) (uuid.UUID, error)
	GetWebsites(ctx context.Context) (*models.WebsiteList, error)
	GetWebsite(ctx context.Context, userID uuid.UUID) (*[]models.Website, error)
	GetIndexedPages(ctx context.Context, websiteID uuid.UUID) (*models.IndexedPages, error)
	DeleteWebsite(ctx context.Context, websiteID uuid.UUID) error
	UpdateWebsiteStatus(websiteID uuid.UUID, status string) error
}

type WebsiteServ struct {
	repo WebsiteStore
}

func NewWebsiteService(repo WebsiteStore) *WebsiteServ {
	return &WebsiteServ{repo: repo}
}

func (s *WebsiteServ) CreateWebsite(ctx context.Context, website *models.Website) (uuid.UUID, error) {
	id := uuid.New()
	website.ID = id
	return id, s.repo.CreateWebsite(ctx, website)
}

func (s *WebsiteServ) GetWebsites(ctx context.Context) (*models.WebsiteList, error) {
	return s.repo.GetWebsites(ctx)
}

func (s *WebsiteServ) GetWebsite(ctx context.Context, userID uuid.UUID) (*[]models.Website, error) {
	return s.repo.GetWebsite(ctx, userID)
}

func (s *WebsiteServ) GetIndexedPages(ctx context.Context, websiteID uuid.UUID) (*models.IndexedPages, error) {
	return s.repo.GetIndexedPages(ctx, websiteID)
}

func (s *WebsiteServ) DeleteWebsite(ctx context.Context, websiteID uuid.UUID) error {
	return s.repo.DeleteWebsite(ctx, websiteID)
}

func (s *WebsiteServ) UpdateWebsiteStatus(websiteID uuid.UUID, status string) error {
	return s.repo.UpdateWebsiteStatus(websiteID, status)
}
