package postgres

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectUserRepository struct {
	db *gorm.DB
}

func NewProjectUserRepositoryImpl(db *gorm.DB) *projectUserRepository {
	return &projectUserRepository{db: db}
}

func (r *projectUserRepository) DeleteProject(ctx context.Context, proj_id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ?", proj_id).Delete(&models.ProjectUser{}).Error
}

func (r *projectUserRepository) Add(ctx context.Context, projectUser models.ProjectUser) error {
	return r.db.WithContext(ctx).Create(&projectUser).Error
}

func (r *projectUserRepository) Delete(ctx context.Context, proj_id, user_id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", proj_id, user_id).Delete(&models.ProjectUser{}).Error
}

func (r *projectUserRepository) GetRole(ctx context.Context, proj_id, user_id uuid.UUID) (*string, error) {
	var project models.ProjectUser
	if err := r.db.WithContext(ctx).First(&project, "project_id = ? AND user_id = ?", proj_id, user_id).Error; err != nil {
		return nil, err
	}
	return &project.Role, nil
}
