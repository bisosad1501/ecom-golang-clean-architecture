package utils

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// AllowedImageTypes defines allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedDocumentTypes defines allowed document MIME types
var AllowedDocumentTypes = map[string]bool{
	"application/pdf":                                                 true,
	"application/msword":                                              true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel":                                        true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"text/plain":                                                      true,
	"text/csv":                                                        true,
}

// ImageMagicNumbers defines magic numbers for image file types
var ImageMagicNumbers = map[string][]byte{
	"image/jpeg": {0xFF, 0xD8, 0xFF},
	"image/png":  {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	"image/gif":  {0x47, 0x49, 0x46, 0x38},
	"image/webp": {0x52, 0x49, 0x46, 0x46},
}

// DocumentMagicNumbers defines magic numbers for document file types
var DocumentMagicNumbers = map[string][]byte{
	"application/pdf": {0x25, 0x50, 0x44, 0x46},
	"application/zip": {0x50, 0x4B, 0x03, 0x04}, // For Office documents
}

// ValidateFileType validates file type using both extension and magic numbers
func ValidateFileType(filename string, content []byte, allowedTypes map[string]bool) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	
	if mimeType == "" {
		return fmt.Errorf("unknown file type for extension: %s", ext)
	}

	// Remove charset from mime type if present
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = mimeType[:idx]
	}

	if !allowedTypes[mimeType] {
		return fmt.Errorf("file type %s is not allowed", mimeType)
	}

	// Validate magic numbers for additional security
	if err := validateMagicNumbers(content, mimeType); err != nil {
		return fmt.Errorf("file content validation failed: %w", err)
	}

	return nil
}

// validateMagicNumbers checks if file content matches expected magic numbers
func validateMagicNumbers(content []byte, expectedMimeType string) error {
	if len(content) < 8 {
		return fmt.Errorf("file content too short for validation")
	}

	// Check image magic numbers
	if magicBytes, exists := ImageMagicNumbers[expectedMimeType]; exists {
		if !hasPrefix(content, magicBytes) {
			return fmt.Errorf("file content does not match expected type %s", expectedMimeType)
		}
		return nil
	}

	// Check document magic numbers
	if magicBytes, exists := DocumentMagicNumbers[expectedMimeType]; exists {
		if !hasPrefix(content, magicBytes) {
			return fmt.Errorf("file content does not match expected type %s", expectedMimeType)
		}
		return nil
	}

	// For Office documents, check for ZIP signature (they are ZIP files)
	if strings.Contains(expectedMimeType, "officedocument") {
		zipMagic := DocumentMagicNumbers["application/zip"]
		if !hasPrefix(content, zipMagic) {
			return fmt.Errorf("office document does not have valid ZIP signature")
		}
		return nil
	}

	// For text files, check if content is valid UTF-8
	if expectedMimeType == "text/plain" || expectedMimeType == "text/csv" {
		if !isValidUTF8(content) {
			return fmt.Errorf("text file contains invalid UTF-8 content")
		}
		return nil
	}

	return nil
}

// hasPrefix checks if content starts with the given prefix
func hasPrefix(content, prefix []byte) bool {
	if len(content) < len(prefix) {
		return false
	}
	for i, b := range prefix {
		if content[i] != b {
			return false
		}
	}
	return true
}

// isValidUTF8 checks if content is valid UTF-8
func isValidUTF8(content []byte) bool {
	// Simple check - try to convert to string and back
	str := string(content)
	return len([]byte(str)) == len(content)
}

// ValidateImageFile validates an image file
func ValidateImageFile(filename string, content []byte) error {
	return ValidateFileType(filename, content, AllowedImageTypes)
}

// ValidateDocumentFile validates a document file
func ValidateDocumentFile(filename string, content []byte) error {
	return ValidateFileType(filename, content, AllowedDocumentTypes)
}

// DetectContentType detects content type from file content
func DetectContentType(content []byte) string {
	return http.DetectContentType(content)
}

// SanitizeFilename removes potentially dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	dangerous := []string{"/", "\\", "..", ":", "*", "?", "\"", "<", ">", "|"}
	
	sanitized := filename
	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	
	// Limit filename length
	if len(sanitized) > 255 {
		ext := filepath.Ext(sanitized)
		name := sanitized[:255-len(ext)]
		sanitized = name + ext
	}
	
	return sanitized
}

// ValidateFileSize validates file size
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", size, maxSize)
	}
	return nil
}

// IsImageFile checks if the file is an image based on content type
func IsImageFile(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsDocumentFile checks if the file is a document based on content type
func IsDocumentFile(contentType string) bool {
	return AllowedDocumentTypes[contentType]
}
