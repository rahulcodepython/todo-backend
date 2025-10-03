package config

import (
	"log"     // Import the "log" package for logging messages, especially fatal errors during configuration loading.
	"os"      // Import the "os" package to access environment variables.
	"strconv" // Import the "strconv" package to convert string values from environment variables to other types (e.g., int).
	"time"    // Import the "time" package to handle time-related operations, specifically for JWT expiration duration.

	"github.com/joho/godotenv" // Import the "godotenv" package to load environment variables from a .env file.
)

// ServerConfig defines the structure for server-related configuration settings.
type ServerConfig struct {
	// Port specifies the port number on which the server will listen for incoming requests.
	Port string
	// Host specifies the host address (e.g., IP address or domain name) on which the server will bind.
	Host string
}

// DatabaseConfig defines the structure for database connection configuration settings.
type DatabaseConfig struct {
	// DBHost specifies the hostname or IP address of the database server.
	DBHost string
	// DBPort specifies the port number on which the database server is listening.
	DBPort int
	// DBUser specifies the username used to authenticate with the database.
	DBUser string
	// DBPassword specifies the password used to authenticate with the database.
	DBPassword string
	// DBName specifies the name of the database to connect to.
	DBName string
	// DBSSLMode specifies the SSL mode for the database connection (e.g., "disable", "require").
	DBSSLMode string
}

// JWTConfig defines the structure for JSON Web Token (JWT) related configuration settings.
type JWTConfig struct {
	// SecretKey specifies the secret key used for signing and verifying JWTs.
	SecretKey string
	// Expires specifies the duration for which a generated JWT will be valid.
	Expires time.Duration
}

// CORSConfig defines the structure for Cross-Origin Resource Sharing (CORS) related configuration settings.
type CORSConfig struct {
	// CorsOrigins specifies a comma-separated list of allowed origins for CORS requests.
	CorsOrigins string
}

// Config is the main configuration struct that aggregates all other configuration types.
// It holds all the application-wide settings loaded from environment variables.
type Config struct {
	// Environment specifies the current operating environment of the application (e.g., "dev", "prod").
	Environment string
	// Server holds the server-specific configuration settings.
	Server ServerConfig
	// Database holds the database-specific configuration settings.
	Database DatabaseConfig
	// JWT holds the JWT-specific configuration settings.
	JWT JWTConfig
	// CORS holds the CORS-specific configuration settings.
	CORS CORSConfig
}

// HandleMissingEnvValues retrieves an environment variable's value or returns a default if it's not set.
// It also logs a warning if the environment variable is missing and a default is used.
//
// Parameters:
// - envName: The name of the environment variable to retrieve.
// - defaultValue: The default value to use if the environment variable is not found or is empty.
//
// Returns:
// - A string containing the environment variable's value or the provided default value.
func HandleMissingEnvValues(envName string, defaultValue string) string {
	// Get the value of the environment variable specified by `envName`.
	envValue := os.Getenv(envName)
	// Check if the retrieved environment variable value is empty.
	if envValue == "" {
		// If it's empty, log a message indicating that the environment variable is missing and a default is being used.
		log.Printf("%s is missing, default value is set.", envName)
		// Return the provided `defaultValue`.
		return defaultValue
	}
	// If the environment variable has a value, return it.
	return envValue
}

// LoadConfig loads all application configurations from environment variables, potentially from a .env file.
// It parses and validates these values, providing defaults where necessary, and returns a pointer to a Config struct.
//
// Returns:
// - *Config: A pointer to the fully populated Config struct.
func LoadConfig() *Config {
	// Load environment variables from a .env file. This allows for easy configuration management in development.
	err := godotenv.Load()
	// Check if there was an error loading the .env file.
	if err != nil {
		// If an error occurs, log a fatal message and exit the application, as configuration is critical.
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve the database port from environment variables, using "5432" as a default.
	dbPort, err := strconv.Atoi(HandleMissingEnvValues("DB_PORT", "5432"))
	// Check if there was an error converting the DB_PORT string to an integer.
	if err != nil {
		// Log a fatal error if the DB_PORT environment variable cannot be converted to an integer.
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	// Retrieve the JWT expiry hours from environment variables, using "24" as a default.
	// This value determines how long JWTs will be valid.
	expiry, err := strconv.Atoi(HandleMissingEnvValues("JWT_EXPIRY_HOURS", "24"))
	// Check if there was an error converting the JWT_EXPIRY_HOURS string to an integer.
	if err != nil {
		// Log a fatal error if the JWT_EXPIRY_HOURS environment variable cannot be converted to an integer.
		log.Fatalf("Error parsing JWT_EXPIRY_HOURS: %v", err)
	}

	// Return a new Config struct populated with values from environment variables or their defaults.
	return &Config{
		// Set the application environment, using "dev" as default if "ENV" is not specified.
		Environment: HandleMissingEnvValues("ENV", "dev"),
		// Populate the ServerConfig struct.
		Server: ServerConfig{
			// Set the server port, using "8000" as default if "PORT" is not specified.
			Port: HandleMissingEnvValues("PORT", "8000"),
			// Set the server host, using "localhost" as default if "HOST" is not specified.
			Host: HandleMissingEnvValues("HOST", "localhost"),
		},
		// Populate the DatabaseConfig struct.
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
		// Populate the JWTConfig struct.
		JWT: JWTConfig{
			// Set the JWT secret key, using a strong default if "JWT_SECRET_KEY" is not specified.
			SecretKey: HandleMissingEnvValues("JWT_SECRET_KEY", "vCYKhw6zTyXIt7ckaKNnv7KarP2wzhZegyoxLLiK6MGKTnVo9z"),
			// Calculate the JWT expiration duration based on the parsed hours.
			Expires: time.Hour * time.Duration(expiry),
		},
		// Populate the CORSConfig struct.
		CORS: CORSConfig{
			// Set the allowed CORS origins, using "http://localhost:3000" as default.
			CorsOrigins: HandleMissingEnvValues("CORS_ORIGINS", "http://localhost:3000"),
		},
	}
}
