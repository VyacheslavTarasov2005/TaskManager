package grpc

import (
	"net"
	interceptor "project-service/internal/delivery/grpc/interceptors"
	"project-service/internal/delivery/grpc/pb"
	"project-service/internal/service/interfaces"

	"google.golang.org/grpc"
)

func SetupServer(projectrService interfaces.ProjectService) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor()),
	)

	pb.RegisterProjectServiceServer(grpcServer, NewProjectServer(projectrService))

	return grpcServer
}

func StartGRPCServer(srv *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return srv.Serve(lis)
}
