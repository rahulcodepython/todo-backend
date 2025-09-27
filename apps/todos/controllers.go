package todos

import (
	"database/sql"
	"log"
	"math"
	"time"

	"github.com/rahulcodepython/todo-backend/apps/users"
	"github.com/rahulcodepython/todo-backend/backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TodosController struct {
	DB *sql.DB
}

func NewTodosController(db *sql.DB) *TodosController {
	return &TodosController{DB: db}
}

// Helper to get user from context
func getUserFromCtx(c *fiber.Ctx) (*users.User, error) {
	user, ok := c.Locals(utils.CtxUserKey).(*users.User)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return user, nil
}

// CreateTodo creates a new todo item.
func (tc *TodosController) CreateTodo(c *fiber.Ctx) error {
	user, err := getUserFromCtx(c)
	if err != nil {
		return err
	}

	body := new(CreateTodoInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	now := time.Now()
	todo := &Todo{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Todo:      body.Todo,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = tc.DB.Exec(CreateTodoQuery, todo.ID, todo.UserID, todo.Todo, todo.Completed, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create todo"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": todo})
}

// GetAllTodos retrieves all todo items for the user.
func (tc *TodosController) GetAllTodos(c *fiber.Ctx) error {
	user, err := getUserFromCtx(c)
	if err != nil {
		return err
	}

	var rows *sql.Rows

	completed := c.Query("completed")
	completedFilter := c.QueryBool("completed")

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

	err = tc.DB.QueryRow(CountTodosByUserIDQuery, user.ID).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting todos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to retrieve todo count"})
	}

	if totalItems == 0 {
		// If there are no items, return an empty response immediately.
		return c.JSON(fiber.Map{
			"success": true,
			"data": PaginatedTodoResponse{
				Data:       []TodoResponse{},
				Count:      0,
				TotalItems: 0,
				TotalPages: 0,
				Page:       page,
				Limit:      limit,
			},
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * limit

	if completed == "" {
		log.Println("Getting all todos")
		rows, err = tc.DB.Query(GetTodosByUserIDQuery, user.ID, limit, offset)
	} else {
		log.Println("Getting filtered todos")
		rows, err = tc.DB.Query(GetTodosByUserIDFilteredComplitionQuery, user.ID, completedFilter, limit, offset)
	}

	if err != nil {
		log.Printf("Error getting todos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to retrieve todos"})
	}
	defer rows.Close()

	var todos []TodoResponse
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return err
		}
		todos = append(todos, TodoResponse{
			ID:        todo.ID,
			Todo:      todo.Todo,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt.Format(time.RFC3339),
			UpdatedAt: todo.UpdatedAt.Format(time.RFC3339),
		})
	}

	paginatedTodoResponse := PaginatedTodoResponse{
		Data:       todos,
		Count:      len(todos),
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}

	return c.JSON(fiber.Map{"success": true, "data": paginatedTodoResponse})
}

// UpdateTodo updates a todo item.
func (tc *TodosController) UpdateTodo(c *fiber.Ctx) error {
	user, err := getUserFromCtx(c)
	if err != nil {
		return err
	}

	todoID := c.Params("id")
	log.Println(todoID) // just for debugging
	body := new(UpdateTodoInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	log.Println(body) // just for debugging

	var todo Todo
	row := tc.DB.QueryRow(UpdateTodoQuery, body.Todo, time.Now(), todoID, user.ID)
	err = row.Scan(&todo.ID, &todo.Todo, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Todo not found"})
		}

		log.Printf("Error updating todo: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update todo"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Todo updated successfully", "data": todo})
}

// CompleteTodo toggles the completed status of a todo.
func (tc *TodosController) CompleteTodo(c *fiber.Ctx) error {
	user, err := getUserFromCtx(c)
	if err != nil {
		return err
	}

	todoID := c.Params("id")
	res, err := tc.DB.Exec(UpdateTodoStatusQuery, time.Now(), todoID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update todo status"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Todo not found or you don't have permission"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Todo status updated"})
}

// DeleteTodo deletes a todo item.
func (tc *TodosController) DeleteTodo(c *fiber.Ctx) error {
	user, err := getUserFromCtx(c)
	if err != nil {
		return err
	}

	todoID := c.Params("id")
	res, err := tc.DB.Exec(DeleteTodoQuery, todoID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete todo"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Todo not found or you don't have permission",
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Todo deleted successfully", "data": fiber.Map{"todo_id": todoID}})
}
