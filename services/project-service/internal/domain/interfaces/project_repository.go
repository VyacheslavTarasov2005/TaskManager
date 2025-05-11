package interfaces

import (
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Add(project models.Project) error
	GetByID(id uuid.UUID) (*models.Project, error)
	GetByOwner(owner uuid.UUID) (*models.Project, error)
	GetAll() ([]*models.Project, error)
	Update(project models.Project) error
}
