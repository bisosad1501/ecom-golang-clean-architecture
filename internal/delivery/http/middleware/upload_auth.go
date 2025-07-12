package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicUploadAuthMiddleware requires either JWT token or API key for public uploads
func PublicUploadAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for JWT token first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			// Use existing JWT middleware logic
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != "" {
				// Validate JWT token (simplified - in production use proper JWT validation)
				c.Set("authenticated", true)
				c.Next()
				return
			}
		}

		// Check for API key (for programmatic access)
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// Validate API key (simplified - in production use proper API key validation)
			if isValidAPIKey(apiKey) {
				c.Set("authenticated", true)
				c.Set("api_key_auth", true)
				c.Next()
				return
			}
		}

		// Check for session-based authentication (for guest uploads with session)
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID != "" && len(sessionID) > 10 {
			// Allow uploads with valid session ID
			c.Set("authenticated", true)
			c.Set("session_auth", true)
			c.Set("session_id", sessionID)
			c.Next()
			return
		}

		// No valid authentication found
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required for file uploads. Please provide a valid JWT token, API key, or session ID.",
		})
		c.Abort()
	}
}

// isValidAPIKey validates API key (simplified implementation)
func isValidAPIKey(apiKey string) bool {
	// In production, this would check against a database or cache
	// For now, accept any key that looks like a valid format
	return len(apiKey) >= 32 && len(apiKey) <= 64
}

// RequireAuthenticationMiddleware ensures user is authenticated for sensitive operations
func RequireAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated via any method
		if authenticated, exists := c.Get("authenticated"); !exists || !authenticated.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// FileUploadSecurityMiddleware adds additional security headers and checks
func FileUploadSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Check content length
		contentLength := c.Request.ContentLength
		if contentLength > 10*1024*1024 { // 10MB limit
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File size exceeds maximum allowed size of 10MB",
			})
			c.Abort()
			return
		}
		
		// Check content type
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid content type. Expected multipart/form-data",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
