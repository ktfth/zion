// Package providers fornece uma interface comum para diferentes provedores de IA
package providers

// Provider define a interface comum para todos os provedores de IA
type Provider interface {
	// Name retorna o nome do provedor
	Name() string

	// GenerateContent faz a chamada à API do provedor e retorna a resposta
	GenerateContent(prompt string) (string, error)
}

// Factory é uma função que cria um novo provedor baseado em configurações
type Factory func(config map[string]string) (Provider, error)
