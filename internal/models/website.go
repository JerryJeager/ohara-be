package models

import "github.com/google/uuid"

type Website struct {
	ID     uuid.UUID `json:"id"`
	Url    string    `json:"url" binding:"required"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

type WebsiteList []Website
