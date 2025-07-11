package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/domain/entities"
	pkgErrors "ecom-golang-clean-architecture/pkg/errors"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// getErrorStatusCode returns appropriate HTTP status code for domain errors
func getErrorStatusCode(err error) int {
	// Check if it's an AppError first
	if appErr := pkgErrors.GetAppError(err); appErr != nil {
		return appErr.StatusCode
	}

	// Fallback to legacy error handling
	switch err {
	case entities.ErrUserNotFound,
		 entities.ErrProductNotFound,
		 entities.ErrCategoryNotFound,
		 entities.ErrCartNotFound,
		 entities.ErrCartItemNotFound,
		 entities.ErrOrderNotFound,
		 entities.ErrPaymentNotFound,
		 entities.ErrNotFound:
		return http.StatusNotFound

	case entities.ErrUserAlreadyExists,
		 entities.ErrCategoryExists,
		 entities.ErrConflict:
		return http.StatusConflict

	case entities.ErrInvalidCredentials,
		 entities.ErrUserNotActive,
		 entities.ErrUnauthorized:
		return http.StatusUnauthorized

	case entities.ErrForbidden:
		return http.StatusForbidden

	case entities.ErrInvalidInput,
		 entities.ErrInvalidQuantity,
		 entities.ErrInvalidProductData,
		 entities.ErrInvalidOrderStatus,
		 entities.ErrInvalidPaymentAmount,
		 entities.ErrInvalidRefundAmount,
		 entities.ErrValidationFailed:
		return http.StatusBadRequest

	case entities.ErrProductNotAvailable,
		 entities.ErrInsufficientStock,
		 entities.ErrOrderCannotBeCancelled,
		 entities.ErrOrderCannotBeRefunded,
		 entities.ErrOrderAlreadyPaid,
		 entities.ErrRefundAmountExceedsPayment,
		 entities.ErrPaymentAlreadyProcessed:
		return http.StatusUnprocessableEntity

	case entities.ErrPaymentFailed:
		return http.StatusPaymentRequired

	default:
		return http.StatusInternalServerError
	}
}
