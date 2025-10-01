// Package users provides handlers and logic for user-related operations
// such as registration, login, logout, and profile management.
package users

import (
	// Standard library package for database/sql operations.
	"database/sql"
	// Standard library package for handling time-related data.
	"time"

	// The web framework used for building the API.
	"github.com/gofiber/fiber/v2"
	// Package for generating UUIDs, specifically V7 which is time-ordered.
	"github.com/google/uuid"
	// Internal package for application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// Internal package for standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
	// Internal package for utility functions like password hashing and token creation.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// UserControl is a struct that holds dependencies for user-related handlers.
// It acts as a controller, centralizing the configuration and database connection.
type UserControl struct {
	// cfg holds the application's configuration, like JWT secrets and timings.
	cfg *config.Config
	// db is the database connection pool used for all user-related queries.
	db *sql.DB
}

// NewUserControl is a constructor function that creates and returns a new instance of UserControl.
// This pattern is used for dependency injection, making the controller easier to test and manage.
func NewUserControl(cfg *config.Config, db *sql.DB) *UserControl {
	// Returns a pointer to a new UserControl struct initialized with the provided config and db connection.
	return &UserControl{
		cfg: cfg,
		db:  db,
	}
}

// BodyParser is a generic helper function to parse the request body into a given struct.
// Using generics [T any] allows this function to work with any struct type (e.g., registerUserRequest, loginUserRequest).
func BodyParser[T any](c *fiber.Ctx, body *T) error {
	// c.BodyParser attempts to parse the incoming request body and unmarshal it into the 'body' struct.
	if err := c.BodyParser(body); err != nil {
		// If parsing fails, it sends a standardized bad request response to the client.
		// This centralizes the parsing and error handling logic.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}
	// If parsing is successful, return nil to indicate no error occurred.
	return nil
}

// CreateNewJWTAndUpdateUser generates a new JWT, saves it to the database,
// and associates it with the provided user account.
func CreateNewJWTAndUpdateUser(user User, uc *UserControl, c *fiber.Ctx) (JWT, error) {
	// Generate a new JWT using the utility function, passing the user's ID and app config.
	jwtToken := utils.CreateToken(user.ID.String(), uc.cfg)
	// Generate a new time-ordered UUID (V7) for the JWT's primary key in the database.
	tokenId, _ := uuid.NewV7() // Ignoring the error as V7 UUID generation is highly reliable.

	// Create a JWT struct to hold the new token information.
	jwt := JWT{
		// The unique ID for this specific JWT record.
		ID: tokenId,
		// The actual JWT string that will be sent to the client.
		Token: jwtToken.Token,
		// The expiration timestamp of the token.
		ExpiresAt: jwtToken.ExpiresAt,
	}

	// Execute a SQL query to insert the new JWT record and update the user's foreign key reference to it.
	// CreateNewJWT_UpdateUserRowQuery is expected to be a constant holding the SQL string.
	_, err := uc.db.Exec(CreateNewJWT_UpdateUserRowQuery, jwt.ID, jwt.Token, jwt.ExpiresAt, user.ID)
	// Check if the database execution resulted in an error.
	if err != nil {
		// If there's an error, return an empty JWT struct and the error itself.
		return JWT{}, err
	}

	// If successful, return the newly created JWT struct and no error.
	return jwt, nil
}

// RegisterUserController handles the logic for new user registration.
func (uc *UserControl) RegisterUserController(c *fiber.Ctx) error {
	// Allocate memory for a new 'registerUserRequest' struct to hold the request body data.
	body := new(registerUserRequest)

	// Parse the request body into the 'body' struct. Note: This helper sends a response on error.
	BodyParser(c, body)

	// Validate that all required fields are present in the request.
	if body.Name == "" || body.Email == "" || body.Password == "" {
		// If any field is missing, send a 400 Bad Request response.
		return response.BadResponse(c, "All fields are required")
	}

	// Declare a variable to store the count of users with the same email.
	var count int

	// Query the database to check if a user with the given email already exists.
	// Scan assigns the result of the query (the count) to the 'count' variable.
	err := uc.db.QueryRow(CheckUniqueEmailQuery, body.Email).Scan(&count)
	// Check for any errors during the database query.
	if err != nil {
		// If an error occurs, return a 500 Internal Server Error response.
		return response.InternelServerError(c, err, "Error checking unique email")
	}

	// If the count is greater than 0, it means the email is already in use.
	if count > 0 {
		// Return a 400 Bad Request response indicating the email is taken.
		return response.BadResponse(c, "This email already is ready used. Try something new!")
	}

	// Generate a new unique, time-ordered ID for the new user.
	userId, _ := uuid.NewV7() // Ignoring the error as V7 UUID generation is highly reliable.
	// Create a new User struct with the data from the request body.
	user := User{
		ID:        userId,
		Name:      body.Name,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now(), // Set the creation timestamp to the current time.
		UpdatedAt: time.Now(), // Set the update timestamp to the current time.
	}

	// Encrypt the user's plaintext password using a secure hashing algorithm (e.g., bcrypt).
	encryptedPassword, err := utils.EncryptPassword(user.Password)
	// Check if the password encryption failed.
	if err != nil {
		// If so, return a 500 Internal Server Error.
		return response.InternelServerError(c, err, "Error encrypting password")
	}
	// Replace the plaintext password in the user struct with its encrypted version.
	user.Password = encryptedPassword

	// Execute the SQL query to insert the new user's data into the database.
	// Note: user.Image and the JWT foreign key are initially set to nil.
	_, err = uc.db.Exec(CreateUserQuery, user.ID, user.Name, user.Email, user.Image, user.Password, nil, user.CreatedAt, user.UpdatedAt)
	// Check if the database insert operation failed.
	if err != nil {
		// If so, return a 500 Internal Server Error.
		return response.InternelServerError(c, err, "Error creating user")
	}

	// After successfully creating the user, generate their first JWT for authentication.
	jwt, err := CreateNewJWTAndUpdateUser(user, uc, c)
	// Check if JWT creation failed.
	if err != nil {
		// If so, return a 500 Internal Server Error.
		return response.InternelServerError(c, err, "Error creating JWT token")
	}

	// Create the response payload struct, which combines user and JWT data.
	responseUser := register_loginUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt), // Format time for a consistent response.
		UpdatedAt: utils.ParseTime(user.UpdatedAt), // Format time for a consistent response.
		Token:     jwt.Token,                       // The JWT string.
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),  // The formatted expiration time.
	}

	// Send a 200 OK response to the client with a success message and the user data.
	return response.OKResponse(c, "User registered successfully", responseUser)
}

// LoginUserController handles the authentication logic for existing users.
func (uc *UserControl) LoginUserController(c *fiber.Ctx) error {
	// Allocate memory for a new 'loginUserRequest' struct.
	body := new(loginUserRequest)

	// Parse the request body into the 'body' struct.
	BodyParser(c, body)

	// Validate that both email and password are provided.
	if body.Email == "" || body.Password == "" {
		// If not, return a 400 Bad Request response.
		return response.BadResponse(c, "All fields are required")
	}

	// Declare a variable 'user' to hold the user data fetched from the database.
	var user User
	// Declare a variable 'jwt' to hold JWT data.
	var jwt JWT

	// Query the database to find a user by their email address.
	err := uc.db.QueryRow(GetUserProfileByEmailQuery, body.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Password, &user.JWT, &user.CreatedAt, &user.UpdatedAt)
	// Check if an error occurred during the query.
	if err != nil {
		// If the error is `sql.ErrNoRows`, it means no user with that email was found.
		if err == sql.ErrNoRows {
			// Return a 404 Not Found response.
			return response.NotFound(c, err, "User not found")
		}
		// For any other database error, return a 500 Internal Server Error.
		return response.InternelServerError(c, err, "Error checking user")
	}

	// Compare the provided plaintext password with the hashed password stored in the database.
	passwordMatched := utils.CompareEncryptedPassword(user.Password, body.Password)
	// If the passwords do not match.
	if !passwordMatched {
		// Return a 401 Unauthorized response. The error 'err' is nil here, but passed for consistency.
		return response.UnauthorizedAccess(c, err, "Invalid credentials")
	}

	// Check if the user already has a JWT associated with their account. `user.JWT` is likely a sql.NullString or similar.
	if !user.JWT.Valid {
		// If no valid JWT exists, create a new one for this login session.
		jwt, err = CreateNewJWTAndUpdateUser(user, uc, c)
		// Handle potential errors during JWT creation.
		if err != nil {
			return response.InternelServerError(c, err, "Error creating JWT token")
		}
	} else {
		// If a JWT already exists, fetch its details along with user info.
		err = uc.db.QueryRow(GetUserLoginInfoQuery, user.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &jwt.ID, &jwt.Token, &jwt.ExpiresAt, &user.CreatedAt, &user.UpdatedAt)
		// Handle errors during this query.
		if err != nil {
			// Specifically handle the case where the user/JWT might not be found.
			if err == sql.ErrNoRows {
				return response.NotFound(c, err, "User not found")
			}
			// Handle other database errors.
			return response.InternelServerError(c, err, "Error checking user")
		}

		// Check if the existing token has expired by comparing its expiration time with the current time.
		if jwt.ExpiresAt.Before(time.Now()) {
			// If the token is expired, delete it from the database to clean up.
			_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
			// Handle potential errors during deletion.
			if err != nil {
				return response.InternelServerError(c, err, "Error deleting expired JWT")
			}

			// Since the old token was expired, create a new one for the user.
			jwt, err = CreateNewJWTAndUpdateUser(user, uc, c)
			// Handle potential errors during the new JWT creation.
			if err != nil {
				return response.InternelServerError(c, err, "Error creating JWT token")
			}
		}
	}

	// Construct the response payload with the user's data and the valid JWT.
	responseUser := register_loginUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: utils.ParseTime(user.CreatedAt),
		UpdatedAt: utils.ParseTime(user.UpdatedAt),
		Token:     jwt.Token,
		ExpiresAt: utils.ParseTime(jwt.ExpiresAt),
	}

	// Send a 200 OK response indicating a successful login.
	return response.OKResponse(c, "User logged in successfully", responseUser)
}

// LogoutUserController handles user logout by invalidating their current JWT.
func (uc *UserControl) LogoutUserController(c *fiber.Ctx) error {
	// Retrieve the JWT data from the request's local context.
	// This data is typically placed here by an authentication middleware after validating the token.
	jwt := c.Locals("jwt").(JWT)

	// Execute a SQL query to delete the JWT record from the database using its unique ID.
	// This effectively invalidates the token for any future requests.
	_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
	// Check for any errors during the database deletion.
	if err != nil {
		// If an error occurs, return a 500 Internal Server Error.
		return response.InternelServerError(c, err, "Error deleting JWT")
	}

	// Send a 200 OK response with a success message and no data payload.
	return response.OKResponse(c, "User logged out successfully", nil)
}

// UserProfileController fetches and returns the profile of the currently authenticated user.
func (uc *UserControl) UserProfileController(c *fiber.Ctx) error {
	// Retrieve the validated JWT data from the request's local context, put there by a middleware.
	jwt := c.Locals("jwt").(JWT)

	// Declare a 'user' variable to hold the profile data.
	var user User

	// Query the database to get the user's profile information by joining with the JWT ID.
	err := uc.db.QueryRow(GetUserProfileByJWTQuery, jwt.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.CreatedAt, &user.UpdatedAt)
	// Check for errors during the database query.
	if err != nil {
		// If the user associated with the JWT is not found, it might indicate a data consistency issue.
		if err == sql.ErrNoRows {
			// This case is treated as an internal server error as a valid JWT should always have a corresponding user.
			return response.InternelServerError(c, err, "Error checking user")
		}
		// For other errors, return a 404 Not Found. This logic could be swapped with the above depending on desired behavior.
		return response.NotFound(c, err, "User not found")
	}

	// If the user profile is fetched successfully, send a 200 OK response with the user data.
	return response.OKResponse(c, "User profile fetched successfully", user)
}
