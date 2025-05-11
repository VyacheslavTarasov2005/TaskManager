package cmd

import (
	"fmt"
	"log"
	"project-service/config"
	"project-service/internal/repository/postgres"
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
}
