package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "request_id"

// Logger returns a gin middleware for logging HTTP requests with request ID
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Generate or retrieve request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Create logger with request ID
		requestLogger := logger.With(zap.String("request_id", requestID))

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		status := c.Writer.Status()

		// Use appropriate log level based on status code
		logFields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.Int("size", c.Writer.Size()),
		}

		switch {
		case status >= 500:
			requestLogger.Error("HTTP Request", logFields...)
		case status >= 400:
			requestLogger.Warn("HTTP Request", logFields...)
		default:
			requestLogger.Info("HTTP Request", logFields...)
		}
	}
}

// GetRequestLogger returns a logger with request ID from context
func GetRequestLogger(c *gin.Context, baseLogger *zap.Logger) *zap.Logger {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return baseLogger.With(zap.String("request_id", id))
		}
	}
	return baseLogger
}
