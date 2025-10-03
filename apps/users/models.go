// This file defines the data models for users and JWTs.
package users

// "time" provides functions for working with time. It is used here to define the CreatedAt and UpdatedAt fields.
import (
	"time"

	// "github.com/google/uuid" is a package for working with UUIDs. It is used here to define the ID and JWT fields.
	"github.com/google/uuid"
)

// User represents the structure of a user in the application.
type User struct {
	// ID is the unique identifier for the user.
	// json:"id" specifies that this field should be marshalled to/from a JSON object with the key "id".
	ID uuid.UUID `json:"id"`
	// Name is the user's name.
	// json:"name" specifies that this field should be marshalled to/from a JSON object with the key "name".
	Name string `json:"name"`
	// Email is the user's email address.
	// json:"email" specifies that this field should be marshalled to/from a JSON object with the key "email".
	Email string `json:"email"`
	// Image is the user's profile image.
	// json:"image" specifies that this field should be marshalled to/from a JSON object with the key "image".
	Image string `json:"image"`
	// Password is the user's hashed password.
	// json:"-" specifies that this field should be omitted from JSON serialization.
	Password string `json:"-"`
	// JWT is the user's JSON Web Token.
	// json:"-" specifies that this field should be omitted from JSON serialization.
	JWT uuid.NullUUID `json:"-"`
	// CreatedAt is the time the user was created.
	// json:"created_at" specifies that this field should be marshalled to/from a JSON object with the key "created_at".
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the user was last updated.
	// json:"updated_at" specifies that this field should be marshalled to/from a JSON object with the key "updated_at".
	UpdatedAt time.Time `json:"updated_at"`
}

// JWT represents the structure of a JSON Web Token.
type JWT struct {
	// ID is the unique identifier for the JWT.
	// json:"id" specifies that this field should be marshalled to/from a JSON object with the key "id".
	ID uuid.UUID `json:"id"`
	// Token is the JWT string.
	// json:"token" specifies that this field should be marshalled to/from a JSON object with the key "token".
	Token string `json:"token"`
	// ExpiresAt is the expiration time of the JWT.
	// json:"expires_at" specifies that this field should be marshalled to/from a JSON object with the key "expires_at".
	ExpiresAt time.Time `json:"expires_at"`
}