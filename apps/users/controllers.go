package users

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/todo-backend/apps"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/response"
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

type UserControl struct {
	cfg *config.Config
	db  *sql.DB
}

func NewUserControl(cfg *config.Config, db *sql.DB) *UserControl {
	if db == nil {
		log.Fatal("Database connection is nil in NewUserControl!")
	}
	return &UserControl{
		cfg: cfg,
		db:  db,
	}
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

func ParseJWT(c *fiber.Ctx) (JWT, bool) {
	jwtInterface := c.Locals("jwt")
	if jwtInterface == nil {
		return JWT{}, false
	}

	jwt, ok := jwtInterface.(JWT)

	return jwt, ok
}

// RegisterUserController handles the logic for new user registration.
func (uc *UserControl) RegisterUserController(c *fiber.Ctx) error {
	// Allocate memory for a new 'registerUserRequest' struct to hold the request body data.
	body := new(registerUserRequest)

	// Parse the request body into the 'body' struct. Note: This helper sends a response on error.
	apps.BodyParser(c, body)

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
	apps.BodyParser(c, body)

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
		return response.InternelServerError(c, err, "Error fetching user profile info")
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
		err = uc.db.QueryRow(GetUserJWTInfoQuery, user.JWT).Scan(&jwt.ID, &jwt.Token, &jwt.ExpiresAt)
		// Handle errors during this query.
		if err != nil {
			// Specifically handle the case where the user/JWT might not be found.
			if err == sql.ErrNoRows {
				return response.NotFound(c, err, "User not found")
			}
			// Handle other database errors.
			return response.InternelServerError(c, err, "Error fetching user login info")
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

func (uc *UserControl) LogoutUserController(c *fiber.Ctx) error {
	jwt, ok := ParseJWT(c)
	if !ok {
		return response.InternelServerError(c, nil, "Invalid JWT type in context")
	}

	// Delete JWT from database
	_, err := uc.db.Exec(DeleteJWTByIdQuery, jwt.ID)
	if err != nil {
		return response.InternelServerError(c, err, "Error deleting JWT")
	}

	return response.OKResponse(c, "User logged out successfully", nil)
}

func (uc *UserControl) UserProfileController(c *fiber.Ctx) error {
	jwt, ok := ParseJWT(c)
	if !ok {
		return response.InternelServerError(c, nil, "Invalid JWT type in context")
	}

	var user User
	err := uc.db.QueryRow(GetUserProfileByJWTQuery, jwt.ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err, "User not found")
		}
		return response.InternelServerError(c, err, "Error fetching user profile")
	}

	return response.OKResponse(c, "User profile fetched successfully", user)
}
