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
			Error:   err.Error(),
		})
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Success: false,
			Message: "All fields are required",
		})
	}

	var count int

	result := uc.db.QueryRow(CheckUniqueEmailQuery, body.Email)
	if err := result.Scan(&count); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error checking unique email",
			Error:   err.Error(),
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusNotAcceptable).JSON(response{
			Success: false,
			Message: "This email already is ready used. Try something new!",
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

	encryptedPassword, err := utils.Encrypt(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error encrypting password",
			Error:   err.Error(),
		})
	}

	user.Password = encryptedPassword
	jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)
	tokenId, _ := uuid.NewV7()

	jwt := JWT{
		ID:        tokenId,
		Token:     jwtToken.Token,
		ExpiresAt: jwtToken.ExpiresAt,
	}

	_, err = uc.db.Exec(CreateJWTTokenQuery, jwt.ID.String(), jwt.Token, jwt.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error creating JWT token",
			Error:   err.Error(),
		})
	}

	_, err = uc.db.Exec(CreateUserQuery, user.ID, user.Name, user.Email, user.Image, user.Password, jwt.ID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error creating user",
			Error:   err.Error(),
		})
	}

	responseUser := registerUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		Token:     jwt.Token,
		ExpiresAt: utils.ParseTime(jwtToken.ExpiresAt),
	}

	response := response{
		Success: true,
		Message: "User registered successfully",
		Data:    responseUser,
	}
	return c.JSON(response)
}
