package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	server := fiber.New()

	address := fmt.Sprintf("%s:%s", "localhost", "8000")
	if err := server.Listen(address); err != nil {
		log.Panicf("Server error: %v", err)
	}
}
