package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/razorpay/test/internal/models"
	"github.com/razorpay/test/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler contains the HTTP handlers for the API
type Handler struct {
	deliveryService *service.DeliveryService
}

// NewHandler creates a new HTTP handler
func NewHandler(deliveryService *service.DeliveryService) *Handler {
	return &Handler{
		deliveryService: deliveryService,
	}
}

// GetCampaigns handles GET /campaigns requests
func (h *Handler) GetCampaigns(c *gin.Context) {
	// Parse query parameters
	request := models.DeliveryRequest{
		App:     c.Query("app"),
		Country: c.Query("country"),
		OS:      c.Query("os"),
	}
	fmt.Println(request)

	// Parse optional limit parameter
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			request.Limit = limit
		}
	}

	// Get campaigns from service
	campaigns, err := h.deliveryService.GetCampaigns(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get campaigns",
			"message": err.Error(),
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"count":     len(campaigns),
	})
}

// PostCampaigns handles POST /campaigns requests
func (h *Handler) PostCampaigns(c *gin.Context) {
	var request models.DeliveryRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Get campaigns from service
	campaigns, err := h.deliveryService.GetCampaigns(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get campaigns",
			"message": err.Error(),
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"count":     len(campaigns),
	})
}

// RefreshTargeting handles POST /refresh requests
func (h *Handler) RefreshTargeting(c *gin.Context) {
	if err := h.deliveryService.RefreshTargetingData(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to refresh targeting data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Targeting data refreshed successfully",
	})
}

// HealthCheck handles GET /health requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "targeting-engine",
	})
}
