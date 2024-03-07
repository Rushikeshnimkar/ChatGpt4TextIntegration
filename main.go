package main

import (
	"encoding/json"
	"fmt"
	"gpt/utils"
	"log"
	"net/http"

	"github.com/joho/godotenv" // If you're using environment variables
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define a handler function for the root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var requestBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Get the prompt from the request body
		prompt, ok := requestBody["prompt"]
		if !ok {
			http.Error(w, "Prompt not found in request body", http.StatusBadRequest)
			return
		}

		// Call GetChatGptResponse function from the utils package
		response, err := utils.GetChatGptResponse(prompt)
		if err != nil {
			http.Error(w, "Error getting ChatGPT response", http.StatusInternalServerError)
			return
		}

		// Write the response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"response": "%s"}`, response)
	})

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
