package handlers

import (
	"net/http"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShippingHandler handles shipping-related HTTP requests
type ShippingHandler struct {
	shippingUseCase usecases.ShippingUseCase
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(shippingUseCase usecases.ShippingUseCase) *ShippingHandler {
	return &ShippingHandler{
		shippingUseCase: shippingUseCase,
	}
}

// GetShippingMethods retrieves available shipping methods
func (h *ShippingHandler) GetShippingMethods(c *gin.Context) {
	var req usecases.GetShippingMethodsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	methods, err := h.shippingUseCase.GetShippingMethods(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get shipping methods",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipping methods retrieved successfully",
		Data: methods,
	})
}

// CalculateShippingCost calculates shipping cost
func (h *ShippingHandler) CalculateShippingCost(c *gin.Context) {
	var req usecases.CalculateShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	cost, err := h.shippingUseCase.CalculateShippingCost(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to calculate shipping cost",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipping cost calculated successfully",
		Data: cost,
	})
}

// CreateShipment creates a new shipment
func (h *ShippingHandler) CreateShipment(c *gin.Context) {
	var req usecases.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	shipment, err := h.shippingUseCase.CreateShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create shipment",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Shipment created successfully",
		Data: shipment,
	})
}

// GetShipment retrieves a shipment by ID
func (h *ShippingHandler) GetShipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid shipment ID",
			Details: err.Error(),
		})
		return
	}

	shipment, err := h.shippingUseCase.GetShipment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Shipment not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipment retrieved successfully",
		Data: shipment,
	})
}

// UpdateShipmentStatus updates shipment status
func (h *ShippingHandler) UpdateShipmentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid shipment ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.ShipmentStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	shipment, err := h.shippingUseCase.UpdateShipmentStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update shipment status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipment status updated successfully",
		Data: shipment,
	})
}

// TrackShipment tracks a shipment by tracking number
func (h *ShippingHandler) TrackShipment(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Tracking number is required",
		})
		return
	}

	tracking, err := h.shippingUseCase.TrackShipment(c.Request.Context(), trackingNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Shipment not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipment tracking retrieved successfully",
		Data: tracking,
	})
}

// CreateReturn creates a new return request
func (h *ShippingHandler) CreateReturn(c *gin.Context) {
	var req usecases.CreateReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	returnReq, err := h.shippingUseCase.CreateReturn(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create return",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Return created successfully",
		Data: returnReq,
	})
}

// GetReturn retrieves a return by ID
func (h *ShippingHandler) GetReturn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid return ID",
			Details: err.Error(),
		})
		return
	}

	returnReq, err := h.shippingUseCase.GetReturn(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Return not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Return retrieved successfully",
		Data: returnReq,
	})
}

// ProcessReturn processes a return request
func (h *ShippingHandler) ProcessReturn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid return ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Status entities.ReturnStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	returnReq, err := h.shippingUseCase.ProcessReturn(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process return",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Return processed successfully",
		Data: returnReq,
	})
}

// CalculateDistanceBasedShipping calculates shipping options based on distance
func (h *ShippingHandler) CalculateDistanceBasedShipping(c *gin.Context) {
	var req usecases.DistanceBasedShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	response, err := h.shippingUseCase.CalculateDistanceBasedShipping(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to calculate distance-based shipping",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Distance-based shipping calculated successfully",
		Data: response,
	})
}

// GetShippingZones returns available shipping zones
func (h *ShippingHandler) GetShippingZones(c *gin.Context) {
	zones, err := h.shippingUseCase.GetShippingZones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get shipping zones",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Shipping zones retrieved successfully",
		Data: zones,
	})
}
