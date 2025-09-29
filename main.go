package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/database"
	"github.com/rahulcodepython/todo-backend/backend/router"
)

func main() {
	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	server := fiber.New()

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	go func() {
		if err := server.Listen(address); err != nil {
			log.Panicf("Server error: %v", err)
		}
	}()

	router.Router(server, cfg, db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Gracefully shutting down...")
	_ = server.Shutdown()

	fmt.Println("Running cleanup tasks...")
	_ = db.Close()

	fmt.Println("Fiber was successful shutdown.")
}
