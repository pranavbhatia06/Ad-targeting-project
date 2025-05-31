package models

import "errors"

// Error definitions for the targeting engine
var (
	ErrMissingAppParam     = errors.New("missing app param")
	ErrMissingCountryParam = errors.New("missing country param")
	ErrMissingOSParam      = errors.New("missing os param")
	ErrCampaignNotFound    = errors.New("campaign not found")
	ErrInvalidRuleType     = errors.New("invalid rule type")
	ErrInvalidDimension    = errors.New("invalid dimension")
)
