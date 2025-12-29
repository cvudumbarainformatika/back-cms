package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
	Redis     RedisConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name string
	Env  string
	Port string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Connection      string
	Host            string
	Port            string
	Database        string
	Username        string
	Password        string
	MaxOpenConns    int // Maximum number of open connections to the database
	MaxIdleConns    int // Maximum number of idle connections in the pool
	ConnMaxLifetime int // Maximum lifetime of a connection in seconds
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret                 string
	AccessTokenExpiration  int // in minutes
	RefreshTokenExpiration int // in minutes
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled       bool
	MaxRequests   int
	WindowSeconds int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "Go Gin Starter Kit"),
			Env:  getEnv("APP_ENV", "local"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Connection:      getEnv("DB_CONNECTION", "mysql"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			Database:        getEnv("DB_DATABASE", ""),
			Username:        getEnv("DB_USERNAME", ""),
			Password:        getEnv("DB_PASSWORD", ""),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", ""),
			AccessTokenExpiration:  getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRATION", 15),
			RefreshTokenExpiration: getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRATION", 10080),
		},
		RateLimit: RateLimitConfig{
			Enabled:       getEnvAsBool("RATE_LIMIT_ENABLED", true),
			MaxRequests:   getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 100),
			WindowSeconds: getEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 60),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	if c.Database.Database == "" {
		return fmt.Errorf("DB_DATABASE is required")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("DB_USERNAME is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool retrieves an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsSlice retrieves an environment variable as a slice (comma-separated) or returns a default value
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
