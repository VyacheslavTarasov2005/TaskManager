package main

import (
	"context"
	"fmt"
	"log"
	"user-service/config"
	"user-service/internal/delivery/grpc"
	"user-service/internal/repository/postgres"
	"user-service/internal/repository/redis"
	"user-service/internal/service/implementations"
	"user-service/migrations"
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

	redisClient, err := redis.NewRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}

	if _, err = redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	refreshTokenRepository := redis.NewRefreshTokenRepositoryImpl(redisClient)
	userRepository := postgres.NewUserRepositoryImpl(dbConn)

	authService := implementations.NewAuthServiceImpl(cfg.JwtSecretKey, cfg.AccessTokenTTL, cfg.RefreshTokenTTL,
		refreshTokenRepository)
	userService := implementations.NewUserServiceImpl(userRepository, authService, refreshTokenRepository)

	grpcServer := grpc.SetupServer(userService, authService)

	if err = grpc.StartGRPCServer(grpcServer, "50051"); err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
	}
}
