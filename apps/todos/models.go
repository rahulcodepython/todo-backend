// This file defines the data model for todos.
package todos

// "github.com/google/uuid" is a package for working with UUIDs. It is used here to define the ID field.
import "github.com/google/uuid"

// Todo represents the structure of a todo item in the application.
type Todo struct {
	// ID is the unique identifier for the todo.
	// json:"id" specifies that this field should be marshalled to/from a JSON object with the key "id".
	ID uuid.UUID `json:"id"`
	// Title is the title of the todo.
	// json:"title" specifies that this field should be marshalled to/from a JSON object with the key "title".
	Title string `json:"title"`
	// Completed is the completion status of the todo.
	// json:"completed" specifies that this field should be marshalled to/from a JSON object with the key "completed".
	Completed bool `json:"completed"`
	// Owner is the ID of the user who owns the todo.
	// json:"owner" specifies that this field should be marshalled to/from a JSON object with the key "owner".
	Owner string `json:"owner"`
	// CreatedAt is the time the todo was created.
	// json:"created_at" specifies that this field should be marshalled to/from a JSON object with the key "created_at".
	CreatedAt string `json:"created_at"`
}