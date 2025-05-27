package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	RedisAddress     string
	RedisPassword    string
	RedisDB          int
	JwtSecretKey     []byte
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
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
		RedisDB:          getEnv("REDIS_DB", 0),
		JwtSecretKey:     []byte(getEnv("JWT_SECRET_KEY", "super_secret_key")),
		AccessTokenTTL:   getEnv("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:  getEnv("REFRESH_TOKEN_TTL", 7*24*time.Hour),
	}

	return config
}

func getEnv[T any](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result any
	var err error

	switch any(defaultValue).(type) {
	case string:
		result = value
	case int:
		result, err = strconv.Atoi(value)
	case time.Duration:
		result, err = time.ParseDuration(value)
	default:
		return defaultValue
	}

	if err != nil {
		return defaultValue
	}

	return result.(T)
}
