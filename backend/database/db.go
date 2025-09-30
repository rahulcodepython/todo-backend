package database

import (
	// Import the "database/sql" package, which provides a generic interface for working with SQL databases in Go.
	"database/sql"
	// Import the "fmt" package for formatted I/O, used here for constructing the database connection string.
	"fmt"
	// Import the "log" package for logging messages, especially critical errors during database operations.
	"log"

	// Import the application's "config" package to access database configuration settings.
	"github.com/rahulcodepython/todo-backend/backend/config"

	// Import the PostgreSQL driver. The underscore `_` indicates that the package is imported for its side effects (registering itself with `database/sql`),
	// rather than for direct use of its exported functions or types.
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PingDB checks the connectivity to the database.
// It attempts to ping the database and logs a fatal error if the connection fails,
// otherwise, it logs a success message. This function is crucial for verifying database availability.
func PingDB(db *sql.DB) {
	// Attempt to ping the database to verify the connection is alive and responsive.
	if err := db.Ping(); err != nil {
		// If an error occurs during the ping, log a message indicating the failure.
		log.Println("Unable to ping database")
		// Terminate the application with a fatal error, as database connectivity is often a critical dependency.
		log.Fatal(err)
	}

	// If the ping is successful, log a message indicating that the database is healthy.
	log.Println("Database is healthy.")
}

// createTable initializes the necessary database tables if they do not already exist.
// This function ensures that the application's schema is present and correctly structured upon startup.
func createTable(db *sql.DB) {
	var query string

	// JWT Token Table
	// Define the SQL query to create the 'jwt_tokens' table.
	// This table stores JSON Web Tokens, typically for session management or blacklisting.
	query = `
		CREATE TABLE IF NOT EXISTS jwt_tokens (
		id UUID PRIMARY KEY, -- 'id' column: A universally unique identifier, serving as the primary key for this table.
		token TEXT NOT NULL UNIQUE, -- 'token' column: Stores the actual JWT string, must be unique and not null.
		expires_at TIMESTAMPTZ NOT NULL, -- 'expires_at' column: Stores the timestamp when the JWT expires, ensuring tokens are not used indefinitely.
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() -- 'created_at' column: Records the timestamp when the token was created, defaults to the current time.
		);
	`
	// Execute the SQL query to create the 'jwt_tokens' table.
	_, err := db.Exec(query)
	if err != nil {
		// If an error occurs during table creation, log a message indicating the failure.
		log.Println("Unable to create jwt token table")
		// Terminate the application with a fatal error, as schema initialization is critical.
		log.Fatal(err)
	}
	// Log a success message after the 'jwt_tokens' table has been created or verified.
	log.Println("jwt_tokens table created successfully.")

	// User Table
	// Define the SQL query to create the 'users' table.
	// This table stores user-related information, including authentication details.
	query = `
		CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY, -- 'id' column: A universally unique identifier for the user, serving as the primary key.
		name TEXT NOT NULL, -- 'name' column: Stores the user's full name, cannot be null.
		email TEXT NOT NULL UNIQUE, -- 'email' column: Stores the user's email address, must be unique and not null, used for login.
		image TEXT, -- 'image' column: Stores a URL or path to the user's profile image, can be null.
		password TEXT NOT NULL, -- 'password' column: Stores the hashed password of the user, cannot be null for security.
		jwt UUID UNIQUE, -- 'jwt' column: Stores a reference to a JWT token ID, allowing for a one-to-one relationship with 'jwt_tokens' table.
		created_at TIMESTAMPTZ NOT NULL, -- 'created_at' column: Records the timestamp when the user account was created.
		updated_at TIMESTAMPTZ NOT NULL, -- 'updated_at' column: Records the timestamp of the last update to the user's account.
		CONSTRAINT fk_jwt -- Define a foreign key constraint named 'fk_jwt'.
			FOREIGN KEY(jwt) -- The 'jwt' column in the 'users' table is the foreign key.
			REFERENCES jwt_tokens(id) -- It references the 'id' column in the 'jwt_tokens' table.
			ON DELETE SET NULL -- If a referenced JWT token is deleted, the 'jwt' column in the 'users' table will be set to NULL.
	);
	`
	// Execute the SQL query to create the 'users' table.
	_, err = db.Exec(query)
	if err != nil {
		// If an error occurs during table creation, log a message indicating the failure.
		log.Println("Unable to create user table")
		// Terminate the application with a fatal error, as schema initialization is critical.
		log.Fatal(err)
	}
	// Log a success message after the 'users' table has been created or verified.
	log.Println("users table created successfully.")
}

// ConnectDB establishes a connection to the PostgreSQL database using the provided configuration.
// It constructs a connection string, opens the database connection, pings it to verify connectivity,
// and ensures that necessary tables are created.
//
// Parameters:
// - cfg: A pointer to the application's configuration struct, containing database connection details.
//
// Returns:
// - *sql.DB: A pointer to the established database connection pool.
func ConnectDB(cfg *config.Config) *sql.DB {
	// Construct the database connection string using parameters from the configuration.
	// This string includes host, port, user, password, database name, and SSL mode.
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBName, cfg.Database.DBSSLMode)

	// Open a new database connection using the "postgres" driver and the constructed connection string.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// If an error occurs during connection opening, log a message indicating the failure.
		log.Println("Unable to connect with database")
		// Terminate the application with a fatal error, as database connectivity is essential.
		log.Fatal(err)
	}

	// Ping the database to ensure the connection is active and healthy.
	PingDB(db)
	// Create necessary tables if they don't exist, ensuring the database schema is ready.
	createTable(db)

	// Return the established database connection pool.
	return db
}
