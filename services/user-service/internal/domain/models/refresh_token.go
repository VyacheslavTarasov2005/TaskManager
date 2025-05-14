package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	UserID    uuid.UUID     `json:"userId"`
	Token     uuid.UUID     `json:"token"`
	CreatedAt time.Time     `json:"createdAt"`
	ExpiresIn time.Duration `json:"expiresIn"`
}

func NewRefreshToken(userID uuid.UUID) *RefreshToken {
	return &RefreshToken{
		UserID:    userID,
		Token:     uuid.New(),
		CreatedAt: time.Now(),
		ExpiresIn: 7 * 24 * time.Hour,
	}
}
