package api

import (
	"io"
	"net/http"

	"go-api/internal/types"
	"go-api/pkg/groq"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	groqClient *groq.Client
}

func NewHandler(apiKey string) *Handler {
	return &Handler{
		groqClient: groq.NewClient(apiKey),
	}
}

func (h *Handler) HandleChatCompletions(c echo.Context) error {
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

	// Make request to Groq API
	resp, err := h.groqClient.ChatCompletion(&chatReq)
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
