-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE file_uploads (
    id VARCHAR(36) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    object_key VARCHAR(500) NOT NULL UNIQUE,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL,
    uploaded_by VARCHAR(36),
    upload_type ENUM('admin', 'user', 'public') NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uploaded_by (uploaded_by),
    INDEX idx_upload_type (upload_type),
    INDEX idx_category (category),
    INDEX idx_object_key (object_key),
    INDEX idx_created_at (created_at)
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS file_uploads;
