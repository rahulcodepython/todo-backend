package utils

const (
	UserTableName   = "users"
	UserTABLESCHEMA = "(id, name, email, image, password, jwt, created_at, updated_at)"

	JWTTableName   = "jwt_tokens"
	JWTTABLESCHEMA = "(id, token, expires_at)"
)
