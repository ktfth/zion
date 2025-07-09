package core

import (
	"fmt"
	"time"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	// Name returns the name of the AI provider
	Name() string

	// GenerateContent generates content based on the given prompt
	GenerateContent(prompt string) (string, error)
}

// ProviderConfig holds configuration for AI providers
type ProviderConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

// Validate checks if the provider configuration is valid
func (c *ProviderConfig) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		BaseDelay:  2 * time.Second,
	}
}
