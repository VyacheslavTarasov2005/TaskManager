package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID  `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	OwnerID   uuid.UUID  `gorm:"not null" json:"owner_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func NewProject(name string, owner_id uuid.UUID) *Project {
	return &Project{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   owner_id,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
}
