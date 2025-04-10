package plugins

import (
	"fmt"
	"strings"
)

// HelloWorldPlugin é um plugin de exemplo que adiciona uma mensagem de boas-vindas
// e modifica o prompt para incluir requisitos adicionais
type HelloWorldPlugin struct{}

// Name retorna o nome do plugin
func (p *HelloWorldPlugin) Name() string {
	return "HelloWorld"
}

// Execute é chamado quando o plugin é executado
func (p *HelloWorldPlugin) Execute() error {
	fmt.Println("HelloWorld plugin: Olá, mundo!")
	return nil
}

// GetHooks retorna os hooks que o plugin implementa
func (p *HelloWorldPlugin) GetHooks() map[ScaffoldHook]interface{} {
	return map[ScaffoldHook]interface{}{
		BeforeGeneration: p.beforeGeneration,
		ModifyPrompt:     p.modifyPrompt,
		AfterGeneration:  p.afterGeneration,
	}
}

// beforeGeneration é executado antes da geração do scaffold
func (p *HelloWorldPlugin) beforeGeneration(ctx *ScaffoldContext) error {
	fmt.Printf("HelloWorld plugin: Iniciando geração para projeto '%s' em %s\n", 
		ctx.ProjectName, ctx.Language)
	return nil
}

// modifyPrompt modifica o prompt para incluir requisitos adicionais
func (p *HelloWorldPlugin) modifyPrompt(prompt string) string {
	// Adiciona requisitos específicos do HelloWorld plugin
	additionalRequirements := "\n\nAdicional do HelloWorld Plugin:\n" +
		"1. Adicione um arquivo hello.md com uma mensagem de boas-vindas\n" +
		"2. Inclua comentários explicativos no código\n"
	
	// Insere os requisitos adicionais antes da linha IMPORTANTE
	if idx := strings.Index(prompt, "IMPORTANTE:"); idx != -1 {
		return prompt[:idx] + additionalRequirements + prompt[idx:]
	}
	
	// Se não encontrar o marcador, apenas adiciona ao final
	return prompt + additionalRequirements
}

// afterGeneration é executado após a geração do scaffold
func (p *HelloWorldPlugin) afterGeneration(ctx *ScaffoldContext) error {
	fmt.Printf("HelloWorld plugin: Geração concluída para projeto '%s'\n", ctx.ProjectName)
	return nil
}

// Inicializa e registra o plugin automaticamente
func init() {
	RegisterPlugin(&HelloWorldPlugin{})
}
