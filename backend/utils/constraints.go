// This file defines constant values used throughout the application.
package utils

// const is a keyword that declares a constant value.
const (
	// UserTableName is the name of the users table in the database.
	UserTableName = "users"
	// UserTableSchema is the schema of the users table in the database.
	UserTableSchema = "id, name, email, image, password, jwt, created_at, updated_at"

	// JWTTableName is the name of the jwt_tokens table in the database.
	JWTTableName = "jwt_tokens"
	// JWTTableSchema is the schema of the jwt_tokens table in the database.
	JWTTableSchema = "id, token, expires_at"

	// TodoTableName is the name of the todos table in the database.
	TodoTableName = "todos"
	// TodoTableSchema is the schema of the todos table in the database.
	TodoTableSchema = "id, title, completed, owner, created_at"
)