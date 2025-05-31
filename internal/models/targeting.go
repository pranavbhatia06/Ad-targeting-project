package models

import (
	"strings"
	"time"
)

// RuleType represents the type of targeting rule
type RuleType string

const (
	RuleTypeInclude RuleType = "INCLUDE"
	RuleTypeExclude RuleType = "EXCLUDE"
)

// Dimension represents the targeting dimension
type Dimension string

const (
	DimensionApp     Dimension = "APP"
	DimensionCountry Dimension = "COUNTRY"
	DimensionOS      Dimension = "OS"
)

// TargetingRule represents a targeting rule for a campaign
type TargetingRule struct {
	ID         int64     `json:"id" db:"id"`
	CampaignID string    `json:"campaign_id" db:"campaign_id"`
	RuleType   RuleType  `json:"rule_type" db:"rule_type"`
	Dimension  Dimension `json:"dimension" db:"dimension"`
	Values     []string  `json:"values" db:"values"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// DeliveryRequest represents an incoming delivery request
type DeliveryRequest struct {
	App     string `json:"app"`
	Country string `json:"country"`
	OS      string `json:"os"`
	Limit   int    `json:"limit,omitempty"`
}

// Validate checks if the delivery request has all required fields
func (dr *DeliveryRequest) Validate() error {
	if dr.App == "" {
		return ErrMissingAppParam
	}
	if dr.Country == "" {
		return ErrMissingCountryParam
	}
	if dr.OS == "" {
		return ErrMissingOSParam
	}
	return nil
}

// Normalize converts request values to lowercase for consistent matching
func (dr *DeliveryRequest) Normalize() {
	dr.App = strings.ToLower(dr.App)
	dr.Country = strings.ToLower(dr.Country)
	dr.OS = strings.ToLower(dr.OS)
}

// TargetingEngine represents the targeting logic
type TargetingEngine struct {
	Rules map[string][]TargetingRule // campaignID -> rules
}

// NewTargetingEngine creates a new targeting engine
func NewTargetingEngine() *TargetingEngine {
	return &TargetingEngine{
		Rules: make(map[string][]TargetingRule),
	}
}

// AddRule adds a targeting rule to the engine
func (te *TargetingEngine) AddRule(rule TargetingRule) {
	if te.Rules[rule.CampaignID] == nil {
		te.Rules[rule.CampaignID] = make([]TargetingRule, 0)
	}
	te.Rules[rule.CampaignID] = append(te.Rules[rule.CampaignID], rule)
}

// Matches checks if a delivery request matches the targeting rules for a campaign
func (te *TargetingEngine) Matches(campaignID string, request DeliveryRequest) bool {
	rules, exists := te.Rules[campaignID]
	if !exists {
		return true // No rules means campaign matches all requests
	}

	// Group rules by dimension
	dimensionRules := make(map[Dimension][]TargetingRule)
	for _, rule := range rules {
		dimensionRules[rule.Dimension] = append(dimensionRules[rule.Dimension], rule)
	}

	// Check each dimension
	for dimension, dimRules := range dimensionRules {
		if !te.matchesDimension(dimension, dimRules, request) {
			return false
		}
	}

	return true
}

// matchesDimension checks if request matches rules for a specific dimension
func (te *TargetingEngine) matchesDimension(dimension Dimension, rules []TargetingRule, request DeliveryRequest) bool {
	var requestValue string
	switch dimension {
	case DimensionApp:
		requestValue = strings.ToLower(request.App)
	case DimensionCountry:
		requestValue = strings.ToLower(request.Country)
	case DimensionOS:
		requestValue = strings.ToLower(request.OS)
	}

	hasInclude := false
	hasExclude := false
	includeMatch := false
	excludeMatch := false

	for _, rule := range rules {
		// Normalize rule values for comparison
		normalizedValues := make([]string, len(rule.Values))
		for i, v := range rule.Values {
			normalizedValues[i] = strings.ToLower(v)
		}

		if rule.RuleType == RuleTypeInclude {
			hasInclude = true
			if te.containsValue(normalizedValues, requestValue) {
				includeMatch = true
			}
		} else if rule.RuleType == RuleTypeExclude {
			hasExclude = true
			if te.containsValue(normalizedValues, requestValue) {
				excludeMatch = true
			}
		}
	}

	// If there's an exclude rule and it matches, reject
	if hasExclude && excludeMatch {
		return false
	}

	// If there's an include rule, it must match
	if hasInclude {
		return includeMatch
	}

	// No rules or only exclude rules that don't match
	return true
}

// containsValue checks if a slice contains a specific value
func (te *TargetingEngine) containsValue(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
