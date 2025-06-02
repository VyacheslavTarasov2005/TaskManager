package main

import (
	"fmt"
	"log"
	"project-service/config"
	"project-service/internal/delivery/grpc"
	"project-service/internal/delivery/grpc/auth"
	"project-service/internal/repository/postgres"
	"project-service/internal/service/implementations"
	"project-service/migrations"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := postgres.NewPostgresConnection(cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresDBName, cfg.PostgresPort)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer func() {
		sqlDB, _ := dbConn.DB()
		sqlDB.Close()
	}()
	fmt.Printf("User Service started with config: %+v\n", cfg)

	if err := migrations.Migrate(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	projectRepository := postgres.NewProjectRepositoryImpl(dbConn)
	projectUserRepository := postgres.NewProjectUserRepositoryImpl(dbConn)

	projectService := implementations.NewProjectServiceImpl(projectRepository)
	projectUserService := implementations.NewProjectUserServiceImpl(projectUserRepository)

	authClient, err := auth.NewUserServiceClient("localhost:50052")

	defer authClient.Close()

	grpcServer := grpc.SetupServer(projectService, projectUserService, authClient)

	if err = grpc.StartGRPCServer(grpcServer, "50051"); err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}
	fmt.Println("API started")
}
