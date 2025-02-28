package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-api/internal/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimiterConfig defines the config for RateLimiter middleware.
type RateLimiterConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// RequestsPerSecond is the rate limit per API key
	RequestsPerSecond rate.Limit

	// Burst is the maximum burst size allowed
	Burst int

	// ExpirationTime is how long to keep rate limiters in memory
	ExpirationTime time.Duration

	// LimitersMutex is a mutex for the Limiters map
	LimitersMutex sync.RWMutex

	// Limiters is a map of API keys to rate limiters
	Limiters map[string]*rateLimiterEntry

	// CleanupTicker for removing expired rate limiters
	CleanupTicker *time.Ticker
}

// rateLimiterEntry represents a rate limiter with its last access time
type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// DefaultRateLimiter returns a middleware that limits requests based on API key
// Each API key is allowed 20 requests per second with a burst of 30
func DefaultRateLimiter() echo.MiddlewareFunc {
	// Initialize config with default values
	config := &RateLimiterConfig{
		Skipper:           middleware.DefaultSkipper,
		RequestsPerSecond: rate.Limit(20),
		Burst:             30,
		ExpirationTime:    30 * time.Minute,
		Limiters:          make(map[string]*rateLimiterEntry),
	}

	// Start cleanup goroutine
	config.CleanupTicker = time.NewTicker(10 * time.Minute)
	go func() {
		for range config.CleanupTicker.C {
			config.cleanup()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Extract API key from Authorization header
			apiKey := extractAPIKey(c)
			if apiKey == "" {
				// If there's no API key, the auth middleware will handle it
				// Just pass through for now
				return next(c)
			}

			// Get or create rate limiter for this API key
			limiter := config.getLimiter(apiKey)

			// Check if request allowed
			if !limiter.Allow() {
				// Return custom rate limit error response
				return c.JSON(http.StatusTooManyRequests, types.ErrorResponse{
					Error: struct {
						Message string      `json:"message"`
						Type    string      `json:"type"`
						Param   string      `json:"param,omitempty"`
						Code    interface{} `json:"code,omitempty"`
					}{
						Message: "Rate limit exceeded for your API key. Please try again later.",
						Type:    "rate_limit_error",
						Code:    http.StatusTooManyRequests,
					},
				})
			}

			// Add rate limit headers
			limit := strconv.FormatFloat(float64(config.RequestsPerSecond), 'f', 2, 64)
			remaining := strconv.Itoa(config.Burst - int(limiter.Tokens()))
			reset := strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10)

			c.Response().Header().Set("X-RateLimit-Limit", limit)
			c.Response().Header().Set("X-RateLimit-Remaining", remaining)
			c.Response().Header().Set("X-RateLimit-Reset", reset)

			return next(c)
		}
	}
}

// getLimiter gets or creates a rate limiter for the given API key
func (config *RateLimiterConfig) getLimiter(apiKey string) *rate.Limiter {
	config.LimitersMutex.RLock()
	entry, exists := config.Limiters[apiKey]
	config.LimitersMutex.RUnlock()

	if exists {
		// Update last seen time
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	// Create new limiter
	limiter := rate.NewLimiter(config.RequestsPerSecond, config.Burst)

	// Store in map
	config.LimitersMutex.Lock()
	config.Limiters[apiKey] = &rateLimiterEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	config.LimitersMutex.Unlock()

	return limiter
}

// cleanup removes expired rate limiters
func (config *RateLimiterConfig) cleanup() {
	cutoff := time.Now().Add(-config.ExpirationTime)

	config.LimitersMutex.Lock()
	defer config.LimitersMutex.Unlock()

	for apiKey, entry := range config.Limiters {
		if entry.lastSeen.Before(cutoff) {
			delete(config.Limiters, apiKey)
		}
	}
}

// extractAPIKey extracts the API key from the Authorization header
func extractAPIKey(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
