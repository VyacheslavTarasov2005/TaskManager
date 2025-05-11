package models

import "github.com/google/uuid"

type RefreshToken struct {
	UserID uuid.UUID
	Token  uuid.UUID
}
