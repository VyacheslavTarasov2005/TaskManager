package migrations

import (
	"project-service/internal/domain/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Project{})
}
