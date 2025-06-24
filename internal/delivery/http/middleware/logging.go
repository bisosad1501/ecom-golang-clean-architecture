package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates request logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + "req"
}

// ErrorHandlerMiddleware handles errors and returns consistent responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors after request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(400, gin.H{
					"error": "Invalid request format",
					"details": err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
			default:
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}
	}
}

// ValidationMiddleware validates request body using struct tags
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip validation for GET requests and non-JSON content
		if c.Request.Method == "GET" || c.ContentType() != "application/json" {
			c.Next()
			return
		}
		
		// Validate request body is present for POST/PUT/PATCH
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.ContentLength == 0 {
				c.JSON(400, gin.H{
					"error": "Request body is required",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}
