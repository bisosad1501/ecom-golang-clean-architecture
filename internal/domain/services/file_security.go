package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// FileSecurityService provides file security validation
type FileSecurityService interface {
	ValidateFileContent(file multipart.File, header *multipart.FileHeader) error
	ScanForMalware(file multipart.File) error
	ValidateImageContent(file multipart.File) error
	ValidateDocumentContent(file multipart.File) error
}

type fileSecurityService struct{}

// NewFileSecurityService creates a new file security service
func NewFileSecurityService() FileSecurityService {
	return &fileSecurityService{}
}

// ValidateFileContent performs comprehensive file content validation
func (s *fileSecurityService) ValidateFileContent(file multipart.File, header *multipart.FileHeader) error {
	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content: %w", err)
	}
	
	// Reset file pointer again
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Detect actual content type
	actualContentType := http.DetectContentType(buffer[:n])
	declaredContentType := header.Header.Get("Content-Type")
	
	// Validate content type matches extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if err := s.validateContentTypeExtensionMatch(actualContentType, ext); err != nil {
		return err
	}
	
	// Validate declared vs actual content type
	if declaredContentType != "" && !s.isContentTypeCompatible(declaredContentType, actualContentType) {
		return fmt.Errorf("declared content type %s does not match actual content type %s", declaredContentType, actualContentType)
	}
	
	// Check for suspicious content
	if err := s.checkSuspiciousContent(buffer[:n]); err != nil {
		return err
	}
	
	// Perform specific validation based on file type
	if strings.HasPrefix(actualContentType, "image/") {
		return s.ValidateImageContent(file)
	} else if s.isDocumentType(actualContentType) {
		return s.ValidateDocumentContent(file)
	}
	
	return nil
}

// validateContentTypeExtensionMatch validates that content type matches file extension
func (s *fileSecurityService) validateContentTypeExtensionMatch(contentType, ext string) error {
	validMappings := map[string][]string{
		"image/jpeg":    {".jpg", ".jpeg"},
		"image/png":     {".png"},
		"image/gif":     {".gif"},
		"image/webp":    {".webp"},
		"application/pdf": {".pdf"},
		"text/plain":    {".txt"},
		"application/msword": {".doc"},
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {".docx"},
	}
	
	for ct, extensions := range validMappings {
		if strings.HasPrefix(contentType, ct) {
			for _, validExt := range extensions {
				if ext == validExt {
					return nil
				}
			}
			return fmt.Errorf("file extension %s does not match content type %s", ext, contentType)
		}
	}
	
	return fmt.Errorf("unsupported content type: %s", contentType)
}

// isContentTypeCompatible checks if declared and actual content types are compatible
func (s *fileSecurityService) isContentTypeCompatible(declared, actual string) bool {
	// Handle common variations
	compatibleTypes := map[string][]string{
		"image/jpeg": {"image/jpeg", "image/jpg"},
		"image/jpg":  {"image/jpeg", "image/jpg"},
		"image/png":  {"image/png"},
		"image/gif":  {"image/gif"},
		"image/webp": {"image/webp"},
		"application/pdf": {"application/pdf"},
		"text/plain": {"text/plain"},
	}
	
	if compatible, exists := compatibleTypes[declared]; exists {
		for _, ct := range compatible {
			if strings.HasPrefix(actual, ct) {
				return true
			}
		}
	}
	
	return strings.HasPrefix(actual, declared)
}

// checkSuspiciousContent checks for suspicious patterns in file content
func (s *fileSecurityService) checkSuspiciousContent(content []byte) error {
	// Check for executable signatures
	suspiciousSignatures := [][]byte{
		{0x4D, 0x5A},                   // PE executable (MZ)
		{0x7F, 0x45, 0x4C, 0x46},       // ELF executable
		{0xCA, 0xFE, 0xBA, 0xBE},       // Mach-O executable
		{0x50, 0x4B, 0x03, 0x04},       // ZIP file (could contain executables)
	}
	
	for _, signature := range suspiciousSignatures {
		if bytes.HasPrefix(content, signature) {
			return fmt.Errorf("suspicious file signature detected")
		}
	}
	
	// Check for script tags in images (potential XSS)
	suspiciousStrings := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
	}
	
	contentStr := strings.ToLower(string(content))
	for _, suspicious := range suspiciousStrings {
		if strings.Contains(contentStr, suspicious) {
			return fmt.Errorf("suspicious content detected: potential script injection")
		}
	}
	
	return nil
}

// ValidateImageContent validates image file content
func (s *fileSecurityService) ValidateImageContent(file multipart.File) error {
	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Read first 1KB to check image headers
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read image content: %w", err)
	}
	
	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	content := buffer[:n]
	
	// Check for valid image headers
	validHeaders := [][]byte{
		{0xFF, 0xD8, 0xFF},                   // JPEG
		{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG
		{0x47, 0x49, 0x46, 0x38},             // GIF
		{0x52, 0x49, 0x46, 0x46},             // WEBP (RIFF)
	}
	
	hasValidHeader := false
	for _, header := range validHeaders {
		if bytes.HasPrefix(content, header) {
			hasValidHeader = true
			break
		}
	}
	
	if !hasValidHeader {
		return fmt.Errorf("invalid image file header")
	}
	
	return nil
}

// ValidateDocumentContent validates document file content
func (s *fileSecurityService) ValidateDocumentContent(file multipart.File) error {
	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Read first 1KB to check document headers
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read document content: %w", err)
	}
	
	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	content := buffer[:n]
	
	// Check for valid document headers
	validHeaders := [][]byte{
		{0x25, 0x50, 0x44, 0x46},             // PDF
		{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}, // MS Office (DOC)
		{0x50, 0x4B, 0x03, 0x04},             // ZIP-based (DOCX)
	}
	
	hasValidHeader := false
	for _, header := range validHeaders {
		if bytes.HasPrefix(content, header) {
			hasValidHeader = true
			break
		}
	}
	
	// For text files, just check if it's valid UTF-8
	if !hasValidHeader {
		// Check if it's a text file
		if s.isValidUTF8(content) {
			return nil
		}
		return fmt.Errorf("invalid document file header")
	}
	
	return nil
}

// ScanForMalware performs basic malware scanning
func (s *fileSecurityService) ScanForMalware(file multipart.File) error {
	// This is a basic implementation
	// In production, you would integrate with a proper antivirus service
	return s.ValidateFileContent(file, &multipart.FileHeader{})
}

// isDocumentType checks if content type is a document type
func (s *fileSecurityService) isDocumentType(contentType string) bool {
	documentTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain",
	}
	
	for _, docType := range documentTypes {
		if strings.HasPrefix(contentType, docType) {
			return true
		}
	}
	
	return false
}

// isValidUTF8 checks if content is valid UTF-8
func (s *fileSecurityService) isValidUTF8(content []byte) bool {
	// Simple check for valid UTF-8
	for i := 0; i < len(content); {
		if content[i] < 0x80 {
			i++
			continue
		}
		
		// Multi-byte character
		if content[i] < 0xC0 {
			return false
		}
		
		var size int
		if content[i] < 0xE0 {
			size = 2
		} else if content[i] < 0xF0 {
			size = 3
		} else if content[i] < 0xF8 {
			size = 4
		} else {
			return false
		}
		
		if i+size > len(content) {
			return false
		}
		
		for j := 1; j < size; j++ {
			if content[i+j] < 0x80 || content[i+j] >= 0xC0 {
				return false
			}
		}
		
		i += size
	}
	
	return true
}
