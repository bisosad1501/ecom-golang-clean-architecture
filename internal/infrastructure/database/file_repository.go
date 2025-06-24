package database

import (
	"context"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *gorm.DB) repositories.FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) CreateFileUpload(ctx context.Context, fileUpload *entities.FileUpload) error {
	return r.db.WithContext(ctx).Create(fileUpload).Error
}

func (r *fileRepository) GetFileUploadByID(ctx context.Context, id string) (*entities.FileUpload, error) {
	var fileUpload entities.FileUpload
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&fileUpload).Error
	if err != nil {
		return nil, err
	}
	return &fileUpload, nil
}

func (r *fileRepository) GetFileUploadByObjectKey(ctx context.Context, objectKey string) (*entities.FileUpload, error) {
	var fileUpload entities.FileUpload
	err := r.db.WithContext(ctx).Where("object_key = ?", objectKey).First(&fileUpload).Error
	if err != nil {
		return nil, err
	}
	return &fileUpload, nil
}

func (r *fileRepository) GetFileUploadsByUser(ctx context.Context, userID string, limit, offset int) ([]*entities.FileUpload, error) {
	var fileUploads []*entities.FileUpload
	err := r.db.WithContext(ctx).
		Where("uploaded_by = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&fileUploads).Error
	return fileUploads, err
}

func (r *fileRepository) GetFileUploadsByTypeAndCategory(ctx context.Context, uploadType entities.FileUploadType, category string, limit, offset int) ([]*entities.FileUpload, error) {
	var fileUploads []*entities.FileUpload
	query := r.db.WithContext(ctx).Where("upload_type = ?", uploadType)
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&fileUploads).Error
	return fileUploads, err
}

func (r *fileRepository) DeleteFileUpload(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.FileUpload{}).Error
}

func (r *fileRepository) UpdateFileUpload(ctx context.Context, fileUpload *entities.FileUpload) error {
	return r.db.WithContext(ctx).Save(fileUpload).Error
}

func (r *fileRepository) FileExists(ctx context.Context, objectKey string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.FileUpload{}).Where("object_key = ?", objectKey).Count(&count).Error
	return count > 0, err
}

func (r *fileRepository) GetFileCountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.FileUpload{}).Where("uploaded_by = ?", userID).Count(&count).Error
	return count, err
}

func (r *fileRepository) GetFileCountByType(ctx context.Context, uploadType entities.FileUploadType) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.FileUpload{}).Where("upload_type = ?", uploadType).Count(&count).Error
	return count, err
}
