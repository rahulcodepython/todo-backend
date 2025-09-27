package users

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Omit from JSON responses
	CreatedAt time.Time `json:"created_at"`
}
