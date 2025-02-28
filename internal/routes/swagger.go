package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegisterSwaggerRoutes registers all Swagger documentation routes
func RegisterSwaggerRoutes(e *echo.Echo) {
	// Set up Swagger handler with options
	swaggerHandler := echoSwagger.EchoWrapHandler(
		echoSwagger.URL("/swagger/doc.json"),
		echoSwagger.DocExpansion("list"),
		echoSwagger.DeepLinking(true),
	)

	// Main Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Additional Swagger route for better browser compatibility
	e.GET("/swagger/index.html", swaggerHandler)
}
