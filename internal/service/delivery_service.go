package service

import (
	"context"
	"fmt"
	"log"

	"github.com/razorpay/test/internal/models"
	"github.com/razorpay/test/internal/repository"
)

// DeliveryService handles campaign delivery logic
type DeliveryService struct {
	campaignRepo *repository.CampaignRepository
	engine       *models.TargetingEngine
}

// NewDeliveryService creates a new delivery service
func NewDeliveryService(campaignRepo *repository.CampaignRepository) *DeliveryService {
	return &DeliveryService{
		campaignRepo: campaignRepo,
		engine:       models.NewTargetingEngine(),
	}
}

// GetCampaigns retrieves campaigns that match the delivery request
func (s *DeliveryService) GetCampaigns(ctx context.Context, request models.DeliveryRequest) ([]models.CampaignResponse, error) {
	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Normalize request for consistent matching
	request.Normalize()

	// Load campaigns and targeting rules
	if err := s.loadTargetingData(ctx); err != nil {
		return nil, fmt.Errorf("failed to load targeting data: %w", err)
	}

	// Get active campaigns
	campaigns, err := s.campaignRepo.GetActiveCampaigns(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}

	// Filter campaigns based on targeting rules
	var matchedCampaigns []models.CampaignResponse
	for _, campaign := range campaigns {
		if campaign.IsActive() && s.engine.Matches(campaign.ID, request) {
			matchedCampaigns = append(matchedCampaigns, campaign.ToResponse())
		}
	}

	return matchedCampaigns, nil
}

// loadTargetingData loads targeting rules into the engine
func (s *DeliveryService) loadTargetingData(ctx context.Context) error {
	rules, err := s.campaignRepo.GetTargetingRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to get targeting rules: %w", err)
	}

	// Clear existing rules and reload
	s.engine = models.NewTargetingEngine()

	for _, rule := range rules {
		s.engine.AddRule(rule)
	}

	log.Printf("Loaded %d targeting rules", len(rules))
	return nil
}

// RefreshTargetingData refreshes the targeting engine with latest data
func (s *DeliveryService) RefreshTargetingData(ctx context.Context) error {
	// Invalidate cache to force fresh data load
	if err := s.campaignRepo.InvalidateCache(ctx); err != nil {
		log.Printf("Warning: failed to invalidate cache: %v", err)
	}

	return s.loadTargetingData(ctx)
}
