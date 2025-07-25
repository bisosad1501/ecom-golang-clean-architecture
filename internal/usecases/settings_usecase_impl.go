package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"ecom-golang-clean-architecture/internal/domain/services"

	"github.com/google/uuid"
)

type settingsUseCaseImpl struct {
	settingsRepo repositories.SettingsRepository
	cache        services.SettingsCache
	auditService services.AuditService
	validator    *SettingsValidator
}

// NewSettingsUseCase creates a new settings use case
func NewSettingsUseCase(settingsRepo repositories.SettingsRepository) SettingsUseCase {
	return &settingsUseCaseImpl{
		settingsRepo: settingsRepo,
		cache:        nil, // Optional
		auditService: nil, // Optional
		validator:    NewSettingsValidator(),
	}
}

// NewSettingsUseCaseWithServices creates a new settings use case with cache and audit services
func NewSettingsUseCaseWithServices(settingsRepo repositories.SettingsRepository, cache services.SettingsCache, auditService services.AuditService) SettingsUseCase {
	return &settingsUseCaseImpl{
		settingsRepo: settingsRepo,
		cache:        cache,
		auditService: auditService,
		validator:    NewSettingsValidator(),
	}
}

// General Settings
func (uc *settingsUseCaseImpl) GetGeneralSettings(ctx context.Context) (*GeneralSettingsResponse, error) {
	var settings []*entities.Setting
	var err error

	// Try to get from cache first
	if uc.cache != nil {
		settings, err = uc.cache.GetSettingsByCategory(ctx, "general")
		if err == nil && settings != nil {
			// Cache hit - use cached data
		} else {
			// Cache miss - get from database
			settings, err = uc.settingsRepo.GetSettingsByCategory(ctx, "general")
			if err != nil {
				return nil, fmt.Errorf("failed to get general settings: %w", err)
			}
			// Cache the result
			uc.cache.SetSettingsByCategory(ctx, "general", settings)
		}
	} else {
		// No cache - get directly from database
		settings, err = uc.settingsRepo.GetSettingsByCategory(ctx, "general")
		if err != nil {
			return nil, fmt.Errorf("failed to get general settings: %w", err)
		}
	}

	response := &GeneralSettingsResponse{
		StoreName:        "My Store",
		StoreDescription: "Welcome to our store",
		StoreEmail:       "admin@store.com",
		StorePhone:       "",
		StoreAddress:     "",
		Currency:         "USD",
		Timezone:         "UTC",
		Language:         "en",
		DateFormat:       "Y-m-d",
		TimeFormat:       "H:i:s",
		UpdatedAt:        time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "store_name":
			response.StoreName = setting.Value
		case "store_description":
			response.StoreDescription = setting.Value
		case "store_email":
			response.StoreEmail = setting.Value
		case "store_phone":
			response.StorePhone = setting.Value
		case "store_address":
			response.StoreAddress = setting.Value
		case "currency":
			response.Currency = setting.Value
		case "timezone":
			response.Timezone = setting.Value
		case "language":
			response.Language = setting.Value
		case "date_format":
			response.DateFormat = setting.Value
		case "time_format":
			response.TimeFormat = setting.Value
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateGeneralSettings(ctx context.Context, req UpdateGeneralSettingsRequest) error {
	// Validate request
	if err := uc.validator.ValidateGeneralSettings(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get current settings for audit logging
	currentSettings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "general")
	if err != nil {
		return fmt.Errorf("failed to get current settings: %w", err)
	}

	// Create map of current values for comparison
	currentValues := make(map[string]interface{})
	for _, setting := range currentSettings {
		currentValues[setting.Key] = setting.Value
	}

	// Prepare new settings
	settings := []*entities.Setting{
		{Key: "store_name", Value: req.StoreName, Type: "string", Category: "general"},
		{Key: "store_description", Value: req.StoreDescription, Type: "string", Category: "general"},
		{Key: "store_email", Value: req.StoreEmail, Type: "string", Category: "general"},
		{Key: "store_phone", Value: req.StorePhone, Type: "string", Category: "general"},
		{Key: "store_address", Value: req.StoreAddress, Type: "string", Category: "general"},
		{Key: "currency", Value: req.Currency, Type: "string", Category: "general"},
		{Key: "timezone", Value: req.Timezone, Type: "string", Category: "general"},
		{Key: "language", Value: req.Language, Type: "string", Category: "general"},
		{Key: "date_format", Value: req.DateFormat, Type: "string", Category: "general"},
		{Key: "time_format", Value: req.TimeFormat, Type: "string", Category: "general"},
	}

	// Track changes for audit
	var changes []entities.SettingsChangeDetail
	newValues := make(map[string]interface{})
	for _, setting := range settings {
		newValues[setting.Key] = setting.Value
		if currentVal, exists := currentValues[setting.Key]; !exists || currentVal != setting.Value {
			changes = append(changes, entities.SettingsChangeDetail{
				Key:      setting.Key,
				OldValue: currentVal,
				NewValue: setting.Value,
				Type:     setting.Type,
			})
		}
	}

	// Update settings in database
	if err := uc.settingsRepo.SetSettings(ctx, settings); err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	// Invalidate cache
	if uc.cache != nil {
		uc.cache.InvalidateCategory(ctx, "general")
		// Also cache the new settings
		uc.cache.SetSettingsByCategory(ctx, "general", settings)
	}

	// Log audit trail (if audit service is available and there are changes)
	if uc.auditService != nil && len(changes) > 0 {
		// Note: In a real implementation, you would get user info from context
		// For now, we'll use a placeholder
		userID := uuid.New()             // This should come from authentication context
		userEmail := "admin@example.com" // This should come from authentication context

		uc.auditService.LogSettingsChange(
			ctx,
			userID,
			userEmail,
			"general",
			changes,
			"", // IP address - should come from request context
			"", // User agent - should come from request context
			"", // Session ID - should come from request context
		)
	}

	return nil
}

// Store Configuration
func (uc *settingsUseCaseImpl) GetStoreConfig(ctx context.Context) (*StoreConfigResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "store")
	if err != nil {
		return nil, err
	}

	response := &StoreConfigResponse{
		MaintenanceMode:    false,
		MaintenanceMessage: "We are currently under maintenance. Please check back later.",
		AllowRegistration:  true,
		RequireEmailVerify: true,
		AutoApproveReviews: false,
		EnableWishlist:     true,
		EnableCompare:      true,
		EnableReviews:      true,
		EnableCoupons:      true,
		EnableInventory:    true,
		EnableMultiVendor:  false,
		UpdatedAt:          time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "maintenance_mode":
			response.MaintenanceMode = setting.Value == "true"
		case "maintenance_message":
			response.MaintenanceMessage = setting.Value
		case "allow_registration":
			response.AllowRegistration = setting.Value == "true"
		case "require_email_verify":
			response.RequireEmailVerify = setting.Value == "true"
		case "auto_approve_reviews":
			response.AutoApproveReviews = setting.Value == "true"
		case "enable_wishlist":
			response.EnableWishlist = setting.Value == "true"
		case "enable_compare":
			response.EnableCompare = setting.Value == "true"
		case "enable_reviews":
			response.EnableReviews = setting.Value == "true"
		case "enable_coupons":
			response.EnableCoupons = setting.Value == "true"
		case "enable_inventory":
			response.EnableInventory = setting.Value == "true"
		case "enable_multi_vendor":
			response.EnableMultiVendor = setting.Value == "true"
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateStoreConfig(ctx context.Context, req UpdateStoreConfigRequest) error {
	settings := []*entities.Setting{
		{Key: "maintenance_mode", Value: fmt.Sprintf("%t", req.MaintenanceMode), Type: "boolean", Category: "store"},
		{Key: "maintenance_message", Value: req.MaintenanceMessage, Type: "string", Category: "store"},
		{Key: "allow_registration", Value: fmt.Sprintf("%t", req.AllowRegistration), Type: "boolean", Category: "store"},
		{Key: "require_email_verify", Value: fmt.Sprintf("%t", req.RequireEmailVerify), Type: "boolean", Category: "store"},
		{Key: "auto_approve_reviews", Value: fmt.Sprintf("%t", req.AutoApproveReviews), Type: "boolean", Category: "store"},
		{Key: "enable_wishlist", Value: fmt.Sprintf("%t", req.EnableWishlist), Type: "boolean", Category: "store"},
		{Key: "enable_compare", Value: fmt.Sprintf("%t", req.EnableCompare), Type: "boolean", Category: "store"},
		{Key: "enable_reviews", Value: fmt.Sprintf("%t", req.EnableReviews), Type: "boolean", Category: "store"},
		{Key: "enable_coupons", Value: fmt.Sprintf("%t", req.EnableCoupons), Type: "boolean", Category: "store"},
		{Key: "enable_inventory", Value: fmt.Sprintf("%t", req.EnableInventory), Type: "boolean", Category: "store"},
		{Key: "enable_multi_vendor", Value: fmt.Sprintf("%t", req.EnableMultiVendor), Type: "boolean", Category: "store"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// Payment Settings
func (uc *settingsUseCaseImpl) GetPaymentSettings(ctx context.Context) (*PaymentSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "payment")
	if err != nil {
		return nil, err
	}

	response := &PaymentSettingsResponse{
		DefaultCurrency:    "USD",
		AcceptedCurrencies: []string{"USD", "EUR", "GBP"},
		StripeEnabled:      false,
		StripePublicKey:    "",
		PayPalEnabled:      false,
		PayPalClientID:     "",
		PayPalSandbox:      true,
		CashOnDelivery:     true,
		BankTransfer:       false,
		BankDetails:        make(map[string]interface{}),
		UpdatedAt:          time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "default_currency":
			response.DefaultCurrency = setting.Value
		case "accepted_currencies":
			json.Unmarshal([]byte(setting.Value), &response.AcceptedCurrencies)
		case "stripe_enabled":
			response.StripeEnabled = setting.Value == "true"
		case "stripe_public_key":
			response.StripePublicKey = setting.Value
		case "paypal_enabled":
			response.PayPalEnabled = setting.Value == "true"
		case "paypal_client_id":
			response.PayPalClientID = setting.Value
		case "paypal_sandbox":
			response.PayPalSandbox = setting.Value == "true"
		case "cash_on_delivery":
			response.CashOnDelivery = setting.Value == "true"
		case "bank_transfer":
			response.BankTransfer = setting.Value == "true"
		case "bank_details":
			json.Unmarshal([]byte(setting.Value), &response.BankDetails)
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdatePaymentSettings(ctx context.Context, req UpdatePaymentSettingsRequest) error {
	acceptedCurrenciesJSON, _ := json.Marshal(req.AcceptedCurrencies)
	bankDetailsJSON, _ := json.Marshal(req.BankDetails)

	settings := []*entities.Setting{
		{Key: "default_currency", Value: req.DefaultCurrency, Type: "string", Category: "payment"},
		{Key: "accepted_currencies", Value: string(acceptedCurrenciesJSON), Type: "json", Category: "payment"},
		{Key: "stripe_enabled", Value: fmt.Sprintf("%t", req.StripeEnabled), Type: "boolean", Category: "payment"},
		{Key: "stripe_public_key", Value: req.StripePublicKey, Type: "string", Category: "payment"},
		{Key: "stripe_secret_key", Value: req.StripeSecretKey, Type: "string", Category: "payment", IsEncrypted: true},
		{Key: "paypal_enabled", Value: fmt.Sprintf("%t", req.PayPalEnabled), Type: "boolean", Category: "payment"},
		{Key: "paypal_client_id", Value: req.PayPalClientID, Type: "string", Category: "payment"},
		{Key: "paypal_client_secret", Value: req.PayPalClientSecret, Type: "string", Category: "payment", IsEncrypted: true},
		{Key: "paypal_sandbox", Value: fmt.Sprintf("%t", req.PayPalSandbox), Type: "boolean", Category: "payment"},
		{Key: "cash_on_delivery", Value: fmt.Sprintf("%t", req.CashOnDelivery), Type: "boolean", Category: "payment"},
		{Key: "bank_transfer", Value: fmt.Sprintf("%t", req.BankTransfer), Type: "boolean", Category: "payment"},
		{Key: "bank_details", Value: string(bankDetailsJSON), Type: "json", Category: "payment"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// Email Settings
func (uc *settingsUseCaseImpl) GetEmailSettings(ctx context.Context) (*EmailSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "email")
	if err != nil {
		return nil, err
	}

	response := &EmailSettingsResponse{
		SMTPHost:     "localhost",
		SMTPPort:     587,
		SMTPUsername: "",
		SMTPSecurity: "tls",
		FromEmail:    "noreply@store.com",
		FromName:     "My Store",
		ReplyToEmail: "",
		UpdatedAt:    time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "smtp_host":
			response.SMTPHost = setting.Value
		case "smtp_port":
			if port := setting.Value; port != "" {
				fmt.Sscanf(port, "%d", &response.SMTPPort)
			}
		case "smtp_username":
			response.SMTPUsername = setting.Value
		case "smtp_security":
			response.SMTPSecurity = setting.Value
		case "from_email":
			response.FromEmail = setting.Value
		case "from_name":
			response.FromName = setting.Value
		case "reply_to_email":
			response.ReplyToEmail = setting.Value
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateEmailSettings(ctx context.Context, req UpdateEmailSettingsRequest) error {
	settings := []*entities.Setting{
		{Key: "smtp_host", Value: req.SMTPHost, Type: "string", Category: "email"},
		{Key: "smtp_port", Value: fmt.Sprintf("%d", req.SMTPPort), Type: "number", Category: "email"},
		{Key: "smtp_username", Value: req.SMTPUsername, Type: "string", Category: "email"},
		{Key: "smtp_password", Value: req.SMTPPassword, Type: "string", Category: "email", IsEncrypted: true},
		{Key: "smtp_security", Value: req.SMTPSecurity, Type: "string", Category: "email"},
		{Key: "from_email", Value: req.FromEmail, Type: "string", Category: "email"},
		{Key: "from_name", Value: req.FromName, Type: "string", Category: "email"},
		{Key: "reply_to_email", Value: req.ReplyToEmail, Type: "string", Category: "email"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

func (uc *settingsUseCaseImpl) TestEmailSettings(ctx context.Context, req TestEmailRequest) error {
	// TODO: Implement email testing logic
	// This would send a test email using the current SMTP settings
	return fmt.Errorf("email testing not implemented yet")
}

// Tax Settings
func (uc *settingsUseCaseImpl) GetTaxSettings(ctx context.Context) (*TaxSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "tax")
	if err != nil {
		return nil, err
	}

	response := &TaxSettingsResponse{
		TaxEnabled:       true,
		TaxIncluded:      false,
		DefaultTaxRate:   0.0,
		TaxCalculation:   "exclusive",
		TaxRates:         []TaxRateResponse{},
		TaxClasses:       []TaxClassResponse{},
		TaxExemptions:    []string{},
		DigitalTaxRate:   0.0,
		ShippingTaxable:  false,
		TaxDisplayFormat: "excluding",
		TaxRounding:      "round_nearest",
		CustomTaxRules:   make(map[string]interface{}),
		UpdatedAt:        time.Now(),
	}

	// Get tax rates and classes
	taxRates, _ := uc.settingsRepo.GetTaxRates(ctx, repositories.TaxRateFilters{IsActive: &[]bool{true}[0]})
	for _, rate := range taxRates {
		response.TaxRates = append(response.TaxRates, TaxRateResponse{
			ID:          rate.ID,
			Name:        rate.Name,
			Rate:        rate.Rate,
			Country:     rate.Country,
			State:       rate.State,
			City:        rate.City,
			ZipCode:     rate.ZipCode,
			TaxClass:    rate.TaxClass,
			Priority:    rate.Priority,
			Compound:    rate.Compound,
			IsActive:    rate.IsActive,
			Description: rate.Description,
			CreatedAt:   rate.CreatedAt,
			UpdatedAt:   rate.UpdatedAt,
		})
	}

	taxClasses, _ := uc.settingsRepo.GetTaxClasses(ctx)
	for _, class := range taxClasses {
		response.TaxClasses = append(response.TaxClasses, TaxClassResponse{
			ID:          class.ID,
			Name:        class.Name,
			Description: class.Description,
			IsActive:    class.IsActive,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		})
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "tax_enabled":
			response.TaxEnabled = setting.Value == "true"
		case "tax_included":
			response.TaxIncluded = setting.Value == "true"
		case "default_tax_rate":
			fmt.Sscanf(setting.Value, "%f", &response.DefaultTaxRate)
		case "tax_calculation":
			response.TaxCalculation = setting.Value
		case "tax_exemptions":
			json.Unmarshal([]byte(setting.Value), &response.TaxExemptions)
		case "digital_tax_rate":
			fmt.Sscanf(setting.Value, "%f", &response.DigitalTaxRate)
		case "shipping_taxable":
			response.ShippingTaxable = setting.Value == "true"
		case "tax_display_format":
			response.TaxDisplayFormat = setting.Value
		case "tax_rounding":
			response.TaxRounding = setting.Value
		case "custom_tax_rules":
			json.Unmarshal([]byte(setting.Value), &response.CustomTaxRules)
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateTaxSettings(ctx context.Context, req UpdateTaxSettingsRequest) error {
	taxExemptionsJSON, _ := json.Marshal(req.TaxExemptions)
	customTaxRulesJSON, _ := json.Marshal(req.CustomTaxRules)

	settings := []*entities.Setting{
		{Key: "tax_enabled", Value: fmt.Sprintf("%t", req.TaxEnabled), Type: "boolean", Category: "tax"},
		{Key: "tax_included", Value: fmt.Sprintf("%t", req.TaxIncluded), Type: "boolean", Category: "tax"},
		{Key: "default_tax_rate", Value: fmt.Sprintf("%.2f", req.DefaultTaxRate), Type: "number", Category: "tax"},
		{Key: "tax_calculation", Value: req.TaxCalculation, Type: "string", Category: "tax"},
		{Key: "tax_exemptions", Value: string(taxExemptionsJSON), Type: "json", Category: "tax"},
		{Key: "digital_tax_rate", Value: fmt.Sprintf("%.2f", req.DigitalTaxRate), Type: "number", Category: "tax"},
		{Key: "shipping_taxable", Value: fmt.Sprintf("%t", req.ShippingTaxable), Type: "boolean", Category: "tax"},
		{Key: "tax_display_format", Value: req.TaxDisplayFormat, Type: "string", Category: "tax"},
		{Key: "tax_rounding", Value: req.TaxRounding, Type: "string", Category: "tax"},
		{Key: "custom_tax_rules", Value: string(customTaxRulesJSON), Type: "json", Category: "tax"},
	}

	err := uc.settingsRepo.SetSettings(ctx, settings)
	if err != nil {
		return err
	}

	// Handle tax rates
	for _, rateReq := range req.TaxRates {
		rate := &entities.TaxRate{
			ID:          rateReq.ID,
			Name:        rateReq.Name,
			Rate:        rateReq.Rate,
			Country:     rateReq.Country,
			State:       rateReq.State,
			City:        rateReq.City,
			ZipCode:     rateReq.ZipCode,
			TaxClass:    rateReq.TaxClass,
			Priority:    rateReq.Priority,
			Compound:    rateReq.Compound,
			IsActive:    rateReq.IsActive,
			Description: rateReq.Description,
		}

		if rateReq.ID == uuid.Nil {
			err = uc.settingsRepo.CreateTaxRate(ctx, rate)
		} else {
			err = uc.settingsRepo.UpdateTaxRate(ctx, rate)
		}
		if err != nil {
			return err
		}
	}

	// Handle tax classes
	for _, classReq := range req.TaxClasses {
		class := &entities.TaxClass{
			ID:          classReq.ID,
			Name:        classReq.Name,
			Description: classReq.Description,
			IsActive:    classReq.IsActive,
		}

		if classReq.ID == uuid.Nil {
			err = uc.settingsRepo.CreateTaxClass(ctx, class)
		} else {
			err = uc.settingsRepo.UpdateTaxClass(ctx, class)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Shipping Settings
func (uc *settingsUseCaseImpl) GetShippingSettings(ctx context.Context) (*ShippingSettingsResponse, error) {
	_, err := uc.settingsRepo.GetSettingsByCategory(ctx, "shipping")
	if err != nil {
		return nil, err
	}

	response := &ShippingSettingsResponse{
		ShippingEnabled:     true,
		FreeShippingEnabled: false,
		FreeShippingAmount:  0.0,
		DefaultShippingCost: 10.0,
		ShippingCalculation: "flat_rate",
		ShippingZones:       []ShippingZoneResponse{},
		ShippingMethods:     []SettingsShippingMethodResponse{},
		ShippingClasses:     []ShippingClassResponse{},
		PackagingOptions:    []PackagingOptionResponse{},
		DimensionUnit:       "cm",
		WeightUnit:          "kg",
		UpdatedAt:           time.Now(),
	}

	// Get shipping data
	zones, _ := uc.settingsRepo.GetShippingZones(ctx)
	for _, zone := range zones {
		var countries, states, zipCodes []string
		json.Unmarshal([]byte(zone.Countries), &countries)
		json.Unmarshal([]byte(zone.States), &states)
		json.Unmarshal([]byte(zone.ZipCodes), &zipCodes)

		response.ShippingZones = append(response.ShippingZones, ShippingZoneResponse{
			ID:          zone.ID,
			Name:        zone.Name,
			Countries:   countries,
			States:      states,
			ZipCodes:    zipCodes,
			IsActive:    zone.IsActive,
			Priority:    zone.SortOrder, // Use SortOrder instead of Priority
			Description: zone.Description,
			CreatedAt:   zone.CreatedAt,
			UpdatedAt:   zone.UpdatedAt,
		})
	}

	methods, _ := uc.settingsRepo.GetShippingMethods(ctx, nil)
	for _, method := range methods {
		response.ShippingMethods = append(response.ShippingMethods, SettingsShippingMethodResponse{
			ID:          method.ID,
			Name:        method.Name,
			Type:        string(method.Type), // Convert ShippingMethodType to string
			Cost:        method.BaseCost,     // Use BaseCost instead of Cost
			MinAmount:   0.0,                 // Default value since not in entity
			MaxAmount:   0.0,                 // Default value since not in entity
			ZoneID:      uuid.Nil,            // Default value since not in entity
			IsActive:    method.IsActive,
			Settings:    make(map[string]interface{}), // Empty map since no Settings field
			Description: method.Description,
			CreatedAt:   method.CreatedAt,
			UpdatedAt:   method.UpdatedAt,
		})
	}

	classes, _ := uc.settingsRepo.GetShippingClasses(ctx)
	for _, class := range classes {
		response.ShippingClasses = append(response.ShippingClasses, ShippingClassResponse{
			ID:          class.ID,
			Name:        class.Name,
			Slug:        class.Slug,
			Description: class.Description,
			Cost:        class.Cost,
			IsActive:    class.IsActive,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		})
	}

	options, _ := uc.settingsRepo.GetPackagingOptions(ctx)
	for _, option := range options {
		response.PackagingOptions = append(response.PackagingOptions, PackagingOptionResponse{
			ID:          option.ID,
			Name:        option.Name,
			Length:      option.Length,
			Width:       option.Width,
			Height:      option.Height,
			Weight:      option.Weight,
			Cost:        option.Cost,
			IsDefault:   option.IsDefault,
			IsActive:    option.IsActive,
			Description: option.Description,
			CreatedAt:   option.CreatedAt,
			UpdatedAt:   option.UpdatedAt,
		})
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateShippingSettings(ctx context.Context, req UpdateShippingSettingsRequest) error {
	settings := []*entities.Setting{
		{Key: "shipping_enabled", Value: fmt.Sprintf("%t", req.ShippingEnabled), Type: "boolean", Category: "shipping"},
		{Key: "free_shipping_enabled", Value: fmt.Sprintf("%t", req.FreeShippingEnabled), Type: "boolean", Category: "shipping"},
		{Key: "free_shipping_amount", Value: fmt.Sprintf("%.2f", req.FreeShippingAmount), Type: "number", Category: "shipping"},
		{Key: "default_shipping_cost", Value: fmt.Sprintf("%.2f", req.DefaultShippingCost), Type: "number", Category: "shipping"},
		{Key: "shipping_calculation", Value: req.ShippingCalculation, Type: "string", Category: "shipping"},
		{Key: "dimension_unit", Value: req.DimensionUnit, Type: "string", Category: "shipping"},
		{Key: "weight_unit", Value: req.WeightUnit, Type: "string", Category: "shipping"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// SEO Settings
func (uc *settingsUseCaseImpl) GetSEOSettings(ctx context.Context) (*SEOSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "seo")
	if err != nil {
		return nil, err
	}

	response := &SEOSettingsResponse{
		SiteTitle:        "My Store",
		SiteDescription:  "Welcome to our online store",
		SiteKeywords:     []string{},
		MetaRobots:       "index,follow",
		CanonicalURL:     "",
		OGTitle:          "",
		OGDescription:    "",
		OGImage:          "",
		TwitterCard:      "summary",
		TwitterSite:      "",
		TwitterCreator:   "",
		GoogleAnalytics:  "",
		GoogleTagManager: "",
		FacebookPixel:    "",
		CustomMeta:       make(map[string]string),
		Sitemap: SitemapSettings{
			Enabled:           true,
			IncludeProducts:   true,
			IncludeCategories: true,
			IncludePages:      true,
			ExcludeURLs:       []string{},
			ChangeFreq:        "weekly",
			Priority:          0.8,
		},
		Robots: RobotsSettings{
			Content:     "",
			Disallow:    []string{},
			Allow:       []string{},
			Crawldelay:  0,
			SitemapURLs: []string{},
		},
		UpdatedAt: time.Now(),
	}

	// Map settings to response (simplified for brevity)
	for _, setting := range settings {
		switch setting.Key {
		case "site_title":
			response.SiteTitle = setting.Value
		case "site_description":
			response.SiteDescription = setting.Value
		case "site_keywords":
			json.Unmarshal([]byte(setting.Value), &response.SiteKeywords)
		case "meta_robots":
			response.MetaRobots = setting.Value
		case "google_analytics":
			response.GoogleAnalytics = setting.Value
		case "facebook_pixel":
			response.FacebookPixel = setting.Value
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateSEOSettings(ctx context.Context, req UpdateSEOSettingsRequest) error {
	keywordsJSON, _ := json.Marshal(req.SiteKeywords)
	customMetaJSON, _ := json.Marshal(req.CustomMeta)
	sitemapJSON, _ := json.Marshal(req.Sitemap)
	robotsJSON, _ := json.Marshal(req.Robots)

	settings := []*entities.Setting{
		{Key: "site_title", Value: req.SiteTitle, Type: "string", Category: "seo"},
		{Key: "site_description", Value: req.SiteDescription, Type: "string", Category: "seo"},
		{Key: "site_keywords", Value: string(keywordsJSON), Type: "json", Category: "seo"},
		{Key: "meta_robots", Value: req.MetaRobots, Type: "string", Category: "seo"},
		{Key: "canonical_url", Value: req.CanonicalURL, Type: "string", Category: "seo"},
		{Key: "og_title", Value: req.OGTitle, Type: "string", Category: "seo"},
		{Key: "og_description", Value: req.OGDescription, Type: "string", Category: "seo"},
		{Key: "og_image", Value: req.OGImage, Type: "string", Category: "seo"},
		{Key: "twitter_card", Value: req.TwitterCard, Type: "string", Category: "seo"},
		{Key: "twitter_site", Value: req.TwitterSite, Type: "string", Category: "seo"},
		{Key: "twitter_creator", Value: req.TwitterCreator, Type: "string", Category: "seo"},
		{Key: "google_analytics", Value: req.GoogleAnalytics, Type: "string", Category: "seo"},
		{Key: "google_tag_manager", Value: req.GoogleTagManager, Type: "string", Category: "seo"},
		{Key: "facebook_pixel", Value: req.FacebookPixel, Type: "string", Category: "seo"},
		{Key: "custom_meta", Value: string(customMetaJSON), Type: "json", Category: "seo"},
		{Key: "sitemap", Value: string(sitemapJSON), Type: "json", Category: "seo"},
		{Key: "robots", Value: string(robotsJSON), Type: "json", Category: "seo"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// Security Settings
func (uc *settingsUseCaseImpl) GetSecuritySettings(ctx context.Context) (*SecuritySettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "security")
	if err != nil {
		return nil, err
	}

	response := &SecuritySettingsResponse{
		TwoFactorEnabled:  false,
		SessionTimeout:    30,
		PasswordPolicy:    "medium",
		MaxLoginAttempts:  5,
		LockoutDuration:   15,
		RequireSSL:        true,
		IPWhitelist:       []string{},
		IPBlacklist:       []string{},
		EnableCaptcha:     false,
		CaptchaProvider:   "recaptcha",
		EnableAuditLog:    true,
		DataRetentionDays: 365,
		EncryptionEnabled: true,
		BackupEncryption:  true,
		UpdatedAt:         time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "two_factor_enabled":
			response.TwoFactorEnabled = setting.Value == "true"
		case "session_timeout":
			fmt.Sscanf(setting.Value, "%d", &response.SessionTimeout)
		case "password_policy":
			response.PasswordPolicy = setting.Value
		case "max_login_attempts":
			fmt.Sscanf(setting.Value, "%d", &response.MaxLoginAttempts)
		case "lockout_duration":
			fmt.Sscanf(setting.Value, "%d", &response.LockoutDuration)
		case "require_ssl":
			response.RequireSSL = setting.Value == "true"
		case "ip_whitelist":
			json.Unmarshal([]byte(setting.Value), &response.IPWhitelist)
		case "ip_blacklist":
			json.Unmarshal([]byte(setting.Value), &response.IPBlacklist)
		case "enable_captcha":
			response.EnableCaptcha = setting.Value == "true"
		case "captcha_provider":
			response.CaptchaProvider = setting.Value
		case "enable_audit_log":
			response.EnableAuditLog = setting.Value == "true"
		case "data_retention_days":
			fmt.Sscanf(setting.Value, "%d", &response.DataRetentionDays)
		case "encryption_enabled":
			response.EncryptionEnabled = setting.Value == "true"
		case "backup_encryption":
			response.BackupEncryption = setting.Value == "true"
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateSecuritySettings(ctx context.Context, req UpdateSecuritySettingsRequest) error {
	ipWhitelistJSON, _ := json.Marshal(req.IPWhitelist)
	ipBlacklistJSON, _ := json.Marshal(req.IPBlacklist)

	settings := []*entities.Setting{
		{Key: "two_factor_enabled", Value: fmt.Sprintf("%t", req.TwoFactorEnabled), Type: "boolean", Category: "security"},
		{Key: "session_timeout", Value: fmt.Sprintf("%d", req.SessionTimeout), Type: "number", Category: "security"},
		{Key: "password_policy", Value: req.PasswordPolicy, Type: "string", Category: "security"},
		{Key: "max_login_attempts", Value: fmt.Sprintf("%d", req.MaxLoginAttempts), Type: "number", Category: "security"},
		{Key: "lockout_duration", Value: fmt.Sprintf("%d", req.LockoutDuration), Type: "number", Category: "security"},
		{Key: "require_ssl", Value: fmt.Sprintf("%t", req.RequireSSL), Type: "boolean", Category: "security"},
		{Key: "ip_whitelist", Value: string(ipWhitelistJSON), Type: "json", Category: "security"},
		{Key: "ip_blacklist", Value: string(ipBlacklistJSON), Type: "json", Category: "security"},
		{Key: "enable_captcha", Value: fmt.Sprintf("%t", req.EnableCaptcha), Type: "boolean", Category: "security"},
		{Key: "captcha_provider", Value: req.CaptchaProvider, Type: "string", Category: "security"},
		{Key: "enable_audit_log", Value: fmt.Sprintf("%t", req.EnableAuditLog), Type: "boolean", Category: "security"},
		{Key: "data_retention_days", Value: fmt.Sprintf("%d", req.DataRetentionDays), Type: "number", Category: "security"},
		{Key: "encryption_enabled", Value: fmt.Sprintf("%t", req.EncryptionEnabled), Type: "boolean", Category: "security"},
		{Key: "backup_encryption", Value: fmt.Sprintf("%t", req.BackupEncryption), Type: "boolean", Category: "security"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// Notification Settings
func (uc *settingsUseCaseImpl) GetNotificationSettings(ctx context.Context) (*NotificationSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "notifications")
	if err != nil {
		return nil, err
	}

	response := &NotificationSettingsResponse{
		EmailNotifications:  true,
		SMSNotifications:    false,
		PushNotifications:   false,
		OrderNotifications:  true,
		StockNotifications:  true,
		ReviewNotifications: true,
		SecurityAlerts:      true,
		MarketingEmails:     false,
		NewsletterEnabled:   true,
		UpdatedAt:           time.Now(),
	}

	// Map settings to response
	for _, setting := range settings {
		switch setting.Key {
		case "email_notifications":
			response.EmailNotifications = setting.Value == "true"
		case "sms_notifications":
			response.SMSNotifications = setting.Value == "true"
		case "push_notifications":
			response.PushNotifications = setting.Value == "true"
		case "order_notifications":
			response.OrderNotifications = setting.Value == "true"
		case "stock_notifications":
			response.StockNotifications = setting.Value == "true"
		case "review_notifications":
			response.ReviewNotifications = setting.Value == "true"
		case "security_alerts":
			response.SecurityAlerts = setting.Value == "true"
		case "marketing_emails":
			response.MarketingEmails = setting.Value == "true"
		case "newsletter_enabled":
			response.NewsletterEnabled = setting.Value == "true"
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateNotificationSettings(ctx context.Context, req UpdateNotificationSettingsRequest) error {
	settings := []*entities.Setting{
		{Key: "email_notifications", Value: fmt.Sprintf("%t", req.EmailNotifications), Type: "boolean", Category: "notifications"},
		{Key: "sms_notifications", Value: fmt.Sprintf("%t", req.SMSNotifications), Type: "boolean", Category: "notifications"},
		{Key: "push_notifications", Value: fmt.Sprintf("%t", req.PushNotifications), Type: "boolean", Category: "notifications"},
		{Key: "order_notifications", Value: fmt.Sprintf("%t", req.OrderNotifications), Type: "boolean", Category: "notifications"},
		{Key: "stock_notifications", Value: fmt.Sprintf("%t", req.StockNotifications), Type: "boolean", Category: "notifications"},
		{Key: "review_notifications", Value: fmt.Sprintf("%t", req.ReviewNotifications), Type: "boolean", Category: "notifications"},
		{Key: "security_alerts", Value: fmt.Sprintf("%t", req.SecurityAlerts), Type: "boolean", Category: "notifications"},
		{Key: "marketing_emails", Value: fmt.Sprintf("%t", req.MarketingEmails), Type: "boolean", Category: "notifications"},
		{Key: "newsletter_enabled", Value: fmt.Sprintf("%t", req.NewsletterEnabled), Type: "boolean", Category: "notifications"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}

// Integration Settings
func (uc *settingsUseCaseImpl) GetIntegrationSettings(ctx context.Context) (*IntegrationSettingsResponse, error) {
	settings, err := uc.settingsRepo.GetSettingsByCategory(ctx, "integrations")
	if err != nil {
		return nil, err
	}

	response := &IntegrationSettingsResponse{
		GoogleAnalytics: GoogleAnalyticsSettings{
			Enabled:      false,
			TrackingID:   "",
			Enhanced:     false,
			Demographics: false,
			Advertising:  false,
			Ecommerce:    false,
		},
		FacebookPixel: FacebookPixelSettings{
			Enabled:  false,
			PixelID:  "",
			Advanced: false,
		},
		MailChimp: MailChimpSettings{
			Enabled:     false,
			APIKey:      "",
			ListID:      "",
			DoubleOptin: true,
		},
		Zapier: ZapierSettings{
			Enabled: false,
			APIKey:  "",
		},
		Webhooks:  []WebhookSettings{},
		UpdatedAt: time.Now(),
	}

	// Map settings to response (simplified)
	for _, setting := range settings {
		switch setting.Key {
		case "google_analytics":
			json.Unmarshal([]byte(setting.Value), &response.GoogleAnalytics)
		case "facebook_pixel":
			json.Unmarshal([]byte(setting.Value), &response.FacebookPixel)
		case "mailchimp":
			json.Unmarshal([]byte(setting.Value), &response.MailChimp)
		case "zapier":
			json.Unmarshal([]byte(setting.Value), &response.Zapier)
		case "webhooks":
			json.Unmarshal([]byte(setting.Value), &response.Webhooks)
		}
		if setting.UpdatedAt.After(response.UpdatedAt) {
			response.UpdatedAt = setting.UpdatedAt
		}
	}

	return response, nil
}

func (uc *settingsUseCaseImpl) UpdateIntegrationSettings(ctx context.Context, req UpdateIntegrationSettingsRequest) error {
	googleAnalyticsJSON, _ := json.Marshal(req.GoogleAnalytics)
	facebookPixelJSON, _ := json.Marshal(req.FacebookPixel)
	mailChimpJSON, _ := json.Marshal(req.MailChimp)
	zapierJSON, _ := json.Marshal(req.Zapier)
	webhooksJSON, _ := json.Marshal(req.Webhooks)

	settings := []*entities.Setting{
		{Key: "google_analytics", Value: string(googleAnalyticsJSON), Type: "json", Category: "integrations"},
		{Key: "facebook_pixel", Value: string(facebookPixelJSON), Type: "json", Category: "integrations"},
		{Key: "mailchimp", Value: string(mailChimpJSON), Type: "json", Category: "integrations"},
		{Key: "zapier", Value: string(zapierJSON), Type: "json", Category: "integrations"},
		{Key: "webhooks", Value: string(webhooksJSON), Type: "json", Category: "integrations"},
	}

	return uc.settingsRepo.SetSettings(ctx, settings)
}
