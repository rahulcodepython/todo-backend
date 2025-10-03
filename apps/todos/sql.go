package todos

import (
	"fmt" // Import the "fmt" package for formatted I/O, used here for constructing SQL query strings.

	"github.com/rahulcodepython/todo-backend/backend/utils" // Import the application's utility package, specifically for database table name constants.
)

var CreateTodoQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5)", utils.TodoTableName, utils.TodoTableSchema)

var GetTodosQuery = fmt.Sprintf("SELECT %s FROM %s WHERE owner = $1", utils.TodoTableSchema, utils.TodoTableName)

var UpdateTodoTitleQuery = fmt.Sprintf("UPDATE %s SET title = $1 WHERE id = $2 returning %s", utils.TodoTableName, utils.TodoTableSchema)

var UpdateTodoCompletedQuery = fmt.Sprintf("UPDATE %s SET completed = $1 WHERE id = $2 returning %s", utils.TodoTableName, utils.TodoTableSchema)

var DeleteTodoQuery = fmt.Sprintf("DELETE FROM %s WHERE id = $1", utils.TodoTableName)

var GetTodoUserQuery = fmt.Sprintf("SELECT owner FROM %s WHERE id = $1", utils.TodoTableName)
