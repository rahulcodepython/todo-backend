package todos

import "time"

type Todo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Todo      string    `json:"todo"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
