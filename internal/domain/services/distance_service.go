package services

import (
	"context"
	"fmt"
	"math"
	"strings"
)

// DistanceService handles distance calculations for shipping
type DistanceService interface {
	CalculateDistance(ctx context.Context, fromLat, fromLng, toLat, toLng float64) (float64, error)
	CalculateDistanceByAddress(ctx context.Context, fromAddress, toAddress string) (float64, error)
	GetShippingZoneByDistance(ctx context.Context, distance float64) (string, error)
	ValidateShippingDistance(ctx context.Context, distance float64, methodID string) (bool, error)
	GetShippingZones() []ShippingZoneInfo
}

type distanceService struct {
	// Could integrate with external APIs like Google Maps, MapBox, etc.
	maxShippingDistance float64
}

// NewDistanceService creates a new distance service
func NewDistanceService() DistanceService {
	return &distanceService{
		maxShippingDistance: 1000.0, // 1000 km max shipping distance
	}
}

// CalculateDistance calculates distance between two coordinates using Haversine formula
func (s *distanceService) CalculateDistance(ctx context.Context, fromLat, fromLng, toLat, toLng float64) (float64, error) {
	// Validate coordinates
	if !isValidLatitude(fromLat) || !isValidLatitude(toLat) {
		return 0, fmt.Errorf("invalid latitude values")
	}
	if !isValidLongitude(fromLng) || !isValidLongitude(toLng) {
		return 0, fmt.Errorf("invalid longitude values")
	}

	// Haversine formula
	const earthRadiusKm = 6371.0

	// Convert degrees to radians
	lat1Rad := toRadians(fromLat)
	lng1Rad := toRadians(fromLng)
	lat2Rad := toRadians(toLat)
	lng2Rad := toRadians(toLng)

	// Calculate differences
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c

	return distance, nil
}

// CalculateDistanceByAddress calculates distance between two addresses
func (s *distanceService) CalculateDistanceByAddress(ctx context.Context, fromAddress, toAddress string) (float64, error) {
	if fromAddress == "" || toAddress == "" {
		return 0, fmt.Errorf("addresses cannot be empty")
	}

	// If same address, return 0
	if fromAddress == toAddress {
		return 0, nil
	}

	// For demo purposes, use realistic distance mapping based on common city pairs
	// In production, this would use a real geocoding service like Google Maps API
	distance := s.getRealisticDistanceByAddress(fromAddress, toAddress)

	if distance > s.maxShippingDistance {
		distance = s.maxShippingDistance
	}

	return distance, nil
}

// getRealisticDistanceByAddress returns realistic distances for demo purposes
func (s *distanceService) getRealisticDistanceByAddress(fromAddress, toAddress string) float64 {
	// Extract city/state info for realistic distance calculation
	fromLower := strings.ToLower(fromAddress)
	toLower := strings.ToLower(toAddress)

	// Same city/state - local delivery
	if s.isSameCity(fromLower, toLower) {
		return 5.0 + (float64(len(fromAddress)%10) * 0.5) // 5-10 km within city
	}

	// Same state - regional delivery
	if s.isSameState(fromLower, toLower) {
		return 50.0 + (float64(len(fromAddress)%50) * 2.0) // 50-150 km within state
	}

	// Cross-country major routes
	if s.isCrossCountry(fromLower, toLower) {
		return 3000.0 + (float64(len(fromAddress)%100) * 10.0) // 3000-4000 km cross-country
	}

	// Default regional distance
	return 200.0 + (float64(len(fromAddress)%30) * 5.0) // 200-350 km default
}

// isSameCity checks if addresses are in the same city
func (s *distanceService) isSameCity(addr1, addr2 string) bool {
	cities := []string{"new york", "los angeles", "chicago", "houston", "phoenix", "philadelphia", "san antonio", "san diego", "dallas", "san jose"}
	for _, city := range cities {
		if strings.Contains(addr1, city) && strings.Contains(addr2, city) {
			return true
		}
	}
	return false
}

// isSameState checks if addresses are in the same state
func (s *distanceService) isSameState(addr1, addr2 string) bool {
	states := []string{"ny", "ca", "tx", "fl", "il", "pa", "oh", "ga", "nc", "mi"}
	for _, state := range states {
		if strings.Contains(addr1, state) && strings.Contains(addr2, state) {
			return true
		}
	}
	return false
}

// isCrossCountry checks if this is a cross-country shipment
func (s *distanceService) isCrossCountry(addr1, addr2 string) bool {
	eastCoast := []string{"new york", "ny", "philadelphia", "pa", "boston", "ma", "miami", "fl"}
	westCoast := []string{"los angeles", "ca", "san francisco", "seattle", "wa", "portland", "or"}

	isAddr1East := false
	isAddr1West := false
	isAddr2East := false
	isAddr2West := false

	for _, city := range eastCoast {
		if strings.Contains(addr1, city) {
			isAddr1East = true
		}
		if strings.Contains(addr2, city) {
			isAddr2East = true
		}
	}

	for _, city := range westCoast {
		if strings.Contains(addr1, city) {
			isAddr1West = true
		}
		if strings.Contains(addr2, city) {
			isAddr2West = true
		}
	}

	return (isAddr1East && isAddr2West) || (isAddr1West && isAddr2East)
}

// GetShippingZoneByDistance determines shipping zone based on distance
func (s *distanceService) GetShippingZoneByDistance(ctx context.Context, distance float64) (string, error) {
	if distance < 0 {
		return "", fmt.Errorf("distance cannot be negative")
	}

	switch {
	case distance <= 10:
		return "local", nil
	case distance <= 50:
		return "regional", nil
	case distance <= 200:
		return "national", nil
	case distance <= 500:
		return "extended", nil
	case distance <= s.maxShippingDistance:
		return "international", nil
	default:
		return "", fmt.Errorf("distance %.2f km exceeds maximum shipping distance %.2f km", distance, s.maxShippingDistance)
	}
}

// ValidateShippingDistance validates if shipping is available for the distance
func (s *distanceService) ValidateShippingDistance(ctx context.Context, distance float64, methodID string) (bool, error) {
	if distance < 0 {
		return false, fmt.Errorf("distance cannot be negative")
	}

	// Check maximum shipping distance
	if distance > s.maxShippingDistance {
		return false, nil
	}

	// Method-specific validation
	switch methodID {
	case "express":
		return distance <= 100, nil // Express only within 100km
	case "standard":
		return distance <= 500, nil // Standard within 500km
	case "economy":
		return distance <= s.maxShippingDistance, nil // Economy for all distances
	default:
		return distance <= 200, nil // Default limit
	}
}

// Helper functions

func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func isValidLongitude(lng float64) bool {
	return lng >= -180 && lng <= 180
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// DistanceCalculationRequest represents a distance calculation request
type DistanceCalculationRequest struct {
	FromLatitude  *float64 `json:"from_latitude"`
	FromLongitude *float64 `json:"from_longitude"`
	FromAddress   string   `json:"from_address"`
	ToLatitude    *float64 `json:"to_latitude"`
	ToLongitude   *float64 `json:"to_longitude"`
	ToAddress     string   `json:"to_address"`
}

// DistanceCalculationResponse represents a distance calculation response
type DistanceCalculationResponse struct {
	Distance     float64 `json:"distance_km"`
	Zone         string  `json:"shipping_zone"`
	IsShippable  bool    `json:"is_shippable"`
	EstimatedCost float64 `json:"estimated_cost"`
}

// ShippingZoneInfo represents shipping zone information
type ShippingZoneInfo struct {
	Zone        string  `json:"zone"`
	MaxDistance float64 `json:"max_distance_km"`
	BaseCost    float64 `json:"base_cost"`
	CostPerKm   float64 `json:"cost_per_km"`
}

// GetShippingZones returns available shipping zones
func (s *distanceService) GetShippingZones() []ShippingZoneInfo {
	return []ShippingZoneInfo{
		{Zone: "local", MaxDistance: 10, BaseCost: 5.0, CostPerKm: 0.5},
		{Zone: "regional", MaxDistance: 50, BaseCost: 10.0, CostPerKm: 0.3},
		{Zone: "national", MaxDistance: 200, BaseCost: 15.0, CostPerKm: 0.2},
		{Zone: "extended", MaxDistance: 500, BaseCost: 25.0, CostPerKm: 0.15},
		{Zone: "international", MaxDistance: 1000, BaseCost: 50.0, CostPerKm: 0.1},
	}
}
