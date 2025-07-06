package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ktfth/zion/ai"
	"github.com/ktfth/zion/plugins"

	"github.com/spf13/cobra"
)

var language string
var projectName string
var description string
var aiProvider string
var apiKey string
var model string
var maxRetries int

// scaffoldCmd define o comando "scaffold".
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Gera a estrutura de um projeto com a ajuda de AI",
	Long: `Gera a estrutura completa de um projeto utilizando inteligência artificial.
	
O comando inclui um sistema de retry inteligente que tenta até 3 vezes (configurável) 
caso a geração inicial falhe, garantindo maior robustez na criação dos projetos.

Exemplos:
  zion scaffold -l go -n meu-projeto -d "API REST com PostgreSQL"
  zion scaffold -l typescript -n webapp -d "App React com TypeScript" -r 5
  zion scaffold -l python -n ml-project -d "Projeto de Machine Learning" --retries 2`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		fmt.Printf("\n🚀 Iniciando geração do projeto\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📦 Projeto: %s\n", projectName)
		fmt.Printf("🔧 Linguagem: %s\n", language)
		fmt.Printf("📝 Descrição: %s\n", description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Lista plugins ativos
		pluginsList := plugins.ListPlugins()
		if len(pluginsList) > 0 {
			fmt.Printf("🔌 Plugins ativos: %v\n\n", pluginsList)
		}

		fmt.Print("🤖 Gerando estrutura com IA...")

		// Preparar configurações do provider
		var response string
		var err error
		var attempt int

		// Lógica de retry para geração do projeto
		for attempt = 1; attempt <= maxRetries; attempt++ {
			if attempt > 1 {
				fmt.Printf("\n🔄 Tentativa %d/%d...", attempt, maxRetries)
			}

			if aiProvider != "" || apiKey != "" || model != "" {
				// Usar configurações customizadas
				response, err = ai.GenerateProjectScaffoldingWithProvider(language, projectName, description, pluginsList, aiProvider, apiKey, model)
			} else {
				// Usar configurações padrão
				response, err = ai.GenerateProjectScaffolding(language, projectName, description, pluginsList)
			}

			if err == nil {
				break
			}

			if attempt < maxRetries {
				fmt.Printf("\n⚠️  Erro na tentativa %d: %v", attempt, err)
				fmt.Printf("\n🔄 Tentando novamente em 2 segundos...")
				time.Sleep(2 * time.Second)
			}
		}

		if err != nil {
			fmt.Printf("\n❌ Falha após %d tentativas:\n%v\n", maxRetries, err)
			if response != "" {
				fmt.Printf("\nResposta da API:\n%s\n", response)
			}
			os.Exit(1)
		}
		fmt.Println(" ✅")

		fmt.Print("📂 Criando estrutura do projeto...")

		// Lógica de retry para criação da estrutura
		var createAttempt int
		for createAttempt = 1; createAttempt <= maxRetries; createAttempt++ {
			if createAttempt > 1 {
				fmt.Printf("\n🔄 Tentativa de criação %d/%d...", createAttempt, maxRetries)
			}

			err = ai.ExtractAndCreateProject(projectName, response)
			if err == nil {
				break
			}

			if createAttempt < maxRetries {
				fmt.Printf("\n⚠️  Erro na tentativa %d: %v", createAttempt, err)
				fmt.Printf("\n🔄 Tentando novamente em 1 segundo...")
				time.Sleep(1 * time.Second)
			}
		}

		if err != nil {
			fmt.Printf("\n⚠️  Erro ao criar estrutura padrão após %d tentativas, tentando método alternativo...\n", maxRetries)
			err = ai.SaveRawResponse(projectName, response)
			if err != nil {
				fmt.Printf("\n❌ Erro ao salvar resposta:\n%v\n", err)
				os.Exit(1)
			}
			fmt.Println("💡 Resposta salva em README.md no diretório do projeto.")
		}
		fmt.Println(" ✅")

		// Executa plugins
		if len(pluginsList) > 0 {
			fmt.Print("🔌 Executando plugins...")
			plugins.ExecutePlugins()
			fmt.Println(" ✅")
		}

		elapsedTime := time.Since(startTime)

		fmt.Printf("\n✨ Projeto criado com sucesso! ✨\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📁 Local: %s\n", projectName)
		fmt.Printf("⏱️  Tempo total: %.2f segundos\n", elapsedTime.Seconds())
		if attempt > 1 {
			fmt.Printf("🔄 Tentativas de geração: %d/%d\n", attempt, maxRetries)
		}
		if createAttempt > 1 {
			fmt.Printf("🔄 Tentativas de criação: %d/%d\n", createAttempt, maxRetries)
		}
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		fmt.Printf("💡 Para começar a desenvolver:\n")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   Consulte o README.md para instruções detalhadas\n\n")
	},
}

func init() {
	// Configura flags para o comando scaffold
	scaffoldCmd.Flags().StringVarP(&language, "language", "l", "", "Linguagem para o scaffold (ex: go, python, etc)")
	scaffoldCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	scaffoldCmd.Flags().StringVarP(&description, "description", "d", "", "Descrição objetiva da estrutura desejada")
	scaffoldCmd.Flags().StringVarP(&aiProvider, "provider", "p", "", "Provider de IA (gemini, gpt, openrouter)")
	scaffoldCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API Key do provider")
	scaffoldCmd.Flags().StringVarP(&model, "model", "m", "", "Modelo específico do provider")
	scaffoldCmd.Flags().IntVarP(&maxRetries, "retries", "r", 3, "Número máximo de tentativas em caso de falha (padrão: 3)")
	scaffoldCmd.MarkFlagRequired("language")
	scaffoldCmd.MarkFlagRequired("name")

	// Registra o comando scaffold no comando raiz
	rootCmd.AddCommand(scaffoldCmd)
}
