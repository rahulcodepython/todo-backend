package users

import (
	"fmt"

	"github.com/rahulcodepython/todo-backend/backend/utils"
)

var CreateJWTTokenQuery = fmt.Sprintf("INSERT INTO %s %s VALUES ($1, $2, $3)", utils.JWTTableName, utils.JWTTableSchema)

var CreateUserQuery = fmt.Sprintf("INSERT INTO %s %s VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", utils.UserTableName, utils.UserTABLESCHEMA)

var CheckUniqueEmailQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", utils.UserTableName)

var CheckUserExistsByEmailGetId_PasswordQuery = fmt.Sprintf("SELECT id, password FROM %s WHERE email = $1", utils.UserTableName)

var GetUserLoginInfoQuery = fmt.Sprintf("SELECT u.id, u.name, u.email, u.image, j.id, j.token, j.expires_at, u.created_at, u.updated_at FROM %s u JOIN %s j ON u.jwt = j.id WHERE u.id = $1", utils.UserTableName, utils.JWTTableName)

var DeleteExpiredJWTQuery = fmt.Sprintf("DELETE FROM %s WHERE id = $1", utils.JWTTableName)

var CreateNewJWT_UpdateUserRowQuery = fmt.Sprintf("WITH new_token AS (INSERT INTO %s %s VALUES ($1, $2, $3) RETURNING id) UPDATE %s SET jwt = (SELECT id FROM new_token) WHERE id = $4", utils.JWTTableName, utils.JWTTableSchema, utils.UserTableName)
