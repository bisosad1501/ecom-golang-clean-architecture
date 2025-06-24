package usecases

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// FileUseCase interface
type FileUseCase interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*FileUploadResponse, error)
	UploadMultipleImages(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*FileUploadResponse, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFileURL(ctx context.Context, filePath string) string
}

type fileUseCase struct {
	uploadPath string
	baseURL    string
}

// NewFileUseCase creates a new file use case
func NewFileUseCase(uploadPath, baseURL string) FileUseCase {
	return &fileUseCase{
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileURL      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	ContentType  string `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// UploadImage uploads a single image file
func (uc *fileUseCase) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (*FileUploadResponse, error) {
	// Validate file type
	if !uc.isValidImageType(header.Header.Get("Content-Type")) {
		return nil, entities.ErrInvalidFileType
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		return nil, entities.ErrFileTooLarge
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Create folder path
	folderPath := filepath.Join(uc.uploadPath, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(folderPath, fileName)
	
	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate relative path for storage
	relativePath := filepath.Join(folder, fileName)
	
	return &FileUploadResponse{
		FileName:    fileName,
		FilePath:    relativePath,
		FileURL:     uc.GetFileURL(ctx, relativePath),
		FileSize:    header.Size,
		ContentType: header.Header.Get("Content-Type"),
		UploadedAt:  time.Now(),
	}, nil
}

// UploadMultipleImages uploads multiple image files
func (uc *fileUseCase) UploadMultipleImages(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*FileUploadResponse, error) {
	responses := make([]*FileUploadResponse, 0, len(files))
	
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue // Skip invalid files
		}
		defer file.Close()

		response, err := uc.UploadImage(ctx, file, fileHeader, folder)
		if err != nil {
			continue // Skip files that fail to upload
		}
		
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return nil, entities.ErrNoValidFiles
	}

	return responses, nil
}

// DeleteFile deletes a file from storage
func (uc *fileUseCase) DeleteFile(ctx context.Context, filePath string) error {
	fullPath := filepath.Join(uc.uploadPath, filePath)
	
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return entities.ErrFileNotFound
	}

	return os.Remove(fullPath)
}

// GetFileURL returns the full URL for a file
func (uc *fileUseCase) GetFileURL(ctx context.Context, filePath string) string {
	return fmt.Sprintf("%s/uploads/%s", uc.baseURL, strings.ReplaceAll(filePath, "\\", "/"))
}

// isValidImageType checks if the content type is a valid image type
func (uc *fileUseCase) isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg", 
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// Image processing utilities
type ImageProcessor interface {
	ResizeImage(ctx context.Context, filePath string, width, height int) error
	CreateThumbnail(ctx context.Context, filePath string, size int) (string, error)
	OptimizeImage(ctx context.Context, filePath string, quality int) error
}

type imageProcessor struct {
	uploadPath string
}

// NewImageProcessor creates a new image processor
func NewImageProcessor(uploadPath string) ImageProcessor {
	return &imageProcessor{
		uploadPath: uploadPath,
	}
}

// ResizeImage resizes an image to specified dimensions
func (ip *imageProcessor) ResizeImage(ctx context.Context, filePath string, width, height int) error {
	// Implementation would use image processing library like imaging or resize
	// For now, return not implemented
	return entities.ErrNotImplemented
}

// CreateThumbnail creates a thumbnail version of an image
func (ip *imageProcessor) CreateThumbnail(ctx context.Context, filePath string, size int) (string, error) {
	// Implementation would create thumbnail
	// For now, return not implemented
	return "", entities.ErrNotImplemented
}

// OptimizeImage optimizes image quality and size
func (ip *imageProcessor) OptimizeImage(ctx context.Context, filePath string, quality int) error {
	// Implementation would optimize image
	// For now, return not implemented
	return entities.ErrNotImplemented
}

// File validation utilities
type FileValidator struct {
	MaxSize      int64
	AllowedTypes []string
}

// NewFileValidator creates a new file validator
func NewFileValidator(maxSize int64, allowedTypes []string) *FileValidator {
	return &FileValidator{
		MaxSize:      maxSize,
		AllowedTypes: allowedTypes,
	}
}

// ValidateFile validates file against rules
func (fv *FileValidator) ValidateFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > fv.MaxSize {
		return entities.ErrFileTooLarge
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	for _, allowedType := range fv.AllowedTypes {
		if contentType == allowedType {
			return nil
		}
	}

	return entities.ErrInvalidFileType
}

// Storage interface for different storage backends
type Storage interface {
	Upload(ctx context.Context, file multipart.File, path string) error
	Delete(ctx context.Context, path string) error
	GetURL(ctx context.Context, path string) string
}

// LocalStorage implements Storage for local file system
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage
func NewLocalStorage(basePath, baseURL string) Storage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Upload uploads file to local storage
func (ls *LocalStorage) Upload(ctx context.Context, file multipart.File, path string) error {
	fullPath := filepath.Join(ls.basePath, path)
	
	// Create directory if not exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	return err
}

// Delete deletes file from local storage
func (ls *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(ls.basePath, path)
	return os.Remove(fullPath)
}

// GetURL returns URL for local storage
func (ls *LocalStorage) GetURL(ctx context.Context, path string) string {
	return fmt.Sprintf("%s/uploads/%s", ls.baseURL, strings.ReplaceAll(path, "\\", "/"))
}

// CloudStorage interface for cloud storage providers
type CloudStorage interface {
	Storage
	GeneratePresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)
}

// S3Storage implements CloudStorage for AWS S3
type S3Storage struct {
	bucket string
	region string
}

// NewS3Storage creates a new S3 storage
func NewS3Storage(bucket, region string) CloudStorage {
	return &S3Storage{
		bucket: bucket,
		region: region,
	}
}

// Upload uploads file to S3
func (s3 *S3Storage) Upload(ctx context.Context, file multipart.File, path string) error {
	// Implementation would use AWS SDK
	return entities.ErrNotImplemented
}

// Delete deletes file from S3
func (s3 *S3Storage) Delete(ctx context.Context, path string) error {
	// Implementation would use AWS SDK
	return entities.ErrNotImplemented
}

// GetURL returns S3 URL
func (s3 *S3Storage) GetURL(ctx context.Context, path string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3.bucket, s3.region, path)
}

// GeneratePresignedURL generates presigned URL for S3
func (s3 *S3Storage) GeneratePresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	// Implementation would use AWS SDK
	return "", entities.ErrNotImplemented
}
