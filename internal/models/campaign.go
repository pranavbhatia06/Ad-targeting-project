package models

import (
	"time"
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusActive   CampaignStatus = "ACTIVE"
	CampaignStatusInactive CampaignStatus = "INACTIVE"
)

// Campaign represents an advertisement campaign
type Campaign struct {
	ID        string         `json:"cid" db:"id"`
	Name      string         `json:"name" db:"name"`
	ImageURL  string         `json:"img" db:"image_url"`
	CTA       string         `json:"cta" db:"cta"`
	Status    CampaignStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// IsActive returns true if the campaign is active
func (c *Campaign) IsActive() bool {
	return c.Status == CampaignStatusActive
}

// CampaignResponse represents the response format for delivery API
type CampaignResponse struct {
	CID string `json:"cid"`
	Img string `json:"img"`
	CTA string `json:"cta"`
}

// ToResponse converts Campaign to CampaignResponse
func (c *Campaign) ToResponse() CampaignResponse {
	return CampaignResponse{
		CID: c.ID,
		Img: c.ImageURL,
		CTA: c.CTA,
	}
}
