package interfaces

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

// отредачить
type UserService interface {
	Create(ctx context.Context, projectName string, userId uuid.UUID) (*uuid.UUID, error)
	Get(ctx context.Context, projectId uuid.UUID) (*models.Project, error)
	Update(ctx context.Context, projectId uuid.UUID, newName string) (*models.Project, error)
	Delete(ctx context.Context, projectId uuid.UUID) error
}
