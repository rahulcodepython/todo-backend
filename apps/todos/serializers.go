package todos

import "github.com/google/uuid"

type Create_UpdateTodoRequest struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
}

type CompleteTodoRequest struct {
	Completed *bool `json:"completed" validate:"required"`
}

type TodoResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt string    `json:"created_at"`
}

type PaginatedTodoResponse struct {
	Results    []TodoResponse `json:"results"`
	Count      int            `json:"count"`
	TotalItems int64          `json:"total_items"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}
