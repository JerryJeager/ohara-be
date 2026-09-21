package models

import (
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	Url               string    `json:"url" binding:"required"`
	Name              string    `json:"name"`
	Status            string    `json:"status"`
	IsLocalDevEnabled bool      `json:"is_local_dev_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type WebsiteList []Website
