package interfaces

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	CreateToken(ctx context.Context, userID uuid.UUID) (*string, *uuid.UUID, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
	RefreshToken(ctx context.Context, refreshToken uuid.UUID) (*string, *uuid.UUID, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}
