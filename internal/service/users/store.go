package users

import (
	"context"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStore interface {
	SaveUser(ctx context.Context, user *models.User) error
	CreateUser(ctx context.Context, user *models.User, otp *models.Otp) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserOtp(ctx context.Context, userID uuid.UUID) (*models.Otp, error)
	VerifyUser(ctx context.Context, userID uuid.UUID) error
}

type UserRepo struct {
	client *gorm.DB
}

func NewUserRepo(client *gorm.DB) *UserRepo {
	return &UserRepo{client: client}
}

func (r *UserRepo) SaveUser(ctx context.Context, user *models.User) error {
	return r.client.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User, otp *models.Otp) error {
	err := r.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		if otp != nil {
			if err := tx.WithContext(ctx).Create(otp).Error; err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.client.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.client.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserOtp(ctx context.Context, userID uuid.UUID) (*models.Otp, error) {
	var otp models.Otp
	if err := r.client.WithContext(ctx).Where("user_id = ?", userID).First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *UserRepo) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.client.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_verified", true).Error; err != nil {
		return err
	}
	return nil
}
