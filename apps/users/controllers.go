package users

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	userId, _ := uuid.NewV7()

	user := User{
		ID:        userId,
		Name:      body.Name,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)

	encryptedPassword, err := utils.Encrypt(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error encrypting password",
		})
	}
	user.Password = encryptedPassword

	tokenId, _ := uuid.NewV7()

	_, err = uc.db.Exec(CreateJWTTokenQuery, tokenId.String(), jwtToken.Token, jwtToken.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error SQL Query: creating JWT token",
		})
	}

	_, err = uc.db.Exec(CreateUserQuery, user.ID, user.Name, user.Email, user.Image, user.Password, jwtToken.Token, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error SQL Query: creating user",
		})
	}

	responseUser := registerUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		Token:     jwtToken.Token,
		ExpiresAt: utils.ParseTime(jwtToken.ExpiresAt),
	}

	response := response{
		Success: true,
		Message: "User registered successfully",
		Data:    responseUser,
	}
	return c.JSON(response)
}
