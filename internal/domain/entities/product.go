package entities

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDraft    ProductStatus = "draft"
)

// Product represents a product in the system
type Product struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string        `json:"name" gorm:"not null;index" validate:"required"`
	Description string        `json:"description" gorm:"type:text"`
	SKU         string        `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Price       float64       `json:"price" gorm:"not null" validate:"required,gt=0"`
	ComparePrice *float64     `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice   *float64      `json:"cost_price" validate:"omitempty,gt=0"`
	Stock       int           `json:"stock" gorm:"default:0" validate:"min=0"`
	Weight      *float64      `json:"weight" validate:"omitempty,gt=0"`
	Dimensions  *Dimensions   `json:"dimensions" gorm:"embedded"`
	CategoryID  uuid.UUID     `json:"category_id" gorm:"type:uuid;index"`
	Status      ProductStatus `json:"status" gorm:"default:'draft'" validate:"required"`
	IsDigital   bool          `json:"is_digital" gorm:"default:false"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Category    Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Images      []ProductImage `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	Tags        []ProductTag  `json:"tags,omitempty" gorm:"many2many:product_tag_associations;"`
	Reviews     []Review      `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
	Suppliers   []Supplier    `json:"suppliers,omitempty" gorm:"many2many:supplier_products;"`
}

// TableName returns the table name for Product entity
func (Product) TableName() string {
	return "products"
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length" validate:"gt=0"`
	Width  float64 `json:"width" validate:"gt=0"`
	Height float64 `json:"height" validate:"gt=0"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	URL       string    `json:"url" gorm:"not null" validate:"required,url"`
	AltText   string    `json:"alt_text"`
	Position  int       `json:"position" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ProductImage entity
func (ProductImage) TableName() string {
	return "product_images"
}

// ProductTag represents a product tag
type ProductTag struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ProductTag entity
func (ProductTag) TableName() string {
	return "tags"
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.Stock > 0
}

// HasDiscount checks if the product has a discount
func (p *Product) HasDiscount() bool {
	return p.ComparePrice != nil && *p.ComparePrice > p.Price
}

// GetDiscountPercentage calculates the discount percentage
func (p *Product) GetDiscountPercentage() float64 {
	if !p.HasDiscount() {
		return 0
	}
	return ((*p.ComparePrice - p.Price) / *p.ComparePrice) * 100
}

// CanReduceStock checks if stock can be reduced by the given quantity
func (p *Product) CanReduceStock(quantity int) bool {
	return p.Stock >= quantity
}

// ReduceStock reduces the product stock
func (p *Product) ReduceStock(quantity int) error {
	if !p.CanReduceStock(quantity) {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	return nil
}

// IncreaseStock increases the product stock
func (p *Product) IncreaseStock(quantity int) {
	p.Stock += quantity
}
