package usecases

// PaginationInfo represents pagination information
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// NewPaginationInfo creates pagination info
func NewPaginationInfo(offset, limit int, total int64) *PaginationInfo {
	// Prevent division by zero
	if limit <= 0 {
		limit = 20 // default limit
	}

	currentPage := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginationInfo{
		CurrentPage: currentPage,
		PerPage:     limit,
		TotalPages:  totalPages,
		TotalItems:  total,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
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


