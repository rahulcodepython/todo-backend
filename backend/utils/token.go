package utils

import (
	"time" // Import the "time" package to handle time-related operations, such as setting token expiration.

	"github.com/golang-jwt/jwt/v5"                           // Import the "jwt" package from "github.com/golang-jwt/jwt/v5" for JSON Web Token (JWT) creation and signing.
	"github.com/rahulcodepython/todo-backend/backend/config" // Import the "config" package to access application-wide configurations, particularly JWT settings.
)

type Token struct {
	Token     string    `json:"token"`      // Define the 'Token' field, which will hold the actual JWT string, and tag it for JSON serialization.
	ExpiresAt time.Time `json:"expires_at"` // Define the 'ExpiresAt' field, a time.Time object indicating when the token expires, and tag it for JSON serialization.
}

// CreateToken generates a new JWT for a given user ID and application configuration.
// It constructs a JWT with specific claims (user ID, expiration, issued at time) and signs it.
//
// Parameters:
// - userId: A string representing the unique identifier of the user for whom the token is being created.
// - cfg: A pointer to the application's configuration struct, containing JWT secret and expiration settings.
//
// Returns:
// - *Token: A pointer to a Token struct containing the generated JWT string and its expiration time, or nil if an error occurs during token signing.
func CreateToken(userId string, cfg *config.Config) *Token {
	// Initialize a new Token struct.
	token := Token{
		Token:     "",                              // Initialize the token string as empty; it will be populated after signing.
		ExpiresAt: time.Now().Add(cfg.JWT.Expires), // Calculate the token's expiration time by adding the configured JWT expiry duration to the current time.
	}

	// Create a map of claims to be included in the JWT payload.
	claims := jwt.MapClaims{
		"user_id": userId,                                 // Set the "user_id" claim to the provided user's identifier.
		"exp":     time.Now().Add(cfg.JWT.Expires).Unix(), // Set the "exp" (expiration time) claim as a Unix timestamp.
		"iat":     time.Now().Unix(),                      // Set the "iat" (issued at time) claim as a Unix timestamp.
	}

	// Create a new JWT token instance using the HS256 signing method and the defined claims.
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenClaims.SignedString([]byte(cfg.JWT.SecretKey)) // Sign the token using the configured JWT secret key, converting it to a byte slice.
	if err != nil {
		return nil
	}

	token.Token = tokenString
	return &token
}
