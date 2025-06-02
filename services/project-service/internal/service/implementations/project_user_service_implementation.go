package implementations

import (
	"context"
	"project-service/internal/domain/interfaces"
	"project-service/internal/domain/models"
	"project-service/internal/service/errors"

	"github.com/google/uuid"
)

type projectUserServiceImpl struct {
	repo interfaces.ProjectUserRepository
}

func NewProjectUserServiceImpl(repo interfaces.ProjectUserRepository) *projectUserServiceImpl {
	return &projectUserServiceImpl{repo: repo}
}

func (s *projectUserServiceImpl) DelteByProject(ctx context.Context, proj_id uuid.UUID) error {
	err := s.repo.DeleteProject(ctx, proj_id)
	if err != nil {
		return errors.ApplicationError{
			StatusCode: 500,
			Code:       "InternalServerError",
			Errors: map[string]string{
				"message": "Failed to delete users from project",
			},
		}
	}
	return nil
}

func (s *projectUserServiceImpl) AddToProject(ctx context.Context, user_id, project_id uuid.UUID, role string) error {
	e_role, err := s.repo.GetRole(ctx, project_id, user_id)
	if err != nil {
		return errors.ApplicationError{
			StatusCode: 500,
			Code:       "InternalServerError",
			Errors: map[string]string{
				"message": "Failed to check user role",
			},
		}
	}
	if e_role != nil {
		return errors.ApplicationError{
			StatusCode: 409,
			Code:       "UserAlreadyInProject",
			Errors: map[string]string{
				"message": "User is already added to the project",
			},
		}
	}

	newProjectUser, err := models.NewProjectUser(user_id, project_id, role)
	if err != nil {
		return errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRole",
			Errors: map[string]string{
				"message": "Invalid user role: must be OWNER or MEMBER",
			},
		}
	}

	err = s.repo.Add(ctx, *newProjectUser)
	if err != nil {
		return errors.ApplicationError{
			StatusCode: 500,
			Code:       "InternalServerError",
			Errors: map[string]string{
				"message": "Failed to add user to project",
			},
		}
	}
	return nil
}

func (s *projectUserServiceImpl) DeleteFromProject(ctx context.Context, proj_id, user_id uuid.UUID) error {
	err := s.repo.Delete(ctx, proj_id, user_id)
	if err != nil {
		return errors.ApplicationError{
			StatusCode: 404,
			Code:       "UserNotFound",
			Errors: map[string]string{
				"message": "User not found in the project",
			},
		}
	}
	return nil
}

func (s *projectUserServiceImpl) GetUserRole(ctx context.Context, proj_id, user_id uuid.UUID) (*models.UserRole, error) {
	e_role, err := s.repo.GetRole(ctx, proj_id, user_id)
	if err != nil {
		return nil, errors.ApplicationError{
			StatusCode: 500,
			Code:       "InternalServerError",
			Errors: map[string]string{
				"message": "Failed to retrieve user role",
			},
		}
	}
	if e_role == nil {
		return nil, errors.ApplicationError{
			StatusCode: 404,
			Code:       "UserNotInProject",
			Errors: map[string]string{
				"message": "User is not a member of this project",
			},
		}
	}
	return e_role, nil
}
