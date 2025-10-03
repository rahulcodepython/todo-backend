// This file defines the SQL queries used for user-related database operations.
package users

// "fmt" provides functions for formatted I/O. It is used here to construct the SQL queries.
import (
	"fmt"

	// "github.com/rahulcodepython/todo-backend/backend/utils" is a local package that provides constant values for table names and schemas.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// CreateUserQuery is the SQL query to insert a new user into the database.
var CreateUserQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", utils.UserTableName, utils.UserTableSchema)

// CheckUniqueEmailQuery is the SQL query to check if an email is unique.
var CheckUniqueEmailQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", utils.UserTableName)

// GetUserProfileByEmailQuery is the SQL query to retrieve a user's profile by email.
var GetUserProfileByEmailQuery = fmt.Sprintf("SELECT %s FROM %s WHERE email = $1", utils.UserTableSchema, utils.UserTableName)

// GetUserJWTInfoQuery is the SQL query to retrieve a user's JWT information by user ID.
var GetUserJWTInfoQuery = fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", utils.JWTTableSchema, utils.JWTTableName)

// DeleteJWTByIdQuery is the SQL query to delete a JWT by its ID.
var DeleteJWTByIdQuery = fmt.Sprintf("DELETE FROM %s WHERE id = $1", utils.JWTTableName)

// CreateNewJWT_UpdateUserRowQuery is the SQL query to create a new JWT and update the user's row with the new JWT.
var CreateNewJWT_UpdateUserRowQuery = fmt.Sprintf("WITH new_token AS (INSERT INTO %s (%s) VALUES ($1, $2, $3) RETURNING id) UPDATE %s SET jwt = (SELECT id FROM new_token) WHERE id = $4", utils.JWTTableName, utils.JWTTableSchema, utils.UserTableName)

// GetUserProfileByJWTQuery is the SQL query to retrieve a user's profile by JWT.
var GetUserProfileByJWTQuery = fmt.Sprintf("SELECT %s FROM %s WHERE jwt = $1", utils.UserTableSchema, utils.UserTableName)