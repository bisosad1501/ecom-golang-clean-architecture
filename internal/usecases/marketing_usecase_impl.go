package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
)

type marketingUseCaseImpl struct {
	marketingRepo repositories.MarketingRepository
}

// NewMarketingUseCase creates a new marketing use case
func NewMarketingUseCase(marketingRepo repositories.MarketingRepository) MarketingUseCase {
	return &marketingUseCaseImpl{
		marketingRepo: marketingRepo,
	}
}

// Email Templates
func (uc *marketingUseCaseImpl) GetEmailTemplates(ctx context.Context, req GetEmailTemplatesRequest) (*GetEmailTemplatesResponse, error) {
	filters := repositories.EmailTemplateFilters{
		Search:    req.Search,
		Type:      req.Type,
		IsActive:  req.IsActive,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    (req.Page - 1) * req.Limit,
	}

	templates, err := uc.marketingRepo.GetEmailTemplates(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.marketingRepo.CountEmailTemplates(ctx, filters)
	if err != nil {
		return nil, err
	}

	var templateResponses []EmailTemplateResponse
	for _, template := range templates {
		var variables []EmailVariableResponse
		// Convert map to slice for response
		for key, value := range template.Variables {
			variables = append(variables, EmailVariableResponse{
				Name:        key,
				Description: fmt.Sprintf("%v", value),
				Type:        "string",
				Required:    false,
			})
		}

		templateResponses = append(templateResponses, EmailTemplateResponse{
			ID:            template.ID,
			Name:          template.Name,
			Subject:       template.Subject,
			Type:          string(template.Type), // Convert EmailType to string
			Content:       template.BodyHTML,     // Use BodyHTML as Content
			PlainText:     template.BodyText,     // Use BodyText as PlainText
			IsActive:      template.IsActive,
			Variables:     variables,
			UsageCount:    0,        // Default value since not in entity
			CreatedBy:     uuid.Nil, // Default value since not in entity
			CreatedByName: "Admin",  // TODO: Get actual user name
			CreatedAt:     template.CreatedAt,
			UpdatedAt:     template.UpdatedAt,
		})
	}

	return &GetEmailTemplatesResponse{
		Templates: templateResponses,
		Total:     total,
		Pagination: &PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		},
	}, nil
}

func (uc *marketingUseCaseImpl) GetEmailTemplate(ctx context.Context, id uuid.UUID) (*EmailTemplateResponse, error) {
	template, err := uc.marketingRepo.GetEmailTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	var variables []EmailVariableResponse
	// Convert map to slice for response
	for key, value := range template.Variables {
		variables = append(variables, EmailVariableResponse{
			Name:        key,
			Description: fmt.Sprintf("%v", value),
			Type:        "string",
			Required:    false,
		})
	}

	return &EmailTemplateResponse{
		ID:            template.ID,
		Name:          template.Name,
		Subject:       template.Subject,
		Type:          string(template.Type), // Convert EmailType to string
		Content:       template.BodyHTML,     // Use BodyHTML as Content
		PlainText:     template.BodyText,     // Use BodyText as PlainText
		IsActive:      template.IsActive,
		Variables:     variables,
		UsageCount:    0,        // Default value since not in entity
		CreatedBy:     uuid.Nil, // Default value since not in entity
		CreatedByName: "Admin",  // TODO: Get actual user name
		CreatedAt:     template.CreatedAt,
		UpdatedAt:     template.UpdatedAt,
	}, nil
}

func (uc *marketingUseCaseImpl) CreateEmailTemplate(ctx context.Context, req CreateEmailTemplateRequest) (*EmailTemplateResponse, error) {
	// Convert variables slice to map
	variablesMap := make(map[string]interface{})
	for _, v := range req.Variables {
		variablesMap[v.Name] = v.Description
	}

	template := &entities.EmailTemplate{
		Name:        req.Name,
		Subject:     req.Subject,
		Type:        entities.EmailType(req.Type), // Convert string to EmailType
		BodyHTML:    req.Content,                  // Use Content as BodyHTML
		BodyText:    req.PlainText,                // Use PlainText as BodyText
		IsActive:    req.IsActive,
		Variables:   variablesMap,
		Description: "Created via API", // Default description
	}

	err := uc.marketingRepo.CreateEmailTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return uc.GetEmailTemplate(ctx, template.ID)
}

func (uc *marketingUseCaseImpl) UpdateEmailTemplate(ctx context.Context, id uuid.UUID, req UpdateEmailTemplateRequest) (*EmailTemplateResponse, error) {
	template, err := uc.marketingRepo.GetEmailTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Subject != "" {
		template.Subject = req.Subject
	}
	if req.Type != "" {
		template.Type = entities.EmailType(req.Type) // Convert string to EmailType
	}
	if req.Content != "" {
		template.BodyHTML = req.Content // Use Content as BodyHTML
	}
	template.BodyText = req.PlainText // Use PlainText as BodyText
	template.IsActive = req.IsActive

	if req.Variables != nil {
		// Convert variables slice to map
		variablesMap := make(map[string]interface{})
		for _, v := range req.Variables {
			variablesMap[v.Name] = v.Description
		}
		template.Variables = variablesMap
	}

	err = uc.marketingRepo.UpdateEmailTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return uc.GetEmailTemplate(ctx, id)
}

func (uc *marketingUseCaseImpl) DeleteEmailTemplate(ctx context.Context, id uuid.UUID) error {
	return uc.marketingRepo.DeleteEmailTemplate(ctx, id)
}

// Campaigns
func (uc *marketingUseCaseImpl) GetCampaigns(ctx context.Context, req GetCampaignsRequest) (*GetCampaignsResponse, error) {
	filters := repositories.CampaignFilters{
		Search:    req.Search,
		Status:    req.Status,
		Type:      req.Type,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    (req.Page - 1) * req.Limit,
	}

	campaigns, err := uc.marketingRepo.GetCampaigns(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.marketingRepo.CountCampaigns(ctx, filters)
	if err != nil {
		return nil, err
	}

	var campaignResponses []CampaignResponse
	for _, campaign := range campaigns {
		var targetAudience TargetAudienceResponse
		var goals []CampaignGoalResponse
		var metrics CampaignMetrics

		json.Unmarshal([]byte(campaign.TargetAudience), &targetAudience)
		json.Unmarshal([]byte(campaign.Goals), &goals)
		json.Unmarshal([]byte(campaign.Metrics), &metrics)

		campaignResponses = append(campaignResponses, CampaignResponse{
			ID:             campaign.ID,
			Name:           campaign.Name,
			Description:    campaign.Description,
			Type:           campaign.Type,
			Status:         campaign.Status,
			Budget:         campaign.Budget,
			Spent:          campaign.Spent,
			StartDate:      campaign.StartDate,
			EndDate:        campaign.EndDate,
			TargetAudience: targetAudience,
			Goals:          goals,
			Metrics:        metrics,
			CreatedBy:      campaign.CreatedBy,
			CreatedByName:  "Admin", // TODO: Get actual user name
			CreatedAt:      campaign.CreatedAt,
			UpdatedAt:      campaign.UpdatedAt,
		})
	}

	return &GetCampaignsResponse{
		Campaigns: campaignResponses,
		Total:     total,
		Pagination: &PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		},
	}, nil
}

func (uc *marketingUseCaseImpl) GetCampaign(ctx context.Context, id uuid.UUID) (*CampaignResponse, error) {
	campaign, err := uc.marketingRepo.GetCampaign(ctx, id)
	if err != nil {
		return nil, err
	}

	var targetAudience TargetAudienceResponse
	var goals []CampaignGoalResponse
	var metrics CampaignMetrics

	json.Unmarshal([]byte(campaign.TargetAudience), &targetAudience)
	json.Unmarshal([]byte(campaign.Goals), &goals)
	json.Unmarshal([]byte(campaign.Metrics), &metrics)

	return &CampaignResponse{
		ID:             campaign.ID,
		Name:           campaign.Name,
		Description:    campaign.Description,
		Type:           campaign.Type,
		Status:         campaign.Status,
		Budget:         campaign.Budget,
		Spent:          campaign.Spent,
		StartDate:      campaign.StartDate,
		EndDate:        campaign.EndDate,
		TargetAudience: targetAudience,
		Goals:          goals,
		Metrics:        metrics,
		CreatedBy:      campaign.CreatedBy,
		CreatedByName:  "Admin", // TODO: Get actual user name
		CreatedAt:      campaign.CreatedAt,
		UpdatedAt:      campaign.UpdatedAt,
	}, nil
}

func (uc *marketingUseCaseImpl) CreateCampaign(ctx context.Context, req CreateCampaignRequest) (*CampaignResponse, error) {
	targetAudienceJSON, _ := json.Marshal(req.TargetAudience)
	goalsJSON, _ := json.Marshal(req.Goals)
	settingsJSON, _ := json.Marshal(req.Settings)

	campaign := &entities.Campaign{
		Name:           req.Name,
		Description:    req.Description,
		Type:           req.Type,
		Status:         req.Status,
		Budget:         req.Budget,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		TargetAudience: string(targetAudienceJSON),
		Goals:          string(goalsJSON),
		Settings:       string(settingsJSON),
		CreatedBy:      req.CreatedBy,
	}

	err := uc.marketingRepo.CreateCampaign(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return uc.GetCampaign(ctx, campaign.ID)
}

func (uc *marketingUseCaseImpl) UpdateCampaign(ctx context.Context, id uuid.UUID, req UpdateCampaignRequest) (*CampaignResponse, error) {
	// Implementation similar to CreateCampaign but for updates
	return nil, fmt.Errorf("not implemented")
}

func (uc *marketingUseCaseImpl) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	return uc.marketingRepo.DeleteCampaign(ctx, id)
}

func (uc *marketingUseCaseImpl) LaunchCampaign(ctx context.Context, id uuid.UUID) error {
	campaign, err := uc.marketingRepo.GetCampaign(ctx, id)
	if err != nil {
		return err
	}

	campaign.Status = "active"
	return uc.marketingRepo.UpdateCampaign(ctx, campaign)
}

func (uc *marketingUseCaseImpl) PauseCampaign(ctx context.Context, id uuid.UUID) error {
	campaign, err := uc.marketingRepo.GetCampaign(ctx, id)
	if err != nil {
		return err
	}

	campaign.Status = "paused"
	return uc.marketingRepo.UpdateCampaign(ctx, campaign)
}

func (uc *marketingUseCaseImpl) GetCampaignAnalytics(ctx context.Context, id uuid.UUID, period string) (*CampaignAnalyticsResponse, error) {
	// TODO: Implement campaign analytics
	return &CampaignAnalyticsResponse{
		CampaignID: id,
		Period:     period,
		Metrics: CampaignMetrics{
			Impressions:       1000,
			Clicks:            50,
			Conversions:       5,
			Revenue:           500.0,
			CTR:               5.0,
			ConversionRate:    10.0,
			CostPerClick:      2.0,
			CostPerConversion: 20.0,
			ROI:               400.0,
			ROAS:              5.0,
		},
		Timeline:             []struct{}{},
		TopPerformingContent: []struct{}{},
	}, nil
}
