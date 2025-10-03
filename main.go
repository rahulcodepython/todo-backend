package main

import (
	// The "fmt" package provides functions for formatted I/O, such as printing to the console.
	"fmt"
	// The "log" package implements a simple logging package, used here for reporting server errors.
	"log"
	// The "os" package provides a platform-independent interface to operating system functionality,
	// such as environment variables and process signals.
	"os"
	// The "os/signal" package implements access to incoming operating system signals,
	// allowing the application to respond to events like Ctrl+C.
	"os/signal"
	// The "syscall" package provides a low-level interface to operating system primitives,
	// used here to specify the SIGTERM signal for graceful shutdown.
	"syscall"

	// "github.com/gofiber/fiber/v2" is a fast, unopinionated, and flexible web framework for Go,
	// used to build the HTTP server and define API routes.
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/backend/config" handles loading application configurations
	// from environment variables or a configuration file, centralizing settings management.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/database" manages the database connection and
	// provides functions for interacting with the database.
	"github.com/rahulcodepython/todo-backend/backend/database"
	// "github.com/rahulcodepython/todo-backend/backend/router" is responsible for setting up
	// and registering all the application's API routes and middleware.
	"github.com/rahulcodepython/todo-backend/backend/router"
)

// main is the entry point of the application. Execution begins here.
func main() {
	// Load application configuration. This function reads environment variables and
	// potentially a .env file to populate the Config struct with settings for the server, database, etc.
	cfg := config.LoadConfig()

	// Establish a connection to the database using the loaded configuration.
	// This function typically returns a database connection pool or a single connection.
	db := database.ConnectDB(cfg)

	// Create a new Fiber application instance. This initializes the web server framework.
	server := fiber.New()

	// Register all application routes and middleware with the Fiber server.
	// This function typically sets up API endpoints, authentication, and other request processing logic.
	router.Router(server, cfg, db)

	// Construct the server address string from the configuration, combining the host and port.
	// For example, if Host is "0.0.0.0" and Port is "8080", address will be "0.0.0.0:8080".
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Start the Fiber server in a new goroutine. This allows the main goroutine to continue
	// executing, specifically to set up signal handling for graceful shutdown.
	go func() {
		// Listen for incoming HTTP requests on the specified address.
		if err := server.Listen(address); err != nil {
			// If the server encounters an error that prevents it from starting or continuing,
			// log the error and terminate the application using log.Panicf.
			log.Panicf("Server error: %v", err)
		}
	}()

	// Create a buffered channel of type os.Signal with a capacity of 1.
	// This channel will be used to receive notifications about operating system signals.
	c := make(chan os.Signal, 1)
	// Register the channel 'c' to receive notifications for specific OS signals:
	// os.Interrupt (typically Ctrl+C) and syscall.SIGTERM (a termination signal).
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// Block the main goroutine until a signal is received on the channel 'c'.
	// This effectively keeps the application running until an interrupt or terminate signal is sent.
	<-c

	// Print a message to the console indicating that the application is starting its graceful shutdown process.
	fmt.Println("Gracefully shutting down...")
	// Attempt to gracefully shut down the Fiber server. This allows ongoing requests to complete
	// and prevents new connections, ensuring a clean exit. The error is ignored with '_'.
	_ = server.Shutdown()

	// Print a message indicating that cleanup tasks, such as closing database connections, are being performed.
	fmt.Println("Running cleanup tasks...")
	// Attempt to close the database connection. This releases database resources and ensures
	// that no open connections are left behind. The error is ignored with '_'.
	_ = db.Close()

	// Print a final message confirming that the Fiber application has been successfully shut down.
	fmt.Println("Fiber was successful shutdown.")
}
