package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://js.stripe.com; "+
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
			"font-src 'self' https://fonts.gstatic.com; "+
			"img-src 'self' data: https:; "+
			"connect-src 'self' https://api.stripe.com; "+
			"frame-src https://js.stripe.com https://hooks.stripe.com; "+
			"object-src 'none'; "+
			"base-uri 'self'")

		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", 
			"geolocation=(), "+
			"microphone=(), "+
			"camera=(), "+
			"payment=(self), "+
			"usb=(), "+
			"magnetometer=(), "+
			"gyroscope=(), "+
			"accelerometer=()")

		// HSTS (only for HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// RequestSizeLimitMiddleware limits request body size
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{
				"error": "Request entity too large",
			})
			c.Abort()
			return
		}
		c.Next()
	})
}

// NoSniffMiddleware prevents MIME type sniffing
func NoSniffMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}

// AntiClickjackingMiddleware prevents clickjacking attacks
func AntiClickjackingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Next()
	}
}
