package implementations

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
	"user-service/internal/domain/interfaces"
	"user-service/internal/domain/models"
	"user-service/internal/service/errors"
	serviceInterfaces "user-service/internal/service/interfaces"
	"user-service/pkg/utils"
)

type UserServiceImpl struct {
	userRepository         interfaces.UserRepository
	authService            serviceInterfaces.AuthService
	refreshTokenRepository interfaces.RefreshTokenRepository
}

func NewUserServiceImpl(
	userRepository interfaces.UserRepository,
	authService serviceInterfaces.AuthService,
	refreshTokenRepository interfaces.RefreshTokenRepository,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository:         userRepository,
		authService:            authService,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, name, email, password string) (*string, *uuid.UUID, error) {
	sameEmailUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if sameEmailUser != nil {
		return nil, nil, errors.ApplicationError{
			StatusCode: 409,
			Code:       "EmailConflict",
			Errors: map[string]string{
				"message": "Email address is already in use",
			},
		}
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user := models.NewUser(name, email, hashedPassword)

	accessToken, refreshToken, err := s.authService.CreateToken(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	if err := s.userRepository.Add(ctx, *user); err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (*string, *uuid.UUID, error) {
	invalidCredentialsError := errors.ApplicationError{
		StatusCode: 401,
		Code:       "InvalidCredentials",
		Errors: map[string]string{
			"message": "Invalid email or password",
		},
	}
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, invalidCredentialsError
	}

	if err := comparePasswords(user.Password, password); err != nil {
		return nil, nil, invalidCredentialsError
	}

	if user.IsDeleted {
		return nil, nil, errors.ApplicationError{
			StatusCode: 403,
			Code:       "DeletedUser",
			Errors: map[string]string{
				"message": "User is deleted",
			},
		}
	}

	accessToken, refreshToken, err := s.authService.CreateToken(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors: map[string]string{
				"message": "User not found",
			},
		}
	}

	if user.IsDeleted {
		return &models.User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsDeleted: user.IsDeleted,
			Name:      "DELETED",
			Email:     "DELETED",
			Password:  "DELETED",
		}, nil
	}

	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) (*models.User,
	error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors: map[string]string{
				"message": "User not found",
			},
		}
	}

	if user.IsDeleted {
		return nil, errors.ApplicationError{
			StatusCode: 403,
			Code:       "DeletedUser",
			Errors: map[string]string{
				"message": "User is deleted",
			},
		}
	}

	sameEmailUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if sameEmailUser != nil && sameEmailUser.ID != user.ID {
		return nil, errors.ApplicationError{
			StatusCode: 409,
			Code:       "EmailConflict",
			Errors: map[string]string{
				"message": "Email address is already in use",
			},
		}
	}

	user.Name = name
	user.Email = email
	user.UpdatedAt = utils.Ptr(time.Now())

	err = s.userRepository.Update(ctx, *user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors: map[string]string{
				"message": "User not found",
			},
		}
	}

	if user.IsDeleted {
		return errors.ApplicationError{
			StatusCode: 403,
			Code:       "DeletedUser",
			Errors: map[string]string{
				"message": "User is deleted",
			},
		}
	}

	if err = comparePasswords(user.Password, oldPassword); err != nil {
		return errors.ApplicationError{
			StatusCode: 403,
			Code:       "IncorrectPassword",
			Errors: map[string]string{
				"message": "Incorrect old password",
			},
		}
	}

	if err = s.refreshTokenRepository.DeleteAllByUserID(ctx, userID); err != nil {
		return err
	}

	user.Password = newPassword
	user.UpdatedAt = utils.Ptr(time.Now())

	err = s.userRepository.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors: map[string]string{
				"message": "User not found",
			},
		}
	}

	if user.IsDeleted {
		return errors.ApplicationError{
			StatusCode: 409,
			Code:       "DeletedUser",
			Errors: map[string]string{
				"message": "User is deleted",
			},
		}
	}

	if err = s.refreshTokenRepository.DeleteAllByUserID(ctx, userID); err != nil {
		return err
	}

	user.IsDeleted = true
	user.UpdatedAt = utils.Ptr(time.Now())

	if err := s.userRepository.Update(ctx, *user); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePasswords(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
