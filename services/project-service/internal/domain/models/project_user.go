package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProject struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	ProjectID uuid.UUID `gorm:"primaryKey"`
	AddedTime time.Time `gorm:"not null"`
}

func NewUserProject(user_id, project_id uuid.UUID) *UserProject {
	return &UserProject{
		UserID:    user_id,
		ProjectID: project_id,
		AddedTime: time.Now(),
	}
}
