package users

import (
	"time" // Import the "time" package to handle time-related operations, such as creation and update timestamps.

	"github.com/google/uuid" // Import the "uuid" package for generating and handling UUIDs, used for unique identifiers.
)

// User represents the structure of a user in the application, mapping to the 'users' table in the database.
type User struct {
	ID        uuid.UUID     `json:"id"`         // ID is the unique identifier for the user, serialized as "id" in JSON.
	Name      string        `json:"name"`       // Name is the user's full name, serialized as "name" in JSON.
	Email     string        `json:"email"`      // Email is the user's email address, serialized as "email" in JSON.
	Image     string        `json:"image"`      // Image is an optional URL or path to the user's profile image, serialized as "image" in JSON.
	Password  string        `json:"-"`          // Password is the user's hashed password, omitted from JSON serialization for security.
	JWT       uuid.NullUUID `json:"-"`          // JWT is a nullable UUID referencing the associated JWT token, omitted from JSON for security and internal use.
	CreatedAt time.Time     `json:"created_at"` // CreatedAt is the timestamp when the user account was created, serialized as "created_at" in JSON.
	UpdatedAt time.Time     `json:"updated_at"` // UpdatedAt is the timestamp when the user account was last updated, serialized as "updated_at" in JSON.
}

// JWT represents the structure of a JSON Web Token stored in the application, mapping to the 'jwt_tokens' table.
type JWT struct {
	ID        uuid.UUID `json:"id"`         // ID is the unique identifier for the JWT token, serialized as "id" in JSON.
	Token     string    `json:"token"`      // Token is the actual JWT string, serialized as "token" in JSON.
	ExpiresAt time.Time `json:"expires_at"` // ExpiresAt is the timestamp when the JWT token expires, serialized as "expires_at" in JSON.
}
