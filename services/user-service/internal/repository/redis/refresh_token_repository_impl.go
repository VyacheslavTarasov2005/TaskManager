package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
	"user-service/internal/domain/models"
)

type RefreshTokenRepositoryImpl struct {
	client *redis.Client
	ttl    time.Duration
}

const (
	tokenKeyPrefix      = "token:"
	userTokensKeyPrefix = "user_tokens:"
)

func NewRefreshTokenRepositoryImpl(client *redis.Client) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{
		client: client,
	}
}

func (r *RefreshTokenRepositoryImpl) Add(ctx context.Context, token models.RefreshToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	pipe := r.client.TxPipeline()

	pipe.Set(ctx, r.tokenKey(token.Token), data, token.ExpiresIn)

	pipe.SAdd(ctx, r.userTokensKey(token.UserID), token.Token.String())
	pipe.Expire(ctx, r.userTokensKey(token.UserID), token.ExpiresIn)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepositoryImpl) GetByToken(ctx context.Context, token uuid.UUID) (*models.RefreshToken, error) {
	data, err := r.client.Get(ctx, r.tokenKey(token)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var refreshToken models.RefreshToken
	if err := json.Unmarshal(data, &refreshToken); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	tokens, err := r.client.SMembers(ctx, r.userTokensKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	pipe := r.client.Pipeline()
	for _, tokenStr := range tokens {
		token := uuid.MustParse(tokenStr)
		pipe.Del(ctx, r.tokenKey(token))
	}

	pipe.Del(ctx, r.userTokensKey(userID))

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepositoryImpl) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	refreshToken, err := r.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if refreshToken == nil {
		return nil
	}

	pipe := r.client.TxPipeline()

	pipe.Del(ctx, r.tokenKey(token))

	pipe.SRem(ctx, r.userTokensKey(refreshToken.UserID), token.String())

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepositoryImpl) tokenKey(token uuid.UUID) string {
	return fmt.Sprintf("%s%s", tokenKeyPrefix, token.String())
}

func (r *RefreshTokenRepositoryImpl) userTokensKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", userTokensKeyPrefix, userID.String())
}
