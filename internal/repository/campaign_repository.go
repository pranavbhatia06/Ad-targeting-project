package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/razorpay/test/internal/models"
	"github.com/razorpay/test/pkg/cache"

	"github.com/lib/pq"
)

// CampaignRepository handles campaign data operations
type CampaignRepository struct {
	db    *sql.DB
	cache *cache.RedisCache
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *sql.DB, cache *cache.RedisCache) *CampaignRepository {
	return &CampaignRepository{
		db:    db,
		cache: cache,
	}
}

// GetActiveCampaigns retrieves all active campaigns with caching
func (r *CampaignRepository) GetActiveCampaigns(ctx context.Context) ([]models.Campaign, error) {
	// Try cache first
	if r.cache != nil {
		campaigns, err := r.cache.GetCampaigns(ctx)
		if err == nil && campaigns != nil {
			return campaigns, nil
		}
	}

	// Fallback to database
	query := `
		SELECT id, name, image_url, cta, status, created_at, updated_at 
		FROM campaigns 
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, models.CampaignStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		err := rows.Scan(
			&campaign.ID,
			&campaign.Name,
			&campaign.ImageURL,
			&campaign.CTA,
			&campaign.Status,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaigns: %w", err)
	}

	// Cache the results
	if r.cache != nil {
		_ = r.cache.SetCampaigns(ctx, campaigns)
	}

	return campaigns, nil
}

// GetTargetingRules retrieves all targeting rules with caching
func (r *CampaignRepository) GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error) {
	// Try cache first
	if r.cache != nil {
		rules, err := r.cache.GetTargetingRules(ctx)
		if err == nil && rules != nil {
			return rules, nil
		}
	}

	// Fallback to database
	query := `
		SELECT id, campaign_id, rule_type, dimension, values, created_at, updated_at
		FROM targeting_rules
		ORDER BY campaign_id, dimension
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query targeting rules: %w", err)
	}
	defer rows.Close()

	var rules []models.TargetingRule
	for rows.Next() {
		var rule models.TargetingRule
		err := rows.Scan(
			&rule.ID,
			&rule.CampaignID,
			&rule.RuleType,
			&rule.Dimension,
			pq.Array(&rule.Values),
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan targeting rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating targeting rules: %w", err)
	}

	// Cache the results
	if r.cache != nil {
		_ = r.cache.SetTargetingRules(ctx, rules)
	}

	return rules, nil
}

// InvalidateCache clears the cache
func (r *CampaignRepository) InvalidateCache(ctx context.Context) error {
	if r.cache != nil {
		return r.cache.InvalidateCache(ctx)
	}
	return nil
}
