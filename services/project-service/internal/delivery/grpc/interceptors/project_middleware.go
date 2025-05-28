package interceptor

import (
	"context"
	"errors"
	"project-service/internal/delivery/grpc/auth"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDKey = contextKey("user_id")

type AuthResponse struct {
	UserID string `json:"user_id"`
}

var publicMethods = map[string]bool{
	"/project.ProjectService/Create":            true,
	"/project.ProjectService/GetProject":        true,
	"/project.ProjectService/GetProjectsByUser": true,
	"/project.ProjectService/Update":            true,
	"/project.ProjectService/Delete":            true,
}

func AuthInterceptor(cli *auth.UserServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := cli.GetClaims(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, userIDKey, userID)

		return handler(ctx, req)
	}
}

func UserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("user_id not found in context")
	}
	return userID, nil
}
