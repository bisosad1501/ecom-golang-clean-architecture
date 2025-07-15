package handlers

import (
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InventoryHandler handles inventory-related HTTP requests
type InventoryHandler struct {
	inventoryUseCase usecases.InventoryUseCase
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryUseCase usecases.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
	}
}

// GetInventory gets inventory by product and warehouse ID
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid product ID",
			Details: err.Error(),
		})
		return
	}

	warehouseIDStr := c.Param("warehouseId")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid warehouse ID",
			Details: err.Error(),
		})
		return
	}

	inventory, err := h.inventoryUseCase.GetInventory(c.Request.Context(), productID, warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Inventory not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Inventory retrieved successfully",
		Data: inventory,
	})
}

// GetInventories gets inventories with pagination
func (h *InventoryHandler) GetInventories(c *gin.Context) {
	warehouseIDStr := c.Query("warehouse_id")
	if warehouseIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "warehouse_id query parameter is required",
			Details: "Please provide a warehouse_id",
		})
		return
	}

	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid warehouse ID",
			Details: err.Error(),
		})
		return
	}

	// Parse and validate pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate and normalize pagination
	page, limit, err = usecases.ValidateAndNormalizePagination(page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	search := c.Query("search")

	req := usecases.GetInventoriesRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	inventories, err := h.inventoryUseCase.GetWarehouseInventories(c.Request.Context(), warehouseID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get inventories",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Inventories retrieved successfully",
		Data: inventories,
	})
}

// UpdateInventory updates inventory
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid inventory ID",
			Details: err.Error(),
		})
		return
	}

	var req usecases.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	inventory, err := h.inventoryUseCase.UpdateInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update inventory",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Inventory updated successfully",
		Data: inventory,
	})
}

// AdjustStock adjusts inventory stock
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var req usecases.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	movement, err := h.inventoryUseCase.AdjustStock(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to adjust stock",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Stock adjusted successfully",
		Data: movement,
	})
}

// GetLowStockItems gets low stock items
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := usecases.GetLowStockItemsRequest{
		Page:  page,
		Limit: limit,
	}

	items, err := h.inventoryUseCase.GetLowStockItems(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get low stock items",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Low stock items retrieved successfully",
		Data: items,
	})
}

// TransferStock transfers stock between warehouses
func (h *InventoryHandler) TransferStock(c *gin.Context) {
	var req usecases.TransferStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err := h.inventoryUseCase.TransferStock(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to transfer stock",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Stock transferred successfully",
		Data: nil,
	})
}

// GetMovements gets inventory movements
func (h *InventoryHandler) GetMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := usecases.GetMovementsRequest{
		Page:  page,
		Limit: limit,
	}

	movements, err := h.inventoryUseCase.GetMovements(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get movements",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Movements retrieved successfully",
		Data: movements,
	})
}

// RecordMovement records an inventory movement
func (h *InventoryHandler) RecordMovement(c *gin.Context) {
	var req usecases.RecordMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	movement, err := h.inventoryUseCase.RecordMovement(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to record movement",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Movement recorded successfully",
		Data: movement,
	})
}

// GetStockAlerts gets stock alerts
func (h *InventoryHandler) GetStockAlerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	alertType := c.Query("type")
	status := c.Query("status")

	req := usecases.GetAlertsRequest{
		Page:   page,
		Limit:  limit,
		Type:   alertType,
		Status: status,
	}

	alerts, err := h.inventoryUseCase.GetStockAlerts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get stock alerts",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Stock alerts retrieved successfully",
		Data: alerts,
	})
}

// ResolveAlert resolves a stock alert
func (h *InventoryHandler) ResolveAlert(c *gin.Context) {
	alertIDStr := c.Param("id")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid alert ID",
			Details: err.Error(),
		})
		return
	}

	var req struct {
		Resolution string `json:"resolution" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// For now, use a system user ID - in real implementation, get from auth context
	systemUserID := uuid.New()

	err = h.inventoryUseCase.ResolveAlert(c.Request.Context(), alertID, req.Resolution, systemUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to resolve alert",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Alert resolved successfully",
	})
}

// GetOutOfStockItems gets out of stock items
func (h *InventoryHandler) GetOutOfStockItems(c *gin.Context) {
	warehouseIDStr := c.Query("warehouse_id")
	var warehouseID *uuid.UUID

	if warehouseIDStr != "" {
		id, err := uuid.Parse(warehouseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid warehouse ID",
				Details: err.Error(),
			})
			return
		}
		warehouseID = &id
	}

	items, err := h.inventoryUseCase.GetOutOfStockItems(c.Request.Context(), warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get out of stock items",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Out of stock items retrieved successfully",
		Data: items,
	})
}
