package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ktfth/zion/ai"
	"github.com/spf13/cobra"
)

// testUltimateGoalCmd representa o comando para testar o Ultimate Goal Controller
var testUltimateGoalCmd = &cobra.Command{
	Use:   "test-ultimate-goal",
	Short: "Testa o sistema de Ultimate Goal Focus",
	Long: `Testa o sistema de Ultimate Goal Focus que condiciona a geração 
baseada no objetivo final do prompt, eliminando arquivos e recursos desnecessários.

Exemplos:
  zion test-ultimate-goal "criar uma API REST simples apenas para CRUD de usuários"
  zion test-ultimate-goal "desenvolver um CLI básico para gerenciar arquivos"
  zion test-ultimate-goal "implementar somente um servidor HTTP mínimo"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Erro: Descrição do projeto é obrigatória")
			fmt.Println("Uso: zion test-ultimate-goal \"sua descrição aqui\"")
			os.Exit(1)
		}

		description := args[0]

		fmt.Printf("🎯 TESTE DO SISTEMA ULTIMATE GOAL FOCUS\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		fmt.Printf("📝 Descrição do projeto: %s\n\n", description)

		// 1. Testar o Ultimate Goal Controller
		fmt.Printf("1️⃣ Testando análise do objetivo...\n")

		goalController := ai.NewUltimateGoalController(description)

		fmt.Printf("   🎯 Objetivo identificado: %s\n", goalController.Goal)
		fmt.Printf("   🧠 Intenção do usuário: %s\n", goalController.Intent)
		fmt.Printf("   📐 Escopo determinado: %s\n", goalController.Scope)
		fmt.Printf("   🔥 Prioridade: %d/10\n", goalController.Priority)
		fmt.Printf("   🔍 Palavras-chave: %v\n", goalController.Keywords)

		if len(goalController.RequiredFiles) > 0 {
			fmt.Printf("   ✅ Arquivos obrigatórios: %v\n", goalController.RequiredFiles)
		}

		if len(goalController.RequiredDirs) > 0 {
			fmt.Printf("   ✅ Diretórios obrigatórios: %v\n", goalController.RequiredDirs)
		}

		if len(goalController.ExcludedFiles) > 0 {
			fmt.Printf("   🚫 Arquivos excluídos: %v\n", goalController.ExcludedFiles)
		}

		if len(goalController.ExcludedDirs) > 0 {
			fmt.Printf("   🚫 Diretórios excluídos: %v\n", goalController.ExcludedDirs)
		}

		fmt.Printf("\n2️⃣ Testando geração do prompt focado no objetivo...\n")

		// 2. Testar geração de prompt
		basePrompt := "Crie um projeto com a seguinte estrutura..."
		focusedPrompt := goalController.BuildGoalFocusedPrompt(basePrompt)

		fmt.Printf("   📏 Tamanho do prompt: %d caracteres\n", len(focusedPrompt))
		fmt.Printf("   🎯 Prompt contém foco no objetivo: %v\n",
			contains(focusedPrompt, "ULTIMATE GOAL FOCUS"))
		fmt.Printf("   📋 Prompt contém regras críticas: %v\n",
			contains(focusedPrompt, "REGRAS CRÍTICAS"))

		// 3. Testar análise de conformidade
		fmt.Printf("\n3️⃣ Testando análise de conformidade...\n")

		goalAnalysis := goalController.GetGoalSummary()
		fmt.Printf("   📊 Confiança na análise: %.1f%%\n", goalAnalysis.Confidence*100)
		fmt.Printf("   🎯 Objetivo principal: %s\n", goalAnalysis.PrimaryGoal)
		fmt.Printf("   📁 Componentes-chave: %v\n", goalAnalysis.KeyComponents)

		// 4. Testar filtro de conteúdo
		fmt.Printf("\n4️⃣ Testando filtro de conteúdo...\n")

		// Criar conteúdo de exemplo
		sampleContent := `{
			"structure": {
				"directories": ["src", "tests", "docs", "docker", "examples"],
				"files": {
					"main.go": {"content": "package main"},
					"database.go": {"content": "package database"},
					"docker-compose.yml": {"content": "version: '3'"},
					"README.md": {"content": "# Project"},
					"test_helpers.go": {"content": "package test"},
					"swagger.yaml": {"content": "swagger: '2.0'"}
				}
			}
		}`

		filteredContent, err := goalController.FilterGeneratedContent(sampleContent)
		if err != nil {
			fmt.Printf("   ❌ Erro no filtro: %v\n", err)
		} else {
			// Contar arquivos antes e depois
			var original, filtered map[string]interface{}
			json.Unmarshal([]byte(sampleContent), &original)
			json.Unmarshal([]byte(filteredContent), &filtered)

			originalFiles := len(original["structure"].(map[string]interface{})["files"].(map[string]interface{}))
			filteredFiles := len(filtered["structure"].(map[string]interface{})["files"].(map[string]interface{}))

			fmt.Printf("   📊 Arquivos antes: %d\n", originalFiles)
			fmt.Printf("   📊 Arquivos depois: %d\n", filteredFiles)
			fmt.Printf("   🎯 Filtro aplicado: %v\n", originalFiles != filteredFiles)
		}

		// 5. Testar integração com Adaptive Controller
		fmt.Printf("\n5️⃣ Testando integração com Adaptive Controller...\n")

		adaptiveController := ai.NewAdaptiveInstructionController("api", "go", description)

		// Testar prompt adaptativo
		adaptivePrompt := adaptiveController.BuildAdaptivePrompt(basePrompt)
		fmt.Printf("   📏 Tamanho do prompt adaptativo: %d caracteres\n", len(adaptivePrompt))
		fmt.Printf("   🎯 Integração com Ultimate Goal: %v\n",
			contains(adaptivePrompt, "ULTIMATE GOAL FOCUS"))

		// Testar validação
		isCompliant, issues := adaptiveController.ValidateGoalCompliance(sampleContent)
		fmt.Printf("   ✅ Conteúdo está em conformidade: %v\n", isCompliant)
		if !isCompliant {
			fmt.Printf("   ⚠️  Problemas encontrados: %d\n", len(issues))
			for _, issue := range issues {
				fmt.Printf("      • %s\n", issue)
			}
		}

		// 6. Testar análise final
		fmt.Printf("\n6️⃣ Resumo final da análise...\n")

		if analysis := adaptiveController.GetGoalAnalysis(); analysis != nil {
			analysisJson, _ := json.MarshalIndent(analysis, "   ", "  ")
			fmt.Printf("   📋 Análise completa:\n%s\n", string(analysisJson))
		}

		fmt.Printf("\n🎉 TESTE CONCLUÍDO COM SUCESSO!\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		fmt.Printf("\n💡 Para usar o sistema em produção:\n")
		fmt.Printf("   zion scaffold -l go -n meu-projeto -d \"%s\"\n", description)
		fmt.Printf("   O sistema automaticamente aplicará o filtro de Ultimate Goal!\n")
	},
}

func contains(text, substring string) bool {
	return len(text) > 0 && len(substring) > 0 &&
		(text == substring || findSubstring(text, substring))
}

func findSubstring(text, substring string) bool {
	if len(substring) > len(text) {
		return false
	}

	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(testUltimateGoalCmd)
}
