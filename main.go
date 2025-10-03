// This file is the main entry point for the todo-backend application.
// It initializes the configuration, database connection, and the Fiber web server.
// It also handles graceful shutdown of the application.
package main

// "fmt" provides functions for formatted I/O. It is used here to print messages to the console.
import (
	"fmt"
	// "log" provides a simple logging package. It is used here to log fatal server errors.
	"log"
	// "os" provides a platform-independent interface to operating system functionality. It is used here to receive signals for graceful shutdown.
	"os"
	// "os/signal" provides functions for handling incoming signals from the operating system. It is used here to listen for interrupt and terminate signals.
	"os/signal"
	// "syscall" provides a low-level interface to operating system primitives. It is used here to specify the SIGTERM signal.
	"syscall"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create the HTTP server and define API routes.
	"github.com/gofiber/fiber/v2"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that handles loading application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/database" is a local package that manages the database connection.
	"github.com/rahulcodepython/todo-backend/backend/database"
	// "github.com/rahulcodepython/todo-backend/backend/router" is a local package that sets up the application's API routes.
	"github.com/rahulcodepython/todo-backend/backend/router"
)

// main is the entry point of the application.
// It initializes the server, database, and router, and then starts the server.
// It also includes logic for graceful shutdown.
func main() {
	// cfg is a variable that holds the application configuration.
	// config.LoadConfig() is called to load the configuration from environment variables or a .env file.
	cfg := config.LoadConfig()

	// db is a variable that holds the database connection.
	// database.ConnectDB() is called to establish a connection to the database using the loaded configuration.
	db := database.ConnectDB(cfg)

	// server is a new instance of a Fiber application.
	// fiber.New() creates a new Fiber server.
	server := fiber.New()

	// router.Router() is called to set up all the application routes and middleware.
	// It takes the Fiber server, configuration, and database connection as arguments.
	router.Router(server, cfg, db)

	// address is a string that represents the server address.
	// It is constructed by combining the server host and port from the configuration.
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// A new goroutine is started to run the Fiber server.
	// This allows the main goroutine to continue and handle graceful shutdown.
	go func() {
		// server.Listen() starts the HTTP server and listens for incoming requests on the specified address.
		if err := server.Listen(address); err != nil {
			// If an error occurs while starting the server, log the error and panic.
			log.Panicf("Server error: %v", err)
		}
	}()

	// c is a channel that will receive operating system signals.
	// It has a buffer size of 1.
	c := make(chan os.Signal, 1)
	// signal.Notify() registers the given channel to receive notifications of the specified signals.
	// In this case, it listens for os.Interrupt (Ctrl+C) and syscall.SIGTERM.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// This is a blocking call that waits for a signal to be received on the channel c.
	<-c

	// A message is printed to the console to indicate that the server is shutting down.
	fmt.Println("Gracefully shutting down...")
	// server.Shutdown() gracefully shuts down the server without interrupting any active connections.
	_ = server.Shutdown()

	// A message is printed to the console to indicate that cleanup tasks are running.
	fmt.Println("Running cleanup tasks...")
	// db.Close() closes the database connection.
	_ = db.Close()

	// A message is printed to the console to indicate that the server has shut down successfully.
	fmt.Println("Fiber was successful shutdown.")
}