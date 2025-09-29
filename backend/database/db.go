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

func ConnectDB(cfg *config.Config) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBName, cfg.Database.DBSSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Unable to connect with database")
		log.Fatal(err)
	}

	PingDB(db)

	return db
}
