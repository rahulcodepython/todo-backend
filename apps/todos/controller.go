package todos

import (
	"database/sql"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/config"
	"github.com/rahulcodepython/todo-backend/backend/response"
	"github.com/rahulcodepython/todo-backend/backend/utils"
)

type TodoController struct {
	cfg *config.Config
	db  *sql.DB
}

func NewTodoControl(cfg *config.Config, db *sql.DB) *TodoController {
	return &TodoController{
		cfg: cfg,
		db:  db,
	}
}

func MatchCurrentUserWithTodoOwner(tc *TodoController, todoId uuid.UUID, currentUserId uuid.UUID) (bool, error) {
	var userId uuid.UUID

	err := tc.db.QueryRow(GetTodoUserQuery, todoId).Scan(&userId)
	if err != nil {
		return false, err
	}

	return userId == currentUserId, nil
}

func (tc *TodoController) CreateTodoController(c *fiber.Ctx) error {
	user := c.Locals("user").(users.User)

	body := new(Create_UpdateTodoRequest)
	if err := c.BodyParser(body); err != nil {
		// If parsing fails, it sends a standardized bad request response to the client.
		// This centralizes the parsing and error handling logic.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	if body.Title == "" {
		// If any field is missing, send a 400 Bad Request response.
		return response.BadResponse(c, "Title is required")
	}

	todoId, _ := uuid.NewV7()

	todo := Todo{
		ID:        todoId,
		Title:     body.Title,
		Completed: false,
		Owner:     user.ID.String(),
		CreatedAt: utils.ParseTime(user.CreatedAt),
	}

	_, err := tc.db.Exec(CreateTodoQuery, todo.ID, todo.Title, todo.Completed, todo.Owner, todo.CreatedAt)
	if err != nil {
		return response.BadInternalResponse(c, err, "Unable to create todo")
	}

	todoResponse := TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}

	return response.OKCreatedResponse(c, "Todo created successfully", todoResponse)
}
func (tc *TodoController) GetTodosController(c *fiber.Ctx) error {
	user := c.Locals("user").(users.User)

	completedQuery := c.Query("completed")
	completed := c.QueryBool("completed")

	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	var totalItems int64
	var err error

	if completedQuery == "" {
		err = tc.db.QueryRow(CountTodosByUserQuery, user.ID).Scan(&totalItems)
	} else {
		err = tc.db.QueryRow(CountTodosByUserFilteredByCompletedQuery, user.ID, completed).Scan(&totalItems)
	}
	if err != nil {
		return response.InternelServerError(c, err, "Failed to retrieve todo count")
	}

	if totalItems == 0 {
		// If there are no items, return an empty response immediately.
		return response.OKResponse(c, "Todos fetched successfully", PaginatedTodoResponse{
			Results:    []TodoResponse{},
			Count:      0,
			TotalItems: 0,
			TotalPages: 0,
			Page:       page,
			Limit:      limit,
		})
	}

	var todos []TodoResponse
	var rows *sql.Rows

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * limit

	if completedQuery == "" {
		rows, err = tc.db.Query(GetTodosByUserQuery, user.ID, limit, offset)
	} else {
		rows, err = tc.db.Query(GetTodosByUserFilteredByCompletedQuery, user.ID, completed, limit, offset)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to retrieve todos"})
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
		if err != nil {
			return response.InternelServerError(c, err, "Unable to get todos")
		}

		todos = append(todos, TodoResponse{
			ID:        todo.ID,
			Title:     todo.Title,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
		})
	}

	paginatedTodoResponse := PaginatedTodoResponse{
		Results:    todos,
		Count:      len(todos),
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}

	return response.OKResponse(c, "Todo fetched successfully", paginatedTodoResponse)
}
func (tc *TodoController) UpdateTodoController(c *fiber.Ctx) error {
	user := c.Locals("user").(users.User)

	todoId := c.Params("id")
	if todoId == "" {
		return response.BadResponse(c, "Todo id is required")
	}

	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	if !matchedCurrentUserWithTodoOwner {
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	body := new(Create_UpdateTodoRequest)
	if err := c.BodyParser(body); err != nil {
		// If parsing fails, it sends a standardized bad request response to the client.
		// This centralizes the parsing and error handling logic.
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	if body.Title == "" {
		return response.BadResponse(c, "Title is required")
	}

	var todo Todo

	err = tc.db.QueryRow(UpdateTodoTitleQuery, body.Title, todoId).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
	if err != nil {
		return response.InternelServerError(c, err, "Unable to update todo")
	}

	todoResponse := TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}

	return response.OKResponse(c, "Todo updated successfully", todoResponse)
}
func (tc *TodoController) DeleteTodoController(c *fiber.Ctx) error {
	user := c.Locals("user").(users.User)

	todoId := c.Params("id")
	if todoId == "" {
		return response.BadResponse(c, "Todo id is required")
	}

	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	if !matchedCurrentUserWithTodoOwner {
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	_, err = tc.db.Exec(DeleteTodoQuery, todoId)
	if err != nil {
		return response.InternelServerError(c, err, "Unable to delete todo")
	}

	return response.OKResponse(c, "Todo deleted successfully", fiber.Map{"todo_id": todoId})
}
func (tc *TodoController) CompleteTodoController(c *fiber.Ctx) error {
	user := c.Locals("user").(users.User)

	todoId := c.Params("id")
	if todoId == "" {
		return response.BadResponse(c, "Todo id is required")
	}

	matchedCurrentUserWithTodoOwner, err := MatchCurrentUserWithTodoOwner(tc, uuid.MustParse(todoId), user.ID)
	if !matchedCurrentUserWithTodoOwner {
		return response.UnauthorizedAccess(c, err, "You are not authorized to update this todo")
	}

	body := new(CompleteTodoRequest)
	if err := c.BodyParser(body); err != nil {
		return response.BadInternalResponse(c, err, "Invalid request body")
	}

	var todo Todo

	err = tc.db.QueryRow(UpdateTodoCompletedQuery, body.Completed, todoId).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Owner, &todo.CreatedAt)
	if err != nil {
		return response.InternelServerError(c, err, "Unable to update todo")
	}

	todoResponse := TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
	}

	return response.OKResponse(c, "Todo updated successfully", todoResponse)
}
