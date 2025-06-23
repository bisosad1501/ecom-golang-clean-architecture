package entities

import "errors"

// Domain errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive     = errors.New("user is not active")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")

	// Product errors
	ErrProductNotFound     = errors.New("product not found")
	ErrProductNotAvailable = errors.New("product not available")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrInvalidProductData  = errors.New("invalid product data")

	// Category errors
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")

	// Cart errors
	ErrCartNotFound    = errors.New("cart not found")
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrInvalidQuantity = errors.New("invalid quantity")

	// Order errors
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
	ErrOrderCannotBeRefunded  = errors.New("order cannot be refunded")
	ErrInvalidOrderStatus     = errors.New("invalid order status")
	ErrOrderAlreadyPaid       = errors.New("order already paid")

	// Payment errors
	ErrPaymentNotFound             = errors.New("payment not found")
	ErrPaymentFailed               = errors.New("payment failed")
	ErrInvalidPaymentAmount        = errors.New("invalid payment amount")
	ErrInvalidRefundAmount         = errors.New("invalid refund amount")
	ErrRefundAmountExceedsPayment  = errors.New("refund amount exceeds payment amount")
	ErrPaymentAlreadyProcessed     = errors.New("payment already processed")

	// General errors
	ErrInvalidInput    = errors.New("invalid input")
	ErrInternalError   = errors.New("internal server error")
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrValidationFailed = errors.New("validation failed")
)
