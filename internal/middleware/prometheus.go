package middleware

import (
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
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(apiKeyRequests)
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

			// Process the request
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
