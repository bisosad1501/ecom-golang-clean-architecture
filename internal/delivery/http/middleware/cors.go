package middleware

import (
	"net/http"
	"strings"

	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates CORS middleware
func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowedOrigin := getAllowedOrigin(origin, cfg.AllowedOrigins)
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)

			// Only allow credentials for specific origins, not wildcard
			if allowedOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getAllowedOrigin checks if the origin is allowed and returns the appropriate value
func getAllowedOrigin(origin string, allowedOrigins []string) string {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return "*"
		}
		if allowed == origin {
			return origin
		}
	}
	return ""
}

// isOriginAllowed checks if the origin is in the allowed list (kept for backward compatibility)
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	return getAllowedOrigin(origin, allowedOrigins) != ""
}
