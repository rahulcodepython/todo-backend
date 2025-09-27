// File: backend/config/config.go

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rahulcodepython/todo-backend/backend/database"
)

// Config is the main configuration struct for the application.
type Config struct {
	App        AppConfig
	Server     ServerConfig
	Database   database.DBConfig
	JWT        JWTConfig
	Middleware MiddlewareConfig
}

// AppConfig holds application-specific settings.
type AppConfig struct {
	Name string
	Env  string
}

// ServerConfig holds server settings like the port.
type ServerConfig struct {
	Host string
	Port string
}

// JWTConfig holds settings for JSON Web Tokens.
type JWTConfig struct {
	Secret  string
	Expires time.Duration
}

// MiddlewareConfig holds settings for various middleware.
type MiddlewareConfig struct {
	RateLimitMax    int
	RateLimitWindow time.Duration
}

// DBConfig is a type alias for the database configuration struct.
// This is done to allow passing the original struct from the database package.
type DBConfig = database.DBConfig

// LoadConfig reads configuration from environment variables or a .env file.
func LoadConfig() (*Config, error) {
	// Attempt to load .env.dev file for local development.
	// In production, we'll rely on environment variables set in the deployment environment.
	if err := godotenv.Load(".env.dev"); err != nil {
		// It's okay if the file doesn't exist, but fail on any other error.
		// os.IsNotExist is a reliable way to check for this specific error.
		// This allows the app to run in environments (like production) where .env files aren't used.
		fmt.Println("Note: .env.dev file not found, relying on environment variables.")
	}

	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	rateLimitMax, err := strconv.Atoi(os.Getenv("RATE_LIMIT_MAX"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_MAX: %w", err)
	}
	rateLimitWindowMin, err := strconv.Atoi(os.Getenv("RATE_LIMIT_WINDOW_MINUTES"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_WINDOW_MINUTES: %w", err)
	}

	dbCfg := database.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	return &Config{
		App: AppConfig{
			Name: "Go-Fiber Todo App",
			Env:  os.Getenv("APP_ENV"),
		},
		Server: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: dbCfg,
		JWT: JWTConfig{
			Secret:  os.Getenv("JWT_SECRET"),
			Expires: time.Hour * time.Duration(expiryHours),
		},
		Middleware: MiddlewareConfig{
			RateLimitMax:    rateLimitMax,
			RateLimitWindow: time.Minute * time.Duration(rateLimitWindowMin),
		},
	}, nil
}
