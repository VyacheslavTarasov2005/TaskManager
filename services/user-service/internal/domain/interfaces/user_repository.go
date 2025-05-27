package interfaces

import (
	"github.com/google/uuid"
	"user-service/internal/domain/models"
)

type UserRepository interface {
	Add(user models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user models.User) error
}
