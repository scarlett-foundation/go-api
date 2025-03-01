package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// httpRequestsTotal counts total HTTP requests
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by status code, method, and path",
		},
		[]string{"status", "method", "path"},
	)

	// httpRequestDuration measures request duration
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// apiKeyRequests counts requests by API key
	apiKeyRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_key_requests_total",
			Help: "Total number of requests by API key (masked)",
		},
		[]string{"api_key"},
	)

	// tokenUsagePrompt tracks the number of prompt tokens used by API key
	tokenUsagePrompt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_usage_prompt_total",
			Help: "Total number of prompt tokens used by API key",
		},
		[]string{"api_key"},
	)

	// tokenUsageCompletion tracks the number of completion tokens used by API key
	tokenUsageCompletion = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_usage_completion_total",
			Help: "Total number of completion tokens used by API key",
		},
		[]string{"api_key"},
	)

	// tokenUsageTotal tracks the total number of tokens used by API key
	tokenUsageTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_usage_total",
			Help: "Total number of tokens used by API key",
		},
		[]string{"api_key"},
	)
)

// responseBodyWriter is a custom response writer that captures the response body
type responseBodyWriter struct {
	io.Writer
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (r *responseBodyWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

// Write captures the response body
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	return r.Writer.Write(b)
}

// responseBody represents the structure of a chat completion response
type responseBody struct {
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(apiKeyRequests)
	prometheus.MustRegister(tokenUsagePrompt)
	prometheus.MustRegister(tokenUsageCompletion)
	prometheus.MustRegister(tokenUsageTotal)
}

// PrometheusMiddleware returns a middleware function that collects Prometheus metrics
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Capture API key from the request
			// The API key is in the Authorization header in the format "Bearer <token>"
			apiKey := "unknown"
			auth := c.Request().Header.Get("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				key := auth[7:]
				if len(key) > 8 {
					// Mask the API key for privacy - keep first and last 4 chars
					apiKey = key[:4] + "..." + key[len(key)-4:]
				} else {
					apiKey = "short_key"
				}
			}

			// Use a custom response writer to capture the response body
			// Only do this for chat completions endpoint to avoid overhead on other endpoints
			if c.Path() == "/chat/completions" && c.Request().Method == "POST" {
				resBody := new(bytes.Buffer)
				writer := &responseBodyWriter{
					Writer:         io.MultiWriter(resBody, c.Response().Writer),
					ResponseWriter: c.Response().Writer,
					statusCode:     http.StatusOK,
				}

				// Replace the response writer
				c.Response().Writer = writer

				// Process the request
				err := next(c)

				// Only try to parse the response body if the status code is 200 OK
				if writer.statusCode == http.StatusOK {
					var response responseBody
					if jsonErr := json.Unmarshal(resBody.Bytes(), &response); jsonErr == nil {
						// Record token usage metrics
						if response.Usage.PromptTokens > 0 {
							tokenUsagePrompt.WithLabelValues(apiKey).Add(float64(response.Usage.PromptTokens))
						}
						if response.Usage.CompletionTokens > 0 {
							tokenUsageCompletion.WithLabelValues(apiKey).Add(float64(response.Usage.CompletionTokens))
						}
						if response.Usage.TotalTokens > 0 {
							tokenUsageTotal.WithLabelValues(apiKey).Add(float64(response.Usage.TotalTokens))
						}
					}
				}

				// Record metrics after the request is processed
				duration := time.Since(start).Seconds()
				status := writer.statusCode
				method := c.Request().Method
				path := c.Request().URL.Path

				// Record HTTP metrics
				httpRequestsTotal.WithLabelValues(strconv.Itoa(status), method, path).Inc()
				httpRequestDuration.WithLabelValues(method, path).Observe(duration)

				// Record API key usage
				if apiKey != "unknown" {
					apiKeyRequests.WithLabelValues(apiKey).Inc()
				}

				return err
			}

			// Standard processing for all other requests
			err := next(c)

			// Record metrics after the request is processed
			duration := time.Since(start).Seconds()
			status := c.Response().Status
			method := c.Request().Method
			path := c.Request().URL.Path

			// Record HTTP metrics
			httpRequestsTotal.WithLabelValues(strconv.Itoa(status), method, path).Inc()
			httpRequestDuration.WithLabelValues(method, path).Observe(duration)

			// Record API key usage
			if apiKey != "unknown" {
				apiKeyRequests.WithLabelValues(apiKey).Inc()
			}

			return err
		}
	}
}

// RegisterPrometheusHandler registers the Prometheus metrics endpoint
func RegisterPrometheusHandler(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
