// This file provides functionality for creating JSON Web Tokens (JWTs).
package utils

// "time" provides functions for working with time. It is used here to set the expiration time of the JWT.
import (
	"time"

	// "github.com/golang-jwt/jwt/v5" is a package for creating and signing JWTs.
	"github.com/golang-jwt/jwt/v5"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to application configuration, including JWT settings.
	"github.com/rahulcodepython/todo-backend/backend/config"
)

// Token represents the structure of a JWT.
type Token struct {
	// Token is the JWT string.
	// json:"token" specifies that this field should be marshalled to/from a JSON object with the key "token".
	Token string `json:"token"`
	// ExpiresAt is the time when the token expires.
	// json:"expires_at" specifies that this field should be marshalled to/from a JSON object with the key "expires_at".
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateToken generates a new JWT for a given user ID.
// It takes a user ID and the application configuration as input.
// It returns a pointer to a Token struct containing the JWT and its expiration time, or nil if an error occurs.
//
// @param userId string - The ID of the user for whom the token is being created.
// @param cfg *config.Config - A pointer to the application's configuration struct.
// @return *Token - A pointer to a Token struct, or nil if an error occurs.
func CreateToken(userId string, cfg *config.Config) *Token {
	// token is a new instance of the Token struct.
	token := Token{
		// The Token field is initialized as an empty string.
		Token: "",
		// The ExpiresAt field is set to the current time plus the configured JWT expiration duration.
		ExpiresAt: time.Now().Add(cfg.JWT.Expires),
	}

	// claims is a map that holds the JWT claims.
	claims := jwt.MapClaims{
		// "user_id" is a claim that stores the user's ID.
		"user_id": userId,
		// "exp" is a claim that stores the expiration time of the token as a Unix timestamp.
		"exp": time.Now().Add(cfg.JWT.Expires).Unix(),
		// "iat" is a claim that stores the time the token was issued as a Unix timestamp.
		"iat": time.Now().Unix(),
	}

	// tokenClaims is a new JWT token with the specified signing method and claims.
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString is the signed JWT string.
	// tokenClaims.SignedString() signs the token with the configured JWT secret key.
	tokenString, err := tokenClaims.SignedString([]byte(cfg.JWT.SecretKey))
	// This checks if an error occurred while signing the token.
	if err != nil {
		// If an error occurs, return nil.
		return nil
	}

	// The Token field of the token struct is set to the signed token string.
	token.Token = tokenString
	// A pointer to the token struct is returned.
	return &token
}