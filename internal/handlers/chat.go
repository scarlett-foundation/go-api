package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"go-api/internal/types"

	"github.com/labstack/echo/v4"
)

const groqEndpoint = "https://api.groq.com/openai/v1/chat/completions"

// HandleChatCompletions handles the chat completions endpoint
func HandleChatCompletions(c echo.Context) error {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
	var chatReq types.ChatRequest
	if err := c.Bind(&chatReq); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
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
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
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
