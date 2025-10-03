package todos

import "github.com/google/uuid"

type Todo struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Owner     string    `json:"owner"`
	CreatedAt string    `json:"created_at"`
}
