package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AbandonedCartHandler struct {
	abandonedCartUseCase usecases.AbandonedCartUseCase
}

func NewAbandonedCartHandler(abandonedCartUseCase usecases.AbandonedCartUseCase) *AbandonedCartHandler {
	return &AbandonedCartHandler{
		abandonedCartUseCase: abandonedCartUseCase,
	}
}

// GetAbandonedCarts gets list of abandoned carts
func (h *AbandonedCartHandler) GetAbandonedCarts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	carts, err := h.abandonedCartUseCase.GetAbandonedCarts(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get abandoned carts",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Abandoned carts retrieved successfully",
		Data:    carts,
	})
}

// GetAbandonedCartStats gets abandoned cart statistics
func (h *AbandonedCartHandler) GetAbandonedCartStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)

	since := time.Now().AddDate(0, 0, -days)

	stats, err := h.abandonedCartUseCase.GetAbandonedCartStats(c.Request.Context(), since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get abandoned cart stats",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Abandoned cart stats retrieved successfully",
		Data:    stats,
	})
}

// ProcessAbandonedCarts processes abandoned carts and sends reminder emails
func (h *AbandonedCartHandler) ProcessAbandonedCarts(c *gin.Context) {
	// First detect abandoned carts
	err := h.abandonedCartUseCase.DetectAbandonedCarts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to detect abandoned carts",
			Details: err.Error(),
		})
		return
	}

	// Then send emails
	err = h.abandonedCartUseCase.SendAbandonedCartEmails(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send abandoned cart emails",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Abandoned carts processed successfully",
		Data:    nil,
	})
}

// SendReminderEmail sends reminder email for specific abandoned cart
func (h *AbandonedCartHandler) SendReminderEmail(c *gin.Context) {
	cartID := c.Param("id")
	if cartID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Cart ID is required",
			Details: "Cart ID parameter is missing",
		})
		return
	}

	// For now, just trigger the general abandoned cart email process
	err := h.abandonedCartUseCase.SendAbandonedCartEmails(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send reminder email",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reminder email sent successfully",
		Data:    nil,
	})
}
