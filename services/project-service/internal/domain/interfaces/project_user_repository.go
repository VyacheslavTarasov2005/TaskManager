package interfaces

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type ProjectUserRepository interface {
	DeleteProject(ctx context.Context, proj_id uuid.UUID) error
	Add(ctx context.Context, projectUser models.ProjectUser) error
	Delete(ctx context.Context, proj_id, user_id uuid.UUID) error
	GetRole(ctx context.Context, proj_id, user_id uuid.UUID) (*string, error)
}
