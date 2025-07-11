package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents an error code
type ErrorCode string

const (
	// User error codes
	ErrCodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeUserNotActive     ErrorCode = "USER_NOT_ACTIVE"
	ErrCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden         ErrorCode = "FORBIDDEN"

	// Product error codes
	ErrCodeProductNotFound     ErrorCode = "PRODUCT_NOT_FOUND"
	ErrCodeProductNotAvailable ErrorCode = "PRODUCT_NOT_AVAILABLE"
	ErrCodeInsufficientStock   ErrorCode = "INSUFFICIENT_STOCK"
	ErrCodeInvalidProductData  ErrorCode = "INVALID_PRODUCT_DATA"

	// Order error codes
	ErrCodeOrderNotFound        ErrorCode = "ORDER_NOT_FOUND"
	ErrCodeOrderCannotBeCancelled ErrorCode = "ORDER_CANNOT_BE_CANCELLED"
	ErrCodeOrderCannotBeRefunded  ErrorCode = "ORDER_CANNOT_BE_REFUNDED"
	ErrCodeInvalidOrderStatus     ErrorCode = "INVALID_ORDER_STATUS"
	ErrCodeOrderAlreadyPaid       ErrorCode = "ORDER_ALREADY_PAID"

	// Payment error codes
	ErrCodePaymentNotFound             ErrorCode = "PAYMENT_NOT_FOUND"
	ErrCodePaymentFailed               ErrorCode = "PAYMENT_FAILED"
	ErrCodeInvalidPaymentAmount        ErrorCode = "INVALID_PAYMENT_AMOUNT"
	ErrCodeInvalidRefundAmount         ErrorCode = "INVALID_REFUND_AMOUNT"
	ErrCodeRefundAmountExceedsPayment  ErrorCode = "REFUND_AMOUNT_EXCEEDS_PAYMENT"
	ErrCodePaymentAlreadyProcessed     ErrorCode = "PAYMENT_ALREADY_PROCESSED"

	// Cart error codes
	ErrCodeCartNotFound    ErrorCode = "CART_NOT_FOUND"
	ErrCodeCartItemNotFound ErrorCode = "CART_ITEM_NOT_FOUND"
	ErrCodeInvalidQuantity ErrorCode = "INVALID_QUANTITY"

	// General error codes
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeConflict         ErrorCode = "CONFLICT"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	// Concurrency error codes
	ErrCodeConcurrencyConflict ErrorCode = "CONCURRENCY_CONFLICT"
	ErrCodeResourceLocked      ErrorCode = "RESOURCE_LOCKED"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Cause      error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithCause adds a cause to the error
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getDefaultStatusCode(code),
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getDefaultStatusCode(code),
		Cause:      err,
	}
}

// getDefaultStatusCode returns the default HTTP status code for an error code
func getDefaultStatusCode(code ErrorCode) int {
	switch code {
	case ErrCodeUserNotFound, ErrCodeProductNotFound, ErrCodeOrderNotFound,
		 ErrCodePaymentNotFound, ErrCodeCartNotFound, ErrCodeCartItemNotFound,
		 ErrCodeNotFound:
		return http.StatusNotFound

	case ErrCodeUserAlreadyExists, ErrCodeConflict:
		return http.StatusConflict

	case ErrCodeInvalidCredentials, ErrCodeUserNotActive, ErrCodeUnauthorized:
		return http.StatusUnauthorized

	case ErrCodeForbidden:
		return http.StatusForbidden

	case ErrCodeInvalidInput, ErrCodeInvalidQuantity, ErrCodeInvalidProductData,
		 ErrCodeInvalidOrderStatus, ErrCodeInvalidPaymentAmount, ErrCodeInvalidRefundAmount,
		 ErrCodeValidationFailed:
		return http.StatusBadRequest

	case ErrCodeProductNotAvailable, ErrCodeInsufficientStock, ErrCodeOrderCannotBeCancelled,
		 ErrCodeOrderCannotBeRefunded, ErrCodeOrderAlreadyPaid, ErrCodeRefundAmountExceedsPayment,
		 ErrCodePaymentAlreadyProcessed:
		return http.StatusUnprocessableEntity

	case ErrCodePaymentFailed:
		return http.StatusPaymentRequired

	case ErrCodeConcurrencyConflict, ErrCodeResourceLocked:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error chain
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// Common error constructors
func UserNotFound() *AppError {
	return New(ErrCodeUserNotFound, "User not found")
}

func UserAlreadyExists() *AppError {
	return New(ErrCodeUserAlreadyExists, "User already exists")
}

func InvalidCredentials() *AppError {
	return New(ErrCodeInvalidCredentials, "Invalid credentials")
}

func ProductNotFound() *AppError {
	return New(ErrCodeProductNotFound, "Product not found")
}

func InsufficientStock() *AppError {
	return New(ErrCodeInsufficientStock, "Insufficient stock")
}

func OrderNotFound() *AppError {
	return New(ErrCodeOrderNotFound, "Order not found")
}

func CartNotFound() *AppError {
	return New(ErrCodeCartNotFound, "Cart not found")
}

func InvalidInput(message string) *AppError {
	return New(ErrCodeInvalidInput, message)
}

func InternalError(message string) *AppError {
	return New(ErrCodeInternalError, message)
}

func ConcurrencyConflict(message string) *AppError {
	return New(ErrCodeConcurrencyConflict, message)
}
