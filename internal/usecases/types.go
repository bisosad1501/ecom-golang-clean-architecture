package usecases

import (
	"time"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/google/uuid"
)

// Shared Response Types

// ProductResponse represents product response
type ProductResponse struct {
	ID           uuid.UUID                 `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	SKU          string                    `json:"sku"`
	Price        float64                   `json:"price"`
	ComparePrice *float64                  `json:"compare_price"`
	CostPrice    *float64                  `json:"cost_price"`
	Stock        int                       `json:"stock"`
	Weight       *float64                  `json:"weight"`
	Dimensions   *DimensionsResponse       `json:"dimensions"`
	Category     *ProductCategoryResponse  `json:"category"`
	Images       []ProductImageResponse    `json:"images"`
	Tags         []ProductTagResponse      `json:"tags"`
	Status       entities.ProductStatus    `json:"status"`
	IsDigital    bool                      `json:"is_digital"`
	IsAvailable  bool                      `json:"is_available"`
	HasDiscount  bool                      `json:"has_discount"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

type DimensionsResponse struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ProductCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
}

type ProductImageResponse struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	AltText  string    `json:"alt_text"`
	Position int       `json:"position"`
}

type ProductTagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
