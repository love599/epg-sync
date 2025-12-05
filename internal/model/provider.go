package model

import "time"

type ProviderConfig struct {
	ID         string         `yaml:"id"`
	Name       string         `yaml:"name"`
	BaseURL    string         `yaml:"base_url"`
	Enabled    bool           `yaml:"enabled"`
	Priority   int            `yaml:"priority"`
	Timeout    time.Duration  `yaml:"timeout"`
	RateLimit  int            `yaml:"rate_limit"`
	MaxRetries int            `yaml:"max_retries"`
	Settings   map[string]any `yaml:"settings"`
}
type ProviderHealth struct {
	ProviderID string    `json:"provider_id"`
	Healthy    bool      `json:"healthy"`
	Message    string    `json:"message"`
	CheckedAt  time.Time `json:"checked_at"`
	Latency    int64     `json:"latency"`
	StatusCode int       `json:"status_code"`
}
