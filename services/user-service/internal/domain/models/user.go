package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
	IsDeleted bool   `gorm:"not null"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Password  string `gorm:"not null"`
}
