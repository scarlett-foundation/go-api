package routes

import (
	"go-api/internal/handlers"
	"go-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all chat-related routes
// @title Chat API Routes
// @description Routes for chat functionality
// @Security BearerAuth
func RegisterRoutes(e *echo.Echo) {
	// Register chat routes
	// Chat completions endpoint for interacting with Groq API
	e.POST("/chat/completions", handlers.HandleChatCompletions, middleware.APIKeyAuth())
}
