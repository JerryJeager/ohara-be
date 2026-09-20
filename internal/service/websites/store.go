package websites

import (
	"context"

	"github.com/JerryJeager/ohara-be/config"
	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebsiteStore interface {
	CreateWebsite(ctx context.Context, website *models.Website) error
	GetWebsites(ctx context.Context) (*models.WebsiteList, error)
	DeleteWebsite(ctx context.Context, websiteID uuid.UUID) error
	GetWebsite(ctx context.Context, websiteID uuid.UUID) (*models.Website, error)
	UpdateWebsiteStatus(websiteID uuid.UUID, status string) error
}

type WebsiteRepo struct {
	client *gorm.DB
}

func NewWebsiteRepo(client *gorm.DB) *WebsiteRepo {
	return &WebsiteRepo{client: client}
}

func (r *WebsiteRepo) CreateWebsite(ctx context.Context, website *models.Website) error {
	return r.client.WithContext(ctx).Create(website).Error
}

func (r *WebsiteRepo) GetWebsites(ctx context.Context) (*models.WebsiteList, error) {
	var websiteList models.WebsiteList
	if err := r.client.WithContext(ctx).Find(&websiteList).Error; err != nil {
		return nil, err
	}
	return &websiteList, nil
}

func (r *WebsiteRepo) GetWebsite(ctx context.Context, websiteID uuid.UUID) (*models.Website, error) {
	var website models.Website
	if err := r.client.WithContext(ctx).First(&website, "id = ?", websiteID).Error; err != nil {
		return nil, err
	}
	return &website, nil
}

func (r *WebsiteRepo) DeleteWebsite(ctx context.Context, websiteID uuid.UUID) error {
	return r.client.WithContext(ctx).Delete(&models.Website{}, websiteID).Error
}

func (r *WebsiteRepo) UpdateWebsiteStatus(websiteID uuid.UUID, status string) error {
	return r.client.Model(&models.Website{}).Where("id = ?", websiteID).Update("status", status).Error
}

func UpdateWebsiteStatus(websiteID uuid.UUID, status string) error {
	return config.Session.Model(&models.Website{}).Where("id = ?", websiteID).Update("status", status).Error
}
