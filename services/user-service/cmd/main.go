package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"user-service/config"
	"user-service/internal/repository/postgres"
	"user-service/internal/repository/redis"
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

	redisDB, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to parse Redis DB: %v", err)
	}

	redisClient, err := redis.NewRedisClient(cfg.RedisAddress, cfg.RedisPassword, redisDB)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
}
