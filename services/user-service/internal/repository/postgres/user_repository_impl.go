package postgres

import (
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

func (repo *UserRepositoryImpl) Add(user models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepositoryImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (repo *UserRepositoryImpl) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (repo *UserRepositoryImpl) Update(user models.User) error {
	return repo.db.Save(user).Error
}
