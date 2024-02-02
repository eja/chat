package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int `json:"index"`
		Message      Message
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func makeOpenAIAPIRequest(model string, messages []Message) (string, error) {
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
	req.Header.Set("Authorization", "Bearer "+chatOptions.OpenaiToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var response ChatCompletionResponse
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

/*
func main() {
	model := "gpt-3.5-turbo"
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Who won the world series in 2020?"},
		{Role: "assistant", Content: "The Los Angeles Dodgers won the World Series in 2020."},
		{Role: "user", Content: "Where was it played?"},
	}

	content, err := makeOpenAIAPIRequest(model, messages)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Assistant's Response:", content)
}
*/
