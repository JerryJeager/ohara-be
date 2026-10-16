package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Website struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	Url               string          `json:"url" binding:"required"`
	Name              string          `json:"name"`
	Contact           string          `json:"contact" binding:"required"`
	Summary           string          `json:"summary"`
	Tone              string          `json:"tone"`
	Topics            *Topics         `json:"topics"`
	Status            string          `json:"status"`
	IsLocalDevEnabled bool            `json:"is_local_dev_enabled"`
	AllowedOrigins    *AllowedOrigins `json:"allowed_origins"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type WebsiteList []Website

type PageUrl struct {
	PageUrl string `json:"page_url"`
}

type IndexedPages []PageUrl

type IsLocalDevEnabled struct {
	IsLocalDevEnabled bool `json:"is_local_dev_enabled" binding:"required"`
}

type SiteProfile struct {
	Summary string   `json:"summary"`
	Tone    string   `json:"tone"`
	Topics  []string `json:"topics"`
}

type Topics []string

func (a Topics) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Topics) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, returned %T", value)
	}
	return json.Unmarshal(bytes, a)
}

type AllowedOrigins []string

func (a AllowedOrigins) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AllowedOrigins) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, returned %T", value)
	}
	return json.Unmarshal(bytes, a)
}

func (w *Website) MarshalJSON() ([]byte, error) {
	website := map[string]interface{}{
		"id":              w.ID,
		"user_id":         w.UserID,
		"url":             w.Url,
		"name":            w.Name,
		"contact":         w.Contact,
		"summary":         w.Summary,
		"tone":            w.Tone,
		"topics":          w.Topics,
		"status":          w.Status,
		"allowed_origins": w.AllowedOrigins,
		"created_at":      w.CreatedAt,
		"updated_at":      w.UpdatedAt,
	}

	return json.Marshal(website)
}
