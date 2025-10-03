// This file defines the configuration for the application.
package config

// "log" provides a simple logging package. It is used here to log messages related to configuration.
import (
	"log"
	// "os" provides a platform-independent interface to operating system functionality. It is used here to access environment variables.
	"os"
	// "strconv" provides functions for converting strings to other types. It is used here to convert the database port and JWT expiry to integers.
	"strconv"
	// "time" provides functions for working with time. It is used here to set the JWT expiration duration.
	"time"

	// "github.com/joho/godotenv" is a package for loading environment variables from a .env file.
	"github.com/joho/godotenv"
)

// ServerConfig defines the structure for server-related configuration.
type ServerConfig struct {
	// Port is the port on which the server will listen.
	Port string
	// Host is the host of the server.
	Host string
}

// DatabaseConfig defines the structure for database-related configuration.
type DatabaseConfig struct {
	// DBHost is the host of the database.
	DBHost string
	// DBPort is the port of the database.
	DBPort int
	// DBUser is the username for the database.
	DBUser string
	// DBPassword is the password for the database.
	DBPassword string
	// DBName is the name of the database.
	DBName string
	// DBSSLMode is the SSL mode for the database connection.
	DBSSLMode string
}

// JWTConfig defines the structure for JWT-related configuration.
type JWTConfig struct {
	// SecretKey is the secret key used for signing JWTs.
	SecretKey string
	// Expires is the duration for which a JWT is valid.
	Expires time.Duration
}

// CORSConfig defines the structure for CORS-related configuration.
type CORSConfig struct {
	// CorsOrigins is a comma-separated list of allowed origins for CORS requests.
	CorsOrigins string
}

// Config is the main configuration struct that aggregates all other configuration types.
type Config struct {
	// Environment is the environment in which the application is running.
	Environment string
	// Server holds the server-specific configuration.
	Server ServerConfig
	// Database holds the database-specific configuration.
	Database DatabaseConfig
	// JWT holds the JWT-specific configuration.
	JWT JWTConfig
	// CORS holds the CORS-specific configuration.
	CORS CORSConfig
}

// HandleMissingEnvValues retrieves the value of an environment variable or returns a default value if it is not set.
// It takes the name of the environment variable and a default value as input.
//
// @param envName string - The name of the environment variable.
// @param defaultValue string - The default value to be returned if the environment variable is not set.
// @return string - The value of the environment variable or the default value.
func HandleMissingEnvValues(envName string, defaultValue string) string {
	// envValue is the value of the environment variable.
	envValue := os.Getenv(envName)
	// This checks if the environment variable is empty.
	if envValue == "" {
		// If the environment variable is empty, a warning is logged.
		log.Printf("%s is missing, default value is set.", envName)
		// The default value is returned.
		return defaultValue
	}
	// The value of the environment variable is returned.
	return envValue
}

// LoadConfig loads the application configuration from environment variables.
// It returns a pointer to a Config struct.
//
// @return *Config - A pointer to the Config struct.
func LoadConfig() *Config {
	// err is the result of loading the .env file.
	err := godotenv.Load()
	// This checks if an error occurred while loading the .env file.
	if err != nil {
		// If an error occurs, a fatal error is logged.
		log.Fatalf("Error loading .env file: %v", err)
	}

	// dbPort is the port of the database.
	dbPort, err := strconv.Atoi(HandleMissingEnvValues("DB_PORT", "5432"))
	// This checks if an error occurred while converting the database port to an integer.
	if err != nil {
		// If an error occurs, a fatal error is logged.
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	// expiry is the JWT expiration duration in hours.
	expiry, err := strconv.Atoi(HandleMissingEnvValues("JWT_EXPIRY_HOURS", "24"))
	// This checks if an error occurred while converting the JWT expiry to an integer.
	if err != nil {
		// If an error occurs, a fatal error is logged.
		log.Fatalf("Error parsing JWT_EXPIRY_HOURS: %v", err)
	}

	// A pointer to a new Config struct is returned.
	return &Config{
		// The Environment field is set to the value of the "ENV" environment variable, or "dev" if it is not set.
		Environment: HandleMissingEnvValues("ENV", "dev"),
		// The Server field is populated with the server configuration.
		Server: ServerConfig{
			// The Port field is set to the value of the "PORT" environment variable, or "8000" if it is not set.
			Port: HandleMissingEnvValues("PORT", "8000"),
			// The Host field is set to the value of the "HOST" environment variable, or "localhost" if it is not set.
			Host: HandleMissingEnvValues("HOST", "localhost"),
		},
		// The Database field is populated with the database configuration.
		Database: DatabaseConfig{
			// The DBHost field is set to the value of the "DB_HOST" environment variable, or "localhost" if it is not set.
			DBHost: HandleMissingEnvValues("DB_HOST", "localhost"),
			// The DBPort field is set to the value of the dbPort variable.
			DBPort: dbPort,
			// The DBUser field is set to the value of the "DB_USER" environment variable, or "postgres" if it is not set.
			DBUser: HandleMissingEnvValues("DB_USER", "postgres"),
			// The DBPassword field is set to the value of the "DB_PASSWORD" environment variable, or "postgres" if it is not set.
			DBPassword: HandleMissingEnvValues("DB_PASSWORD", "postgres"),
			// The DBName field is set to the value of the "DB_NAME" environment variable, or "postgres" if it is not set.
			DBName:    HandleMissingEnvValues("DB_NAME", "postgres"), // The DBSSLMode field is set to the value of the `DB_SSLMODE` environment variable, or `disable` if it is not set.
			DBSSLMode: HandleMissingEnvValues("DB_SSLMODE", "disable"),
		},
		// The JWT field is populated with the JWT configuration.
		JWT: JWTConfig{
			// The SecretKey field is set to the value of the "JWT_SECRET_KEY" environment variable, or a default value if it is not set.
			SecretKey: HandleMissingEnvValues("JWT_SECRET_KEY", "vCYKhw6zTyXIt7ckaKNnv7KarP2wzhZegyoxLLiK6MGKTnVo9z"),
			// The Expires field is set to the JWT expiration duration.
			Expires: time.Hour * time.Duration(expiry),
		},
		// The CORS field is populated with the CORS configuration.
		CORS: CORSConfig{
			// The CorsOrigins field is set to the value of the "CORS_ORIGINS" environment variable, or "http://localhost:3000" if it is not set.
			CorsOrigins: HandleMissingEnvValues("CORS_ORIGINS", "http://localhost:3000"),
		},
	}
}
