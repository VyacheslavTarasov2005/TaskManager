package implementations

import (
	"project-service/internal/domain/interfaces"
)

type projectUserServiceImpl struct {
	repo interfaces.ProjectUserRepository
}

func NewProjectUserServiceImpl(repo interfaces.ProjectUserRepository) *projectUserServiceImpl {
	return &projectUserServiceImpl{repo: repo}
}
