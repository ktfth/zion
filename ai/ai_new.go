package ai

import (
	"fmt"

	"github.com/ktfth/zion/ai/providers"
	"github.com/ktfth/zion/config"
	"github.com/ktfth/zion/plugins"
)

type ScaffoldResponse struct {
	Structure struct {
		Directories []string               `json:"directories"`
		Files       map[string]interface{} `json:"files"`
	} `json:"structure"`
}

// GenerateProjectScaffolding gera uma estrutura de projeto com base na linguagem, nome e descrição fornecidos
func GenerateProjectScaffolding(language, projectName, description string, registeredPlugins []string) (string, error) {
	// Substituir SamplePlugin por CorePlugin em registeredPlugins
	for i, plugin := range registeredPlugins {
		if plugin == "SamplePlugin" {
			registeredPlugins[i] = "CorePlugin"
		}
	}

	// Criar o contexto de scaffold para os plugins
	ctx := &plugins.ScaffoldContext{
		ProjectName: projectName,
		Language:    language,
		Description: description,
	}

	// Executar o hook BeforeGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.BeforeGeneration, ctx)

	// Construir o prompt usando o novo sistema melhorado
	prompt := buildImprovedPrompt(language, projectName, description)

	// Executar o hook ModifyPrompt para todos os plugins
	ctx.Prompt = prompt
	ctx = plugins.ExecuteHook(plugins.ModifyPrompt, ctx)
	prompt = ctx.Prompt

	// Obter o provedor de IA configurado
	cfg := config.LoadConfig()
	provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, cfg.GetAIConfig())
	if err != nil {
		return "", fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Gerar conteúdo usando o provedor
	fmt.Printf("🤖 Usando provedor de IA: %s\n", provider.Name())
	response, err := provider.GenerateContent(prompt)
	if err != nil {
		return "", err
	}

	// Atualizar a resposta no contexto
	ctx.Response = response

	// Executar o hook AfterGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.AfterGeneration, ctx)

	// Obter a resposta possivelmente modificada pelos plugins
	response = ctx.Response

	// Processar a resposta antes de retornar
	if response != "" {
		processedResponse, err := processScaffoldResponse(response)
		if err != nil {
			return "", fmt.Errorf("erro ao processar resposta: %v", err)
		}
		response = processedResponse
	}

	return response, nil
}

// GenerateProjectScaffoldingWithProvider gera uma estrutura de projeto com configurações customizadas de provider
func GenerateProjectScaffoldingWithProvider(language, projectName, description string, registeredPlugins []string, providerName, apiKey, model string) (string, error) {
	// Substituir SamplePlugin por CorePlugin em registeredPlugins
	for i, plugin := range registeredPlugins {
		if plugin == "SamplePlugin" {
			registeredPlugins[i] = "CorePlugin"
		}
	}

	// Criar o contexto de scaffold para os plugins
	ctx := &plugins.ScaffoldContext{
		ProjectName: projectName,
		Language:    language,
		Description: description,
	}

	// Executar o hook BeforeGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.BeforeGeneration, ctx)

	// Construir o prompt usando o novo sistema melhorado
	prompt := buildImprovedPrompt(language, projectName, description)

	// Executar o hook ModifyPrompt para todos os plugins
	ctx.Prompt = prompt
	ctx = plugins.ExecuteHook(plugins.ModifyPrompt, ctx)
	prompt = ctx.Prompt

	// Obter configurações do provider
	cfg := config.LoadConfig()

	// Usar provider customizado se especificado
	if providerName != "" {
		cfg.AIProvider = providerName
	}

	// Obter configuração do AI
	aiConfig := cfg.GetAIConfig()

	// Override da API key se especificado
	if apiKey != "" {
		aiConfig["api_key"] = apiKey
	}

	// Override do modelo se especificado
	if model != "" {
		aiConfig["model"] = model
	}

	// Criar o provider
	provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, aiConfig)
	if err != nil {
		return "", fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Gerar conteúdo usando o provedor
	fmt.Printf("🤖 Usando provedor de IA: %s\n", provider.Name())
	response, err := provider.GenerateContent(prompt)
	if err != nil {
		return "", err
	}

	// Atualizar a resposta no contexto
	ctx.Response = response

	// Executar o hook AfterGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.AfterGeneration, ctx)

	// Obter a resposta possivelmente modificada pelos plugins
	response = ctx.Response

	// Processar a resposta antes de retornar
	if response != "" {
		processedResponse, err := processScaffoldResponse(response)
		if err != nil {
			return "", fmt.Errorf("erro ao processar resposta: %v", err)
		}
		response = processedResponse
	}

	return response, nil
}

// Rest of the existing functions remain the same...
