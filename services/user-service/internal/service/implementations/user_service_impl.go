package implementations

import (
	"context"
	"github.com/google/uuid"
	"user-service/internal/domain/interfaces"
	"user-service/internal/domain/models"
	"user-service/internal/service/errors"
	serviceInterfaces "user-service/internal/service/interfaces"
)

type UserServiceImpl struct {
	userRepository interfaces.UserRepository
	authService    serviceInterfaces.AuthService
}

func NewUserServiceImpl(
	userRepository interfaces.UserRepository,
	authService serviceInterfaces.AuthService,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		authService:    authService,
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

	user := models.NewUser(name, email, password)

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

	if user.Password != password {
		return nil, nil, invalidCredentialsError
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

	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) error {
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

	sameEmailUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if sameEmailUser != nil && sameEmailUser.ID != user.ID {
		return errors.ApplicationError{
			StatusCode: 409,
			Code:       "EmailConflict",
			Errors: map[string]string{
				"message": "Email address is already in use",
			},
		}
	}

	user.Name = name
	user.Email = email
	err = s.userRepository.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
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

	if user.Password != oldPassword {
		return errors.ApplicationError{
			StatusCode: 403,
			Code:       "IncorrectPassword",
			Errors: map[string]string{
				"message": "Incorrect old password",
			},
		}
	}

	user.Password = newPassword

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

	user.IsDeleted = true
	if err := s.userRepository.Update(ctx, *user); err != nil {
		return err
	}

	return nil
}
