package cmd

import (
	"fmt"
	"log"
	"project-service/config"
	"project-service/internal/delivery/grpc"
	"project-service/internal/repository/postgres"
	"project-service/internal/service/implementations"
	"project-service/migrations"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("User Service started with config: %+v\n", cfg)

	dbConn, err := postgres.NewPostgresConnection(cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresDBName, cfg.PostgresPort)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	if err := migrations.Migrate(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	projectRepository := postgres.NewProjectRepositoryImpl(dbConn)

	projectService := implementations.NewProjectServiceImpl(projectRepository)

	grpcServer := grpc.SetupServer(projectService)

	if err = grpc.StartGRPCServer(grpcServer, "50051"); err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}
}
