package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user-service/internal/domain/models"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) Add(ctx context.Context, user models.User) error {
	return repo.db.WithContext(ctx).Create(&user).Error
}

func (repo *UserRepositoryImpl) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := repo.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (repo *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (repo *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (repo *UserRepositoryImpl) Update(ctx context.Context, user models.User) error {
	return repo.db.WithContext(ctx).Save(user).Error
}
