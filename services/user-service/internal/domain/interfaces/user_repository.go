package interfaces

import (
	"context"
	"github.com/google/uuid"
	"user-service/internal/domain/models"
)

type UserRepository interface {
	Add(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user models.User) error
}
