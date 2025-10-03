// This file contains a test script for making API requests to the todo-backend.
// It includes functionality to send POST requests to create new todo items,
// complete with JSON data and JWT authentication.
package main

// "bytes" provides functions for manipulating byte slices. It is used here to create an io.Reader from the JSON request body.
import (
	"bytes"
	// "encoding/json" provides functions for encoding and decoding JSON data. It is used here to marshal the Todo struct into a JSON byte slice.
	"encoding/json"
	// "fmt" provides functions for formatted I/O. It is used here for logging and printing output to the console.
	"fmt"
	// "io" provides basic interfaces to I/O primitives. It is used here to read the HTTP response body.
	"io"
	// "log" provides a simple logging package. It is used here to log fatal errors.
	"log"
	// "net/http" provides HTTP client and server implementations. It is used here to create and send HTTP requests.
	"net/http"
)

// Todo represents the data structure for a single to-do item in the application.
type Todo struct {
	// Title is the string containing the main task description.
	// json:"title" specifies that this field should be marshalled to/from a JSON object with the key "title".
	Title string `json:"title"`
	// Completed is a boolean flag indicating if the task is finished.
	// json:"completed" specifies that this field should be marshalled to/from a JSON object with the key "completed".
	Completed bool `json:"completed"`
}

// api_request sends a POST request to the API to create a new todo item.
// It constructs the request with a JSON body and an authorization token,
// sends the request, and processes the response.
//
// @param id int - An integer used to generate a unique title for the todo item.
func api_request(id int) {
	// apiURL stores the URL of the local Fiber API endpoint for creating todos.
	apiURL := "http://127.0.0.1:8000/api/v1/todos/create"

	// authToken stores the JWT token for authentication.
	// IMPORTANT: This should be replaced with a valid token from the /login endpoint.
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTk2MTIzMDAsImlhdCI6MTc1OTUyNTkwMCwidXNlcl9pZCI6IjAxOTlhYmVhLTUzMDgtNzBhZS05ZThmLWJlNDNjODA3NWUzNyJ9.UfXJB0UT60eCkigRLK2o7HW3V9l3d0i4KL2UuGLXqO4"

	// newTodo is an instance of the Todo struct, representing the data to be sent in the request body.
	newTodo := Todo{
		// Title is dynamically generated using the id parameter.
		Title: fmt.Sprintf("Complete the task %d", id),
		// Completed is set to true if the id is even, and false otherwise.
		Completed: id%2 == 0,
	}

	// jsonData is a byte slice that will hold the JSON representation of the newTodo struct.
	// json.Marshal is used to convert the struct into JSON.
	jsonData, err := json.Marshal(newTodo)
	// This checks if an error occurred during the JSON marshalling process.
	if err != nil {
		// If an error occurs, log the error and exit the program.
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// req is a new HTTP request object.
	// http.NewRequest creates a POST request to the specified apiURL with the JSON data as the body.
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	// This checks if an error occurred while creating the new HTTP request.
	if err != nil {
		// If an error occurs, log the error and exit the program.
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// This sets the "Content-Type" header of the request to "application/json".
	req.Header.Set("Content-Type", "application/json")

	// This sets the "Authorization" header of the request for JWT authentication.
	// The token is prefixed with "Bearer ".
	req.Header.Set("Authorization", "Bearer "+authToken)

	// client is a new HTTP client.
	client := &http.Client{}

	// This prints a message to the console indicating that the request is being sent.
	fmt.Println("Sending request to API...")
	// resp holds the HTTP response from the server.
	// client.Do sends the HTTP request.
	resp, err := client.Do(req)
	// This checks if an error occurred while sending the request.
	if err != nil {
		// If an error occurs, log the error and exit the program.
		log.Fatalf("Error sending request: %v", err)
	}
	// This defers the closing of the response body until the function returns.
	defer resp.Body.Close()

	// body is a byte slice that will hold the content of the HTTP response body.
	// io.ReadAll reads the entire response body.
	body, err := io.ReadAll(resp.Body)
	// This checks if an error occurred while reading the response body.
	if err != nil {
		// If an error occurs, log the error and exit the program.
		log.Fatalf("Error reading response body: %v", err)
	}

	// This prints a separator line to the console.
	fmt.Println("---------------------------------")
	// This prints the HTTP status code of the response.
	fmt.Println("Status Code:", resp.StatusCode)
	// This prints the response body as a string.
	fmt.Println("Response Body:", string(body))
	// This prints a separator line to the console.
	fmt.Println("---------------------------------")

	// This checks if the HTTP status code of the response is 201 (Created).
	if resp.StatusCode == http.StatusCreated {
		// If the status code is 201, print a success message.
		fmt.Println("✅ Successfully created a new todo!")
	} else {
		// If the status code is not 201, print a failure message.
		fmt.Println("❌ Failed to create todo. Check the status code and response body for details.")
	}
}

// main is the entry point of the program.
// It calls the api_request function in a loop to send multiple requests.
func main() {
	// This loop iterates from 1 to 50.
	for i := 1; i <= 50; i++ {
		// In each iteration, call the api_request function with the current value of i.
		api_request(i)
	}
}