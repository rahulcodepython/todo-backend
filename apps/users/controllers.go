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

	err := uc.db.QueryRow(CheckUniqueEmailQuery, body.Email).Scan(&count)
	if err != nil {
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

	encryptedPassword, err := utils.EncryptPassword(user.Password)
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

	responseUser := register_loginUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		Token:     jwt.Token,
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),
	}

	response := response{
		Success: true,
		Message: "User registered successfully",
		Data:    responseUser,
	}
	return c.JSON(response)
}

func (uc *UserControl) LoginUserController(c *fiber.Ctx) error {
	body := new(loginUserRequest)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Success: false,
			Message: "All fields are required",
		})
	}

	var user User
	var jwt JWT

	err := uc.db.QueryRow(GetUserProfileByEmailQuery, body.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Password, &user.JWT, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(response{
				Success: false,
				Message: "User not found + 1",
				Error:   err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error checking user",
			Error:   err.Error(),
		})
	}

	passwordMatched := utils.CompareEncryptedPassword(user.Password, body.Password)
	if !passwordMatched {
		return c.Status(fiber.StatusUnauthorized).JSON(response{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	if !user.JWT.Valid {
		jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)
		tokenId, _ := uuid.NewV7()

		jwt.ID = tokenId
		jwt.Token = jwtToken.Token
		jwt.ExpiresAt = jwtToken.ExpiresAt

		_, err = uc.db.Exec(CreateNewJWT_UpdateUserRowQuery, jwt.ID, jwt.Token, jwt.ExpiresAt, user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{
				Success: false,
				Message: "Error creating new JWT and updating the user",
				Error:   err.Error(),
			})
		}
	} else {
		err = uc.db.QueryRow(GetUserLoginInfoQuery, user.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &jwt.ID, &jwt.Token, &jwt.ExpiresAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(response{
					Success: false,
					Message: "User not found + 2",
					Error:   err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(response{
				Success: false,
				Message: "Error finding user",
				Error:   err.Error(),
			})
		}

		if jwt.ExpiresAt.Before(time.Now()) {
			jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)
			oldJWTId := jwt.ID
			tokenId, _ := uuid.NewV7()

			jwt.ID = tokenId
			jwt.Token = jwtToken.Token
			jwt.ExpiresAt = jwtToken.ExpiresAt

			_, err := uc.db.Exec(DeleteJWTByIdQuery, oldJWTId)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response{
					Success: false,
					Message: "Error deleting expired JWT",
					Error:   err.Error(),
				})
			}

			_, err = uc.db.Exec(CreateNewJWT_UpdateUserRowQuery, jwt.ID, jwt.Token, jwt.ExpiresAt, user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response{
					Success: false,
					Message: "Error creating new JWT and updating the user",
					Error:   err.Error(),
				})
			}
		}
	}

	responseUser := register_loginUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		Token:     jwt.Token,
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),
	}

	return c.Status(fiber.StatusOK).JSON(response{
		Success: true,
		Message: "User logged in successfully",
		Data:    responseUser,
	})
}

func (uc *UserControl) LogoutUserController(c *fiber.Ctx) error {
	jwt := c.Locals("jwt").(JWT)

	_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{
			Success: false,
			Message: "Error deleting JWT",
			Error:   err.Error(),
		})
	}

	return c.JSON(response{
		Success: true,
		Message: "User logged out successfully",
	})
}

func (uc *UserControl) UserProfileController(c *fiber.Ctx) error {
	jwt := c.Locals("jwt").(JWT)

	var user User

	err := uc.db.QueryRow(GetUserProfileByJWTQuery, jwt.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusInternalServerError).JSON(response{
				Success: false,
				Message: "Internal Server Error",
				Error:   err.Error(),
			})
		}

		return c.Status(fiber.StatusNotFound).JSON(response{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response{
		Success: true,
		Message: "User profile fetched successfully",
		Data:    user,
	})
}
