package config

import (
	"os"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	RedisAddress     string
	RedisPassword    string
	RedisDB          string
	JwtSecretKey     []byte
}

func LoadConfig() *Config {
	config := &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "123456"),
		PostgresDBName:   getEnv("POSTGRES_DB_NAME", "UsersDB"),
		RedisAddress:     getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnv("REDIS_DB", "0"),
		JwtSecretKey:     []byte(getEnv("JWT_SECRET_KEY", "super_secret_key")),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
