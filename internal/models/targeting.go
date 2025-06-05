package models

import (
	"fmt"
	"strings"
	"time"
)

// State-Country mappings for validation
var stateCountryMap = map[string]string{
	// India states
	"delhi":            "india",
	"maharashtra":      "india",
	"karnataka":        "india",
	"tamil nadu":       "india",
	"telangana":        "india",
	"andhra pradesh":   "india",
	"gujarat":          "india",
	"rajasthan":        "india",
	"uttar pradesh":    "india",
	"west bengal":      "india",
	"punjab":           "india",
	"haryana":          "india",
	"kerala":           "india",
	"odisha":           "india",
	"bihar":            "india",
	"jharkhand":        "india",
	"assam":            "india",
	"himachal pradesh": "india",
	"uttarakhand":      "india",
	"goa":              "india",

	// US states
	"california":     "us",
	"new york":       "us",
	"texas":          "us",
	"florida":        "us",
	"illinois":       "us",
	"pennsylvania":   "us",
	"ohio":           "us",
	"georgia":        "us",
	"north carolina": "us",
	"michigan":       "us",
	"new jersey":     "us",
	"virginia":       "us",
	"washington":     "us",
	"arizona":        "us",
	"massachusetts":  "us",
	"tennessee":      "us",
	"indiana":        "us",
	"missouri":       "us",
	"maryland":       "us",
	"wisconsin":      "us",

	// Canada provinces
	"ontario":                   "canada",
	"quebec":                    "canada",
	"british columbia":          "canada",
	"alberta":                   "canada",
	"manitoba":                  "canada",
	"saskatchewan":              "canada",
	"nova scotia":               "canada",
	"new brunswick":             "canada",
	"newfoundland and labrador": "canada",
	"prince edward island":      "canada",

	// UK countries/regions
	"england":          "uk",
	"scotland":         "uk",
	"wales":            "uk",
	"northern ireland": "uk",
	"london":           "uk",

	// Australia states
	"new south wales":              "australia",
	"victoria":                     "australia",
	"queensland":                   "australia",
	"western australia":            "australia",
	"south australia":              "australia",
	"tasmania":                     "australia",
	"northern territory":           "australia",
	"australian capital territory": "australia",
}

// Country aliases for flexible matching
var countryAliases = map[string]string{
	"in":             "india",
	"ind":            "india",
	"us":             "us",
	"usa":            "us",
	"united states":  "us",
	"ca":             "canada",
	"can":            "canada",
	"uk":             "uk",
	"gb":             "uk",
	"united kingdom": "uk",
	"au":             "australia",
	"aus":            "australia",
}

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
	DimensionsState  Dimension = "STATE"
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
	State   string `json:"state"`
	Limit   int    `json:"limit,omitempty"`
}

// Validate checks if the delivery request has all required fields and validates state-country combination
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

	// Validate state-country combination if both are provided
	if dr.State != "" && dr.Country != "" {
		if err := dr.validateStateCountry(); err != nil {
			return err
		}
	}

	return nil
}

// validateStateCountry validates that the state belongs to the specified country
func (dr *DeliveryRequest) validateStateCountry() error {
	stateName := strings.ToLower(strings.TrimSpace(dr.State))
	countryName := strings.ToLower(strings.TrimSpace(dr.Country))

	// Normalize country name using aliases
	if alias, exists := countryAliases[countryName]; exists {
		countryName = alias
	}

	// Check if state exists in our mapping
	expectedCountry, stateExists := stateCountryMap[stateName]
	if !stateExists {
		return fmt.Errorf("unknown state: '%s'", dr.State)
	}

	// Check if state belongs to the specified country
	if expectedCountry != countryName {
		return fmt.Errorf("invalid state-country combination: state '%s' does not belong to country '%s' (belongs to '%s')",
			dr.State, dr.Country, expectedCountry)
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
	case DimensionsState:
		requestValue = strings.ToLower(request.State)
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
