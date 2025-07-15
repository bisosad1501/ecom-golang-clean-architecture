package usecases

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Constants for pagination
const (
	DefaultLimit = 20
	MaxLimit     = 100
	MinLimit     = 1

	// Ecommerce-specific pagination limits
	ProductsPerPage     = 12  // Standard grid layout (3x4 or 4x3)
	OrdersPerPage       = 10  // Order history
	ReviewsPerPage      = 5   // Product reviews
	NotificationsPerPage = 15 // User notifications
	SearchResultsPerPage = 20 // Search results
	WishlistPerPage     = 12  // Wishlist items
	AdminUsersPerPage   = 25  // Admin user management
	AdminOrdersPerPage  = 20  // Admin order management

	// Large dataset limits
	MaxSearchResults    = 1000 // Maximum search results to prevent performance issues
	MaxOrderHistory     = 500  // Maximum order history per user
)

// PaginationRequest represents pagination request parameters
type PaginationRequest struct {
	Page  int `json:"page" form:"page" binding:"min=1"`
	Limit int `json:"limit" form:"limit" binding:"min=1,max=100"`
}

// PaginationInfo represents standardized pagination information
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`

	// Enhanced ecommerce fields
	StartIndex int    `json:"start_index"` // 1-based index of first item on current page
	EndIndex   int    `json:"end_index"`   // 1-based index of last item on current page
	NextPage   *int   `json:"next_page"`   // Next page number (null if no next page)
	PrevPage   *int   `json:"prev_page"`   // Previous page number (null if no previous page)

	// SEO and UX fields
	CanonicalURL string `json:"canonical_url,omitempty"` // SEO canonical URL
	PageSizes    []int  `json:"page_sizes,omitempty"`    // Available page sizes for UX

	// Cursor pagination support
	NextCursor *string `json:"next_cursor,omitempty"` // Cursor for next page
	PrevCursor *string `json:"prev_cursor,omitempty"` // Cursor for previous page
	UseCursor  bool    `json:"use_cursor"`            // Whether cursor pagination is being used

	// Performance and caching fields
	CacheKey    string `json:"cache_key"`    // Cache key for this pagination result
	CacheTTL    int    `json:"cache_ttl"`    // Cache TTL in seconds
	IsCached    bool   `json:"is_cached"`    // Whether this result is from cache
	QueryTime   int    `json:"query_time"`   // Query execution time in milliseconds
}

// EcommercePaginationContext provides context for business logic
type EcommercePaginationContext struct {
	EntityType    string `json:"entity_type"`    // "products", "orders", "reviews", etc.
	UserID        string `json:"user_id,omitempty"`
	CategoryID    string `json:"category_id,omitempty"`
	SearchQuery   string `json:"search_query,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	FilterApplied bool   `json:"filter_applied"`
	Period        string `json:"period,omitempty"` // For time-based data: daily, weekly, monthly, etc.
}

// ValidateAndNormalizePagination validates and normalizes pagination parameters
func ValidateAndNormalizePagination(page, limit int) (int, int, error) {
	// Validate and normalize page
	if page < 1 {
		page = 1
	}

	// Validate and normalize limit
	if limit < MinLimit {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return page, limit, nil
}

// ValidateAndNormalizePaginationForEntity validates pagination with entity-specific defaults
func ValidateAndNormalizePaginationForEntity(page, limit int, entityType string) (int, int, error) {
	// Validate page
	if page < 1 {
		page = 1
	}

	// Set entity-specific default limit if not provided
	if limit <= 0 {
		switch entityType {
		case "products":
			limit = ProductsPerPage
		case "orders":
			limit = OrdersPerPage
		case "reviews":
			limit = ReviewsPerPage
		case "notifications":
			limit = NotificationsPerPage
		case "search":
			limit = SearchResultsPerPage
		case "wishlist":
			limit = WishlistPerPage
		case "admin_users":
			limit = AdminUsersPerPage
		case "admin_orders":
			limit = AdminOrdersPerPage
		default:
			limit = DefaultLimit
		}
	}

	// Validate limit bounds
	if limit < MinLimit {
		limit = MinLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return page, limit, nil
}

// NewPaginationInfo creates standardized pagination info from page-based parameters
func NewPaginationInfo(page, limit int, total int64) *PaginationInfo {
	// Validate inputs
	page, limit, _ = ValidateAndNormalizePagination(page, limit)

	// Calculate total pages (ceiling division)
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	// Ensure page doesn't exceed total pages
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate start and end indices (1-based)
	startIndex := (page-1)*limit + 1
	endIndex := page * limit
	if endIndex > int(total) {
		endIndex = int(total)
	}
	if total == 0 {
		startIndex = 0
		endIndex = 0
	}

	// Calculate next and previous page numbers
	var nextPage, prevPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	// Standard page sizes for ecommerce
	pageSizes := []int{12, 24, 48, 96}

	return &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		StartIndex: startIndex,
		EndIndex:   endIndex,
		NextPage:   nextPage,
		PrevPage:   prevPage,
		PageSizes:  pageSizes,
	}
}

// NewEcommercePaginationInfoWithCursor creates enhanced pagination info with cursor support
func NewEcommercePaginationInfoWithCursor(page, limit int, total int64, ctx *EcommercePaginationContext, firstID, lastID string, timestamp int64) *PaginationInfo {
	pagination := NewEcommercePaginationInfo(page, limit, total, ctx)

	// Add cursor support if enabled
	if pagination.UseCursor {
		if pagination.HasNext && lastID != "" {
			nextCursor := GenerateCursor(lastID, timestamp)
			pagination.NextCursor = &nextCursor
		}
		if pagination.HasPrev && firstID != "" {
			prevCursor := GenerateCursor(firstID, timestamp)
			pagination.PrevCursor = &prevCursor
		}
	}

	return pagination
}

// NewPaginationInfoFromOffset creates pagination info from offset-based parameters
func NewPaginationInfoFromOffset(offset, limit int, total int64) *PaginationInfo {
	// Validate inputs
	if limit <= 0 {
		limit = DefaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	// Convert offset to page
	page := (offset / limit) + 1

	return NewPaginationInfo(page, limit, total)
}

// NewEcommercePaginationInfo creates enhanced pagination with business context
func NewEcommercePaginationInfo(page, limit int, total int64, ctx *EcommercePaginationContext) *PaginationInfo {
	// Use entity-specific validation
	if ctx != nil {
		page, limit, _ = ValidateAndNormalizePaginationForEntity(page, limit, ctx.EntityType)
	} else {
		page, limit, _ = ValidateAndNormalizePagination(page, limit)
	}

	// Create base pagination
	pagination := NewPaginationInfo(page, limit, total)

	// Add business logic enhancements
	if ctx != nil {
		// Adjust page sizes based on entity type
		switch ctx.EntityType {
		case "products":
			pagination.PageSizes = []int{12, 24, 48, 96} // Grid-friendly sizes
		case "orders":
			pagination.PageSizes = []int{10, 20, 50}     // List-friendly sizes
		case "reviews":
			pagination.PageSizes = []int{5, 10, 20}      // Smaller sizes for detailed content
		case "search":
			pagination.PageSizes = []int{20, 40, 60}     // Search result sizes
		default:
			pagination.PageSizes = []int{10, 20, 50, 100}
		}

		// Check if cursor pagination should be used
		pagination.UseCursor = ShouldUseCursorPagination(total, ctx.EntityType)

		// Generate cache key
		cacheParams := map[string]interface{}{
			"page":  page,
			"limit": limit,
		}
		if ctx.SearchQuery != "" {
			cacheParams["query"] = ctx.SearchQuery
		}
		if ctx.CategoryID != "" {
			cacheParams["category"] = ctx.CategoryID
		}
		pagination.CacheKey = GenerateCacheKey(ctx.EntityType, ctx.UserID, cacheParams)
		pagination.CacheTTL = GetCacheTTL(ctx.EntityType)

		// Generate canonical URL for SEO (if needed)
		if ctx.SearchQuery != "" || ctx.CategoryID != "" {
			// This would be implemented based on your URL structure
			// pagination.CanonicalURL = generateCanonicalURL(ctx)
		}
	}

	return pagination
}

// ToOffset converts page-based pagination to offset
func (p *PaginationInfo) ToOffset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

// ValidatePage checks if the requested page is valid
func (p *PaginationInfo) ValidatePage() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if p.Page > p.TotalPages && p.TotalPages > 0 {
		return fmt.Errorf("page %d exceeds total pages %d", p.Page, p.TotalPages)
	}
	return nil
}

// IsFirstPage checks if current page is the first page
func (p *PaginationInfo) IsFirstPage() bool {
	return p.Page == 1
}

// IsLastPage checks if current page is the last page
func (p *PaginationInfo) IsLastPage() bool {
	return p.Page == p.TotalPages
}

// GetPageRange returns a range of page numbers for pagination UI
func (p *PaginationInfo) GetPageRange(maxPages int) []int {
	if maxPages <= 0 {
		maxPages = 5 // Default to 5 pages
	}

	start := p.Page - maxPages/2
	if start < 1 {
		start = 1
	}

	end := start + maxPages - 1
	if end > p.TotalPages {
		end = p.TotalPages
		start = end - maxPages + 1
		if start < 1 {
			start = 1
		}
	}

	pages := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// ShouldShowPagination determines if pagination should be displayed
func (p *PaginationInfo) ShouldShowPagination() bool {
	return p.TotalPages > 1
}

// GetItemsDisplayText returns text like "Showing 1-12 of 150 items"
func (p *PaginationInfo) GetItemsDisplayText() string {
	if p.Total == 0 {
		return "No items found"
	}
	if p.Total == 1 {
		return "Showing 1 item"
	}
	return fmt.Sprintf("Showing %d-%d of %d items", p.StartIndex, p.EndIndex, p.Total)
}

// Performance optimization functions

// ShouldUseCursorPagination determines if cursor-based pagination should be used
func ShouldUseCursorPagination(total int64, entityType string) bool {
	// Use cursor pagination for large datasets to improve performance
	switch entityType {
	case "products":
		return total > 10000 // Large product catalogs
	case "orders":
		return total > 5000  // Large order history
	case "search":
		return total > 1000  // Large search results
	default:
		return total > 10000
	}
}

// GenerateCursor creates a cursor for pagination
func GenerateCursor(id string, timestamp int64) string {
	// Simple cursor format: base64(id:timestamp)
	cursor := fmt.Sprintf("%s:%d", id, timestamp)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}

// ParseCursor parses a cursor to extract ID and timestamp
func ParseCursor(cursor string) (string, int64, error) {
	if cursor == "" {
		return "", 0, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", 0, fmt.Errorf("invalid cursor format")
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid cursor format")
	}

	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid cursor timestamp")
	}

	return parts[0], timestamp, nil
}

// GenerateCacheKey creates a cache key for pagination results
func GenerateCacheKey(entityType, userID string, params map[string]interface{}) string {
	// Create a deterministic cache key
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	parts = append(parts, "pagination", entityType)
	if userID != "" {
		parts = append(parts, "user", userID)
	}

	for _, k := range keys {
		parts = append(parts, k, fmt.Sprintf("%v", params[k]))
	}

	return strings.Join(parts, ":")
}

// GetCacheTTL returns appropriate cache TTL for entity type
func GetCacheTTL(entityType string) int {
	switch entityType {
	case "products":
		return 300 // 5 minutes for products
	case "orders":
		return 60 // 1 minute for orders (more dynamic)
	case "search":
		return 180 // 3 minutes for search results
	case "analytics":
		return 600 // 10 minutes for analytics
	default:
		return 300 // Default 5 minutes
	}
}

// CalculateOptimalPageSize suggests optimal page size based on entity type and context
func CalculateOptimalPageSize(entityType string, deviceType string, connectionSpeed string) int {
	baseSize := DefaultLimit

	// Adjust based on entity type
	switch entityType {
	case "products":
		baseSize = ProductsPerPage
	case "orders":
		baseSize = OrdersPerPage
	case "reviews":
		baseSize = ReviewsPerPage
	}

	// Adjust based on device type
	switch deviceType {
	case "mobile":
		baseSize = baseSize / 2 // Smaller pages for mobile
	case "tablet":
		baseSize = int(float64(baseSize) * 0.75) // Slightly smaller for tablet
	}

	// Adjust based on connection speed
	switch connectionSpeed {
	case "slow":
		baseSize = baseSize / 2 // Smaller pages for slow connections
	case "fast":
		baseSize = int(float64(baseSize) * 1.5) // Larger pages for fast connections
	}

	// Ensure within bounds
	if baseSize < MinLimit {
		baseSize = MinLimit
	}
	if baseSize > MaxLimit {
		baseSize = MaxLimit
	}

	return baseSize
}

// Common request/response structs for missing types
type TopProductsRequest struct {
	Period string `json:"period"`
	Limit  int    `json:"limit"`
}

type TopCategoriesRequest struct {
	Period string `json:"period"`
	Limit  int    `json:"limit"`
}

type TopProductsResponse struct {
	Products []interface{} `json:"products"`
	Period   string        `json:"period"`
	Total    int64         `json:"total"`
}


