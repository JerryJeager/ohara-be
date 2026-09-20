package models

import (
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	FirstName    string       `json:"first_name" binding:"required"`
	LastName     string       `json:"last_name" binding:"required"`
	Email        string       `json:"email" binding:"required"`
	Password     string       `json:"password" binding:"required"`
	IsVerified   bool         `json:"is_verified"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUserEmail struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Otp    string `json:"otp" binding:"required"`
}

func (user *User) HashPassword() error {
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	password = html.EscapeString(strings.TrimSpace(password))
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

