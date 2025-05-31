package http

import (
	"github.com/razorpay/test/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all HTTP routes for the application
func SetupRoutes(deliveryService *service.DeliveryService) *gin.Engine {
	// Create Gin router with default middleware
	router := gin.Default()

	// Create handler
	handler := NewHandler(deliveryService)

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Campaign delivery endpoints
		v1.GET("/campaigns", handler.GetCampaigns)
		v1.POST("/campaigns", handler.PostCampaigns)

		// Admin endpoints
		v1.POST("/refresh", handler.RefreshTargeting)
	}

	return router
}
