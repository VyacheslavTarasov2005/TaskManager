package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validatable interface {
	Validate() error
}

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if v, ok := req.(Validatable); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
}
