package models

import (
	"crypto/sha256"
	"encoding/hex"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email" binding:"required"`
	ProfilePicture string    `json:"profile_picture"`
	GoogleID       string    `json:"google_id"`
	RefreshToken   string    `json:"refresh_token"`
	CreatedAt      time.Time `json:"created_at"`
}

type GoogleAuthReq struct {
	IDToken string `json:"id_token" binding:"required"`
}

type AccessTokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func VerifyRefreshToken(refreshToken, storedHash string) bool {
	refreshToken = html.EscapeString(strings.TrimSpace(refreshToken))

	hash := sha256.Sum256([]byte(refreshToken))
	incomingHash := hex.EncodeToString(hash[:])

	return incomingHash == storedHash
}

func (user *User) HashRefreshToken() {
	user.RefreshToken = html.EscapeString(strings.TrimSpace(user.RefreshToken))

	hash := sha256.Sum256([]byte(user.RefreshToken))

	user.RefreshToken = hex.EncodeToString(hash[:])
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

func VerifyPassword(password, hashedPassword string) error {
	password = html.EscapeString(strings.TrimSpace(password))
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
