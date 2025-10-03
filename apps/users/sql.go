package users

import (
	"fmt" // Import the "fmt" package for formatted I/O, used here for constructing SQL query strings.

	"github.com/rahulcodepython/todo-backend/backend/utils" // Import the application's utility package, specifically for database table name constants.
)

// CreateUserQuery defines the SQL query to insert a new user into the 'users' table.
// It uses `fmt.Sprintf` to dynamically insert the table name and schema from `utils` constants,
// ensuring consistency and reducing hardcoding. The query expects 8 parameters for user details.
var CreateUserQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", utils.UserTableName, utils.UserTableSchema)

// CheckUniqueEmailQuery defines the SQL query to count users with a specific email address.
// This is used to verify if an email is already registered in the 'users' table,
// ensuring email uniqueness during user registration. It expects 1 parameter for the email.
var CheckUniqueEmailQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", utils.UserTableName)

// GetUserProfileByEmailQuery defines the SQL query to retrieve all columns from the 'users' table
// for a given email address. This query is used to fetch a user's complete profile information
// based on their email, typically during login or when checking for existing users.
// It uses `fmt.Sprintf` to dynamically insert the table name and expects 1 parameter for the email.
var GetUserProfileByEmailQuery = fmt.Sprintf("SELECT %s FROM %s WHERE email = $1", utils.UserTableSchema, utils.UserTableName)

// GetUserLoginInfoQuery defines a SQL query to retrieve detailed user and JWT information by user ID.
// It performs a JOIN operation between the 'users' table (aliased as 'u') and the 'jwt_tokens' table (aliased as 'j')
// on the 'jwt' column of the 'users' table and the 'id' column of the 'jwt_tokens' table.
// This query is crucial for fetching all necessary data for user authentication and session management,
// including the user's profile details and the associated JWT token and its expiration.
// It uses `fmt.Sprintf` to dynamically insert table names and expects 1 parameter for the user ID.
var GetUserJWTInfoQuery = fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", utils.JWTTableSchema, utils.JWTTableName)

// DeleteJWTByIdQuery defines the SQL query to delete a JWT token from the 'jwt_tokens' table by its ID.
// This query is used to invalidate or revoke a specific JWT, typically during logout or when a token expires.
// It uses `fmt.Sprintf` to dynamically insert the table name and expects 1 parameter for the JWT ID.
var DeleteJWTByIdQuery = fmt.Sprintf("DELETE FROM %s WHERE id = $1", utils.JWTTableName)

// CreateNewJWT_UpdateUserRowQuery defines a complex SQL query using a Common Table Expression (CTE) to
// first insert a new JWT token into the 'jwt_tokens' table and then update the 'users' table
// to link the newly created JWT to a specific user. This ensures atomic creation and association.
// It expects 4 parameters: JWT ID, JWT Token string, JWT expiration time, and User ID.
var CreateNewJWT_UpdateUserRowQuery = fmt.Sprintf("WITH new_token AS (INSERT INTO %s (%s) VALUES ($1, $2, $3) RETURNING id) UPDATE %s SET jwt = (SELECT id FROM new_token) WHERE id = $4", utils.JWTTableName, utils.JWTTableSchema, utils.UserTableName)

// GetUserProfileByJWTQuery defines the SQL query to retrieve user profile information based on a JWT ID.
// It performs a JOIN operation between the 'users' table (aliased as 'u') and the 'jwt_tokens' table (aliased as 'j')
// on the 'jwt' column of the 'users' table and the 'id' column of the 'jwt_tokens' table.
// This query is used to fetch user details when only the JWT ID is known, typically after authentication.
// It expects 1 parameter for the JWT ID.
var GetUserProfileByJWTQuery = fmt.Sprintf("SELECT %s FROM %s WHERE jwt = $1", utils.UserTableSchema, utils.UserTableName)
