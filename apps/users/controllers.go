package users

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

type UserControl struct {
	cfg *config.Config
	db  *sql.DB
}

func NewUserControl(cfg *config.Config, db *sql.DB) *UserControl {
	return &UserControl{
		cfg: cfg,
		db:  db,
	}
}

func (uc *UserControl) RegisterUserController(c *fiber.Ctx) error {
	body := new(registerUserRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Success: false,
			Message: "All fields are required",
		})
	}

	user := User{
		Name:      body.Name,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	responseUser := registerUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
	}

	response := response{
		Success: true,
		Message: "User registered successfully",
		Data:    responseUser,
	}
	return c.JSON(response)
}
