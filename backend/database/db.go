package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/rahulcodepython/todo-backend/backend/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func PingDB(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Println("Unable to ping database")
		log.Fatal(err)
	}

	log.Println("Database is healthy.")
}

func createTable(db *sql.DB) {
	var query string

	// JWT Token Table
	query = `
		CREATE TABLE IF NOT EXISTS jwt_tokens (
		id UUID PRIMARY KEY,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Unable to create jwt token table")
		log.Fatal(err)
	}
	log.Println("jwt_tokens table created successfully.")

	// User Table
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
	_, err = db.Exec(query)
	if err != nil {
		log.Println("Unable to create user table")
		log.Fatal(err)
	}
	log.Println("users table created successfully.")
}

func ConnectDB(cfg *config.Config) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBName, cfg.Database.DBSSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Unable to connect with database")
		log.Fatal(err)
	}

	PingDB(db)
	createTable(db)

	return db
}
