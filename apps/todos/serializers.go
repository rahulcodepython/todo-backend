package todos

type CreateTodoInput struct {
	Todo string `json:"todo" validate:"required"`
}

type UpdateTodoInput struct {
	Todo string `json:"todo" validate:"required"`
}

type TodoResponse struct {
	ID        string `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PaginatedTodoResponse struct {
	Data       []TodoResponse `json:"data"`
	Count      int            `json:"count"`
	TotalItems int64          `json:"total_items"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}
