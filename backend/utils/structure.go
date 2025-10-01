package utils

// Response defines the standard structure for API responses across the application.
// This struct is used to ensure consistency in how data, messages, and errors are returned
// to clients, making API consumption predictable and easier to handle.
type Response struct {
	// Success indicates whether the API request was successful.
	// It is a boolean field that will be serialized to JSON as "success".
	Success bool `json:"success"`
	// Message provides a human-readable description of the response.
	// This can be a success message, an informational message, or a general error message.
	Message string `json:"message"`
	// Data holds the actual payload of the response, typically for successful operations.
	// It is an empty interface to allow for any type of data and is omitted from JSON if empty.
	Data interface{} `json:"data,omitempty"`
	// Error holds detailed error information, typically for unsuccessful operations.
	// It is an empty interface to allow for various error structures and is omitted from JSON if empty.
	Error interface{} `json:"error,omitempty"`
}
