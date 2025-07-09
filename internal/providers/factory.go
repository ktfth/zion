package providers

import (
	"fmt"
	"os"

	"github.com/ktfth/zion/internal/core"
)

// ProviderFactory creates AI providers based on configuration
type ProviderFactory struct{}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

// CreateProvider creates an AI provider based on the provider name and configuration
func (f *ProviderFactory) CreateProvider(providerName string, config *core.ProviderConfig) (core.AIProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid provider config: %w", err)
	}

	switch providerName {
	case "gemini":
		return NewGeminiProvider(config.APIKey, config.Model), nil
	case "openai", "gpt":
		return NewOpenAIProvider(config.APIKey, config.Model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// GetDefaultProvider returns the default provider based on available API keys
func (f *ProviderFactory) GetDefaultProvider() (core.AIProvider, error) {
	// Check for available API keys in order of preference
	if geminiKey := os.Getenv("GEMINI_API_KEY"); geminiKey != "" {
		config := &core.ProviderConfig{
			APIKey: geminiKey,
			Model:  "gemini-2.0-flash",
		}
		return f.CreateProvider("gemini", config)
	}

	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		config := &core.ProviderConfig{
			APIKey: openaiKey,
			Model:  "gpt-3.5-turbo",
		}
		return f.CreateProvider("openai", config)
	}

	return nil, fmt.Errorf("no AI provider configured - please set GEMINI_API_KEY or OPENAI_API_KEY")
}
