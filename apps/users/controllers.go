package users

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UsersController struct {
	DB  *sql.DB
	Cfg *config.Config
}

func NewUsersController(db *sql.DB, cfg *config.Config) *UsersController {
	return &UsersController{DB: db, Cfg: cfg}
}

func (uc *UsersController) GenerateToken(user *User) (*AuthResponse, error) {
	token, err := utils.GenerateJWT(user.ID, uc.Cfg.JWT.Secret, uc.Cfg.JWT.Expires)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(uc.Cfg.JWT.Expires)

	_, err = uc.DB.Exec(CreateJWTQuery, uuid.New().String(), user.ID, token, now, expiresAt)
	if err != nil {
		log.Printf("Error storing JWT: %v", err)
		return nil, err
	}

	res := &AuthResponse{
		Token:  token,
		Expiry: expiresAt.Format(time.RFC3339),
	}

	return res, err
}

// Register creates a new user account.
func (uc *UsersController) Register(c *fiber.Ctx) error {
	body := new(RegisterUserInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to hash password"})
	}

	user := &User{
		ID:        uuid.New().String(),
		Name:      body.Name,
		Email:     body.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	_, err = uc.DB.Exec(CreateUserQuery, user.ID, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		// You might want to check for unique constraint violation here specifically
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not create user",
			"error":   err.Error(),
		})
	}

	token, err := uc.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "User registered successfully", "data": token})
}

// Login authenticates a user and returns a JWT.
func (uc *UsersController) Login(c *fiber.Ctx) error {
	body := new(LoginUserInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	var user User
	err := uc.DB.QueryRow(GetUserByEmailQuery, body.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": utils.ErrInvalidCredentials})
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": utils.ErrInvalidCredentials})
	}

	token, err := uc.GenerateToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"success": true, "data": token})
}

// Logout invalidates the user's current JWT.
func (uc *UsersController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := uc.DB.Exec(DeleteJWTByTokenQuery, tokenString)
	if err != nil {
		log.Printf("Error on logout (deleting token): %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Logged out successfully"})
}

// GetProfile returns the profile of the currently authenticated user.
func (uc *UsersController) GetProfile(c *fiber.Ctx) error {
	user, ok := c.Locals(utils.CtxUserKey).(*User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}
