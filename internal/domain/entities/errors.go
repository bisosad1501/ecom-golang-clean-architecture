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

	// Brand errors
	ErrBrandNotFound = errors.New("brand not found")
	ErrBrandExists   = errors.New("brand already exists")

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

	// Address errors
	ErrAddressNotFound = errors.New("address not found")

	// Wishlist errors
	ErrWishlistItemNotFound = errors.New("wishlist item not found")

	// User preference errors
	ErrUserPreferenceNotFound = errors.New("user preference not found")

	// Account verification errors
	ErrAccountVerificationNotFound = errors.New("account verification not found")
	ErrInvalidVerificationCode     = errors.New("invalid verification code")
	ErrVerificationCodeExpired     = errors.New("verification code expired")

	// Password reset errors
	ErrPasswordResetNotFound = errors.New("password reset not found")
	ErrPasswordResetExpired  = errors.New("password reset expired")
	ErrPasswordResetUsed     = errors.New("password reset already used")

	// Review errors
	ErrReviewNotFound = errors.New("review not found")
	ErrReviewVoteNotFound = errors.New("review vote not found")

	// Coupon errors
	ErrCouponNotFound = errors.New("coupon not found")
	ErrCouponCodeExists = errors.New("coupon code already exists")
	ErrCouponInvalid = errors.New("coupon is invalid")
	ErrCouponExpired = errors.New("coupon has expired")
	ErrCouponNotApplicable = errors.New("coupon is not applicable")
	ErrCouponUsageLimitExceeded = errors.New("coupon usage limit exceeded")

	// Promotion errors
	ErrPromotionNotFound = errors.New("promotion not found")

	// Loyalty program errors
	ErrLoyaltyProgramNotFound = errors.New("loyalty program not found")
	ErrInsufficientPoints = errors.New("insufficient loyalty points")

	// General errors
	ErrInvalidInput     = errors.New("invalid input")
	ErrInternalError    = errors.New("internal server error")
	ErrNotFound         = errors.New("resource not found")
	ErrConflict         = errors.New("resource conflict")
	ErrValidationFailed = errors.New("validation failed")
	ErrNotImplemented   = errors.New("feature not implemented")

	// File upload errors
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrFileNotFound    = errors.New("file not found")
	ErrNoValidFiles    = errors.New("no valid files provided")

	// Shipping errors
	ErrShippingMethodNotFound = errors.New("shipping method not found")
	ErrShipmentNotFound       = errors.New("shipment not found")
	ErrReturnNotFound         = errors.New("return not found")
	ErrOrderCannotBeReturned  = errors.New("order cannot be returned")
)
