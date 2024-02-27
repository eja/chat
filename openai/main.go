package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eja/chat/internal/core"
	"net/http"
)

type typeOpenaiResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int `json:"index"`
		Message      TypeMessage
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

type TypeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Request(model string, messages []TypeMessage) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Define the request payload
	payload, err := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": messages,
	})
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+core.Options.OpenaiToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var response typeOpenaiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("Error decoding JSON response: %v", err)
	}

	// Check if there's a valid assistant message
	if len(response.Choices) > 0 {
		assistantMessage := response.Choices[0].Message
		return assistantMessage.Content, nil
	}

	return "", fmt.Errorf("No valid assistant message in the response")
}
