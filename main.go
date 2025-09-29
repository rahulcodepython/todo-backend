package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func main() {
	cfg := config.LoadConfig()

	server := fiber.New()

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := server.Listen(address); err != nil {
		log.Panicf("Server error: %v", err)
	}
}
