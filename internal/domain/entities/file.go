package entities

import (
	"time"
)

// FileUpload represents an uploaded file in the system
type FileUpload struct {
	ID          string    `json:"id" gorm:"primary_key"`
	FileName    string    `json:"fileName" gorm:"not null"`
	OriginalName string    `json:"originalName" gorm:"not null"`
	ObjectKey   string    `json:"objectKey" gorm:"not null;uniqueIndex"`
	FileSize    int64     `json:"fileSize" gorm:"not null"`
	ContentType string    `json:"contentType" gorm:"not null"`
	URL         string    `json:"url" gorm:"not null"`
	
	// Upload context
	UploadedBy   *string       `json:"uploadedBy,omitempty" gorm:"index"`  // User ID if authenticated
	UploadType   FileUploadType `json:"uploadType" gorm:"not null;index"`  // admin, user, public
	Category     string        `json:"category" gorm:"not null;index"`     // images, documents, etc.
	
	// Metadata
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// FileUploadType defines the type of upload
type FileUploadType string

const (
	FileUploadTypeAdmin  FileUploadType = "admin"
	FileUploadTypeUser   FileUploadType = "user"
	FileUploadTypePublic FileUploadType = "public"
)

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	File        interface{} `json:"-"`               // multipart.File
	Header      interface{} `json:"-"`               // *multipart.FileHeader
	Category    string      `json:"category"`        // images, documents, etc.
	UploadType  FileUploadType `json:"uploadType"`   // admin, user, public
	UploadedBy  *string     `json:"uploadedBy,omitempty"` // User ID if authenticated
}

// FileUploadResponse represents the response after successful upload
type FileUploadResponse struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"fileSize"`
	ContentType string    `json:"contentType"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"createdAt"`
}

// FileConfig represents configuration for file uploads
type FileConfig struct {
	MaxFileSize      int64    `json:"maxFileSize"`      // in bytes
	AllowedTypes     []string `json:"allowedTypes"`     // MIME types
	AllowedExtensions []string `json:"allowedExtensions"` // file extensions
}

// DefaultImageConfig returns default configuration for image uploads
func DefaultImageConfig() *FileConfig {
	return &FileConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{
			"image/jpeg",
			"image/jpg", 
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{
			".jpg",
			".jpeg", 
			".png",
			".gif",
			".webp",
		},
	}
}

// DefaultDocumentConfig returns default configuration for document uploads
func DefaultDocumentConfig() *FileConfig {
	return &FileConfig{
		MaxFileSize: 10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"text/plain",
		},
		AllowedExtensions: []string{
			".pdf",
			".doc",
			".docx",
			".txt",
		},
	}
}
