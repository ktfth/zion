package providers

import (
	"testing"

	"github.com/ktfth/zion/internal/core"
)

func TestProviderFactory_CreateProvider(t *testing.T) {
	factory := NewProviderFactory()

	tests := []struct {
		name         string
		providerName string
		config       *core.ProviderConfig
		wantErr      bool
	}{
		{
			name:         "valid gemini provider",
			providerName: "gemini",
			config: &core.ProviderConfig{
				APIKey: "test-key",
				Model:  "gemini-2.0-flash",
			},
			wantErr: false,
		},
		{
			name:         "valid openai provider",
			providerName: "openai",
			config: &core.ProviderConfig{
				APIKey: "test-key",
				Model:  "gpt-3.5-turbo",
			},
			wantErr: false,
		},
		{
			name:         "invalid provider name",
			providerName: "invalid",
			config: &core.ProviderConfig{
				APIKey: "test-key",
			},
			wantErr: true,
		},
		{
			name:         "missing api key",
			providerName: "gemini",
			config: &core.ProviderConfig{
				Model: "gemini-2.0-flash",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := factory.CreateProvider(tt.providerName, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderFactory.CreateProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("ProviderFactory.CreateProvider() returned nil provider")
			}
		})
	}
}

func TestGeminiProvider_Name(t *testing.T) {
	provider := NewGeminiProvider("test-key", "gemini-2.0-flash")
	if provider.Name() != "Gemini" {
		t.Errorf("GeminiProvider.Name() = %v, want %v", provider.Name(), "Gemini")
	}
}

func TestOpenAIProvider_Name(t *testing.T) {
	provider := NewOpenAIProvider("test-key", "gpt-3.5-turbo")
	if provider.Name() != "OpenAI" {
		t.Errorf("OpenAIProvider.Name() = %v, want %v", provider.Name(), "OpenAI")
	}
}
