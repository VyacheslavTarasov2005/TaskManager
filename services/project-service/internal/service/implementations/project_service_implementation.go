package implementations

import (
	"context"
	"fmt"
	"project-service/internal/domain/interfaces"
	"project-service/internal/domain/models"
	"project-service/internal/service/errors"
	"project-service/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type projectServiceImpl struct {
	repo interfaces.ProjectRepository
}

func NewProjectServiceImpl(repo interfaces.ProjectRepository) *projectServiceImpl {
	return &projectServiceImpl{repo: repo}
}

func (s *projectServiceImpl) Create(ctx context.Context, projectName string, userId uuid.UUID) (*uuid.UUID, error) {

	existing, err := s.repo.GetByOwnerAndProjectName(ctx, userId, projectName)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("project with name '%s' already exists for this user", projectName)
	}
	if err != nil {
		return nil, err
	}

	project := models.NewProject(projectName, userId)

	if err := s.repo.Add(ctx, *project); err != nil {
		return nil, err
	}

	return &project.ID, nil
}

func (s *projectServiceImpl) GetProject(ctx context.Context, projectId uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, projectId)
}

func (s *projectServiceImpl) GetProjectsByUser(ctx context.Context, ownerId uuid.UUID, query string) ([]*models.Project, error) {
	return s.repo.GetByOwner(ctx, ownerId, query)
}

func (s *projectServiceImpl) Update(ctx context.Context, userId, projectId uuid.UUID, newName string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, projectId)
	if err != nil {
		return nil, err
	}

	projectSameName, err := s.repo.GetByID(ctx, projectId)
	if err != nil {
		return nil, err
	}
	if projectSameName != nil && projectSameName.Name == newName {
		return nil, errors.ApplicationError{
			StatusCode: 409,
			Code:       "EmailConflict",
			Errors: map[string]string{
				"message": "Project name is already in use",
			},
		}
	}

	project.Name = newName
	project.UpdatedAt = utils.Ptr(time.Now())

	if err := s.repo.Update(ctx, *project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) Delete(ctx context.Context, projectId uuid.UUID) error {
	return s.repo.Delete(ctx, projectId)
}
