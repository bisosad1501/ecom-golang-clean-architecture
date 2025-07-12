package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// SessionValidationMiddleware validates session ID format for guest operations
func SessionValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate session ID for guest operations (when no auth token is present)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// User is authenticated, skip session validation
			c.Next()
			return
		}

		// Check if session ID is provided
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			// No session ID provided, let the handler decide if it's required
			c.Next()
			return
		}

		// Validate session ID format
		if err := validateSessionID(sessionID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid session ID format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// validateSessionID validates the session ID format
func validateSessionID(sessionID string) error {
	// Session ID should be 10-128 characters, alphanumeric with hyphens and underscores
	if len(sessionID) < 10 || len(sessionID) > 128 {
		return fmt.Errorf("session ID must be between 10 and 128 characters")
	}

	// Check format: alphanumeric, hyphens, and underscores only
	validFormat := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validFormat.MatchString(sessionID) {
		return fmt.Errorf("session ID can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

// GuestCartMiddleware ensures guest cart operations have valid session ID
func GuestCartMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to guest operations (no auth token)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			c.Next()
			return
		}

		// For guest cart operations, session ID is required
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required for guest cart operations",
			})
			c.Abort()
			return
		}

		// Validate session ID format
		if err := validateSessionID(sessionID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid session ID format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
