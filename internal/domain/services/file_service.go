package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/storage"
)

// FileService quản lý các operations với file
type FileService interface {
	// UploadFile upload file và trả về thông tin file
	UploadFile(ctx context.Context, req *entities.FileUploadRequest) (*entities.FileUploadResponse, error)
	
	// DeleteFile xóa file
	DeleteFile(ctx context.Context, id string) error
	
	// GetFileURL lấy URL của file
	GetFileURL(objectKey string) string
	
	// GetFileUpload lấy thông tin file upload
	GetFileUpload(ctx context.Context, id string) (*entities.FileUpload, error)
	
	// GetFileUploads lấy danh sách file uploads
	GetFileUploads(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error)
	
	// ValidateFile kiểm tra file có hợp lệ không
	ValidateFile(header *multipart.FileHeader, config *entities.FileConfig) error
}

type fileService struct {
	storageProvider storage.StorageProvider
	fileRepo        repositories.FileRepository
}

// NewFileService tạo file service mới
func NewFileService(storageProvider storage.StorageProvider, fileRepo repositories.FileRepository) FileService {
	return &fileService{
		storageProvider: storageProvider,
		fileRepo:        fileRepo,
	}
}

func (fs *fileService) UploadFile(ctx context.Context, req *entities.FileUploadRequest) (*entities.FileUploadResponse, error) {
	file, ok := req.File.(multipart.File)
	if !ok {
		return nil, fmt.Errorf("invalid file type")
	}
	
	header, ok := req.Header.(*multipart.FileHeader)
	if !ok {
		return nil, fmt.Errorf("invalid file header type")
	}

	// Generate unique filename
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), fileExt)
	
	// Generate object key based on upload type and category
	var objectKey string
	switch req.UploadType {
	case entities.FileUploadTypeAdmin:
		objectKey = fmt.Sprintf("admin/%s/%s", req.Category, fileName)
	case entities.FileUploadTypeUser:
		objectKey = fmt.Sprintf("user/%s/%s", req.Category, fileName)
	case entities.FileUploadTypePublic:
		objectKey = fmt.Sprintf("public/%s/%s", req.Category, fileName)
	default:
		return nil, fmt.Errorf("invalid upload type: %s", req.UploadType)
	}

	// Upload file to storage
	fileURL, err := fs.storageProvider.UploadFile(file, objectKey, header.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create file upload record
	fileUpload := &entities.FileUpload{
		ID:           uuid.New().String(),
		FileName:     fileName,
		OriginalName: header.Filename,
		ObjectKey:    objectKey,
		FileSize:     header.Size,
		ContentType:  header.Header.Get("Content-Type"),
		URL:          fileURL,
		UploadedBy:   req.UploadedBy,
		UploadType:   req.UploadType,
		Category:     req.Category,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := fs.fileRepo.CreateFileUpload(ctx, fileUpload); err != nil {
		// Try to cleanup uploaded file if database save fails
		if deleteErr := fs.storageProvider.DeleteFile(objectKey); deleteErr != nil {
			// Log the cleanup error but don't override the original error
			fmt.Printf("Warning: failed to cleanup uploaded file after database error: %v\n", deleteErr)
		}
		return nil, fmt.Errorf("failed to save file upload record: %w", err)
	}

	return &entities.FileUploadResponse{
		ID:          fileUpload.ID,
		URL:         fileUpload.URL,
		FileName:    fileUpload.FileName,
		FileSize:    fileUpload.FileSize,
		ContentType: fileUpload.ContentType,
		Message:     "File uploaded successfully",
		CreatedAt:   fileUpload.CreatedAt,
	}, nil
}

func (fs *fileService) DeleteFile(ctx context.Context, id string) error {
	// Get file upload record
	fileUpload, err := fs.fileRepo.GetFileUploadByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get file upload record: %w", err)
	}

	// Delete from storage
	if err := fs.storageProvider.DeleteFile(fileUpload.ObjectKey); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	if err := fs.fileRepo.DeleteFileUpload(ctx, id); err != nil {
		return fmt.Errorf("failed to delete file upload record: %w", err)
	}

	return nil
}

func (fs *fileService) GetFileURL(objectKey string) string {
	return fs.storageProvider.GetFileURL(objectKey)
}

func (fs *fileService) GetFileUpload(ctx context.Context, id string) (*entities.FileUpload, error) {
	return fs.fileRepo.GetFileUploadByID(ctx, id)
}

func (fs *fileService) GetFileUploads(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error) {
	return fs.fileRepo.GetFileUploadsByTypeAndCategory(ctx, uploadType, category, limit, offset)
}

func (fs *fileService) ValidateFile(header *multipart.FileHeader, config *entities.FileConfig) error {
	// Check file size
	if header.Size > config.MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", header.Size, config.MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExt := false
	for _, allowedExt := range config.AllowedExtensions {
		if ext == allowedExt {
			validExt = true
			break
		}
	}
	if !validExt {
		return fmt.Errorf("file extension %s is not allowed. Allowed extensions: %v", ext, config.AllowedExtensions)
	}

	// Check content type if available
	contentType := header.Header.Get("Content-Type")
	if contentType != "" {
		validType := false
		for _, allowedType := range config.AllowedTypes {
			if contentType == allowedType {
				validType = true
				break
			}
		}
		if !validType {
			return fmt.Errorf("content type %s is not allowed. Allowed types: %v", contentType, config.AllowedTypes)
		}
	}

	return nil
}
