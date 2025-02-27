package routes

import (
	"go-api/internal/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all chat-related routes
func RegisterRoutes(e *echo.Echo) {
	e.POST("/chat/completions", handlers.HandleChatCompletions)
}
