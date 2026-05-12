package config

import (
	"os"
)

type Config struct {
	JWTSecret          string
	JWTExpirationHours int
	ServerPort         string
	DatabaseURL        string
}

func Load() *Config {
	return &Config{
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpirationHours: 24 * 7, // 7 days
		ServerPort:         getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
