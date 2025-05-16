package grpc

import (
	"google.golang.org/grpc"
	"net"
	"user-service/internal/delivery/grpc/interceptors"
	"user-service/internal/delivery/grpc/pb"
	"user-service/internal/service/interfaces"
)

func SetupServer(userService interfaces.UserService, authService interfaces.AuthService) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthInterceptor(authService),
			interceptors.ValidationInterceptor(),
		),
	)

	pb.RegisterUserServiceServer(srv, NewUserServer(userService, authService))

	return srv
}

func StartGRPCServer(srv *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return srv.Serve(lis)
}
