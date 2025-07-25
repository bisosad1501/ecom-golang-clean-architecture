package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// MarketingUseCase defines marketing-related operations
type MarketingUseCase interface {
	// Email Templates
	GetEmailTemplates(ctx context.Context, req GetEmailTemplatesRequest) (*GetEmailTemplatesResponse, error)
	GetEmailTemplate(ctx context.Context, id uuid.UUID) (*EmailTemplateResponse, error)
	CreateEmailTemplate(ctx context.Context, req CreateEmailTemplateRequest) (*EmailTemplateResponse, error)
	UpdateEmailTemplate(ctx context.Context, id uuid.UUID, req UpdateEmailTemplateRequest) (*EmailTemplateResponse, error)
	DeleteEmailTemplate(ctx context.Context, id uuid.UUID) error

	// Campaigns
	GetCampaigns(ctx context.Context, req GetCampaignsRequest) (*GetCampaignsResponse, error)
	GetCampaign(ctx context.Context, id uuid.UUID) (*CampaignResponse, error)
	CreateCampaign(ctx context.Context, req CreateCampaignRequest) (*CampaignResponse, error)
	UpdateCampaign(ctx context.Context, id uuid.UUID, req UpdateCampaignRequest) (*CampaignResponse, error)
	DeleteCampaign(ctx context.Context, id uuid.UUID) error
	LaunchCampaign(ctx context.Context, id uuid.UUID) error
	PauseCampaign(ctx context.Context, id uuid.UUID) error
	GetCampaignAnalytics(ctx context.Context, id uuid.UUID, period string) (*CampaignAnalyticsResponse, error)
}

// Campaign types
type GetCampaignsRequest struct {
	Search    string `json:"search" form:"search"`
	Status    string `json:"status" form:"status"`
	Type      string `json:"type" form:"type"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

type GetCampaignsResponse struct {
	Campaigns  []CampaignResponse `json:"campaigns"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

type CampaignResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Type            string                   `json:"type"`
	Status          string                   `json:"status"`
	Budget          float64                  `json:"budget"`
	Spent           float64                  `json:"spent"`
	StartDate       time.Time                `json:"start_date"`
	EndDate         *time.Time               `json:"end_date"`
	TargetAudience  TargetAudienceResponse   `json:"target_audience"`
	Goals           []CampaignGoalResponse   `json:"goals"`
	Metrics         CampaignMetrics          `json:"metrics"`
	CreatedBy       uuid.UUID                `json:"created_by"`
	CreatedByName   string                   `json:"created_by_name"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

type CreateCampaignRequest struct {
	Name           string                   `json:"name" validate:"required"`
	Description    string                   `json:"description"`
	Type           string                   `json:"type" validate:"required"`
	Status         string                   `json:"status"`
	Budget         float64                  `json:"budget"`
	StartDate      time.Time                `json:"start_date" validate:"required"`
	EndDate        *time.Time               `json:"end_date"`
	TargetAudience TargetAudienceRequest    `json:"target_audience"`
	Goals          []CampaignGoalRequest    `json:"goals"`
	Settings       map[string]interface{}   `json:"settings"`
	CreatedBy      uuid.UUID                `json:"created_by" validate:"required"`
}

type UpdateCampaignRequest struct {
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Type           string                   `json:"type"`
	Status         string                   `json:"status"`
	Budget         float64                  `json:"budget"`
	StartDate      time.Time                `json:"start_date"`
	EndDate        *time.Time               `json:"end_date"`
	TargetAudience TargetAudienceRequest    `json:"target_audience"`
	Goals          []CampaignGoalRequest    `json:"goals"`
	Settings       map[string]interface{}   `json:"settings"`
}

type TargetAudienceRequest struct {
	AgeRange     *AgeRange `json:"age_range"`
	Gender       []string  `json:"gender"`
	Locations    []string  `json:"locations"`
	Interests    []string  `json:"interests"`
	Behaviors    []string  `json:"behaviors"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type TargetAudienceResponse struct {
	AgeRange     *AgeRange `json:"age_range"`
	Gender       []string  `json:"gender"`
	Locations    []string  `json:"locations"`
	Interests    []string  `json:"interests"`
	Behaviors    []string  `json:"behaviors"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type AgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type CampaignGoalRequest struct {
	Type        string  `json:"type" validate:"required"`
	Target      float64 `json:"target"`
	Description string  `json:"description"`
}

type CampaignGoalResponse struct {
	Type        string  `json:"type"`
	Target      float64 `json:"target"`
	Current     float64 `json:"current"`
	Description string  `json:"description"`
	Progress    float64 `json:"progress"`
}

type CampaignMetrics struct {
	Impressions       int64   `json:"impressions"`
	Clicks            int64   `json:"clicks"`
	Conversions       int64   `json:"conversions"`
	Revenue           float64 `json:"revenue"`
	CTR               float64 `json:"ctr"`
	ConversionRate    float64 `json:"conversion_rate"`
	CostPerClick      float64 `json:"cost_per_click"`
	CostPerConversion float64 `json:"cost_per_conversion"`
	ROI               float64 `json:"roi"`
	ROAS              float64 `json:"roas"`
}

type CampaignAnalyticsResponse struct {
	CampaignID           uuid.UUID       `json:"campaign_id"`
	Period               string          `json:"period"`
	Metrics              CampaignMetrics `json:"metrics"`
	Timeline             []struct{}      `json:"timeline"`
	TopPerformingContent []struct{}      `json:"top_performing_content"`
}
