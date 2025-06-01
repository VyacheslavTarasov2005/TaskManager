package postgres

import (
	"context"
	"project-service/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepositoryImpl(db *gorm.DB) *projectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Add(ctx context.Context, project models.Project) error {
	now := time.Now()
	project.ID = uuid.New()
	project.CreatedAt = now
	return r.db.WithContext(ctx).Create(&project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}
func (r *projectRepository) GetByOwnerAndProjectName(ctx context.Context, owner uuid.UUID, projectName string) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "owner_id = ? AND name = ?", owner, projectName).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByOwner(ctx context.Context, owner uuid.UUID, query string) ([]*models.Project, error) {
	var projects []*models.Project

	db := r.db.WithContext(ctx).Model(&models.Project{}).Where("owner_id = ?", owner)

	if query != "" {
		db = db.Where("name ILIKE ?", "%"+query+"%")
	}

	err := db.
		Order("created_at DESC").
		Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project models.Project) error {
	return r.db.WithContext(ctx).Save(&project).Error
}

func (r *projectRepository) Delete(ctx context.Context, projectId uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", projectId).Error
}
