package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"user-service/internal/service/interfaces"
)

func AuthInterceptor(authService interfaces.AuthService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		claims, err := authService.VerifyToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
		}

		newCtx := context.WithValue(ctx, "user_id", userID)

		return handler(newCtx, req)
	}
}

func isPublicMethod(method string) bool {
	publicMethods := map[string]bool{
		"/user.UserService/Register":       true,
		"/user.UserService/Login":          true,
		"/user.UserService/RefreshToken":   true,
		"/user.UserService/RecoverAccount": true,
	}
	return publicMethods[method]
}
