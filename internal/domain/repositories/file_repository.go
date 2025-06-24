package repositories

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
)

// FileRepository defines the interface for file upload data operations
type FileRepository interface {
	// Create a new file upload record
	CreateFileUpload(ctx context.Context, fileUpload *entities.FileUpload) error
	
	// Get file upload by ID
	GetFileUploadByID(ctx context.Context, id string) (*entities.FileUpload, error)
	
	// Get file upload by object key
	GetFileUploadByObjectKey(ctx context.Context, objectKey string) (*entities.FileUpload, error)
	
	// Get file uploads by user ID
	GetFileUploadsByUser(ctx context.Context, userID string, limit, offset int) ([]*entities.FileUpload, error)
	
	// Get file uploads by type and category
	GetFileUploadsByTypeAndCategory(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error)
	
	// Delete file upload record
	DeleteFileUpload(ctx context.Context, id string) error
	
	// Update file upload record
	UpdateFileUpload(ctx context.Context, fileUpload *entities.FileUpload) error
	
	// Check if file exists
	FileExists(ctx context.Context, objectKey string) (bool, error)
	
	// Get total count of files by user
	GetFileCountByUser(ctx context.Context, userID string) (int64, error)
	
	// Get total count of files by type
	GetFileCountByType(ctx context.Context, uploadType entities.FileUploadType) (int64, error)
}
