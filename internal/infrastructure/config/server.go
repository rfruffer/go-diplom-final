package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	// Server settings
	Port int
	Host string

	// JWT settings
	JWTSigningKey   string
	JWTAccessTTL    time.Duration
	JWTRefreshTTL   time.Duration
	JWTIssuer       string

	// Database settings
	DatabaseDSN string
}

// LoadServerConfig загружает конфигурацию сервера из переменных окружения
func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		Port:          getEnvAsInt("SERVER_PORT", 8080),
		Host:          getEnv("SERVER_HOST", "localhost"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY", ""),
		JWTAccessTTL:  getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL: getEnvAsDuration("JWT_REFRESH_TTL", 24*time.Hour),
		JWTIssuer:     getEnv("JWT_ISSUER", "gophkeeper"),
		DatabaseDSN:   getEnv("DATABASE_DSN", ""),
	}

	// Валидация обязательных параметров
	if cfg.JWTSigningKey == "" {
		return nil, fmt.Errorf("JWT_SIGNING_KEY environment variable is required")
	}

	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is required")
	}

	return cfg, nil
}

// DefaultServerConfig возвращает конфигурацию по умолчанию (для разработки)
// ВНИМАНИЕ: НЕ использовать в production!
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:          8080,
		Host:          "localhost",
		JWTSigningKey: "development-secret-key-change-in-production",
		JWTAccessTTL:  15 * time.Minute,
		JWTRefreshTTL: 24 * time.Hour,
		JWTIssuer:     "gophkeeper-dev",
		DatabaseDSN:   "host=localhost port=5432 user=postgres password=postgres dbname=gophkeeper sslmode=disable",
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDuration получает значение переменной окружения как time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
