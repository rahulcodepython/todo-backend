package todos

type Create_UpdateTodoRequest struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
}

type CompleteTodoRequest struct {
	Completed *bool `json:"completed" validate:"required"`
}

type TodoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}
