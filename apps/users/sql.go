package users

import (
	"fmt"

	"github.com/rahulcodepython/todo-backend/backend/utils"
)

var CreateJWTTokenQuery = fmt.Sprintf("INSERT INTO %s %s VALUES ($1, $2, $3)", utils.JWTTableName, utils.JWTTABLESCHEMA)

var CreateUserQuery = fmt.Sprintf("INSERT INTO %s %s VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", utils.UserTableName, utils.UserTABLESCHEMA)

var CheckUniqueEmailQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", utils.UserTableName)
