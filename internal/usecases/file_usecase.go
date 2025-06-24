package usecases

import (
	"context"
	"mime/multipart"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/services"
)

// FileUseCase defines the interface for file upload use cases
type FileUseCase interface {
	// UploadImage uploads an image file
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploadType entities.FileUploadType, userID *string) (*entities.FileUploadResponse, error)
	
	// UploadDocument uploads a document file
	UploadDocument(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploadType entities.FileUploadType, userID *string) (*entities.FileUploadResponse, error)
	
	// DeleteFile deletes a file
	DeleteFile(ctx context.Context, fileID string) error
	
	// GetFileUpload gets file upload info
	GetFileUpload(ctx context.Context, fileID string) (*entities.FileUpload, error)
	
	// GetFileUploads gets list of file uploads
	GetFileUploads(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error)
}

type fileUseCase struct {
	fileService services.FileService
}

// NewFileUseCase creates a new file use case
func NewFileUseCase(fileService services.FileService) FileUseCase {
	return &fileUseCase{
		fileService: fileService,
	}
}

func (uc *fileUseCase) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploadType entities.FileUploadType, userID *string) (*entities.FileUploadResponse, error) {
	// Validate image file
	config := entities.DefaultImageConfig()
	if err := uc.fileService.ValidateFile(header, config); err != nil {
		return nil, err
	}

	// Create upload request
	req := &entities.FileUploadRequest{
		File:       file,
		Header:     header,
		Category:   "images",
		UploadType: uploadType,
		UploadedBy: userID,
	}

	return uc.fileService.UploadFile(ctx, req)
}

func (uc *fileUseCase) UploadDocument(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploadType entities.FileUploadType, userID *string) (*entities.FileUploadResponse, error) {
	// Validate document file
	config := entities.DefaultDocumentConfig()
	if err := uc.fileService.ValidateFile(header, config); err != nil {
		return nil, err
	}

	// Create upload request
	req := &entities.FileUploadRequest{
		File:       file,
		Header:     header,
		Category:   "documents",
		UploadType: uploadType,
		UploadedBy: userID,
	}

	return uc.fileService.UploadFile(ctx, req)
}

func (uc *fileUseCase) DeleteFile(ctx context.Context, fileID string) error {
	return uc.fileService.DeleteFile(ctx, fileID)
}

func (uc *fileUseCase) GetFileUpload(ctx context.Context, fileID string) (*entities.FileUpload, error) {
	return uc.fileService.GetFileUpload(ctx, fileID)
}

func (uc *fileUseCase) GetFileUploads(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error) {
	return uc.fileService.GetFileUploads(ctx, uploadType, category, limit, offset)
}
