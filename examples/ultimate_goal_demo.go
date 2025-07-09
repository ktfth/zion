package main

import (
	"fmt"
	"log"

	"github.com/ktfth/zion/ai"
)

// Demonstração do Sistema Ultimate Goal Focus
func main() {
	fmt.Printf("🎯 DEMONSTRAÇÃO - SISTEMA ULTIMATE GOAL FOCUS\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Cenários de teste
	scenarios := []struct {
		name        string
		description string
		expected    string
	}{
		{
			name:        "API REST Mínima",
			description: "criar uma API REST simples apenas para CRUD de usuários",
			expected:    "Deve gerar apenas arquivos essenciais: main.go, handlers.go, models.go",
		},
		{
			name:        "CLI Básico",
			description: "desenvolver um CLI básico para gerenciar arquivos",
			expected:    "Deve focar em: main.go, cmd.go, cli.go - sem web interface",
		},
		{
			name:        "Servidor HTTP Mínimo",
			description: "implementar somente um servidor HTTP mínimo",
			expected:    "Deve conter apenas: main.go, server.go - sem middlewares complexos",
		},
		{
			name:        "Projeto Completo",
			description: "desenvolver uma aplicação web completa com banco de dados",
			expected:    "Deve incluir: frontend, backend, database, testes, docker",
		},
		{
			name:        "Protótipo Rápido",
			description: "fazer um protótipo rápido de dashboard",
			expected:    "Deve ser mínimo: index.html, app.js, style.css",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("%d️⃣ CENÁRIO: %s\n", i+1, scenario.name)
		fmt.Printf("   📝 Descrição: %s\n", scenario.description)
		fmt.Printf("   🎯 Esperado: %s\n", scenario.expected)
		fmt.Printf("   ⚡ Resultado:\n")

		// Criar controller para análise
		controller := ai.NewUltimateGoalController(scenario.description)

		// Mostrar análise
		fmt.Printf("      • Objetivo: %s\n", controller.Goal)
		fmt.Printf("      • Escopo: %s\n", controller.Scope)
		fmt.Printf("      • Prioridade: %d/10\n", controller.Priority)
		fmt.Printf("      • Arquivos obrigatórios: %v\n", controller.RequiredFiles)
		fmt.Printf("      • Arquivos excluídos: %v\n", controller.ExcludedFiles)

		// Testar geração de prompt
		basePrompt := "Crie um projeto Go com estrutura completa"
		focusedPrompt := controller.BuildGoalFocusedPrompt(basePrompt)

		fmt.Printf("      • Prompt focado: %d chars\n", len(focusedPrompt))
		fmt.Printf("      • Contém foco: %v\n", containsText(focusedPrompt, "ULTIMATE GOAL FOCUS"))

		// Testar filtro
		sampleContent := createSampleContent()
		filteredContent, err := controller.FilterGeneratedContent(sampleContent)
		if err != nil {
			fmt.Printf("      • Erro no filtro: %v\n", err)
		} else {
			originalCount := countFiles(sampleContent)
			filteredCount := countFiles(filteredContent)
			fmt.Printf("      • Arquivos: %d → %d (redução: %d%%)\n",
				originalCount, filteredCount,
				((originalCount-filteredCount)*100)/originalCount)
		}

		// Validar conformidade
		analysis := controller.GetGoalSummary()
		fmt.Printf("      • Confiança: %.1f%%\n", analysis.Confidence*100)

		fmt.Printf("   ✅ Análise concluída\n\n")
	}

	// Demonstração de integração com Adaptive Controller
	fmt.Printf("🔧 INTEGRAÇÃO COM ADAPTIVE CONTROLLER\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	testDescription := "criar uma API REST simples apenas para CRUD de usuários"
	adaptiveController := ai.NewAdaptiveInstructionController("api", "go", testDescription)

	// Testar prompt adaptativo
	basePrompt := "Crie um projeto Go profissional"
	adaptivePrompt := adaptiveController.BuildAdaptivePrompt(basePrompt)

	fmt.Printf("📋 Prompt adaptativo gerado:\n")
	fmt.Printf("   • Tamanho: %d caracteres\n", len(adaptivePrompt))
	fmt.Printf("   • Contém Ultimate Goal: %v\n", containsText(adaptivePrompt, "ULTIMATE GOAL"))
	fmt.Printf("   • Contém regras críticas: %v\n", containsText(adaptivePrompt, "REGRAS CRÍTICAS"))

	// Testar validação
	sampleContent := createSampleContent()
	isCompliant, issues := adaptiveController.ValidateGoalCompliance(sampleContent)

	fmt.Printf("   • Conformidade: %v\n", isCompliant)
	if !isCompliant {
		fmt.Printf("   • Problemas: %d\n", len(issues))
		for _, issue := range issues {
			fmt.Printf("     - %s\n", issue)
		}
	}

	// Mostrar análise final
	if analysis := adaptiveController.GetGoalAnalysis(); analysis != nil {
		fmt.Printf("   • Análise completa disponível: ✅\n")
		fmt.Printf("   • Objetivo: %s\n", analysis.PrimaryGoal)
		fmt.Printf("   • Componentes-chave: %v\n", analysis.KeyComponents)
	}

	fmt.Printf("\n🎉 DEMONSTRAÇÃO CONCLUÍDA COM SUCESSO!\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("\n💡 Para usar na prática:\n")
	fmt.Printf("   go run main.go\n")
	fmt.Printf("   zion scaffold -l go -n meu-projeto -d \"sua descrição específica\"\n")
	fmt.Printf("   zion test-ultimate-goal \"sua descrição aqui\"\n")
}

// containsText verifica se o texto contém a substring
func containsText(text, substring string) bool {
	return len(text) > len(substring) && findInText(text, substring)
}

// findInText encontra substring no texto
func findInText(text, substring string) bool {
	textLen := len(text)
	subLen := len(substring)

	if subLen > textLen {
		return false
	}

	for i := 0; i <= textLen-subLen; i++ {
		if text[i:i+subLen] == substring {
			return true
		}
	}

	return false
}

// createSampleContent cria conteúdo de exemplo para teste
func createSampleContent() string {
	return `{
		"structure": {
			"directories": ["src", "tests", "docs", "docker", "examples", "benchmarks"],
			"files": {
				"main.go": {"content": "package main"},
				"handlers.go": {"content": "package handlers"},
				"models.go": {"content": "package models"},
				"database.go": {"content": "package database"},
				"docker-compose.yml": {"content": "version: '3'"},
				"Dockerfile": {"content": "FROM golang"},
				"README.md": {"content": "# Project"},
				"test_helpers.go": {"content": "package test"},
				"benchmark_test.go": {"content": "package benchmark"},
				"swagger.yaml": {"content": "swagger: '2.0'"},
				"kubernetes.yaml": {"content": "apiVersion: v1"},
				"performance_test.go": {"content": "package performance"}
			}
		}
	}`
}

// countFiles conta o número de arquivos no JSON
func countFiles(content string) int {
	// Implementação simples para contar arquivos
	// Em produção, usaria JSON parsing
	count := 0
	for i := 0; i < len(content); i++ {
		if i < len(content)-8 && content[i:i+8] == "content\"" {
			count++
		}
	}
	return count
}

// init inicializa o exemplo
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
