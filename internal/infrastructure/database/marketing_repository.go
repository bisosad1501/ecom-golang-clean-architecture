package database

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type marketingRepository struct {
	db *gorm.DB
}

// NewMarketingRepository creates a new marketing repository
func NewMarketingRepository(db *gorm.DB) repositories.MarketingRepository {
	return &marketingRepository{db: db}
}

// Email Templates
func (r *marketingRepository) GetEmailTemplates(ctx context.Context, filters repositories.EmailTemplateFilters) ([]*entities.EmailTemplate, error) {
	query := r.db.WithContext(ctx).Model(&entities.EmailTemplate{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR subject ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	// Sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var templates []*entities.EmailTemplate
	err := query.Find(&templates).Error
	return templates, err
}

func (r *marketingRepository) GetEmailTemplate(ctx context.Context, id uuid.UUID) (*entities.EmailTemplate, error) {
	var template entities.EmailTemplate
	err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email template not found")
		}
		return nil, err
	}
	return &template, nil
}

func (r *marketingRepository) CreateEmailTemplate(ctx context.Context, template *entities.EmailTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *marketingRepository) UpdateEmailTemplate(ctx context.Context, template *entities.EmailTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *marketingRepository) DeleteEmailTemplate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.EmailTemplate{}, "id = ?", id).Error
}

func (r *marketingRepository) CountEmailTemplates(ctx context.Context, filters repositories.EmailTemplateFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.EmailTemplate{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR subject ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Campaigns
func (r *marketingRepository) GetCampaigns(ctx context.Context, filters repositories.CampaignFilters) ([]*entities.Campaign, error) {
	query := r.db.WithContext(ctx).Model(&entities.Campaign{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	// Sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var campaigns []*entities.Campaign
	err := query.Find(&campaigns).Error
	return campaigns, err
}

func (r *marketingRepository) GetCampaign(ctx context.Context, id uuid.UUID) (*entities.Campaign, error) {
	var campaign entities.Campaign
	err := r.db.WithContext(ctx).First(&campaign, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, err
	}
	return &campaign, nil
}

func (r *marketingRepository) CreateCampaign(ctx context.Context, campaign *entities.Campaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *marketingRepository) UpdateCampaign(ctx context.Context, campaign *entities.Campaign) error {
	return r.db.WithContext(ctx).Save(campaign).Error
}

func (r *marketingRepository) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Campaign{}, "id = ?", id).Error
}

func (r *marketingRepository) CountCampaigns(ctx context.Context, filters repositories.CampaignFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Campaign{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Promotions
func (r *marketingRepository) GetPromotions(ctx context.Context, filters repositories.PromotionFilters) ([]*entities.Promotion, error) {
	query := r.db.WithContext(ctx).Model(&entities.Promotion{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var promotions []*entities.Promotion
	err := query.Find(&promotions).Error
	return promotions, err
}

func (r *marketingRepository) GetPromotion(ctx context.Context, id uuid.UUID) (*entities.Promotion, error) {
	var promotion entities.Promotion
	err := r.db.WithContext(ctx).First(&promotion, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("promotion not found")
		}
		return nil, err
	}
	return &promotion, nil
}

func (r *marketingRepository) CreatePromotion(ctx context.Context, promotion *entities.Promotion) error {
	return r.db.WithContext(ctx).Create(promotion).Error
}

func (r *marketingRepository) UpdatePromotion(ctx context.Context, promotion *entities.Promotion) error {
	return r.db.WithContext(ctx).Save(promotion).Error
}

func (r *marketingRepository) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Promotion{}, "id = ?", id).Error
}

func (r *marketingRepository) CountPromotions(ctx context.Context, filters repositories.PromotionFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Promotion{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Email Campaigns
func (r *marketingRepository) GetEmailCampaigns(ctx context.Context, filters repositories.EmailCampaignFilters) ([]*entities.EmailCampaign, error) {
	query := r.db.WithContext(ctx).Preload("Template")

	if filters.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var campaigns []*entities.EmailCampaign
	err := query.Find(&campaigns).Error
	return campaigns, err
}

func (r *marketingRepository) GetEmailCampaign(ctx context.Context, id uuid.UUID) (*entities.EmailCampaign, error) {
	var campaign entities.EmailCampaign
	err := r.db.WithContext(ctx).Preload("Template").First(&campaign, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email campaign not found")
		}
		return nil, err
	}
	return &campaign, nil
}

func (r *marketingRepository) CreateEmailCampaign(ctx context.Context, campaign *entities.EmailCampaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *marketingRepository) UpdateEmailCampaign(ctx context.Context, campaign *entities.EmailCampaign) error {
	return r.db.WithContext(ctx).Save(campaign).Error
}

func (r *marketingRepository) DeleteEmailCampaign(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.EmailCampaign{}, "id = ?", id).Error
}

func (r *marketingRepository) CountEmailCampaigns(ctx context.Context, filters repositories.EmailCampaignFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.EmailCampaign{})

	if filters.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Newsletter Subscribers
func (r *marketingRepository) GetNewsletterSubscribers(ctx context.Context, filters repositories.NewsletterSubscriberFilters) ([]*entities.NewsletterSubscriber, error) {
	query := r.db.WithContext(ctx).Model(&entities.NewsletterSubscriber{})

	if filters.Search != "" {
		query = query.Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Segment != "" {
		query = query.Where("segments LIKE ?", "%"+filters.Segment+"%")
	}

	// Sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		direction := "ASC"
		if filters.SortOrder == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, direction)
	}
	query = query.Order(orderBy)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var subscribers []*entities.NewsletterSubscriber
	err := query.Find(&subscribers).Error
	return subscribers, err
}

func (r *marketingRepository) GetNewsletterSubscriber(ctx context.Context, id uuid.UUID) (*entities.NewsletterSubscriber, error) {
	var subscriber entities.NewsletterSubscriber
	err := r.db.WithContext(ctx).First(&subscriber, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("newsletter subscriber not found")
		}
		return nil, err
	}
	return &subscriber, nil
}

func (r *marketingRepository) GetNewsletterSubscriberByEmail(ctx context.Context, email string) (*entities.NewsletterSubscriber, error) {
	var subscriber entities.NewsletterSubscriber
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&subscriber).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("newsletter subscriber not found")
		}
		return nil, err
	}
	return &subscriber, nil
}

func (r *marketingRepository) CreateNewsletterSubscriber(ctx context.Context, subscriber *entities.NewsletterSubscriber) error {
	return r.db.WithContext(ctx).Create(subscriber).Error
}

func (r *marketingRepository) UpdateNewsletterSubscriber(ctx context.Context, subscriber *entities.NewsletterSubscriber) error {
	return r.db.WithContext(ctx).Save(subscriber).Error
}

func (r *marketingRepository) DeleteNewsletterSubscriber(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.NewsletterSubscriber{}, "id = ?", id).Error
}

func (r *marketingRepository) CountNewsletterSubscribers(ctx context.Context, filters repositories.NewsletterSubscriberFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.NewsletterSubscriber{})

	if filters.Search != "" {
		query = query.Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Segment != "" {
		query = query.Where("segments LIKE ?", "%"+filters.Segment+"%")
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Loyalty Programs
func (r *marketingRepository) GetLoyaltyPrograms(ctx context.Context) ([]*entities.LoyaltyProgram, error) {
	var programs []*entities.LoyaltyProgram
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&programs).Error
	return programs, err
}

func (r *marketingRepository) GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*entities.LoyaltyProgram, error) {
	var program entities.LoyaltyProgram
	err := r.db.WithContext(ctx).First(&program, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("loyalty program not found")
		}
		return nil, err
	}
	return &program, nil
}

func (r *marketingRepository) CreateLoyaltyProgram(ctx context.Context, program *entities.LoyaltyProgram) error {
	return r.db.WithContext(ctx).Create(program).Error
}

func (r *marketingRepository) UpdateLoyaltyProgram(ctx context.Context, program *entities.LoyaltyProgram) error {
	return r.db.WithContext(ctx).Save(program).Error
}

func (r *marketingRepository) DeleteLoyaltyProgram(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.LoyaltyProgram{}, "id = ?", id).Error
}

// Referral Programs
func (r *marketingRepository) GetReferralPrograms(ctx context.Context) ([]*entities.ReferralProgram, error) {
	var programs []*entities.ReferralProgram
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&programs).Error
	return programs, err
}

func (r *marketingRepository) GetReferralProgram(ctx context.Context, id uuid.UUID) (*entities.ReferralProgram, error) {
	var program entities.ReferralProgram
	err := r.db.WithContext(ctx).First(&program, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("referral program not found")
		}
		return nil, err
	}
	return &program, nil
}

func (r *marketingRepository) CreateReferralProgram(ctx context.Context, program *entities.ReferralProgram) error {
	return r.db.WithContext(ctx).Create(program).Error
}

func (r *marketingRepository) UpdateReferralProgram(ctx context.Context, program *entities.ReferralProgram) error {
	return r.db.WithContext(ctx).Save(program).Error
}

func (r *marketingRepository) DeleteReferralProgram(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ReferralProgram{}, "id = ?", id).Error
}
