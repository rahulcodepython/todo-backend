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

type CORSConfig struct {
	CorsOrigins string
}

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	CORS        CORSConfig
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
		// Log a fatal error if the DB_PORT environment variable cannot be converted to an integer.
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	// Parse the JWT_EXPIRY_HOURS environment variable into an integer.
	// This value determines how long JWTs will be valid.
	expiry, err := strconv.Atoi(HandleMissingEnvValues("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		// Log a fatal error if the JWT_EXPIRY_HOURS environment variable cannot be converted to an integer.
		log.Fatalf("Error parsing JWT_EXPIRY_HOURS: %v", err)
	}

	return &Config{
		Environment: HandleMissingEnvValues("ENV", "dev"),
		Server: ServerConfig{
			// Set the server port, using "8000" as default if "PORT" is not specified.
			Port: HandleMissingEnvValues("PORT", "8000"),
			// Set the server host, using "localhost" as default if "HOST" is not specified.
			Host: HandleMissingEnvValues("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			// Set the database host, using "localhost" as default.
			DBHost: HandleMissingEnvValues("DB_HOST", "localhost"),
			// Set the database port, using the parsed integer value or its default.
			DBPort: dbPort,
			// Set the database user, using "postgres" as default.
			DBUser: HandleMissingEnvValues("DB_USER", "postgres"),
			// Set the database password, using "postgres" as default.
			DBPassword: HandleMissingEnvValues("DB_PASSWORD", "postgres"),
			// Set the database name, using "postgres" as default.
			DBName: HandleMissingEnvValues("DB_NAME", "postgres"),
			// Set the database SSL mode, using "disable" as default.
			DBSSLMode: HandleMissingEnvValues("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			// Set the JWT secret key, using a strong default if "JWT_SECRET_KEY" is not specified.
			SecretKey: HandleMissingEnvValues("JWT_SECRET_KEY", "vCYKhw6zTyXIt7ckaKNnv7KarP2wzhZegyoxLLiK6MGKTnVo9z"),
			// Calculate the JWT expiration duration based on the parsed hours.
			Expires: time.Hour * time.Duration(expiry),
		},
		CORS: CORSConfig{
			// Set the allowed CORS origins, using "http://localhost:3000" as default.
			CorsOrigins: HandleMissingEnvValues("CORS_ORIGINS", "http://localhost:3000"),
		},
	}
}
