// This file provides functions for connecting to and initializing the database.
package database

// "database/sql" provides a generic SQL interface. It is used here to interact with the database.
import (
	"database/sql"
	// "fmt" provides functions for formatted I/O. It is used here to construct the database connection string.
	"fmt"
	// "log" provides a simple logging package. It is used here to log database-related messages.
	"log"

	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"

	// _ "github.com/lib/pq" is the PostgreSQL driver. The underscore indicates that it is imported for its side effects (registering the driver).
	_ "github.com/lib/pq"
)

// PingDB checks if the database connection is alive.
// It takes a database connection as input.
//
// @param db *sql.DB - The database connection.
func PingDB(db *sql.DB) {
	// db.Ping() verifies a connection to the database is still alive, establishing a connection if necessary.
	if err := db.Ping(); err != nil {
		// If the ping fails, a message is logged.
		log.Println("Unable to ping database")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}

	// If the ping is successful, a success message is logged.
	log.Println("Database is healthy.")
}

// createTable creates the necessary tables in the database if they do not already exist.
// It takes a database connection as input.
//
// @param db *sql.DB - The database connection.
func createTable(db *sql.DB) {
	// query is a variable that will hold the SQL query.
	var query string

	// This is the SQL query to create the jwt_tokens table.
	query = `
		CREATE TABLE IF NOT EXISTS jwt_tokens (
		id UUID PRIMARY KEY,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	// db.Exec() executes a query without returning any rows.
	_, err := db.Exec(query)
	// This checks if an error occurred while creating the table.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to create jwt token table")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}
	// A success message is logged after the table is created.
	log.Println("jwt_tokens table created successfully.")

	// This is the SQL query to create the users table.
	query = `
		CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		image TEXT,
		password TEXT NOT NULL,
		jwt UUID UNIQUE,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		CONSTRAINT fk_jwt
			FOREIGN KEY(jwt)
			REFERENCES jwt_tokens(id)
			ON DELETE SET NULL
		);
	`
	// db.Exec() executes a query without returning any rows.
	_, err = db.Exec(query)
	// This checks if an error occurred while creating the table.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to create user table")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}
	// A success message is logged after the table is created.
	log.Println("users table created successfully.")

	// This is the SQL query to create the todos table.
	query = `
		CREATE TABLE IF NOT EXISTS todos (
		id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		owner UUID NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT fk_owner
			FOREIGN KEY(owner)
			REFERENCES users(id)
			ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(owner);
		`
	// db.Exec() executes a query without returning any rows.
	_, err = db.Exec(query)
	// This checks if an error occurred while creating the table.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to create todos table")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}
	// A success message is logged after the table is created.
	log.Println("todos table created successfully.")
}

// ConnectDB establishes a connection to the database.
// It takes the application configuration as input and returns a database connection.
//
// @param cfg *config.Config - The application configuration.
// @return *sql.DB - The database connection.
func ConnectDB(cfg *config.Config) *sql.DB {
	// connectionString is the connection string for the database.
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBName, cfg.Database.DBSSLMode)

	// db is the database connection.
	// sql.Open() opens a database specified by its database driver name and a driver-specific data source name.
	db, err := sql.Open("postgres", connectionString)
	// This checks if an error occurred while opening the database connection.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to connect with database")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}

	// PingDB() is called to check if the database connection is alive.
	PingDB(db)
	// createTable() is called to create the necessary tables in the database.
	createTable(db)

	// The database connection is returned.
	return db
}