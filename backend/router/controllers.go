package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	DB *sql.DB
}

func NewController(db *sql.DB) *HealthController {
	return &HealthController{DB: db}
}

func (hc *HealthController) HealthCheckUp(c *fiber.Ctx) error {
	// Ping the database to verify the connection is alive.
	// Using PingContext allows the request context to cancel the ping if the client disconnects.
	if err := hc.DB.PingContext(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "Database connection is down.",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Backend is up and running, and the database connection is healthy.",
		"success": true,
	})
}
