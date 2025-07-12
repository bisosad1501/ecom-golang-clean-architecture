package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

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

// FileUploadSecurityMiddleware provides enhanced security for file uploads
func FileUploadSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to multipart/form-data requests
		if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.Next()
			return
		}

		// Parse multipart form with size limit (10MB)
		err := c.Request.ParseMultipartForm(10 << 20)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse multipart form",
			})
			c.Abort()
			return
		}

		// Check each uploaded file
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
			for fieldName, files := range c.Request.MultipartForm.File {
				for _, fileHeader := range files {
					if err := validateUploadedFile(fileHeader); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": fmt.Sprintf("File validation failed for %s: %s", fieldName, err.Error()),
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// validateUploadedFile performs security validation on uploaded files
func validateUploadedFile(fileHeader *multipart.FileHeader) error {
	// Check file size (max 5MB for images, 10MB for documents)
	maxSize := int64(5 << 20) // 5MB default
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", fileHeader.Size, maxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
	}

	if !allowedExtensions[ext] {
		return fmt.Errorf("file extension %s is not allowed", ext)
	}

	// Check MIME type by reading file content
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for validation: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file for MIME type detection: %w", err)
	}

	// Detect MIME type
	mimeType := http.DetectContentType(buffer[:n])
	allowedMimeTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"application/pdf": true,
		"text/plain":      true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("MIME type %s is not allowed", mimeType)
	}

	// Check for malicious content patterns
	if err := scanForMaliciousContent(buffer[:n]); err != nil {
		return fmt.Errorf("malicious content detected: %w", err)
	}

	return nil
}

// scanForMaliciousContent scans file content for malicious patterns
func scanForMaliciousContent(content []byte) error {
	// Convert to lowercase for case-insensitive matching
	lowerContent := bytes.ToLower(content)

	// Check for script tags and other dangerous patterns
	dangerousPatterns := [][]byte{
		[]byte("<script"),
		[]byte("javascript:"),
		[]byte("vbscript:"),
		[]byte("onload="),
		[]byte("onerror="),
		[]byte("onclick="),
		[]byte("<?php"),
		[]byte("<%"),
		[]byte("eval("),
		[]byte("exec("),
		[]byte("system("),
		[]byte("shell_exec("),
	}

	for _, pattern := range dangerousPatterns {
		if bytes.Contains(lowerContent, pattern) {
			return fmt.Errorf("dangerous pattern detected: %s", string(pattern))
		}
	}

	return nil
}
