package config

import (
	"os"
	"strconv"
)

type FileStorageConfig struct {
	Provider    string `json:"provider" env:"FILE_STORAGE_PROVIDER" default:"local"`
	LocalConfig LocalStorageConfig
	S3Config    S3StorageConfig
}

type LocalStorageConfig struct {
	BaseDir    string `json:"base_dir" env:"UPLOAD_DIR" default:"uploads"`
	PublicPath string `json:"public_path" env:"PUBLIC_PATH" default:"/uploads"`
	BaseURL    string `json:"base_url" env:"BASE_URL" default:"http://localhost:8080"`
	MaxSize    int64  `json:"max_size" env:"MAX_FILE_SIZE" default:"5242880"` // 5MB
}

type S3StorageConfig struct {
	Region          string `json:"region" env:"AWS_REGION"`
	Bucket          string `json:"bucket" env:"AWS_S3_BUCKET"`
	AccessKeyID     string `json:"access_key_id" env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secret_access_key" env:"AWS_SECRET_ACCESS_KEY"`
	CDNDomain       string `json:"cdn_domain" env:"AWS_CLOUDFRONT_DOMAIN"`
}

func LoadFileStorageConfig() *FileStorageConfig {
	maxSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)
	if maxSize == 0 {
		maxSize = 5242880 // 5MB default
	}

	return &FileStorageConfig{
		Provider: getEnvOrDefault("FILE_STORAGE_PROVIDER", "local"),
		LocalConfig: LocalStorageConfig{
			BaseDir:    getEnvOrDefault("UPLOAD_DIR", "uploads"),
			PublicPath: getEnvOrDefault("PUBLIC_PATH", "/uploads"),
			BaseURL:    getEnvOrDefault("BASE_URL", "http://localhost:8080"),
			MaxSize:    maxSize,
		},
		S3Config: S3StorageConfig{
			Region:          os.Getenv("AWS_REGION"),
			Bucket:          os.Getenv("AWS_S3_BUCKET"),
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			CDNDomain:       os.Getenv("AWS_CLOUDFRONT_DOMAIN"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
