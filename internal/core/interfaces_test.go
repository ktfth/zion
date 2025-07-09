package core

import (
	"testing"
	"time"
)

func TestProviderConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ProviderConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &ProviderConfig{
				APIKey:      "test-key",
				Model:       "test-model",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: &ProviderConfig{
				Model:       "test-model",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
			wantErr: true,
		},
		{
			name: "empty api key",
			config: &ProviderConfig{
				APIKey:      "",
				Model:       "test-model",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("DefaultRetryConfig().MaxRetries = %v, want %v", config.MaxRetries, 3)
	}

	if config.BaseDelay != 2*time.Second {
		t.Errorf("DefaultRetryConfig().BaseDelay = %v, want %v", config.BaseDelay, 2*time.Second)
	}
}
