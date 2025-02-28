package models

import "go-api/internal/types"

// ChatRequestExample is an example of a chat request
// This will be used by Swagger to show an example in the UI
// @Description Example chat request
type ChatRequestExample struct {
	// Model ID
	Model string `json:"model" example:"deepseek-r1-distill-llama-70b"`
	// Array of messages in the conversation
	Messages []struct {
		// Role of the message sender
		Role string `json:"role" example:"user"`
		// Content of the message
		Content string `json:"content" example:"Hello, how are you?"`
	} `json:"messages"`
	// Sampling temperature
	Temperature float64 `json:"temperature" example:"0.7"`
	// Maximum number of tokens to generate
	MaxTokens int `json:"max_tokens" example:"100"`
	// Whether to stream the response
	Stream bool `json:"stream" example:"false"`
}

// GetChatCompletionExample returns a concrete example for documentation
func GetChatCompletionExample() types.ChatRequest {
	return types.ChatRequest{
		Model: "deepseek-r1-distill-llama-70b",
		Messages: []types.Message{
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   100,
		Stream:      false,
	}
}
