package interfaces

import (
	"context"
	"github.com/google/uuid"
	"user-service/internal/domain/models"
)

type UserService interface {
	Register(ctx context.Context, name, email, password string) (*string, *uuid.UUID, error)
	Login(ctx context.Context, email, password string) (*string, *uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	RecoverAccount(ctx context.Context, email, password string) (*string, *uuid.UUID, error)
}
