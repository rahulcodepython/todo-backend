package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

type JWTConfig struct {
	SecretKey string
	Expires   time.Duration
}

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
}

func HandleMissingEnvValues(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		log.Printf("%s is missing, default value is set.", envName)
		return defaultValue
	}
	return envValue
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbPort, err := strconv.Atoi(HandleMissingEnvValues("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	expiry, err := strconv.Atoi(HandleMissingEnvValues("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		log.Fatalf("Error parsing JWT_EXPIRY_HOURS: %v", err)
	}

	return &Config{
		Environment: HandleMissingEnvValues("ENV", "dev"),
		Server: ServerConfig{
			Port: HandleMissingEnvValues("PORT", "8000"),
			Host: HandleMissingEnvValues("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			DBHost:     HandleMissingEnvValues("DB_HOST", "localhost"),
			DBPort:     dbPort,
			DBUser:     HandleMissingEnvValues("DB_USER", "postgres"),
			DBPassword: HandleMissingEnvValues("DB_PASSWORD", "postgres"),
			DBName:     HandleMissingEnvValues("DB_NAME", "postgres"),
			DBSSLMode:  HandleMissingEnvValues("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: HandleMissingEnvValues("JWT_SECRET_KEY", "vCYKhw6zTyXIt7ckaKNnv7KarP2wzhZegyoxLLiK6MGKTnVo9z"),
			Expires:   time.Hour * time.Duration(expiry),
		},
	}
}
