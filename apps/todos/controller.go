// This file defines the controllers for todo-related operations.
package todos

// "database/sql" provides a generic SQL interface. It is used here to interact with the database.
import (
	"database/sql"
	// "math" provides basic mathematical functions. It is used here to calculate the total number of pages.
	"math"

	// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to define the controllers.
	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid" is a package for working with UUIDs. It is used here to generate and parse UUIDs.
	"github.com/google/uuid"
	// "github.com/rahulcodepython/todo-backend/apps/users" is a local package that contains user-related models.
	"github.com/rahulcodepython/todo-backend/apps/users"
	// "github.com/rahulcodepython/todo-backend/backend/config" is a local package that provides access to the application configuration.
	"github.com/rahulcodepython/todo-backend/backend/config"
	// "github.com/rahulcodepython/todo-backend/backend/response" is a local package that provides standardized API responses.
	"github.com/rahulcodepython/todo-backend/backend/response"
	// "github.com/rahulcodepython/todo-backend/backend/utils" is a local package that provides utility functions.
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

// TodoController is a struct that holds the configuration and database connection.
type TodoController struct {
	// cfg is the application configuration.
	cfg *config.Config
	// db is the database connection.
	db *sql.DB
}

// NewTodoControl creates a new TodoController.
// It takes the application configuration and database connection as input.
//
// @param cfg *config.Config - The application configuration.
// @param db *sql.DB - The database connection.
// @return *TodoController - A pointer to the new TodoController.
func NewTodoControl(cfg *config.Config, db *sql.DB) *TodoController {
	// A new TodoController is returned.
	return &TodoController{
		// The cfg field is set to the application configuration.
		cfg: cfg,
		// The db field is set to the database connection.
		db: db,
	}
}

// MatchCurrentUserWithTodoOwner checks if the current user is the owner of the todo.
// It takes a TodoController, a todo ID, and a current user ID as input.
//
// @param tc *TodoController - The TodoController.
// @param todoId uuid.UUID - The ID of the todo.
// @param currentUserId uuid.UUID - The ID of the current user.
// @return bool - True if the current user is the owner of the todo, false otherwise.
// @return error - An error if one occurred.
func MatchCurrentUserWithTodoOwner(tc *TodoController, todoId uuid.UUID, currentUserId uuid.UUID) (bool, error) {
	// userId is a variable that will hold the ID of the todo's owner.
	var userId uuid.UUID

	// err is the result of querying the database for the todo's owner.
	err := tc.db.QueryRow(GetTodoUserQuery, todoId).Scan(&userId)
	// This checks if an error occurred while querying the database.
	if err != nil {
		// If an error occurs, false and the error are returned.
		return false, err
	}

	// The function returns true if the todo's owner ID matches the current user's ID.
	return userId == currentUserId, nil
}

// CreateTodoController handles the creation of a new todo.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (tc *TodoController) CreateTodoController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(users.User)

	// body is a new Create_UpdateTodoRequest struct.
	body := new(Create_UpdateTodoRequest)
	// This parses the request body into the body struct.
	if err := c.BodyParser(body); err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	// This checks if the title is empty.
	if body.Title == "" {
		// If the title is empty, a bad request response is returned.
		return response.BadResponse(c, "Title is required")
	}

	// todoId is the new UUID for the todo.
	todoId, _ := uuid.NewV7()

	// todo is a new Todo struct.
	todo := Todo{
		// The ID field is set to the new UUID.
		ID: todoId,
		// The Title field is set to the todo's title.
		Title: body.Title,
		// The Completed field is set to false.
		Completed: false,
		// The Owner field is set to the current user's ID.
		Owner: user.ID.String(),
		// The CreatedAt field is set to the user's creation time.
		CreatedAt: utils.ParseTime(user.CreatedAt),
	}

	// _, err is the result of executing the SQL query to create the new todo.
	_, err := tc.db.Exec(CreateTodoQuery, todo.ID, todo.Title, todo.Completed, todo.Owner, todo.CreatedAt)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Unable to create todo")
	}

	// todoResponse is a new TodoResponse struct.
	todoResponse := TodoResponse{
		// The ID field is set to the todo's ID.
		ID: todo.ID,
		// The Title field is set to the todo's title.
		Title: todo.Title,
		// The Completed field is set to the todo's completion status.
		Completed: todo.Completed,
		// The CreatedAt field is set to the todo's creation time.
		CreatedAt: todo.CreatedAt,
	}

	// A created response is returned with a success message and the todo data.
	return response.OKCreatedResponse(c, "Todo created successfully", todoResponse)
}

// GetTodosController handles the retrieval of todos.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (tc *TodoController) GetTodosController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(users.User)

	// completedQuery is the value of the "completed" query parameter.
	completedQuery := c.Query("completed")
	// completed is the boolean value of the "completed" query parameter.
	completed := c.QueryBool("completed")

	// page is the value of the "page" query parameter, with a default of 1.
	page := c.QueryInt("page", 1)
	// This ensures that the page number is at least 1.
	if page <= 0 {
		// If the page number is less than or equal to 0, it is set to 1.
		page = 1
	}

	// limit is the value of the "limit" query parameter, with a default of 10.
	limit := c.QueryInt("limit", 10)
	// This ensures that the limit is at least 1.
	if limit <= 0 {
		// If the limit is less than or equal to 0, it is set to 10.
		limit = 10
	// This ensures that the limit is at most 100.
	} else if limit > 100 {
		// If the limit is greater than 100, it is set to 100.
		limit = 100
	}

	// totalItems is a variable that will hold the total number of todos.
	var totalItems int64
	// err is a variable that will hold any errors that occur.
	var err error

	// This checks if the "completed" query parameter is empty.
	if completedQuery == "" {
		// If it is empty, the total number of todos for the user is retrieved.
		err = tc.db.QueryRow(CountTodosByUserQuery, user.ID).Scan(&totalItems)
	} else {
		// If it is not empty, the total number of todos for the user, filtered by completion status, is retrieved.
		err = tc.db.QueryRow(CountTodosByUserFilteredByCompletedQuery, user.ID, completed).Scan(&totalItems)
	}
	// This checks if an error occurred while querying the database.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Failed to retrieve todo count")
	}

	// This checks if there are no todos.
	if totalItems == 0 {
		// If there are no todos, an OK response is returned with an empty list of todos.
		return response.OKResponse(c, "Todos fetched successfully", PaginatedTodoResponse{
			Results: []TodoResponse{},
			Count: 0,
			TotalItems: 0,
			TotalPages: 0,
			Page: page,
			Limit: limit,
		})
	}

	// todos is a slice that will hold the retrieved todos.
	var todos []TodoResponse
	// rows is a variable that will hold the result of the database query.
	var rows *sql.Rows

	// totalPages is the total number of pages.
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	// This ensures that the page number is not greater than the total number of pages.
	if page > totalPages {
		// If the page number is greater than the total number of pages, it is set to the total number of pages.
		page = totalPages
	}

	// offset is the number of todos to skip.
	offset := (page - 1) * limit

	// This checks if the "completed" query parameter is empty.
	if completedQuery == "" {
		// If it is empty, all todos for the user are retrieved.
		rows, err = tc.db.Query(GetTodosByUserQuery, user.ID, limit, offset)
	} else {
		// If it is not empty, all todos for the user, filtered by completion status, are retrieved.
		rows, err = tc.db.Query(GetTodosByUserFilteredByCompletedQuery, user.ID, completed, limit, offset)
	}

	// This checks if an error occurred while querying the database.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to retrieve todos"})
	}
	// This defers the closing of the rows until the function returns.
	defer rows.Close()

	// This iterates over the rows.
	for rows.Next() {
		// todo is a new Todo struct.
		var todo Todo

		// err is the result of scanning the row into the todo struct.
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
		// This checks if an error occurred while scanning the row.
		if err != nil {
			// If an error occurs, an internal server error response is returned.
			return response.InternelServerError(c, err, "Unable to get todos")
		}

		// The todo is appended to the todos slice.
		todos = append(todos, TodoResponse{
			ID: todo.ID,
			Title: todo.Title,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
		})
	}

	// paginatedTodoResponse is a new PaginatedTodoResponse struct.
	paginatedTodoResponse := PaginatedTodoResponse{
		// The Results field is set to the retrieved todos.
		Results: todos,
		// The Count field is set to the number of retrieved todos.
		Count: len(todos),
		// The TotalItems field is set to the total number of todos.
		TotalItems: totalItems,
		// The TotalPages field is set to the total number of pages.
		TotalPages: totalPages,
		// The Page field is set to the current page number.
		Page: page,
		// The Limit field is set to the number of todos per page.
		Limit: limit,
	}

	// An OK response is returned with a success message and the paginated todo data.
	return response.OKResponse(c, "Todo fetched successfully", paginatedTodoResponse)
}

// UpdateTodoController handles the update of a todo.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (tc *TodoController) UpdateTodoController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(users.User)

	// todoId is the value of the "id" path parameter.
	todoId := c.Params("id")
	// This checks if the todo ID is empty.
	if todoId == "" {
		// If the todo ID is empty, a bad request response is returned.
		return response.BadResponse(c, "Todo id is required")
	}

	// matchedCurrentUserWithTodoOwner is a boolean that indicates whether the current user is the owner of the todo.
	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	// This checks if the current user is not the owner of the todo.
	if !matchedCurrentUserWithTodoOwner {
		// If the current user is not the owner of the todo, an unauthorized access response is returned.
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	// body is a new Create_UpdateTodoRequest struct.
	body := new(Create_UpdateTodoRequest)
	// This parses the request body into the body struct.
	if err := c.BodyParser(body); err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	// This checks if the title is empty.
	if body.Title == "" {
		// If the title is empty, a bad request response is returned.
		return response.BadResponse(c, "Title is required")
	}

	// todo is a new Todo struct.
	var todo Todo

	// err is the result of executing the SQL query to update the todo.
	err = tc.db.QueryRow(UpdateTodoTitleQuery, body.Title, todoId).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Unable to update todo")
	}

	// todoResponse is a new TodoResponse struct.
	todoResponse := TodoResponse{
		// The ID field is set to the todo's ID.
		ID: todo.ID,
		// The Title field is set to the todo's title.
		Title: todo.Title,
		// The Completed field is set to the todo's completion status.
		Completed: todo.Completed,
		// The CreatedAt field is set to the todo's creation time.
		CreatedAt: todo.CreatedAt,
	}

	// An OK response is returned with a success message and the updated todo data.
	return response.OKResponse(c, "Todo updated successfully", todoResponse)
}

// DeleteTodoController handles the deletion of a todo.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (tc *TodoController) DeleteTodoController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(users.User)

	// todoId is the value of the "id" path parameter.
	todoId := c.Params("id")
	// This checks if the todo ID is empty.
	if todoId == "" {
		// If the todo ID is empty, a bad request response is returned.
		return response.BadResponse(c, "Todo id is required")
	}

	// matchedCurrentUserWithTodoOwner is a boolean that indicates whether the current user is the owner of the todo.
	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	// This checks if the current user is not the owner of the todo.
	if !matchedCurrentUserWithTodoOwner {
		// If the current user is not the owner of the todo, an unauthorized access response is returned.
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	// _, err is the result of executing the SQL query to delete the todo.
	_, err = tc.db.Exec(DeleteTodoQuery, todoId)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Unable to delete todo")
	}

	// An OK response is returned with a success message and the deleted todo's ID.
	return response.OKResponse(c, "Todo deleted successfully", fiber.Map{"todo_id": todoId})
}

// CompleteTodoController handles the completion of a todo.
// It takes a Fiber context as input.
//
// @param c *fiber.Ctx - The Fiber context.
// @return error - An error if one occurred.
func (tc *TodoController) CompleteTodoController(c *fiber.Ctx) error {
	// user is the User object retrieved from the local context.
	user := c.Locals("user").(users.User)

	// todoId is the value of the "id" path parameter.
	todoId := c.Params("id")
	// This checks if the todo ID is empty.
	if todoId == "" {
		// If the todo ID is empty, a bad request response is returned.
		return response.BadResponse(c, "Todo id is required")
	}

	// matchedCurrentUserWithTodoOwner is a boolean that indicates whether the current user is the owner of the todo.
	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	// This checks if the current user is not the owner of the todo.
	if !matchedCurrentUserWithTodoOwner {
		// If the current user is not the owner of the todo, an unauthorized access response is returned.
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	// body is a new CompleteTodoRequest struct.
	body := new(CompleteTodoRequest)
	// This parses the request body into the body struct.
	if err := c.BodyParser(body); err != nil {
		// If an error occurs, a bad request response is returned.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	// todo is a new Todo struct.
	var todo Todo

	// err is the result of executing the SQL query to update the todo's completion status.
	err = tc.db.QueryRow(UpdateTodoCompletedQuery, body.Completed, todoId).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
	// This checks if an error occurred while executing the query.
	if err != nil {
		// If an error occurs, an internal server error response is returned.
		return response.InternelServerError(c, err, "Unable to update todo")
	}

	// todoResponse is a new TodoResponse struct.
	todoResponse := TodoResponse{
		// The ID field is set to the todo's ID.
		ID: todo.ID,
		// The Title field is set to the todo's title.
		Title: todo.Title,
		// The Completed field is set to the todo's completion status.
		Completed: todo.Completed,
		// The CreatedAt field is set to the todo's creation time.
		CreatedAt: todo.CreatedAt,
	}

	// An OK response is returned with a success message and the updated todo data.
	return response.OKResponse(c, "Todo updated successfully", todoResponse)
}