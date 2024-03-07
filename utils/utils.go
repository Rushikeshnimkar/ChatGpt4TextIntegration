package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

// ChatGPTRequest represents the JSON request structure for ChatGPT
type ChatGPTRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// GetChatGptResponse sends a prompt to OpenAI's ChatGPT API and returns the response
func GetChatGptResponse(prompt string) (string, error) {
	var (
		ChatGptURL = "https://api.openai.com/v1/chat/completions"
	)

	// Get the OpenAI API key from environment variables
	openAIKey := os.Getenv("OPENAI_KEY")
	if openAIKey == "" {
		return "", errors.New("OPENAI_KEY not found in environment variables")
	}

	// Create the request body struct
	requestData := ChatGPTRequest{
		Model: "gpt-4-turbo-preview",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal request body to JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST", ChatGptURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)

	// Create HTTP client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("non-OK status code received: " + resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return response body as string
	return string(body), nil
}
