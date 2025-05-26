package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetProject struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUpdateProject struct {
	Name    string    `json:"name" binding:"required"`
	OwnerID uuid.UUID `json:"owner_id" binding:"required"`
}
