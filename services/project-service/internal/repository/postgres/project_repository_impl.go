package postgres

import (
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

func (r *projectRepository) Add(project models.Project) error {
	return r.db.Create(&project).Error
}

func (r *projectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByOwner(owner uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "owner_id = ?", owner).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetAll() ([]*models.Project, error) {
	var projects []*models.Project
	if err := r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) Update(project models.Project) error {
	return r.db.Save(&project).Error
}
