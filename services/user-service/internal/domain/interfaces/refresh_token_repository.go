package interfaces

import (
	"context"
	"github.com/google/uuid"
	"user-service/internal/domain/models"
)

type RefreshTokenRepository interface {
	Add(ctx context.Context, token models.RefreshToken) error
	GetByToken(ctx context.Context, token uuid.UUID) (*models.RefreshToken, error)
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByToken(ctx context.Context, token uuid.UUID) error
}
