package interfaces

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Add(ctx context.Context, project models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetByOwner(ctx context.Context, owner uuid.UUID) (*models.Project, error)
	GetAll(ctx context.Context) ([]*models.Project, error)
	Update(ctx context.Context, project models.Project) error
}
