package models

import (
	"time"

	"github.com/google/uuid"
)

type Otp struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid; unique"`
	Otp       string    `json:"otp"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
