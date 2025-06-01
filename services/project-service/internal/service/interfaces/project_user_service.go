package interfaces

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type ProjectUserService interface {
	DelteByProject(ctx context.Context, proj_id uuid.UUID) error
	AddToProject(ctx context.Context, projectUser models.ProjectUser) error
	DeleteFromProject(ctx context.Context, proj_id, user_id uuid.UUID) error
	GetUserRole(ctx context.Context, proj_id, user_id uuid.UUID) (*string, error)
}
