package providers

import (
	"fmt"
	"sync"
)

// DefaultManager is the default provider manager instance
var DefaultManager = NewProviderManager()

// ProviderManager gerencia o registro e seleção de provedores de IA
type ProviderManager struct {
	factories map[string]Factory
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewProviderManager cria uma nova instância do gerenciador
func NewProviderManager() *ProviderManager {
	manager := &ProviderManager{
		factories: make(map[string]Factory),
		providers: make(map[string]Provider),
	}

	// Register default providers
	manager.RegisterFactory("gemini", NewGeminiProvider)
	manager.RegisterFactory("gpt", NewGPTProvider)
	manager.RegisterFactory("openrouter", NewOpenRouterProvider)

	return manager
}

// RegisterFactory registra uma nova fábrica de provedores
func (m *ProviderManager) RegisterFactory(name string, factory Factory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[name] = factory
}

// GetProvider retorna um provedor existente ou cria um novo usando a fábrica apropriada
func (m *ProviderManager) GetProvider(name string, config map[string]string) (Provider, error) {
	m.mu.RLock()
	if provider, exists := m.providers[name]; exists {
		m.mu.RUnlock()
		return provider, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	factory, exists := m.factories[name]
	if !exists {
		return nil, fmt.Errorf("provedor não suportado: %s", name)
	}

	provider, err := factory(config)
	if err != nil {
		return nil, err
	}

	m.providers[name] = provider
	return provider, nil
}

// GetDefaultProvider retorna o provedor padrão (Gemini)
func (m *ProviderManager) GetDefaultProvider(config map[string]string) (Provider, error) {
	return m.GetProvider("gemini", config)
}

// init ensures the DefaultManager is properly initialized with all providers
func init() {
	// DefaultManager is already initialized in the package var declaration
	// Just register the default providers here
	DefaultManager.RegisterFactory("gemini", NewGeminiProvider)
	DefaultManager.RegisterFactory("gpt", NewGPTProvider)
	DefaultManager.RegisterFactory("openrouter", NewOpenRouterProvider)
}
