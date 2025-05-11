package migrations

import (
	"gorm.io/gorm"
	"user-service/internal/domain/models"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
