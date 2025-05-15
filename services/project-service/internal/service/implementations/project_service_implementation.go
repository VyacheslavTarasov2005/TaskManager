package implementations

import (
	"context"
	"project-service/internal/domain/interfaces"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
)

type projectServiceImpl struct {
	repo interfaces.ProjectRepository
}

func NewProjectService(repo interfaces.ProjectRepository) *projectServiceImpl {
	return &projectServiceImpl{repo: repo}
}

func (s *projectServiceImpl) Create(ctx context.Context, projectName string, userId uuid.UUID) (*uuid.UUID, error) {

}

func (s *projectServiceImpl) Get(ctx context.Context, projectId uuid.UUID) (*models.Project, error) {

}

func (s *projectServiceImpl) Update(ctx context.Context, projectId uuid.UUID, newName string) (*models.Project, error) {

}

func (s *projectServiceImpl) Delete(ctx context.Context, projectId uuid.UUID) error {

}
