package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/razorpay/test/internal/config"
	"github.com/razorpay/test/internal/models"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements caching using Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(cfg config.RedisConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisCache{
		client: rdb,
		ttl:    cfg.TTL,
	}
}

// Ping tests the Redis connection
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetCampaigns retrieves campaigns from cache
func (r *RedisCache) GetCampaigns(ctx context.Context) ([]models.Campaign, error) {
	key := "campaigns:active"
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get campaigns from cache: %w", err)
	}

	var campaigns []models.Campaign
	if err := json.Unmarshal([]byte(val), &campaigns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campaigns: %w", err)
	}

	return campaigns, nil
}

// SetCampaigns stores campaigns in cache
func (r *RedisCache) SetCampaigns(ctx context.Context, campaigns []models.Campaign) error {
	key := "campaigns:active"
	data, err := json.Marshal(campaigns)
	if err != nil {
		return fmt.Errorf("failed to marshal campaigns: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set campaigns in cache: %w", err)
	}

	return nil
}

// GetTargetingRules retrieves targeting rules from cache
func (r *RedisCache) GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error) {
	key := "targeting_rules"
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get targeting rules from cache: %w", err)
	}

	var rules []models.TargetingRule
	if err := json.Unmarshal([]byte(val), &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal targeting rules: %w", err)
	}

	return rules, nil
}

// SetTargetingRules stores targeting rules in cache
func (r *RedisCache) SetTargetingRules(ctx context.Context, rules []models.TargetingRule) error {
	key := "targeting_rules"
	data, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("failed to marshal targeting rules: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set targeting rules in cache: %w", err)
	}

	return nil
}

// InvalidateCache clears all cached data
func (r *RedisCache) InvalidateCache(ctx context.Context) error {
	keys := []string{"campaigns:active", "targeting_rules"}
	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("failed to delete cache key %s: %w", key, err)
		}
	}
	return nil
}
