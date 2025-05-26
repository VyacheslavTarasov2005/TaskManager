package implementations

import (
	"context"
	"errors"
	"fmt"
	"project-service/internal/domain/interfaces"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	project := models.Project{
		ID:      uuid.New(),
		Name:    projectName,
		OwnerID: userId,
	}

	if err := s.repo.Add(ctx, project); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &project.ID, nil
}

func (s *projectServiceImpl) GetProject(ctx context.Context, projectId uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, projectId)
}

func (s *projectServiceImpl) GetProjectsByUser(ctx context.Context, ownerId uuid.UUID) ([]*models.Project, error) {
	return s.repo.GetByOwner(ctx, ownerId)
}

func (s *projectServiceImpl) Update(ctx context.Context, projectId uuid.UUID, newName string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	project.Name = newName

	if err := s.repo.Update(ctx, *project); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) Delete(ctx context.Context, projectId uuid.UUID) error {
	return s.repo.Delete(ctx, projectId)
}
