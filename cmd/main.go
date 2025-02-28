package main

// @title Scarlett API
// @version 1.0
// @description A Go API service for the Scarlett Protocol that provides chat completion functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.scarlett.ai/support
// @contact.email help@scarlett.ai

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description API key authentication with Bearer prefix (e.g., "Bearer your-api-key"). The 'Bearer ' prefix is REQUIRED - requests without it will be rejected.

import (
	"log"
	"os"

	// Import swagger docs
	_ "go-api/docs/swagger"
	"go-api/internal/middleware"
	"go-api/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

const (
	// DefaultPort is the port used when no PORT environment variable is set
	DefaultPort = "8082"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())

	// Add rate limiter middleware
	e.Use(middleware.DefaultRateLimiter())

	// Register Swagger documentation routes
	routes.RegisterSwaggerRoutes(e)

	// Register API routes
	routes.RegisterRoutes(e)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	e.Logger.Fatal(e.Start(":" + port))
}
