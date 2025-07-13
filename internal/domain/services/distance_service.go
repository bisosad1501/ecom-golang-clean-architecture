package services

import (
	"context"
	"fmt"
	"math"
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
	// In a real implementation, this would:
	// 1. Geocode addresses to coordinates using Google Maps API or similar
	// 2. Calculate distance using coordinates
	// For now, return a mock distance based on address similarity

	if fromAddress == "" || toAddress == "" {
		return 0, fmt.Errorf("addresses cannot be empty")
	}

	// Mock implementation - in production, integrate with geocoding service
	if fromAddress == toAddress {
		return 0, nil
	}

	// Return mock distance based on string comparison
	// This is just for demonstration - replace with real geocoding
	mockDistance := float64(len(fromAddress)+len(toAddress)) / 10.0
	if mockDistance > s.maxShippingDistance {
		mockDistance = s.maxShippingDistance
	}

	return mockDistance, nil
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
