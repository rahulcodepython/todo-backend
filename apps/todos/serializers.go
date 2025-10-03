// This file defines the serializers for todo-related requests and responses.
package todos

// "github.com/google/uuid" is a package for working with UUIDs. It is used here to define the ID field in the response struct.
import "github.com/google/uuid"

// Create_UpdateTodoRequest defines the structure for a create or update todo request.
type Create_UpdateTodoRequest struct {
	// Title is the title of the todo.
	// json:"title" specifies that this field should be marshalled to/from a JSON object with the key "title".
	// validate:"required,min=3,max=255" specifies that this field is required, has a minimum length of 3, and a maximum length of 255.
	Title string `json:"title" validate:"required,min=3,max=255"`
}

// CompleteTodoRequest defines the structure for a complete todo request.
type CompleteTodoRequest struct {
	// Completed is the completion status of the todo.
	// json:"completed" specifies that this field should be marshalled to/from a JSON object with the key "completed".
	// validate:"required" specifies that this field is required.
	Completed *bool `json:"completed" validate:"required"`
}

// TodoResponse defines the structure for a todo response.
type TodoResponse struct {
	// ID is the unique identifier for the todo.
	// json:"id" specifies that this field should be marshalled to/from a JSON object with the key "id".
	ID uuid.UUID `json:"id"`
	// Title is the title of the todo.
	// json:"title" specifies that this field should be marshalled to/from a JSON object with the key "title".
	Title string `json:"title"`
	// Completed is the completion status of the todo.
	// json:"completed" specifies that this field should be marshalled to/from a JSON object with the key "completed".
	Completed bool `json:"completed"`
	// CreatedAt is the time the todo was created.
	// json:"created_at" specifies that this field should be marshalled to/from a JSON object with the key "created_at".
	CreatedAt string `json:"created_at"`
}

// PaginatedTodoResponse defines the structure for a paginated todo response.
type PaginatedTodoResponse struct {
	// Results is a slice of todos.
	// json:"results" specifies that this field should be marshalled to/from a JSON object with the key "results".
	Results []TodoResponse `json:"results"`
	// Count is the number of todos in the current page.
	// json:"count" specifies that this field should be marshalled to/from a JSON object with the key "count".
	Count int `json:"count"`
	// TotalItems is the total number of todos.
	// json:"total_items" specifies that this field should be marshalled to/from a JSON object with the key "total_items".
	TotalItems int64 `json:"total_items"`
	// TotalPages is the total number of pages.
	// json:"total_pages" specifies that this field should be marshalled to/from a JSON object with the key "total_pages".
	TotalPages int `json:"total_pages"`
	// Page is the current page number.
	// json:"page" specifies that this field should be marshalled to/from a JSON object with the key "page".
	Page int `json:"page"`
	// Limit is the number of todos per page.
	// json:"limit" specifies that this field should be marshalled to/from a JSON object with the key "limit".
	Limit int `json:"limit"`
}