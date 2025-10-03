package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Todo struct defines the shape of our JSON data.
// It must match the structure the API expects.
type Todo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func api_request(id int) {
	// --- 1. Configuration ---

	// The URL of your local Fiber API endpoint
	apiURL := "http://127.0.0.1:8000/api/v1/todos/create"

	// IMPORTANT: Replace this with a valid JWT token you get from your /login endpoint
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTk2MTIzMDAsImlhdCI6MTc1OTUyNTkwMCwidXNlcl9pZCI6IjAxOTlhYmVhLTUzMDgtNzBhZS05ZThmLWJlNDNjODA3NWUzNyJ9.UfXJB0UT60eCkigRLK2o7HW3V9l3d0i4KL2UuGLXqO4"

	// --- 2. Prepare the Request Body ---

	// Create an instance of our Todo struct with the data we want to send
	newTodo := Todo{
		Title:     fmt.Sprintf("Complete the task %d", id),
		Completed: id%2 == 0,
	}

	// Marshal the Go struct into a JSON byte slice
	jsonData, err := json.Marshal(newTodo)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Create a new HTTP request object
	// We use bytes.NewBuffer to create an io.Reader from our JSON byte slice
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// --- 3. Set the Request Headers ---

	// Set the Content-Type header to indicate we are sending JSON
	req.Header.Set("Content-Type", "application/json")

	// Set the Authorization header for JWT authentication
	// The format is "Bearer YOUR_TOKEN"
	req.Header.Set("Authorization", "Bearer "+authToken)

	// --- 4. Send the Request ---

	// Create an HTTP client
	client := &http.Client{}

	// Send the request using the client
	fmt.Println("Sending request to API...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	// CRUCIAL: Always close the response body to prevent resource leaks
	defer resp.Body.Close()

	// --- 5. Process the Response ---

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the results
	fmt.Println("---------------------------------")
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
	fmt.Println("---------------------------------")

	// Check if the request was successful (201 Created is typical for a successful POST)
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Successfully created a new todo!")
	} else {
		fmt.Println("❌ Failed to create todo. Check the status code and response body for details.")
	}
}

func main() {
	for i := 1; i <= 50; i++ {
		api_request(i)
	}
}
