package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectUser struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	ProjectID uuid.UUID `gorm:"primaryKey"`
	Role      string    `gorm:"type:enum('OWNER','MEMBER');not null" json:"role"`
	AddedTime time.Time `gorm:"not null"`
}

func NewUserProject(user_id, project_id uuid.UUID, role string) *ProjectUser {
	return &ProjectUser{
		UserID:    user_id,
		ProjectID: project_id,
		Role:      role,
		AddedTime: time.Now(),
	}
}
