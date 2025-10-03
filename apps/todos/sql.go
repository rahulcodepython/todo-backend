// This file defines the SQL queries used for todo-related database operations.
package todos

// "fmt" provides functions for formatted I/O. It is used here to construct the SQL queries.
import (
	"fmt"

	// "github.com/rahulcodepython/todo-backend/backend/utils" is a local package that provides constant values for table names and schemas.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// CreateTodoQuery is the SQL query to insert a new todo into the database.
var CreateTodoQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5)", utils.TodoTableName, utils.TodoTableSchema)

// GetTodosByUserQuery is the SQL query to retrieve all todos for a specific user.
var GetTodosByUserQuery = fmt.Sprintf("SELECT %s FROM %s WHERE owner = $1 LIMIT $2 OFFSET $3", utils.TodoTableSchema, utils.TodoTableName)

// GetTodosByUserFilteredByCompletedQuery is the SQL query to retrieve all todos for a specific user, filtered by completion status.
var GetTodosByUserFilteredByCompletedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE owner = $1 AND completed = $2 LIMIT $3 OFFSET $4", utils.TodoTableSchema, utils.TodoTableName)

// UpdateTodoTitleQuery is the SQL query to update the title of a todo.
var UpdateTodoTitleQuery = fmt.Sprintf("UPDATE %s SET title = $1 WHERE id = $2 returning %s", utils.TodoTableName, utils.TodoTableSchema)

// UpdateTodoCompletedQuery is the SQL query to update the completion status of a todo.
var UpdateTodoCompletedQuery = fmt.Sprintf("UPDATE %s SET completed = $1 WHERE id = $2 returning %s", utils.TodoTableName, utils.TodoTableSchema)

// DeleteTodoQuery is the SQL query to delete a todo.
var DeleteTodoQuery = fmt.Sprintf("DELETE FROM %s WHERE id = $1", utils.TodoTableName)

// GetTodoUserQuery is the SQL query to retrieve the owner of a todo.
var GetTodoUserQuery = fmt.Sprintf("SELECT owner FROM %s WHERE id = $1", utils.TodoTableName)

// CountTodosByUserQuery is the SQL query to count all todos for a specific user.
var CountTodosByUserQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE owner = $1", utils.TodoTableName)

// CountTodosByUserFilteredByCompletedQuery is the SQL query to count all todos for a specific user, filtered by completion status.
var CountTodosByUserFilteredByCompletedQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE owner = $1 AND completed = $2", utils.TodoTableName)