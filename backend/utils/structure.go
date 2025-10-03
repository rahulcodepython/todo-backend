// This file defines the standard structure for API responses.
package utils

// Response represents the structure of a standard API response.
// It includes fields for success status, a message, data, and an error.
type Response struct {
	// Success indicates whether the API request was successful.
	// json:"success" specifies that this field should be marshalled to/from a JSON object with the key "success".
	Success bool `json:"success"`
	// Message provides a human-readable description of the response.
	// json:"message" specifies that this field should be marshalled to/from a JSON object with the key "message".
	Message string `json:"message"`
	// Data holds the actual payload of the response.
	// It is an empty interface to allow for any type of data.
	// json:"data,omitempty" specifies that this field should be marshalled to/from a JSON object with the key "data", and should be omitted if empty.
	Data interface{} `json:"data,omitempty"`
	// Error holds detailed error information.
	// It is an empty interface to allow for various error structures.
	// json:"error,omitempty" specifies that this field should be marshalled to/from a JSON object with the key "error", and should be omitted if empty.
	Error interface{} `json:"error,omitempty"`
}