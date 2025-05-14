package implementations

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"user-service/internal/domain/interfaces"
	"user-service/internal/domain/models"
	"user-service/internal/service/errors"
)

type AuthServiceImpl struct {
	jwtSecretKey           []byte
	ttl                    time.Duration
	refreshTokenRepository interfaces.RefreshTokenRepository
}

func NewAuthServiceImpl(jwtSecretKey []byte, ttl time.Duration) *AuthServiceImpl {
	return &AuthServiceImpl{
		jwtSecretKey: jwtSecretKey,
		ttl:          ttl,
	}
}

func (s *AuthServiceImpl) CreateToken(ctx context.Context, userID uuid.UUID) (*string, *uuid.UUID, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken := models.NewRefreshToken(userID)
	if err := s.refreshTokenRepository.Add(ctx, *refreshToken); err != nil {
		return nil, nil, err
	}

	return &tokenString, &refreshToken.Token, nil
}

func (s *AuthServiceImpl) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	errorResponse := errors.ApplicationError{
		StatusCode: 401,
		Code:       "Unauthorized",
		Errors: map[string]string{
			"message": "Invalid token",
		},
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorResponse
		}
		return s.jwtSecretKey, nil
	})

	if err != nil {
		return nil, errorResponse
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errorResponse
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken uuid.UUID) (*string, *uuid.UUID, error) {
	dbToken, err := s.refreshTokenRepository.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}
	if dbToken == nil {
		return nil, nil, errors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors: map[string]string{
				"message": "Invalid token",
			},
		}
	}

	if err := s.refreshTokenRepository.DeleteByToken(ctx, dbToken.Token); err != nil {
		return nil, nil, err
	}

	newAccessToken, newRefreshToken, err := s.CreateToken(ctx, dbToken.UserID)
	if err != nil {
		return nil, nil, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := s.refreshTokenRepository.DeleteAllByUserID(ctx, userID); err != nil {
		return err
	}

	return nil
}
