package interceptor

import (
	"context"
	"project-service/internal/delivery/grpc/auth"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthResponse struct {
	UserID string `json:"user_id"`
}

var publicMethods = map[string]bool{
	"/project.ProjectService/Create":            false,
	"/project.ProjectService/GetProject":        true,
	"/project.ProjectService/GetProjectsByUser": false,
	"/project.ProjectService/Update":            false,
	"/project.ProjectService/Delete":            false,
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

		newCtx := context.WithValue(ctx, "user_id", userID)

		return handler(newCtx, req)
	}
}
