package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const groqEndpoint = "https://api.groq.com/openai/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

type ChatRequest struct {
	Messages         []Message `json:"messages"`
	Model            string    `json:"model"`
	Temperature      float64   `json:"temperature,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	N                int       `json:"n,omitempty"`
	User             string    `json:"user,omitempty"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/chat/completions", handleChatCompletions)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func handleChatCompletions(c echo.Context) error {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "GROQ_API_KEY not set",
				Type:    "internal_error",
			},
		})
	}

	// Parse request body
	var chatReq ChatRequest
	if err := c.Bind(&chatReq); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "Invalid request body",
				Type:    "invalid_request_error",
			},
		})
	}

	// Set default values
	if chatReq.N == 0 {
		chatReq.N = 1
	}

	// Prepare request to Groq API
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "Failed to marshal request",
				Type:    "internal_error",
			},
		})
	}

	// Create request to Groq API
	req, err := http.NewRequest("POST", groqEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "Failed to create request",
				Type:    "internal_error",
			},
		})
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "Failed to make request to Groq API",
				Type:    "api_error",
			},
		})
	}
	defer resp.Body.Close()

	// If streaming is requested, stream the response
	if chatReq.Stream {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")

		// Copy the stream directly to the client
		_, err = io.Copy(c.Response().Writer, resp.Body)
		return err
	}

	// For non-streaming responses, just proxy the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: struct {
				Message string      `json:"message"`
				Type    string      `json:"type"`
				Param   string      `json:"param,omitempty"`
				Code    interface{} `json:"code,omitempty"`
			}{
				Message: "Failed to read response from Groq API",
				Type:    "api_error",
			},
		})
	}

	// Return response with same status code and body
	return c.JSONBlob(resp.StatusCode, body)
}
