// This file defines the serializers for user-related requests and responses.
package users

// "github.com/google/uuid" is a package for working with UUIDs. It is used here to define the ID field in the response struct.
import "github.com/google/uuid"

// registerUserRequest defines the structure for a user registration request.
type registerUserRequest struct {
	// Name is the user's name.
	// json:"name" specifies that this field should be marshalled to/from a JSON object with the key "name".
	// validate:"required,min=2,max=100" specifies that this field is required, has a minimum length of 2, and a maximum length of 100.
	Name string `json:"name" validate:"required,min=2,max=100"`
	// Email is the user's email address.
	// json:"email" specifies that this field should be marshalled to/from a JSON object with the key "email".
	// validate:"required,email" specifies that this field is required and must be a valid email address.
	Email string `json:"email" validate:"required,email"`
	// Image is the user's profile image.
	// json:"image" specifies that this field should be marshalled to/from a JSON object with the key "image".
	Image string `json:"image"`
	// Password is the user's password.
	// json:"password" specifies that this field should be marshalled to/from a JSON object with the key "password".
	// validate:"required,min=6" specifies that this field is required and has a minimum length of 6.
	Password string `json:"password" validate:"required,min=6"`
}

// register_loginUserResponse defines the structure for a user registration or login response.
type register_loginUserResponse struct {
	// ID is the user's ID.
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
	// Token is the user's JWT.
	// json:"token,omitempty" specifies that this field should be marshalled to/from a JSON object with the key "token", and should be omitted if empty.
	Token string `json:"token,omitempty"`
	// ExpiresAt is the expiration time of the JWT.
	// json:"expires_at,omitempty" specifies that this field should be marshalled to/from a JSON object with the key "expires_at", and should be omitted if empty.
	ExpiresAt string `json:"expires_at,omitempty"`
	// CreatedAt is the time the user was created.
	// json:"created_at" specifies that this field should be marshalled to/from a JSON object with the key "created_at".
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the user was last updated.
	// json:"updated_at" specifies that this field should be marshalled to/from a JSON object with the key "updated_at".
	UpdatedAt string `json:"updated_at"`
}

// loginUserRequest defines the structure for a user login request.
type loginUserRequest struct {
	// Email is the user's email address.
	// json:"email" specifies that this field should be marshalled to/from a JSON object with the key "email".
	// validate:"required,email" specifies that this field is required and must be a valid email address.
	Email string `json:"email" validate:"required,email"`
	// Password is the user's password.
	// json:"password" specifies that this field should be marshalled to/from a JSON object with the key "password".
	// validate:"required,min=6" specifies that this field is required and has a minimum length of 6.
	Password string `json:"password" validate:"required,min=6"`
}