package usecases

import (
	"ecom-golang-clean-architecture/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

// Shared Response Types

// ProductResponse represents product response
type ProductResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	SKU              string    `json:"sku"`

	// SEO and Metadata
	Slug            string                     `json:"slug"`
	MetaTitle       string                     `json:"meta_title"`
	MetaDescription string                     `json:"meta_description"`
	Keywords        string                     `json:"keywords"`
	Featured        bool                       `json:"featured"`
	Visibility      entities.ProductVisibility `json:"visibility"`

	// Pricing
	Price        float64  `json:"price"`
	ComparePrice *float64 `json:"compare_price"`
	CostPrice    *float64 `json:"cost_price"`

	// Sale Pricing
	SalePrice              *float64   `json:"sale_price"`
	SaleStartDate          *time.Time `json:"sale_start_date"`
	SaleEndDate            *time.Time `json:"sale_end_date"`
	CurrentPrice           float64    `json:"current_price"`
	IsOnSale               bool       `json:"is_on_sale"`
	SaleDiscountPercentage float64    `json:"sale_discount_percentage"`

	// Inventory
	Stock             int                  `json:"stock"`
	LowStockThreshold int                  `json:"low_stock_threshold"`
	TrackQuantity     bool                 `json:"track_quantity"`
	AllowBackorder    bool                 `json:"allow_backorder"`
	StockStatus       entities.StockStatus `json:"stock_status"`
	IsLowStock        bool                 `json:"is_low_stock"`

	// Physical Properties
	Weight     *float64            `json:"weight"`
	Dimensions *DimensionsResponse `json:"dimensions"`

	// Shipping and Tax
	RequiresShipping bool   `json:"requires_shipping"`
	ShippingClass    string `json:"shipping_class"`
	TaxClass         string `json:"tax_class"`
	CountryOfOrigin  string `json:"country_of_origin"`

	// Categorization
	Category *ProductCategoryResponse `json:"category"`
	Brand    *ProductBrandResponse    `json:"brand"`

	// Content
	Images     []ProductImageResponse     `json:"images"`
	Tags       []ProductTagResponse       `json:"tags"`
	Attributes []ProductAttributeResponse `json:"attributes"`
	Variants   []ProductVariantResponse   `json:"variants"`

	// Status and Type
	Status      entities.ProductStatus `json:"status"`
	ProductType entities.ProductType   `json:"product_type"`
	IsDigital   bool                   `json:"is_digital"`
	IsAvailable bool                   `json:"is_available"`
	HasDiscount bool                   `json:"has_discount"`
	HasVariants bool                   `json:"has_variants"`
	MainImage   string                 `json:"main_image"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

type ProductBrandResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	IsActive    bool      `json:"is_active"`
}

type ProductAttributeResponse struct {
	ID          uuid.UUID  `json:"id"`
	AttributeID uuid.UUID  `json:"attribute_id"`
	TermID      *uuid.UUID `json:"term_id"`
	Name        string     `json:"name"`
	Value       string     `json:"value"`
	Position    int        `json:"position"`
}

type ProductVariantResponse struct {
	ID           uuid.UUID                         `json:"id"`
	SKU          string                            `json:"sku"`
	Price        float64                           `json:"price"`
	ComparePrice *float64                          `json:"compare_price"`
	CostPrice    *float64                          `json:"cost_price"`
	Stock        int                               `json:"stock"`
	Weight       *float64                          `json:"weight"`
	Dimensions   *DimensionsResponse               `json:"dimensions"`
	Image        string                            `json:"image"`
	Position     int                               `json:"position"`
	IsActive     bool                              `json:"is_active"`
	Attributes   []ProductVariantAttributeResponse `json:"attributes"`
}

type ProductVariantAttributeResponse struct {
	ID            uuid.UUID `json:"id"`
	AttributeID   uuid.UUID `json:"attribute_id"`
	AttributeName string    `json:"attribute_name"`
	TermID        uuid.UUID `json:"term_id"`
	TermName      string    `json:"term_name"`
	TermValue     string    `json:"term_value"`
}

// Inventory Types

// InventoryResponse represents inventory response
type InventoryResponse struct {
	ID                uuid.UUID          `json:"id"`
	ProductID         uuid.UUID          `json:"product_id"`
	WarehouseID       uuid.UUID          `json:"warehouse_id"`
	QuantityOnHand    int                `json:"quantity_on_hand"`
	QuantityReserved  int                `json:"quantity_reserved"`
	QuantityAvailable int                `json:"quantity_available"`
	ReorderLevel      int                `json:"reorder_level"`
	MaxStockLevel     *int               `json:"max_stock_level"`
	MinStockLevel     *int               `json:"min_stock_level"`
	AverageCost       float64            `json:"average_cost"`
	LastCost          *float64           `json:"last_cost"`
	LastMovementAt    *time.Time         `json:"last_movement_at"`
	LastCountAt       *time.Time         `json:"last_count_at"`
	IsLowStock        bool               `json:"is_low_stock"`
	IsOutOfStock      bool               `json:"is_out_of_stock"`
	IsOverStock       bool               `json:"is_over_stock"`
	IsActive          bool               `json:"is_active"`
	Product           *ProductResponse   `json:"product,omitempty"`
	Warehouse         *WarehouseResponse `json:"warehouse,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

// WarehouseResponse represents warehouse response
type WarehouseResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	Type        string    `json:"type"`
	IsActive    bool      `json:"is_active"`
	IsDefault   bool      `json:"is_default"`
}

// AdjustStockRequest represents adjust stock request
type AdjustStockRequest struct {
	ProductID     uuid.UUID  `json:"product_id" validate:"required"`
	WarehouseID   uuid.UUID  `json:"warehouse_id" validate:"required"`
	QuantityDelta int        `json:"quantity_delta" validate:"required"`
	Reason        string     `json:"reason" validate:"required"`
	Notes         string     `json:"notes"`
	BatchNumber   string     `json:"batch_number"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	AdjustedBy    uuid.UUID  `json:"adjusted_by" validate:"required"`
}

// GetLowStockItemsRequest represents request for getting low stock items
type GetLowStockItemsRequest struct {
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
}

// LowStockItemsResponse represents response for low stock items
type LowStockItemsResponse struct {
	Items      []*InventoryResponse `json:"items"`
	Total      int64                `json:"total"`
	Pagination PaginationResponse   `json:"pagination"`
}

// UpdateInventoryRequest represents update inventory request
type UpdateInventoryRequest struct {
	ProductID      uuid.UUID  `json:"product_id" validate:"required"`
	WarehouseID    uuid.UUID  `json:"warehouse_id" validate:"required"`
	QuantityOnHand *int       `json:"quantity_on_hand"`
	ReorderLevel   *int       `json:"reorder_level"`
	MaxStockLevel  *int       `json:"max_stock_level"`
	MinStockLevel  *int       `json:"min_stock_level"`
	AverageCost    *float64   `json:"average_cost"`
	LastCost       *float64   `json:"last_cost"`
	LastCountAt    *time.Time `json:"last_count_at"`
	UpdatedBy      uuid.UUID  `json:"updated_by" validate:"required"`
}

// GetInventoriesRequest represents request for getting inventories
type GetInventoriesRequest struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	Search         string `json:"search"`
	LowStockOnly   bool   `json:"low_stock_only"`
	OutOfStockOnly bool   `json:"out_of_stock_only"`
}

// InventoriesListResponse represents inventories list response
type InventoriesListResponse struct {
	Items      []*InventoryResponse `json:"items"`
	Total      int64                `json:"total"`
	Pagination PaginationResponse   `json:"pagination"`
}

// MovementResponse represents inventory movement response
type MovementResponse struct {
	ID             uuid.UUID  `json:"id"`
	InventoryID    uuid.UUID  `json:"inventory_id"`
	Type           string     `json:"type"`
	Reason         string     `json:"reason"`
	Quantity       int        `json:"quantity"`
	UnitCost       *float64   `json:"unit_cost"`
	TotalCost      *float64   `json:"total_cost"`
	QuantityBefore int        `json:"quantity_before"`
	QuantityAfter  int        `json:"quantity_after"`
	ReferenceType  *string    `json:"reference_type"`
	ReferenceID    *uuid.UUID `json:"reference_id"`
	Notes          string     `json:"notes"`
	BatchNumber    *string    `json:"batch_number"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
}

// RecordMovementRequest represents record movement request
type RecordMovementRequest struct {
	ProductID     uuid.UUID  `json:"product_id" validate:"required"`
	WarehouseID   uuid.UUID  `json:"warehouse_id" validate:"required"`
	Type          string     `json:"type" validate:"required"`
	Reason        string     `json:"reason" validate:"required"`
	Quantity      int        `json:"quantity" validate:"required"`
	UnitCost      *float64   `json:"unit_cost"`
	ReferenceType *string    `json:"reference_type"`
	ReferenceID   *uuid.UUID `json:"reference_id"`
	Notes         string     `json:"notes"`
	BatchNumber   *string    `json:"batch_number"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	CreatedBy     uuid.UUID  `json:"created_by" validate:"required"`
}

// GetMovementsRequest represents request for getting movements
type GetMovementsRequest struct {
	InventoryID *uuid.UUID `json:"inventory_id"`
	ProductID   *uuid.UUID `json:"product_id"`
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	Type        string     `json:"type"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
}

// MovementsListResponse represents movements list response
type MovementsListResponse struct {
	Items      []*MovementResponse `json:"items"`
	Total      int64               `json:"total"`
	Pagination PaginationResponse  `json:"pagination"`
}

// TransferStockRequest represents transfer stock request
type TransferStockRequest struct {
	ProductID       uuid.UUID `json:"product_id" validate:"required"`
	FromWarehouseID uuid.UUID `json:"from_warehouse_id" validate:"required"`
	ToWarehouseID   uuid.UUID `json:"to_warehouse_id" validate:"required"`
	Quantity        int       `json:"quantity" validate:"required"`
	Reason          string    `json:"reason" validate:"required"`
	Notes           string    `json:"notes"`
	TransferredBy   uuid.UUID `json:"transferred_by" validate:"required"`
}

// GetAlertsRequest represents request for getting alerts
type GetAlertsRequest struct {
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	ProductID   *uuid.UUID `json:"product_id"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
}

// AlertsListResponse represents alerts list response
type AlertsListResponse struct {
	Items      []*AlertResponse   `json:"items"`
	Total      int64              `json:"total"`
	Pagination PaginationResponse `json:"pagination"`
}

// AlertResponse represents alert response
type AlertResponse struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Message     string     `json:"message"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	ProductID   *uuid.UUID `json:"product_id"`
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	TriggeredAt time.Time  `json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	ResolvedBy  *uuid.UUID `json:"resolved_by"`
	Resolution  string     `json:"resolution"`
}

// InventoryReportRequest represents inventory report request
type InventoryReportRequest struct {
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	CategoryID  *uuid.UUID `json:"category_id"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	ReportType  string     `json:"report_type"`
}

// MovementReportRequest represents movement report request
type MovementReportRequest struct {
	WarehouseID  *uuid.UUID `json:"warehouse_id"`
	ProductID    *uuid.UUID `json:"product_id"`
	DateFrom     *time.Time `json:"date_from"`
	DateTo       *time.Time `json:"date_to"`
	MovementType string     `json:"movement_type"`
}

// MovementReportResponse represents movement report response
type MovementReportResponse struct {
	ReportType  string                 `json:"report_type"`
	GeneratedAt time.Time              `json:"generated_at"`
	DateRange   *DateRangeResponse     `json:"date_range,omitempty"`
	Summary     *MovementReportSummary `json:"summary"`
	Items       []*MovementReportItem  `json:"items"`
}

// MovementReportSummary represents movement report summary
type MovementReportSummary struct {
	TotalMovements   int     `json:"total_movements"`
	TotalInbound     int     `json:"total_inbound"`
	TotalOutbound    int     `json:"total_outbound"`
	TotalAdjustments int     `json:"total_adjustments"`
	NetChange        int     `json:"net_change"`
	ValueChange      float64 `json:"value_change"`
}

// MovementReportItem represents movement report item
type MovementReportItem struct {
	Date      time.Time          `json:"date"`
	Product   *ProductResponse   `json:"product"`
	Warehouse *WarehouseResponse `json:"warehouse"`
	Type      string             `json:"type"`
	Reason    string             `json:"reason"`
	Quantity  int                `json:"quantity"`
	UnitCost  *float64           `json:"unit_cost"`
	TotalCost *float64           `json:"total_cost"`
}

// DateRangeResponse represents date range response
type DateRangeResponse struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// PaginationResponse represents pagination response
type PaginationResponse = PaginationInfo

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Note: Notification types are defined in notification_usecase.go to avoid duplication
// Note: Payment types are defined in payment_usecase.go to avoid duplication
