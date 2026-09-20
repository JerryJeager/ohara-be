package users

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/utils"
	"github.com/JerryJeager/ohara-be/internal/utils/emails"
	"github.com/google/uuid"
)

type UserSv interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	VerfiyUserEmail(ctx context.Context, verify *models.VerifyUserEmail) error
	Login(ctx context.Context, user *models.UserLogin) (*models.User, string, error)
}

type UserServ struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserServ {
	return &UserServ{repo: repo}
}

func (s *UserServ) CreateUser(ctx context.Context, user *models.User) (string, error) {

	id := uuid.New()
	user.ID = id

	var otp models.Otp
	otp.ID = uuid.New()
	otp.UserID = id
	otp.Otp = utils.GetOtp()
	otp.ExpiresAt = time.Now().Add(time.Hour * 24 * 5) //expires after five days

	if err := user.HashPassword(); err != nil {
		return "", err
	}

	if err := s.repo.CreateUser(ctx, user, &otp); err != nil {
		return "", err
	}

	go func() {
		if err := emails.SendEmail(user.Email, "Welcome", emails.VerifyEmailTemplate(user.FirstName, otp.Otp)); err != nil {
			log.Printf("failed to send welcome email")
		}
	}()

	return id.String(), nil
}

func (s *UserServ) VerfiyUserEmail(ctx context.Context, verify *models.VerifyUserEmail) error {
	var user *models.User
	var err error
	if verify.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, verify.Email)
		if err != nil {
			user, err = s.repo.GetUserByID(ctx, verify.UserID)
			if err != nil {
				return err
			}
		}
	} else {
		user, err = s.repo.GetUserByID(ctx, verify.UserID)
		if err != nil {
			return err
		}
	}

	otp, err := s.repo.GetUserOtp(ctx, user.ID)
	if err != nil {
		return err
	}

	if otp.ExpiresAt.Compare(otp.CreatedAt) == -1 {
		return errors.New("expired token")
	}
	if otp.Otp != verify.Otp {
		return errors.New("otp is invalid")
	}
	return s.repo.VerifyUser(ctx, user.ID)
}

func (s *UserServ) Login(ctx context.Context, user *models.UserLogin) (*models.User, string, error) {
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)
	u, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, "", err
	}
	if !u.IsVerified {
		return nil, "", errors.New("only verified users can login")
	}

	if err := models.VerifyPassword(user.Password, u.Password); err != nil {
		return nil, "", err
	}
	token, err := utils.GenerateToken(u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}
