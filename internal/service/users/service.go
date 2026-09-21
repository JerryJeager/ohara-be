package users

import (
	"context"
	"errors"
	"log"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/utils"
	"github.com/google/uuid"
)

type UserSv interface {
	AuthUser(ctx context.Context, user *models.User) (*models.User, *models.AccessTokens, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
	RefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) (*models.AccessTokens, error)
}

type UserServ struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserServ {
	return &UserServ{repo: repo}
}

// create new user if user doesn't exist and then login else just login
func (s *UserServ) AuthUser(ctx context.Context, user *models.User) (*models.User, *models.AccessTokens, error) {
	var userID uuid.UUID
	existingUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("user does not exist-> creating new user with email: %s", user.Email)
		userID = uuid.New()
	} else if existingUser.Email == user.Email {
		log.Printf("user exists already->log in user with email: %s", existingUser.Email)
		userID = existingUser.ID
		user = existingUser
	}
	user.ID = userID

	accessToken, err := utils.GenerateToken(user.ID, 1)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := utils.GenerateToken(user.ID, 30)
	if err != nil {
		return nil, nil, err
	}
	user.RefreshToken = refreshToken
	user.HashRefreshToken()
	if err := s.repo.SaveUser(ctx, user); err != nil {
		return nil, nil, err
	}

	return user, &models.AccessTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserServ) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *UserServ) RefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) (*models.AccessTokens, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !models.VerifyRefreshToken(refreshToken, user.RefreshToken) {
		return nil, errors.New("invalid refresh token used")
	}

	accessToken, err := utils.GenerateToken(user.ID, 1)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := utils.GenerateToken(user.ID, 30)
	if err != nil {
		return nil, err
	}

	user.RefreshToken = newRefreshToken
	user.HashRefreshToken()
	if err := s.repo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return &models.AccessTokens{
		RefreshToken: newRefreshToken,
		AccessToken:  accessToken,
	}, nil

}
