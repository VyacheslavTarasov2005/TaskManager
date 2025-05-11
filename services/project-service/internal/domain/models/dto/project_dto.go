package dto

import "time"

type GetProject struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uint      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUpdateProject struct {
	Name    string `json:"name" binding:"required"`
	OwnerID uint   `json:"owner_id" binding:"required"`
}
