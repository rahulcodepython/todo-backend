package utils

const (
	// Context Keys
	CtxUserKey = "user"

	// Error Messages
	ErrMissingAuthHeader  = "missing authorization header"
	ErrInvalidAuthHeader  = "invalid authorization header format"
	ErrInvalidToken       = "invalid or malformed token"
	ErrTokenExpired       = "token has expired"
	ErrUserNotFound       = "user not found"
	ErrInvalidCredentials = "invalid email or password"
)
