package database

import (
	"database/sql"
	"fmt"

	"github.com/razorpay/test/internal/config"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CreateTables creates the necessary database tables
func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS campaigns (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			image_url TEXT NOT NULL,
			cta VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS targeting_rules (
			id SERIAL PRIMARY KEY,
			campaign_id VARCHAR(255) NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
			rule_type VARCHAR(50) NOT NULL,
			dimension VARCHAR(50) NOT NULL,
			values TEXT[] NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);`,
		`CREATE INDEX IF NOT EXISTS idx_targeting_rules_campaign_id ON targeting_rules(campaign_id);`,
		`CREATE INDEX IF NOT EXISTS idx_targeting_rules_dimension ON targeting_rules(dimension);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

// SeedData inserts sample data for testing
func SeedData(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	//defer tx.Rollback()

	// Clear existing data
	if _, err := tx.Exec("DELETE FROM targeting_rules"); err != nil {
		return fmt.Errorf("failed to clear targeting_rules: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM campaigns"); err != nil {
		return fmt.Errorf("failed to clear campaigns: %w", err)
	}

	// Insert sample campaigns
	campaigns := []struct {
		id       string
		name     string
		imageURL string
		cta      string
		status   string
	}{
		{"spotify", "Spotify - Music for everyone", "https://somelink", "Download", "ACTIVE"},
		{"duolingo", "Duolingo: Best way to learn", "https://somelink2", "Install", "ACTIVE"},
		{"subwaysurfer", "Subway Surfer", "https://somelink3", "Play", "ACTIVE"},
	}

	for _, campaign := range campaigns {
		Result, err := tx.Exec(
			"INSERT INTO campaigns (id, name, image_url, cta, status) VALUES ($1, $2, $3, $4, $5);",
			campaign.id, campaign.name, campaign.imageURL, campaign.cta, campaign.status,
		)
		row, _ := Result.RowsAffected()
		fmt.Println("RESULT ", row)
		if err != nil {
			return fmt.Errorf("failed to insert campaign %s: %w", campaign.id, err)
		}
	}

	// Insert sample targeting rules
	rules := []struct {
		campaignID string
		ruleType   string
		dimension  string
		values     []string
	}{
		{"spotify", "INCLUDE", "COUNTRY", []string{"us", "canada"}},
		{"duolingo", "INCLUDE", "OS", []string{"android", "ios"}},
		{"duolingo", "EXCLUDE", "COUNTRY", []string{"us"}},
		{"subwaysurfer", "INCLUDE", "OS", []string{"android"}},
		{"subwaysurfer", "INCLUDE", "APP", []string{"com.gametion.ludokinggame"}},
	}

	for _, rule := range rules {
		_, err := tx.Exec(
			`INSERT INTO targeting_rules (campaign_id, rule_type, dimension, "values") VALUES ($1, $2, $3, $4);`,
			rule.campaignID, rule.ruleType, rule.dimension, pq.Array(rule.values),
		)
		if err != nil {
			return fmt.Errorf("failed to insert targeting rule for campaign %s: %w", rule.campaignID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
