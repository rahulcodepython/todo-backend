// This file defines the controllers for user-related operations.
package users

// "database/sql" provides a generic SQL interface. It is used here to interact with the database.
import (
	"database/sql"
	// "log" provides a simple logging package. It is used here to log fatal errors.
	"log"
	// "time" provides functions for working with time. It is used here to set timestamps.
	"time"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to define the controllers.
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid" is a package for working with UUIDs. It is used here to generate new UUIDs.
	"github.com/google/uuid"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
	// "github.com/rahulcodepython/todo-backend/backend/utils" is a local package that provides utility functions.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// UserControl is a struct that holds the configuration and database connection.
type UserControl struct {
	// cfg is the application configuration.
	cfg *config.Config
	// db is the database connection.
	db *sql.DB
}

// NewUserControl creates a new UserControl.
// It takes the application configuration and database connection as input.
//
// @param cfg *config.Config - The application configuration.
// @param db *sql.DB - The database connection.
// @return *UserControl - A pointer to the new UserControl.
func NewUserControl(cfg *config.Config, db *sql.DB) *UserControl {
	// This checks if the database connection is nil.
	if db == nil {
		// If the database connection is nil, a fatal error is logged.
		log.Fatal("Database connection is nil in NewUserControl!")
	}
	// A new UserControl is returned.
	return &UserControl{
		// The cfg field is set to the application configuration.
		cfg: cfg,
		// The db field is set to the database connection.
		db: db,
	}
}

// CreateNewJWTAndUpdateUser creates a new JWT and updates the user's row with the new JWT.
// It takes a user, a UserControl, and a Fiber context as input.
//
// @param user User - The user for whom the JWT is being created.
// @param uc *UserControl - The UserControl.
// @param c *fiber.Ctx - The Fiber context.
// @return JWT - The new JWT.
// @return error - An error if one occurred.
func CreateNewJWTAndUpdateUser(user User, uc *UserControl, c *fiber.Ctx) (JWT, error) {
	// jwtToken is the new JWT.
	jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)
	// tokenId is the new UUID for the JWT.
	tokenId, _ := uuid.NewV7()

	// jwt is a new JWT struct.
	jwt := JWT{
		// The ID field is set to the new UUID.
		ID: tokenId,
		// The Token field is set to the new JWT string.
		Token: jwtToken.Token,
		// The ExpiresAt field is set to the expiration time of the JWT.
		ExpiresAt: jwtToken.ExpiresAt,
	}

	// _, err is the result of executing the SQL query to create the new JWT and update the user's row.
	_, err := uc.db.Exec(CreateNewJWT_UpdateUserRowQuery, jwt.ID, jwt.Token, jwt.ExpiresAt, user.ID)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an empty JWT and the error are returned.
		return JWT{}, err
	}

	// The new JWT and no error are returned.
	return jwt, nil
}

// RegisterUserController handles user registration.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (uc *UserControl) RegisterUserController(c *fiber.Ctx) error {
	// body is a new registerUserRequest struct.
	body := new(registerUserRequest)
	// This parses the request body into the body struct.
	if err := c.BodyParser(body); err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	// This checks if all required fields are present.
	if body.Name == "" || body.Email == "" || body.Password == "" {
		// If any field is missing, a bad request response is returned.
		return response.BadResponse(c, "All fields are required")
	}

	// count is a variable that will hold the number of users with the same email.
	var count int

	// err is the result of querying the database to check if the email is unique.
	err := uc.db.QueryRow(CheckUniqueEmailQuery, body.Email).Scan(&count)
	// This checks if an error occurred while querying the database.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error checking unique email")
	}

	// This checks if the email is already in use.
	if count > 0 {
		// If the email is already in use, a bad request response is returned.
		return response.BadResponse(c, "This email already is ready used. Try something new!")
	}

	// userId is the new UUID for the user.
	userId, _ := uuid.NewV7()
	// user is a new User struct.
	user := User{
		// The ID field is set to the new UUID.
		ID: userId,
		// The Name field is set to the user's name.
		Name: body.Name,
		// The Email field is set to the user's email address.
		Email: body.Email,
		// The Password field is set to the user's password.
		Password: body.Password,
		// The CreatedAt field is set to the current time.
		CreatedAt: time.Now(),
		// The UpdatedAt field is set to the current time.
		UpdatedAt: time.Now(),
	}

	// encryptedPassword is the user's encrypted password.
	encryptedPassword, err := utils.EncryptPassword(user.Password)
	// This checks if an error occurred while encrypting the password.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error encrypting password")
	}
	// The user's password is replaced with the encrypted password.
	user.Password = encryptedPassword

	// _, err is the result of executing the SQL query to create the new user.
	_, err = uc.db.Exec(CreateUserQuery, user.ID, user.Name, user.Email, user.Image, user.Password, nil, user.CreatedAt, user.UpdatedAt)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error creating user")
	}

	// jwt is the new JWT for the user.
	jwt, err := CreateNewJWTAndUpdateUser(user, uc, c)
	// This checks if an error occurred while creating the JWT.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error creating JWT token")
	}

	// responseUser is a new register_loginUserResponse struct.
	responseUser := register_loginUserResponse{
		// The ID field is set to the user's ID.
		ID: user.ID,
		// The Name field is set to the user's name.
		Name: user.Name,
		// The Email field is set to the user's email address.
		Email: user.Email,
		// The CreatedAt field is set to the user's creation time.
		CreatedAt: utils.ParseTime(user.CreatedAt),
		// The UpdatedAt field is set to the user's last update time.
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		// The Token field is set to the new JWT.
		Token: jwt.Token,
		// The ExpiresAt field is set to the expiration time of the JWT.
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),
	}

	// An OK response is returned with a success message and the user data.
	return response.OKResponse(c, "User registered successfully", responseUser)
}

// LoginUserController handles user login.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (uc *UserControl) LoginUserController(c *fiber.Ctx) error {
	// body is a new loginUserRequest struct.
	body := new(loginUserRequest)
	// This parses the request body into the body struct.
	if err := c.BodyParser(body); err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	// This checks if all required fields are present.
	if body.Email == "" || body.Password == "" {
		// If any field is missing, a bad request response is returned.
		return response.BadResponse(c, "All fields are required")
	}

	// user is a variable that will hold the user's data.
	var user User
	// jwt is a variable that will hold the JWT data.
	var jwt JWT

	// err is the result of querying the database for the user's profile.
	err := uc.db.QueryRow(GetUserProfileByEmailQuery, body.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Password, &user.JWT, &user.CreatedAt, &user.UpdatedAt)
	// This checks if an error occurred while querying the database.
	if err != nil {
		// This checks if the error is sql.ErrNoRows.
		if err == sql.ErrNoRows {
			// If no user is found, a not found response is returned.
			return response.NotFound(c, err, "User not found")
		}
		// For any other error, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error fetching user profile info")
	}

	// passwordMatched is a boolean that indicates whether the passwords match.
	passwordMatched := utils.CompareEncryptedPassword(user.Password, body.Password)
	// This checks if the passwords do not match.
	if !passwordMatched {
		// If the passwords do not match, an unauthorized access response is returned.
		return response.UnauthorizedAccess(c, err, "Invalid credentials")
	}

	// This checks if the user already has a valid JWT.
	if !user.JWT.Valid {
		// If the user does not have a valid JWT, a new one is created.
		jwt, err = CreateNewJWTAndUpdateUser(user, uc, c)
		// This checks if an error occurred while creating the JWT.
		if err != nil {
			// If an error occurs, an internal server error response is returned.
			return response.InternelServerError(c, err, "Error creating JWT token")
		}
	} else {
		// If the user already has a JWT, its information is retrieved from the database.
		err = uc.db.QueryRow(GetUserJWTInfoQuery, user.JWT).Scan(&jwt.ID, &jwt.Token, &jwt.ExpiresAt)
		// This checks if an error occurred while querying the database.
		if err != nil {
			// This checks if the error is sql.ErrNoRows.
			if err == sql.ErrNoRows {
				// If no JWT is found, a not found response is returned.
				return response.NotFound(c, err, "User not found")
			}
			// For any other error, an internal server error response is returned.
			return response.InternelServerError(c, err, "Error fetching user login info")
		}

		// This checks if the JWT has expired.
		if jwt.ExpiresAt.Before(time.Now()) {
			// If the JWT has expired, it is deleted from the database.
			_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
			// This checks if an error occurred while deleting the JWT.
			if err != nil {
				// If an error occurs, an internal server error response is returned.
				return response.InternelServerError(c, err, "Error deleting expired JWT")
			}

			// A new JWT is created for the user.
			jwt, err = CreateNewJWTAndUpdateUser(user, uc, c)
			// This checks if an error occurred while creating the JWT.
			if err != nil {
				// If an error occurs, an internal server error response is returned.
				return response.InternelServerError(c, err, "Error creating JWT token")
			}
		}
	}

	// responseUser is a new register_loginUserResponse struct.
	responseUser := register_loginUserResponse{
		// The ID field is set to the user's ID.
		ID: user.ID,
		// The Name field is set to the user's name.
		Name: user.Name,
		// The Email field is set to the user's email address.
		Email: user.Email,
		// The CreatedAt field is set to the user's creation time.
		CreatedAt: utils.ParseTime(user.CreatedAt),
		// The UpdatedAt field is set to the user's last update time.
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		// The Token field is set to the new JWT.
		Token: jwt.Token,
		// The ExpiresAt field is set to the expiration time of the JWT.
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),
	}

	// An OK response is returned with a success message and the user data.
	return response.OKResponse(c, "User logged in successfully", responseUser)
}

// LogoutUserController handles user logout.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (uc *UserControl) LogoutUserController(c *fiber.Ctx) error {
	// jwt is the JWT object retrieved from the local context.
	jwt := c.Locals("jwt").(JWT)

	// _, err is the result of executing the SQL query to delete the JWT.
	_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Error deleting JWT")
	}

	// An OK response is returned with a success message.
	return response.OKResponse(c, "User logged out successfully", nil)
}

// UserProfileController handles retrieving the user's profile.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (uc *UserControl) UserProfileController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(User)
	// An OK response is returned with a success message and the user data.
	return response.OKResponse(c, "User profile fetched successfully", user)
}