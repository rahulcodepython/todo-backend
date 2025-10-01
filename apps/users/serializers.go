package users

import "github.com/google/uuid" // Import the uuid package for handling UUID types, specifically for user and JWT IDs.

// registerUserRequest defines the structure for the incoming JSON request body when a user registers.
// It includes validation tags for ensuring data integrity and proper formatting.
type registerUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"` // User's name, required, with a minimum length of 2 and maximum of 100 characters.
	Email    string `json:"email" validate:"required,email"`        // User's email, required, and must be a valid email format.
	Image    string `json:"image"`                                  // Optional URL or path to the user's profile image.
	Password string `json:"password" validate:"required,min=6"`     // User's password, required, with a minimum length of 6 characters for security.
}

// register_loginUserResponse defines the structure for the JSON response sent back to the client
// after a successful user registration or login. It includes user details and JWT information.
type register_loginUserResponse struct {
	ID        uuid.UUID `json:"id"`                   // Unique identifier for the user.
	Name      string    `json:"name"`                 // User's name.
	Email     string    `json:"email"`                // User's email address.
	Image     string    `json:"image"`                // User's profile image URL/path.
	Token     string    `json:"token,omitempty"`      // JWT token string, omitted if empty (e.g., for profile fetch without token refresh).
	ExpiresAt string    `json:"expires_at,omitempty"` // Expiration timestamp of the JWT, omitted if empty.
	CreatedAt string    `json:"created_at"`           // Timestamp when the user account was created.
	UpdatedAt string    `json:"updated_at"`           // Timestamp when the user account was last updated.
}

// loginUserRequest defines the structure for the incoming JSON request body when a user logs in.
// It includes validation tags for ensuring data integrity and proper formatting.
type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`    // User's email, required, and must be a valid email format.
	Password string `json:"password" validate:"required,min=6"` // User's password, required, with a minimum length of 6 characters.
}

// userProfileResponse defines the structure for the JSON response sent back to the client
// when fetching a user's profile. It includes user details but typically excludes sensitive info like password and JWT.
// This struct is currently identical to register_loginUserResponse but might diverge if specific fields need to be excluded
// or added for profile viewing purposes.
type userProfileResponse register_loginUserResponse
