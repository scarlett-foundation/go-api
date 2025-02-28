package middleware

import (
	"net/http"
	"os"
	"strings"

	"go-api/internal/types"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

type APIKeys struct {
	Keys []string `yaml:"api_keys"`
}

// APIKeyAuth middleware validates the API key in request headers
func APIKeyAuth() echo.MiddlewareFunc {
	// Load API keys from yaml file
	data, err := os.ReadFile("api-keys.yaml")
	if err != nil {
		panic("Failed to read api-keys.yaml")
	}

	var apiKeys APIKeys
	if err := yaml.Unmarshal(data, &apiKeys); err != nil {
		panic("Failed to parse api-keys.yaml")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get API key from Authorization header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
					Error: struct {
						Message string      `json:"message"`
						Type    string      `json:"type"`
						Param   string      `json:"param,omitempty"`
						Code    interface{} `json:"code,omitempty"`
					}{
						Message: "Missing API key",
						Type:    "unauthorized",
					},
				})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
					Error: struct {
						Message string      `json:"message"`
						Type    string      `json:"type"`
						Param   string      `json:"param,omitempty"`
						Code    interface{} `json:"code,omitempty"`
					}{
						Message: "Invalid Authorization header format",
						Type:    "unauthorized",
					},
				})
			}

			apiKey := parts[1]

			// Check if API key is valid
			isValid := false
			for _, validKey := range apiKeys.Keys {
				if apiKey == validKey {
					isValid = true
					break
				}
			}

			if !isValid {
				return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
					Error: struct {
						Message string      `json:"message"`
						Type    string      `json:"type"`
						Param   string      `json:"param,omitempty"`
						Code    interface{} `json:"code,omitempty"`
					}{
						Message: "Invalid API key",
						Type:    "unauthorized",
					},
				})
			}

			// Valid API key, proceed to the next handler
			return next(c)
		}
	}
}
