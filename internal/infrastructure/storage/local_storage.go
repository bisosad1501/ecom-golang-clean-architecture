package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"ecom-golang-clean-architecture/internal/domain/storage"
)

type LocalFileStorage struct {
	config *config.LocalStorageConfig
}

func NewLocalStorage(cfg *config.LocalStorageConfig) (storage.StorageProvider, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(cfg.BaseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	
	return &LocalFileStorage{
		config: cfg,
	}, nil
}

// Ensure LocalFileStorage implements StorageProvider
var _ storage.StorageProvider = (*LocalFileStorage)(nil)

func (s *LocalFileStorage) UploadFile(file multipart.File, objectKey string, contentType string) (string, error) {
	// Create full file path
	fullPath := filepath.Join(s.config.BaseDir, objectKey)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()
	
	// Reset file pointer to beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}
	
	// Generate public URL
	url := s.GetFileURL(objectKey)
	return url, nil
}

func (s *LocalFileStorage) DeleteFile(objectKey string) error {
	fullPath := filepath.Join(s.config.BaseDir, objectKey)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}
	
	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (s *LocalFileStorage) GetFileURL(objectKey string) string {
	// Clean the object key to ensure proper URL format
	cleanKey := strings.TrimPrefix(objectKey, "/")
	return fmt.Sprintf("%s/%s", s.config.PublicPath, cleanKey)
}

func (s *LocalFileStorage) FileExists(objectKey string) (bool, error) {
	fullPath := filepath.Join(s.config.BaseDir, objectKey)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
