package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	JWTSecret              string
	JWTExpirationHours     int
	ServerPort             string
	DatabaseURL            string
	Environment            string
	RateLimitEnabled       bool
	RateLimitPerMinute     int
	CORSAllowedOrigins     []string
}

func Load() *Config {
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" || jwtSecret == "your-secret-key-change-in-production" {
		log.Fatal("JWT_SECRET must be set to a secure random value in production")
	}

	environment := getEnv("ENVIRONMENT", "development")
	
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	allowedOrigins := strings.Split(corsOrigins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return &Config{
		JWTSecret:          jwtSecret,
		JWTExpirationHours: 24 * 7,
		ServerPort:         getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		Environment:        environment,
		RateLimitEnabled:   getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
		CORSAllowedOrigins: allowedOrigins,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return parsed
	}
	return defaultValue
}
