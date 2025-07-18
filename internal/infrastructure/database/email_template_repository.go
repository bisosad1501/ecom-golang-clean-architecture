package database

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailTemplateRepository struct {
	db *gorm.DB
}

// NewEmailTemplateRepository creates a new email template repository
func NewEmailTemplateRepository(db *gorm.DB) repositories.EmailTemplateRepository {
	return &emailTemplateRepository{db: db}
}

// Create creates a new email template
func (r *emailTemplateRepository) Create(ctx context.Context, template *entities.EmailTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID gets an email template by ID
func (r *emailTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.EmailTemplate, error) {
	var template entities.EmailTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByName gets an email template by name
func (r *emailTemplateRepository) GetByName(ctx context.Context, name string) (*entities.EmailTemplate, error) {
	var template entities.EmailTemplate
	err := r.db.WithContext(ctx).Where("name = ? AND is_active = true", name).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update updates an email template
func (r *emailTemplateRepository) Update(ctx context.Context, template *entities.EmailTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete deletes an email template
func (r *emailTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.EmailTemplate{}, id).Error
}

// List lists email templates with pagination
func (r *emailTemplateRepository) List(ctx context.Context, offset, limit int) ([]*entities.EmailTemplate, error) {
	var templates []*entities.EmailTemplate
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&templates).Error
	return templates, err
}

// GetByType gets email templates by type
func (r *emailTemplateRepository) GetByType(ctx context.Context, emailType entities.EmailType) ([]*entities.EmailTemplate, error) {
	var templates []*entities.EmailTemplate
	err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = true", emailType).
		Order("created_at DESC").
		Find(&templates).Error
	return templates, err
}

// GetActive gets all active email templates
func (r *emailTemplateRepository) GetActive(ctx context.Context) ([]*entities.EmailTemplate, error) {
	var templates []*entities.EmailTemplate
	err := r.db.WithContext(ctx).
		Where("is_active = true").
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}

// GetLatestVersion gets the latest version of a template by name
func (r *emailTemplateRepository) GetLatestVersion(ctx context.Context, name string) (*entities.EmailTemplate, error) {
	var template entities.EmailTemplate
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		Order("version DESC").
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByVersion gets a specific version of a template
func (r *emailTemplateRepository) GetByVersion(ctx context.Context, name string, version int) (*entities.EmailTemplate, error) {
	var template entities.EmailTemplate
	err := r.db.WithContext(ctx).
		Where("name = ? AND version = ?", name, version).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}
