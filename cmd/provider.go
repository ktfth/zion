package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktfth/zion/ai/providers"
	"github.com/ktfth/zion/config"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Gerencia provedores de IA",
	Long:  `Gerencia e configura provedores de IA como Gemini, GPT e OpenRouter.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os provedores disponíveis",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🤖 Provedores de IA disponíveis:\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Lista dos provedores disponíveis
		providersList := []string{"gemini", "gpt", "openrouter"}

		cfg := config.LoadConfig()

		for _, providerName := range providersList {
			status := "❌ Não configurado"

			// Verificar se o provider está configurado
			switch providerName {
			case "gemini":
				if cfg.GeminiAPIKey != "" {
					status = "✅ Configurado"
					if cfg.AIProvider == "gemini" {
						status += " (Ativo)"
					}
				}
			case "gpt":
				if cfg.OpenAIAPIKey != "" {
					status = "✅ Configurado"
					if cfg.AIProvider == "gpt" {
						status += " (Ativo)"
					}
				}
			case "openrouter":
				if cfg.OpenRouterAPIKey != "" {
					status = "✅ Configurado"
					if cfg.AIProvider == "openrouter" {
						status += " (Ativo)"
					}
				}
			}

			fmt.Printf("• %s: %s\n", strings.ToUpper(providerName), status)
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("💡 Para configurar um provider:\n")
		fmt.Printf("   zion provider config <provider>\n")
		fmt.Printf("   zion provider test <provider>\n\n")
	},
}

var configCmd = &cobra.Command{
	Use:   "config [provider]",
	Short: "Mostra informações de configuração para um provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName := strings.ToLower(args[0])

		fmt.Printf("🔧 Configuração do provider %s:\n", strings.ToUpper(providerName))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		switch providerName {
		case "gemini":
			fmt.Printf("Variáveis de ambiente necessárias:\n")
			fmt.Printf("• GEMINI_API_KEY (obrigatória)\n")
			fmt.Printf("• ZION_AI_PROVIDER=gemini (para usar como padrão)\n\n")
			fmt.Printf("Exemplo de configuração:\n")
			fmt.Printf("export GEMINI_API_KEY=\"your-api-key-here\"\n")
			fmt.Printf("export ZION_AI_PROVIDER=\"gemini\"\n\n")
			fmt.Printf("Como obter a API key:\n")
			fmt.Printf("1. Acesse https://makersuite.google.com/app/apikey\n")
			fmt.Printf("2. Crie uma nova API key\n")
			fmt.Printf("3. Configure a variável de ambiente\n")

		case "gpt":
			fmt.Printf("Variáveis de ambiente necessárias:\n")
			fmt.Printf("• OPENAI_API_KEY (obrigatória)\n")
			fmt.Printf("• ZION_AI_PROVIDER=gpt (para usar como padrão)\n\n")
			fmt.Printf("Variáveis opcionais:\n")
			fmt.Printf("• OPENAI_MODEL (padrão: gpt-3.5-turbo)\n")
			fmt.Printf("• OPENAI_MAX_TOKENS (padrão: 2048)\n")
			fmt.Printf("• OPENAI_TEMPERATURE (padrão: 0.7)\n\n")
			fmt.Printf("Exemplo de configuração:\n")
			fmt.Printf("export OPENAI_API_KEY=\"sk-your-api-key-here\"\n")
			fmt.Printf("export ZION_AI_PROVIDER=\"gpt\"\n")
			fmt.Printf("export OPENAI_MODEL=\"gpt-4\"\n\n")
			fmt.Printf("Como obter a API key:\n")
			fmt.Printf("1. Acesse https://platform.openai.com/api-keys\n")
			fmt.Printf("2. Crie uma nova API key\n")
			fmt.Printf("3. Configure a variável de ambiente\n")

		case "openrouter":
			fmt.Printf("Variáveis de ambiente necessárias:\n")
			fmt.Printf("• OPENROUTER_API_KEY (obrigatória)\n")
			fmt.Printf("• ZION_AI_PROVIDER=openrouter (para usar como padrão)\n\n")
			fmt.Printf("Variáveis opcionais:\n")
			fmt.Printf("• OPENROUTER_MODEL (padrão: meta-llama/llama-3.2-3b-instruct:free)\n")
			fmt.Printf("• OPENROUTER_MAX_TOKENS (padrão: 2048)\n")
			fmt.Printf("• OPENROUTER_TEMPERATURE (padrão: 0.7)\n")
			fmt.Printf("• OPENROUTER_BASE_URL (padrão: https://openrouter.ai/api/v1)\n\n")
			fmt.Printf("Exemplo de configuração:\n")
			fmt.Printf("export OPENROUTER_API_KEY=\"sk-or-v1-your-api-key-here\"\n")
			fmt.Printf("export ZION_AI_PROVIDER=\"openrouter\"\n")
			fmt.Printf("export OPENROUTER_MODEL=\"anthropic/claude-3-haiku\"\n\n")
			fmt.Printf("Como obter a API key:\n")
			fmt.Printf("1. Acesse https://openrouter.ai/keys\n")
			fmt.Printf("2. Crie uma conta ou faça login\n")
			fmt.Printf("3. Crie uma nova API key\n")
			fmt.Printf("4. Configure a variável de ambiente\n\n")
			fmt.Printf("Modelos populares disponíveis:\n")
			fmt.Printf("• meta-llama/llama-3.2-3b-instruct:free (gratuito)\n")
			fmt.Printf("• anthropic/claude-3-haiku\n")
			fmt.Printf("• openai/gpt-3.5-turbo\n")
			fmt.Printf("• google/gemini-pro\n")

		default:
			fmt.Printf("❌ Provider não reconhecido: %s\n", providerName)
			fmt.Printf("Providers disponíveis: gemini, gpt, openrouter\n")
			os.Exit(1)
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	},
}

var testCmd = &cobra.Command{
	Use:   "test [provider]",
	Short: "Testa a conexão com um provider específico",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName := strings.ToLower(args[0])

		fmt.Printf("🧪 Testando provider %s...\n", strings.ToUpper(providerName))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Carregar configuração
		cfg := config.LoadConfig()

		// Criar configuração específica do provider
		var aiConfig map[string]string
		switch providerName {
		case "gemini":
			aiConfig = map[string]string{
				"api_key": cfg.GeminiAPIKey,
			}
		case "gpt":
			aiConfig = map[string]string{
				"api_key":     cfg.OpenAIAPIKey,
				"model":       os.Getenv("OPENAI_MODEL"),
				"max_tokens":  os.Getenv("OPENAI_MAX_TOKENS"),
				"temperature": os.Getenv("OPENAI_TEMPERATURE"),
			}
		case "openrouter":
			aiConfig = map[string]string{
				"api_key":     cfg.OpenRouterAPIKey,
				"model":       os.Getenv("OPENROUTER_MODEL"),
				"max_tokens":  os.Getenv("OPENROUTER_MAX_TOKENS"),
				"temperature": os.Getenv("OPENROUTER_TEMPERATURE"),
				"base_url":    os.Getenv("OPENROUTER_BASE_URL"),
			}
		default:
			fmt.Printf("❌ Provider não reconhecido: %s\n", providerName)
			os.Exit(1)
		}

		// Verificar se a API key está configurada
		if aiConfig["api_key"] == "" {
			fmt.Printf("❌ API key não configurada para %s\n", strings.ToUpper(providerName))
			fmt.Printf("💡 Execute: zion provider config %s\n", providerName)
			os.Exit(1)
		}

		// Criar o provider
		provider, err := providers.DefaultManager.GetProvider(providerName, aiConfig)
		if err != nil {
			fmt.Printf("❌ Erro ao criar provider: %v\n", err)
			os.Exit(1)
		}

		// Testar com um prompt simples
		testPrompt := "Responda apenas 'OK' para confirmar que está funcionando."

		fmt.Printf("🔄 Enviando prompt de teste...\n")
		response, err := provider.GenerateContent(testPrompt)
		if err != nil {
			fmt.Printf("❌ Erro no teste: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Teste bem-sucedido!\n")
		fmt.Printf("🤖 Resposta: %s\n", strings.TrimSpace(response))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("💡 Provider %s está funcionando corretamente!\n", strings.ToUpper(providerName))
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Mostra o provider atualmente configurado",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		fmt.Printf("🤖 Provider atual: %s\n", strings.ToUpper(cfg.AIProvider))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Mostrar status do provider atual
		switch cfg.AIProvider {
		case "gemini":
			if cfg.GeminiAPIKey != "" {
				fmt.Printf("✅ Status: Configurado\n")
				fmt.Printf("🔑 API Key: %s***\n", cfg.GeminiAPIKey[:min(8, len(cfg.GeminiAPIKey))])
			} else {
				fmt.Printf("❌ Status: Não configurado\n")
			}
		case "gpt":
			if cfg.OpenAIAPIKey != "" {
				fmt.Printf("✅ Status: Configurado\n")
				fmt.Printf("🔑 API Key: %s***\n", cfg.OpenAIAPIKey[:min(8, len(cfg.OpenAIAPIKey))])
				if model := os.Getenv("OPENAI_MODEL"); model != "" {
					fmt.Printf("🎯 Modelo: %s\n", model)
				}
			} else {
				fmt.Printf("❌ Status: Não configurado\n")
			}
		case "openrouter":
			if cfg.OpenRouterAPIKey != "" {
				fmt.Printf("✅ Status: Configurado\n")
				fmt.Printf("🔑 API Key: %s***\n", cfg.OpenRouterAPIKey[:min(8, len(cfg.OpenRouterAPIKey))])
				if model := os.Getenv("OPENROUTER_MODEL"); model != "" {
					fmt.Printf("🎯 Modelo: %s\n", model)
				}
			} else {
				fmt.Printf("❌ Status: Não configurado\n")
			}
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("💡 Para alterar o provider:\n")
		fmt.Printf("   export ZION_AI_PROVIDER=\"provider_name\"\n")
		fmt.Printf("   Providers disponíveis: gemini, gpt, openrouter\n")
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	// Adicionar subcomandos ao comando provider
	providerCmd.AddCommand(listCmd)
	providerCmd.AddCommand(configCmd)
	providerCmd.AddCommand(testCmd)
	providerCmd.AddCommand(currentCmd)

	// Registrar o comando provider no comando raiz
	rootCmd.AddCommand(providerCmd)
}
