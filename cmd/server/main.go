package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/razorpay/test/internal/config"
	"github.com/razorpay/test/internal/repository"
	"github.com/razorpay/test/internal/service"
	httpTransport "github.com/razorpay/test/internal/transport/http"
	"github.com/razorpay/test/pkg/cache"
	"github.com/razorpay/test/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables and seed data
	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	fmt.Println("CREATE TABLES ERR", err)

	if err := database.SeedData(db); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}
	fmt.Println("SEED TABLES ERR", err)
	// Initialize Redis cache
	redisCache := cache.NewRedisCache(cfg.Redis)
	defer redisCache.Close()

	// Test Redis connection
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// Initialize repository
	campaignRepo := repository.NewCampaignRepository(db, redisCache)

	// Initialize service
	deliveryService := service.NewDeliveryService(campaignRepo)

	// Setup HTTP routes
	router := httpTransport.SetupRoutes(deliveryService)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
