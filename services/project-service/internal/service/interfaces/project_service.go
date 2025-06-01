package interfaces

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, projectName string, userId uuid.UUID) (*uuid.UUID, error)
	GetProject(ctx context.Context, projectId uuid.UUID) (*models.Project, error)
	GetProjectsByUser(ctx context.Context, ownerId uuid.UUID, query string) ([]*models.Project, error)
	Update(ctx context.Context, userId, projectId uuid.UUID, newName string) (*models.Project, error)
	Delete(ctx context.Context, projectId uuid.UUID) error
}
