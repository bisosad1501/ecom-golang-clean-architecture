package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddressHandler handles address-related HTTP requests
type AddressHandler struct {
	addressUseCase usecases.AddressUseCase
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressUseCase usecases.AddressUseCase) *AddressHandler {
	return &AddressHandler{
		addressUseCase: addressUseCase,
	}
}

// CreateAddress handles creating a new address
// @Summary Create a new address
// @Description Create a new address for the authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateAddressRequest true "Create address request"
// @Success 201 {object} usecases.AddressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	var req usecases.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	address, err := h.addressUseCase.CreateAddress(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Address created successfully",
		Data:    address,
	})
}

// GetAddresses handles getting user's addresses
// @Summary Get user's addresses
// @Description Get all addresses for the authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} usecases.AddressResponse
// @Failure 401 {object} ErrorResponse
// @Router /addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addresses, err := h.addressUseCase.GetUserAddresses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: addresses,
	})
}

// GetAddress handles getting a specific address
// @Summary Get address by ID
// @Description Get a specific address by its ID
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} usecases.AddressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /addresses/{id} [get]
func (h *AddressHandler) GetAddress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid address ID",
		})
		return
	}

	address, err := h.addressUseCase.GetAddress(c.Request.Context(), userID, addressID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: address,
	})
}

// UpdateAddress handles updating an address
// @Summary Update address
// @Description Update an existing address
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body usecases.UpdateAddressRequest true "Update address request"
// @Success 200 {object} usecases.AddressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid address ID",
		})
		return
	}

	var req usecases.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	address, err := h.addressUseCase.UpdateAddress(c.Request.Context(), userID, addressID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Address updated successfully",
		Data:    address,
	})
}

// DeleteAddress handles deleting an address
// @Summary Delete address
// @Description Delete an existing address
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid address ID",
		})
		return
	}

	err = h.addressUseCase.DeleteAddress(c.Request.Context(), userID, addressID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Address deleted successfully",
	})
}

// SetDefaultAddress handles setting an address as default
// @Summary Set default address
// @Description Set an address as default for shipping or billing
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body map[string]string true "Address type (shipping/billing/both)"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /addresses/{id}/default [post]
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid address ID",
		})
		return
	}

	var req struct {
		Type string `json:"type" validate:"required,oneof=shipping billing both"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	addressType := entities.AddressType(req.Type)
	err = h.addressUseCase.SetDefaultAddress(c.Request.Context(), userID, addressID, addressType)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Default address set successfully",
	})
}

// GetDefaultAddress handles getting default address
// @Summary Get default address
// @Description Get default address for shipping or billing
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "Address type (shipping/billing)"
// @Success 200 {object} usecases.AddressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /addresses/default [get]
func (h *AddressHandler) GetDefaultAddress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	addressTypeStr := c.Query("type")
	if addressTypeStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Address type is required",
		})
		return
	}

	addressType := entities.AddressType(addressTypeStr)
	if addressType != entities.AddressTypeShipping && addressType != entities.AddressTypeBilling {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid address type. Must be 'shipping' or 'billing'",
		})
		return
	}

	address, err := h.addressUseCase.GetDefaultAddress(c.Request.Context(), userID, addressType)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: address,
	})
}
