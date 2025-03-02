package routes

import (
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegisterSwaggerRoutes registers all Swagger documentation routes
func RegisterSwaggerRoutes(e *echo.Echo) {
	// Get API host from environment
	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = "localhost:8082"
	}

	// Force HTTPS in production, use both in development
	scheme := "http"
	if os.Getenv("ENVIRONMENT") == "production" {
		scheme = "https"
	}

	// Set up Swagger handler with options
	swaggerHandler := echoSwagger.EchoWrapHandler(
		echoSwagger.URL(scheme+"://"+apiHost+"/swagger/doc.json"),
		echoSwagger.DocExpansion("list"),
		echoSwagger.DeepLinking(true),
		echoSwagger.InstanceName("swagger"),
	)

	// Main Swagger endpoint
	e.GET("/swagger/*", swaggerHandler)
	// Additional Swagger route for better browser compatibility
	e.GET("/swagger/index.html", swaggerHandler)
}
