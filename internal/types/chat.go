package types

// Message represents a chat message with role and content
// @Description A message in a chat conversation
type Message struct {
	// Role of the message sender (e.g., user, assistant)
	// example: user
	Role string `json:"role" example:"user"`
	// Content of the message
	// example: Hello, how are you today?
	Content string `json:"content" example:"Hello, how are you today?"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type ErrorResponse struct {
	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   string      `json:"param,omitempty"`
		Code    interface{} `json:"code,omitempty"`
	} `json:"error"`
}

// ChatRequest represents a chat completion request
// @Description Request payload for chat completions
type ChatRequest struct {
	// Array of messages in the conversation
	Messages []Message `json:"messages" example:"[{\"role\":\"user\",\"content\":\"Tell me about artificial intelligence\"}]"`
	// Model ID to use for completion
	Model string `json:"model" example:"deepseek-r1-distill-llama-70b"`
	// Sampling temperature between 0 and 2
	Temperature float64 `json:"temperature,omitempty" example:"0.7"`
	// Maximum number of tokens to generate
	MaxTokens int `json:"max_tokens,omitempty" example:"100"`
	// Nucleus sampling parameter
	TopP float64 `json:"top_p,omitempty" example:"1.0"`
	// Frequency penalty for token generation
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty" example:"0"`
	// Presence penalty for token generation
	PresencePenalty float64 `json:"presence_penalty,omitempty" example:"0"`
	// Whether to stream the response
	Stream bool `json:"stream,omitempty" example:"false"`
	// Sequences to stop generation
	Stop []string `json:"stop,omitempty" example:"[\"END\",\"STOP\"]"`
	// Number of completions to generate
	N int `json:"n,omitempty" example:"1"`
	// Optional user identifier
	User string `json:"user,omitempty" example:"user123"`
}
