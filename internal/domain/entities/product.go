package entities

import (
	"fmt"
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

// ProductVisibility represents the visibility of a product
type ProductVisibility string

const (
	ProductVisibilityVisible ProductVisibility = "visible"
	ProductVisibilityHidden  ProductVisibility = "hidden"
	ProductVisibilityPrivate ProductVisibility = "private"
)

// ProductType represents the type of a product
type ProductType string

const (
	ProductTypeSimple   ProductType = "simple"
	ProductTypeVariable ProductType = "variable"
	ProductTypeGrouped  ProductType = "grouped"
	ProductTypeExternal ProductType = "external"
)

// StockStatus represents the stock status of a product
type StockStatus string

const (
	StockStatusInStock     StockStatus = "in_stock"
	StockStatusOutOfStock  StockStatus = "out_of_stock"
	StockStatusOnBackorder StockStatus = "on_backorder"
	StockStatusLowStock    StockStatus = "low_stock"
)

// Product represents a product in the system
type Product struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string    `json:"name" gorm:"not null;index" validate:"required"`
	Description      string    `json:"description" gorm:"type:text"`
	ShortDescription string    `json:"short_description" gorm:"type:text"`
	SKU              string    `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`

	// SEO and Metadata
	Slug            string            `json:"slug" gorm:"uniqueIndex" validate:"required"`
	MetaTitle       string            `json:"meta_title"`
	MetaDescription string            `json:"meta_description" gorm:"type:text"`
	Keywords        string            `json:"keywords"`
	Featured        bool              `json:"featured" gorm:"default:false"`
	Visibility      ProductVisibility `json:"visibility" gorm:"default:'visible'" validate:"required"`

	// Pricing
	Price        float64  `json:"price" gorm:"not null" validate:"required,gt=0"`
	ComparePrice *float64 `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64 `json:"cost_price" validate:"omitempty,gt=0"`

	// Sale Pricing
	SalePrice     *float64   `json:"sale_price" validate:"omitempty,gt=0"`
	SaleStartDate *time.Time `json:"sale_start_date"`
	SaleEndDate   *time.Time `json:"sale_end_date"`

	// Inventory
	Stock             int         `json:"stock" gorm:"default:0" validate:"min=0"`
	LowStockThreshold int         `json:"low_stock_threshold" gorm:"default:5"`
	TrackQuantity     bool        `json:"track_quantity" gorm:"default:true"`
	AllowBackorder    bool        `json:"allow_backorder" gorm:"default:false"`
	StockStatus       StockStatus `json:"stock_status" gorm:"default:'in_stock'"`

	// Physical Properties
	Weight     *float64    `json:"weight" validate:"omitempty,gt=0"`
	Dimensions *Dimensions `json:"dimensions" gorm:"embedded"`

	// Shipping and Tax
	RequiresShipping bool   `json:"requires_shipping" gorm:"default:true"`
	ShippingClass    string `json:"shipping_class"`
	TaxClass         string `json:"tax_class" gorm:"default:'standard'"`
	CountryOfOrigin  string `json:"country_of_origin"`

	// Categorization - CategoryID removed, use ProductCategory many-to-many as single source of truth
	BrandID    *uuid.UUID `json:"brand_id" gorm:"type:uuid;index"`

	// Status and Type
	Status      ProductStatus `json:"status" gorm:"default:'draft'" validate:"required"`
	ProductType ProductType   `json:"product_type" gorm:"default:'simple'" validate:"required"`
	IsDigital   bool          `json:"is_digital" gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships - Category relationship removed, use ProductCategory many-to-many
	Brand           *Brand                  `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
	Images          []ProductImage          `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	Tags            []ProductTag            `json:"tags,omitempty" gorm:"many2many:product_tag_associations;"`
	Reviews         []Review                `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
	Suppliers       []Supplier              `json:"suppliers,omitempty" gorm:"many2many:supplier_products;"`
	Variants        []ProductVariant        `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Attributes      []ProductAttributeValue `json:"attributes,omitempty" gorm:"foreignKey:ProductID"`
	RelatedProducts []Product               `json:"related_products,omitempty" gorm:"many2many:product_relations;joinForeignKey:ProductID;joinReferences:RelatedProductID"`
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

// Brand represents a product brand
type Brand struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex" validate:"required"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Description string    `json:"description" gorm:"type:text"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website" validate:"omitempty,url"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Computed fields (not stored in database)
	ProductCount int `json:"product_count" gorm:"-"`

	// Relationships
	Products []Product `json:"products,omitempty" gorm:"foreignKey:BrandID"`
}

// TableName returns the table name for Brand entity
func (Brand) TableName() string {
	return "brands"
}

// ProductVariant represents a product variant (e.g., different sizes, colors)
type ProductVariant struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID    uuid.UUID   `json:"product_id" gorm:"type:uuid;not null;index"`
	SKU          string      `json:"sku" gorm:"uniqueIndex;not null" validate:"required"`
	Price        float64     `json:"price" gorm:"not null" validate:"required,gt=0"`
	ComparePrice *float64    `json:"compare_price" validate:"omitempty,gt=0"`
	CostPrice    *float64    `json:"cost_price" validate:"omitempty,gt=0"`
	Stock        int         `json:"stock" gorm:"default:0" validate:"min=0"`
	Weight       *float64    `json:"weight" validate:"omitempty,gt=0"`
	Dimensions   *Dimensions `json:"dimensions" gorm:"embedded"`
	Image        string      `json:"image"`
	Position     int         `json:"position" gorm:"default:0"`
	IsActive     bool        `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Product    Product                   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Attributes []ProductVariantAttribute `json:"attributes,omitempty" gorm:"foreignKey:VariantID"`
}

// TableName returns the table name for ProductVariant entity
func (ProductVariant) TableName() string {
	return "product_variants"
}

// ProductComparison represents a product comparison session
type ProductComparison struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    *uuid.UUID `json:"user_id" gorm:"type:uuid;index"` // Optional for guest users
	SessionID string     `json:"session_id" gorm:"index"`        // For guest users
	Name      string     `json:"name"`                           // Optional comparison name
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Items []ProductComparisonItem `json:"items,omitempty" gorm:"foreignKey:ComparisonID"`
}

// TableName returns the table name for ProductComparison entity
func (ProductComparison) TableName() string {
	return "product_comparisons"
}

// ProductComparisonItem represents an item in a product comparison
type ProductComparisonItem struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ComparisonID uuid.UUID `json:"comparison_id" gorm:"type:uuid;not null;index"`
	ProductID    uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Position     int       `json:"position" gorm:"default:0"` // Order in comparison
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Comparison ProductComparison `json:"comparison,omitempty" gorm:"foreignKey:ComparisonID"`
	Product    Product           `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName returns the table name for ProductComparisonItem entity
func (ProductComparisonItem) TableName() string {
	return "product_comparison_items"
}

// FilterSet represents a saved filter configuration
type FilterSet struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	SessionID   string     `json:"session_id" gorm:"index"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Filters     string     `json:"filters" gorm:"type:jsonb"` // JSON encoded filter parameters
	IsPublic    bool       `json:"is_public" gorm:"default:false"`
	UsageCount  int        `json:"usage_count" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for FilterSet entity
func (FilterSet) TableName() string {
	return "filter_sets"
}

// FilterUsage tracks filter usage analytics
type FilterUsage struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	SessionID   string     `json:"session_id" gorm:"index"`
	FilterType  string     `json:"filter_type" gorm:"index"` // price, category, brand, attribute, etc.
	FilterKey   string     `json:"filter_key" gorm:"index"`  // specific filter identifier
	FilterValue string     `json:"filter_value"`             // filter value used
	ResultCount int        `json:"result_count"`             // number of results returned
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for FilterUsage entity
func (FilterUsage) TableName() string {
	return "filter_usage"
}

// ProductFilterOption represents available filter options for products
type ProductFilterOption struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID  *uuid.UUID `json:"category_id" gorm:"type:uuid;index"`
	FilterType  string     `json:"filter_type" gorm:"index"` // price, brand, attribute, rating, etc.
	FilterKey   string     `json:"filter_key" gorm:"index"`  // attribute_id for attributes, 'brand' for brands
	FilterValue string     `json:"filter_value"`             // the actual filter value
	Label       string     `json:"label"`                    // display label
	Count       int        `json:"count" gorm:"default:0"`   // number of products with this filter
	Position    int        `json:"position" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ProductFilterOption entity
func (ProductFilterOption) TableName() string {
	return "product_filter_options"
}

// ProductAttribute represents a product attribute (e.g., Color, Size)
type ProductAttribute struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex" validate:"required"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Type        string    `json:"type" gorm:"default:'text'"` // text, select, color, image
	Description string    `json:"description" gorm:"type:text"`
	Position    int       `json:"position" gorm:"default:0"`
	IsRequired  bool      `json:"is_required" gorm:"default:false"`
	IsVisible   bool      `json:"is_visible" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Terms []ProductAttributeTerm `json:"terms,omitempty" gorm:"foreignKey:AttributeID"`
}

// TableName returns the table name for ProductAttribute entity
func (ProductAttribute) TableName() string {
	return "product_attributes"
}

// ProductAttributeTerm represents a term/value for an attribute (e.g., Red, Blue for Color)
type ProductAttributeTerm struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AttributeID uuid.UUID `json:"attribute_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Slug        string    `json:"slug" gorm:"not null" validate:"required"`
	Value       string    `json:"value"`
	Color       string    `json:"color"`
	Image       string    `json:"image"`
	Position    int       `json:"position" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Attribute ProductAttribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
}

// TableName returns the table name for ProductAttributeTerm entity
func (ProductAttributeTerm) TableName() string {
	return "product_attribute_terms"
}

// ProductAttributeValue represents the value of an attribute for a specific product
type ProductAttributeValue struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID   uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;index"`
	AttributeID uuid.UUID  `json:"attribute_id" gorm:"type:uuid;not null;index"`
	TermID      *uuid.UUID `json:"term_id" gorm:"type:uuid;index"`
	Value       string     `json:"value"`
	Position    int        `json:"position" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Product   Product               `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Attribute ProductAttribute      `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
	Term      *ProductAttributeTerm `json:"term,omitempty" gorm:"foreignKey:TermID"`
}

// TableName returns the table name for ProductAttributeValue entity
func (ProductAttributeValue) TableName() string {
	return "product_attribute_values"
}

// ProductVariantAttribute represents the attribute values for a specific variant
type ProductVariantAttribute struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VariantID   uuid.UUID `json:"variant_id" gorm:"type:uuid;not null;index"`
	AttributeID uuid.UUID `json:"attribute_id" gorm:"type:uuid;not null;index"`
	TermID      uuid.UUID `json:"term_id" gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Variant   ProductVariant       `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	Attribute ProductAttribute     `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
	Term      ProductAttributeTerm `json:"term,omitempty" gorm:"foreignKey:TermID"`
}

// TableName returns the table name for ProductVariantAttribute entity
func (ProductVariantAttribute) TableName() string {
	return "product_variant_attributes"
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.Stock > 0
}

// HasDiscount checks if the product has any type of discount
func (p *Product) HasDiscount() bool {
	return p.IsOnSale() || p.HasCompareDiscount()
}

// HasCompareDiscount checks if the product has a compare price discount
func (p *Product) HasCompareDiscount() bool {
	return p.ComparePrice != nil && *p.ComparePrice > p.Price
}

// GetDiscountPercentage calculates the effective discount percentage
// Priority: Sale discount > Compare price discount
func (p *Product) GetDiscountPercentage() float64 {
	// Priority 1: Sale discount (if currently on sale)
	if p.IsOnSale() {
		return p.GetSaleDiscountPercentage()
	}

	// Priority 2: Compare price discount
	if p.HasCompareDiscount() {
		return p.GetCompareDiscountPercentage()
	}

	return 0
}

// GetCompareDiscountPercentage calculates discount percentage from compare price
func (p *Product) GetCompareDiscountPercentage() float64 {
	if !p.HasCompareDiscount() {
		return 0
	}
	return ((*p.ComparePrice - p.Price) / *p.ComparePrice) * 100
}

// CanReduceStock checks if stock can be reduced by the given quantity
func (p *Product) CanReduceStock(quantity int) bool {
	if !p.TrackQuantity {
		return true
	}
	if p.AllowBackorder {
		return true
	}
	return p.Stock >= quantity
}

// GetCurrentPrice returns the current effective price (sale price if active, otherwise regular price)
func (p *Product) GetCurrentPrice() float64 {
	if p.IsOnSale() {
		return *p.SalePrice
	}
	return p.Price
}

// IsOnSale checks if the product is currently on sale
func (p *Product) IsOnSale() bool {
	if p.SalePrice == nil || *p.SalePrice <= 0 {
		return false
	}

	now := time.Now()

	// Check sale start date
	if p.SaleStartDate != nil && now.Before(*p.SaleStartDate) {
		return false
	}

	// Check sale end date
	if p.SaleEndDate != nil && now.After(*p.SaleEndDate) {
		return false
	}

	return *p.SalePrice < p.Price
}

// GetSaleDiscountPercentage calculates the sale discount percentage
func (p *Product) GetSaleDiscountPercentage() float64 {
	if !p.IsOnSale() {
		return 0
	}
	return ((p.Price - *p.SalePrice) / p.Price) * 100
}

// GetOriginalPrice returns the price to show as "crossed out" when there's a discount
// Returns the higher price that should be displayed with strikethrough
func (p *Product) GetOriginalPrice() *float64 {
	// Priority 1: If on sale, show regular price as original
	if p.IsOnSale() {
		return &p.Price
	}

	// Priority 2: If has compare price discount, show compare price as original
	if p.HasCompareDiscount() {
		return p.ComparePrice
	}

	// No discount, no original price to show
	return nil
}

// GetDisplayPrice returns the current selling price (what customer pays)
func (p *Product) GetDisplayPrice() float64 {
	return p.GetCurrentPrice()
}

// IsLowStock checks if the product is low on stock
func (p *Product) IsLowStock() bool {
	if !p.TrackQuantity {
		return false
	}
	return p.Stock <= p.LowStockThreshold && p.Stock > 0
}

// UpdateStockStatus updates the stock status based on current stock level
func (p *Product) UpdateStockStatus() {
	if !p.TrackQuantity {
		p.StockStatus = StockStatusInStock
		return
	}

	if p.Stock <= 0 {
		if p.AllowBackorder {
			p.StockStatus = StockStatusOnBackorder
		} else {
			p.StockStatus = StockStatusOutOfStock
		}
	} else if p.IsLowStock() {
		p.StockStatus = StockStatusLowStock
	} else {
		p.StockStatus = StockStatusInStock
	}
}

// IsVisible checks if the product is visible to customers
func (p *Product) IsVisible() bool {
	return p.Status == ProductStatusActive && p.Visibility == ProductVisibilityVisible
}

// HasVariants checks if the product has variants
func (p *Product) HasVariants() bool {
	return p.ProductType == ProductTypeVariable && len(p.Variants) > 0
}

// GetMainImage returns the main product image (first image or empty string)
func (p *Product) GetMainImage() string {
	if len(p.Images) > 0 {
		return p.Images[0].URL
	}
	return ""
}

// ReduceStock reduces the product stock
func (p *Product) ReduceStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInput
	}

	if !p.CanReduceStock(quantity) {
		return ErrInsufficientStock
	}

	p.Stock -= quantity

	// Update stock status based on remaining stock
	p.UpdateStockStatus()

	return nil
}

// IncreaseStock increases the product stock
func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInput
	}

	p.Stock += quantity

	// Update stock status based on new stock level
	p.UpdateStockStatus()

	return nil
}

// Validate validates all product business rules
func (p *Product) Validate() error {
	// Validate required fields
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.SKU == "" {
		return fmt.Errorf("product SKU is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	// Validate compare price
	if p.ComparePrice != nil && *p.ComparePrice <= p.Price {
		return fmt.Errorf("compare price must be greater than regular price")
	}

	// Validate cost price
	if p.CostPrice != nil && *p.CostPrice < 0 {
		return fmt.Errorf("cost price cannot be negative")
	}

	// Validate stock
	if p.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	if p.LowStockThreshold < 0 {
		return fmt.Errorf("low stock threshold cannot be negative")
	}

	// Validate weight
	if p.Weight != nil && *p.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}

	// Validate dimensions
	if p.Dimensions != nil {
		if p.Dimensions.Length <= 0 || p.Dimensions.Width <= 0 || p.Dimensions.Height <= 0 {
			return fmt.Errorf("all dimensions must be greater than 0")
		}
	}

	// Validate sale pricing
	if err := p.ValidateSalePricing(); err != nil {
		return err
	}

	return nil
}

// ValidateSalePricing validates sale pricing business rules
func (p *Product) ValidateSalePricing() error {
	// If sale price is set, validate business rules
	if p.SalePrice != nil {
		// Sale price must be less than regular price
		if *p.SalePrice >= p.Price {
			return fmt.Errorf("sale price must be less than regular price")
		}

		// Sale price must be positive
		if *p.SalePrice <= 0 {
			return fmt.Errorf("sale price must be greater than 0")
		}

		// If both start and end dates are set, start must be before end
		if p.SaleStartDate != nil && p.SaleEndDate != nil {
			if p.SaleStartDate.After(*p.SaleEndDate) {
				return fmt.Errorf("sale start date must be before sale end date")
			}
		}
	}

	return nil
}
