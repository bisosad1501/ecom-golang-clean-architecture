package services

import (
	"context"
	"fmt"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
)

// ShippingCompatibilityService handles shipping method compatibility checks
type ShippingCompatibilityService interface {
	ValidateShippingMethodForAddress(ctx context.Context, method *entities.ShippingMethod, address *entities.Address) error
	GetCompatibleShippingMethods(ctx context.Context, address *entities.Address, weight float64, orderValue float64) ([]*entities.ShippingMethod, error)
	ValidateShippingConstraints(ctx context.Context, method *entities.ShippingMethod, weight float64, dimensions *entities.Dimensions) error
}

type shippingCompatibilityService struct {
	// Could inject external services for address validation, etc.
}

// NewShippingCompatibilityService creates a new shipping compatibility service
func NewShippingCompatibilityService() ShippingCompatibilityService {
	return &shippingCompatibilityService{}
}

// ValidateShippingMethodForAddress validates if a shipping method is compatible with an address
func (s *shippingCompatibilityService) ValidateShippingMethodForAddress(ctx context.Context, method *entities.ShippingMethod, address *entities.Address) error {
	if method == nil {
		return fmt.Errorf("shipping method is required")
	}
	if address == nil {
		return fmt.Errorf("address is required")
	}

	// Check if method is active
	if !method.IsActive {
		return fmt.Errorf("shipping method '%s' is not active", method.Name)
	}

	// Check international shipping restrictions
	if address.IsInternational() {
		if !s.supportsInternationalShipping(method) {
			return fmt.Errorf("shipping method '%s' does not support international shipping", method.Name)
		}
	}

	// Check carrier-specific restrictions
	if err := s.validateCarrierRestrictions(method, address); err != nil {
		return fmt.Errorf("carrier restriction: %w", err)
	}

	// Check address completeness for shipping
	if !address.IsShippingAddress() {
		return fmt.Errorf("address is not configured for shipping")
	}

	return nil
}

// GetCompatibleShippingMethods returns shipping methods compatible with the given address and constraints
func (s *shippingCompatibilityService) GetCompatibleShippingMethods(ctx context.Context, address *entities.Address, weight float64, orderValue float64) ([]*entities.ShippingMethod, error) {
	// This would typically query the repository, but for now we'll return a mock implementation
	// In a real implementation, this would:
	// 1. Query all active shipping methods
	// 2. Filter by address compatibility
	// 3. Filter by weight/value constraints
	// 4. Return compatible methods

	return nil, fmt.Errorf("not implemented - would query repository for compatible methods")
}

// ValidateShippingConstraints validates shipping method constraints
func (s *shippingCompatibilityService) ValidateShippingConstraints(ctx context.Context, method *entities.ShippingMethod, weight float64, dimensions *entities.Dimensions) error {
	if method == nil {
		return fmt.Errorf("shipping method is required")
	}

	// Validate weight constraints
	if weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}

	if method.MaxWeight > 0 && weight > method.MaxWeight {
		return fmt.Errorf("weight %.2f kg exceeds maximum weight %.2f kg for method '%s'", weight, method.MaxWeight, method.Name)
	}

	// Validate dimensions if provided
	if dimensions != nil {
		if err := s.validateDimensions(method, dimensions); err != nil {
			return fmt.Errorf("dimension validation failed: %w", err)
		}
	}

	return nil
}

// supportsInternationalShipping checks if method supports international shipping
func (s *shippingCompatibilityService) supportsInternationalShipping(method *entities.ShippingMethod) bool {
	// Check by carrier
	internationalCarriers := map[string]bool{
		"DHL":     true,
		"FedEx":   true,
		"UPS":     true,
		"USPS":    true,
		"TNT":     true,
	}

	if internationalCarriers[method.Carrier] {
		return true
	}

	// Check by method type
	internationalTypes := map[string]bool{
		"express":     true,
		"premium":     true,
		"international": true,
	}

	return internationalTypes[string(method.Type)]
}

// validateCarrierRestrictions validates carrier-specific restrictions
func (s *shippingCompatibilityService) validateCarrierRestrictions(method *entities.ShippingMethod, address *entities.Address) error {
	carrier := strings.ToUpper(method.Carrier)
	country := strings.ToUpper(address.Country)

	// Example carrier restrictions
	switch carrier {
	case "VIETNAMPOST":
		// VietnamPost only ships within Vietnam
		if country != "VN" && country != "VIETNAM" {
			return fmt.Errorf("VietnamPost only ships within Vietnam")
		}
	case "VIETTELPOST":
		// ViettelPost only ships within Vietnam
		if country != "VN" && country != "VIETNAM" {
			return fmt.Errorf("ViettelPost only ships within Vietnam")
		}
	case "DHL":
		// DHL has global coverage but may have restrictions for certain countries
		restrictedCountries := []string{"NORTH KOREA", "IRAN", "SYRIA"}
		for _, restricted := range restrictedCountries {
			if country == restricted {
				return fmt.Errorf("DHL does not ship to %s", address.Country)
			}
		}
	}

	return nil
}

// validateDimensions validates package dimensions
func (s *shippingCompatibilityService) validateDimensions(method *entities.ShippingMethod, dimensions *entities.Dimensions) error {
	if dimensions.Length <= 0 || dimensions.Width <= 0 || dimensions.Height <= 0 {
		return fmt.Errorf("all dimensions must be greater than 0")
	}

	// Calculate dimensional weight (length * width * height / 5000 for cm³ to kg conversion)
	dimensionalWeight := (dimensions.Length * dimensions.Width * dimensions.Height) / 5000

	// Check if dimensional weight exceeds method limits
	if method.MaxWeight > 0 && dimensionalWeight > method.MaxWeight {
		return fmt.Errorf("dimensional weight %.2f kg exceeds maximum weight %.2f kg for method '%s'", dimensionalWeight, method.MaxWeight, method.Name)
	}

	// Method-specific dimension limits
	switch strings.ToLower(string(method.Type)) {
	case "express":
		// Express typically has smaller size limits
		maxDimension := 100.0 // cm
		if dimensions.Length > maxDimension || dimensions.Width > maxDimension || dimensions.Height > maxDimension {
			return fmt.Errorf("express shipping: no dimension can exceed %.0f cm", maxDimension)
		}
	case "economy":
		// Economy may have larger size limits but restrictions on total size
		totalSize := dimensions.Length + dimensions.Width + dimensions.Height
		if totalSize > 300 {
			return fmt.Errorf("economy shipping: total dimensions cannot exceed 300 cm")
		}
	}

	return nil
}
