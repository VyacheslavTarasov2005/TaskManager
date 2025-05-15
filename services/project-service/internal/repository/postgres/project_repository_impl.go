package postgres

import (
	"context"
	"project-service/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *projectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Add(ctx context.Context, project models.Project) error {
	return r.db.WithContext(ctx).Create(&project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByOwner(ctx context.Context, owner uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "owner_id = ?", owner).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	var projects []*models.Project
	if err := r.db.WithContext(ctx).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project models.Project) error {
	return r.db.WithContext(ctx).Save(&project).Error
}
