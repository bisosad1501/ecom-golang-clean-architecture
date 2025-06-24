package storage

import "mime/multipart"

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	// UploadFile uploads a file to storage and returns the file URL
	UploadFile(file multipart.File, objectKey string, contentType string) (string, error)
	
	// DeleteFile deletes a file from storage
	DeleteFile(objectKey string) error
	
	// GetFileURL returns the URL for accessing a file
	GetFileURL(objectKey string) string
	
	// FileExists checks if a file exists in storage
	FileExists(objectKey string) (bool, error)
}
