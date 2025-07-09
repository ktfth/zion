package ai

import (
	"encoding/json"
	"testing"
)

func TestUltimateGoalController_AnalyzeGoal(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		wantGoal     string
		wantScope    string
		wantPriority int
		wantIntent   string
	}{
		{
			name:         "API REST Mínima",
			description:  "criar uma API REST simples apenas para CRUD de usuários",
			wantGoal:     "uma api rest simples apenas para crud de usuários",
			wantScope:    "minimal",
			wantPriority: 8,
			wantIntent:   "specific_minimal",
		},
		{
			name:         "CLI Básico",
			description:  "desenvolver um CLI básico para gerenciar arquivos",
			wantGoal:     "um cli básico para gerenciar arquivos",
			wantScope:    "minimal",
			wantPriority: 7,
			wantIntent:   "quick_essential",
		},
		{
			name:         "Projeto Completo",
			description:  "desenvolver uma aplicação web completa com banco de dados",
			wantGoal:     "uma aplicação web completa com banco de dados",
			wantScope:    "comprehensive",
			wantPriority: 4,
			wantIntent:   "comprehensive",
		},
		{
			name:         "Protótipo Rápido",
			description:  "fazer um protótipo rápido de dashboard",
			wantGoal:     "um protótipo rápido de dashboard",
			wantScope:    "minimal",
			wantPriority: 8,
			wantIntent:   "prototype",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewUltimateGoalController(tt.description)

			if controller.Goal != tt.wantGoal {
				t.Errorf("Goal = %v, want %v", controller.Goal, tt.wantGoal)
			}

			if controller.Scope != tt.wantScope {
				t.Errorf("Scope = %v, want %v", controller.Scope, tt.wantScope)
			}

			if controller.Priority != tt.wantPriority {
				t.Errorf("Priority = %v, want %v", controller.Priority, tt.wantPriority)
			}

			if controller.Intent != tt.wantIntent {
				t.Errorf("Intent = %v, want %v", controller.Intent, tt.wantIntent)
			}
		})
	}
}

func TestUltimateGoalController_BuildGoalFocusedPrompt(t *testing.T) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")
	basePrompt := "Crie um projeto Go"

	prompt := controller.BuildGoalFocusedPrompt(basePrompt)

	// Verificar se o prompt contém elementos esperados
	expectedElements := []string{
		"ULTIMATE GOAL FOCUS",
		"REGRAS CRÍTICAS",
		"OBJETIVO PRINCIPAL",
		"ELIMINE qualquer componente",
		"MANTENHA laser focus",
	}

	for _, element := range expectedElements {
		if !containsString(prompt, element) {
			t.Errorf("Expected element '%s' not found in prompt", element)
		}
	}

	// Verificar informações específicas
	if !containsString(prompt, controller.Goal) {
		t.Errorf("Goal not found in prompt: %s", controller.Goal)
	}

	if !containsString(prompt, controller.Scope) {
		t.Errorf("Scope not found in prompt: %s", controller.Scope)
	}
}

func TestUltimateGoalController_FilterGeneratedContent(t *testing.T) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")

	// Conteúdo de exemplo com arquivos desnecessários
	sampleContent := `{
		"structure": {
			"directories": ["src", "tests", "docs", "docker", "examples"],
			"files": {
				"main.go": {"content": "package main"},
				"handlers.go": {"content": "package handlers"},
				"models.go": {"content": "package models"},
				"docker-compose.yml": {"content": "version: '3'"},
				"Dockerfile": {"content": "FROM golang"},
				"swagger.yaml": {"content": "swagger: '2.0'"},
				"benchmark_test.go": {"content": "package benchmark"}
			}
		}
	}`

	filteredContent, err := controller.FilterGeneratedContent(sampleContent)
	if err != nil {
		t.Fatalf("FilterGeneratedContent failed: %v", err)
	}

	// Verificar se arquivos desnecessários foram removidos
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(filteredContent), &result); err != nil {
		t.Fatalf("Failed to parse filtered content: %v", err)
	}

	files := result["structure"].(map[string]interface{})["files"].(map[string]interface{})

	// Arquivos essenciais devem estar presentes
	essentialFiles := []string{"main.go", "handlers.go", "models.go"}
	for _, file := range essentialFiles {
		if _, exists := files[file]; !exists {
			t.Errorf("Essential file %s was filtered out", file)
		}
	}

	// Arquivos desnecessários devem ser removidos (se scope for minimal)
	if controller.Scope == "minimal" {
		unnecessaryFiles := []string{"docker-compose.yml", "swagger.yaml", "benchmark_test.go"}
		for _, file := range unnecessaryFiles {
			if _, exists := files[file]; exists {
				t.Errorf("Unnecessary file %s was not filtered out", file)
			}
		}
	}
}

func TestUltimateGoalController_ValidateCompliance(t *testing.T) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")

	// Teste com conteúdo não conforme
	nonCompliantContent := `{
		"structure": {
			"files": {
				"main.go": {"content": "package main"},
				"docker-compose.yml": {"content": "version: '3'"},
				"swagger.yaml": {"content": "swagger: '2.0'"}
			}
		}
	}`

	// Simular validação através do adaptive controller
	adaptiveController := NewAdaptiveInstructionController("api", "go", "criar uma API REST simples apenas para CRUD de usuários")

	isCompliant, issues := adaptiveController.ValidateGoalCompliance(nonCompliantContent)

	if isCompliant {
		t.Error("Expected non-compliant content to fail validation")
	}

	if len(issues) == 0 {
		t.Error("Expected validation issues to be found")
	}

	// Verificar se os problemas específicos foram detectados
	foundDockerIssue := false
	foundSwaggerIssue := false

	for _, issue := range issues {
		if containsString(issue, "docker-compose.yml") {
			foundDockerIssue = true
		}
		if containsString(issue, "swagger.yaml") {
			foundSwaggerIssue = true
		}
	}

	if !foundDockerIssue && controller.Scope == "minimal" {
		t.Error("Expected docker-compose.yml to be flagged as unnecessary")
	}

	if !foundSwaggerIssue && controller.Scope == "minimal" {
		t.Error("Expected swagger.yaml to be flagged as unnecessary")
	}
}

func TestUltimateGoalController_GetGoalSummary(t *testing.T) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")

	summary := controller.GetGoalSummary()

	// Verificar campos obrigatórios
	if summary.PrimaryGoal == "" {
		t.Error("PrimaryGoal should not be empty")
	}

	if summary.Confidence <= 0 || summary.Confidence > 1 {
		t.Errorf("Confidence should be between 0 and 1, got %f", summary.Confidence)
	}

	if len(summary.KeyComponents) == 0 {
		t.Error("KeyComponents should not be empty")
	}

	// Verificar se componentes esperados estão presentes
	expectedComponents := []string{"api", "rest", "crud"}
	found := 0
	for _, expected := range expectedComponents {
		for _, component := range summary.KeyComponents {
			if component == expected {
				found++
				break
			}
		}
	}

	if found == 0 {
		t.Error("Expected to find at least one key component related to API/REST/CRUD")
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        []string
	}{
		{
			name:        "API REST",
			description: "criar uma API REST com PostgreSQL",
			want:        []string{"api", "rest", "postgres"},
		},
		{
			name:        "Frontend React",
			description: "desenvolver frontend com React e TypeScript",
			want:        []string{"frontend", "react", "typescript"},
		},
		{
			name:        "CLI Go",
			description: "fazer um CLI em Go",
			want:        []string{"cli", "go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keywords := extractKeywords(tt.description)

			for _, expectedKeyword := range tt.want {
				found := false
				for _, keyword := range keywords {
					if keyword == expectedKeyword {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected keyword '%s' not found in %v", expectedKeyword, keywords)
				}
			}
		})
	}
}

func TestDetermineOptimalScope(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "Minimal with 'apenas'",
			description: "criar apenas uma API simples",
			want:        "minimal",
		},
		{
			name:        "Comprehensive with 'completo'",
			description: "desenvolver uma aplicação web completa",
			want:        "comprehensive",
		},
		{
			name:        "Balanced default",
			description: "criar uma API para usuários",
			want:        "balanced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := determineOptimalScope(tt.description)
			if scope != tt.want {
				t.Errorf("determineOptimalScope() = %v, want %v", scope, tt.want)
			}
		})
	}
}

// Helper function para testes
func containsString(text, substring string) bool {
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

func BenchmarkUltimateGoalController_BuildPrompt(b *testing.B) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")
	basePrompt := "Crie um projeto Go com estrutura completa"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = controller.BuildGoalFocusedPrompt(basePrompt)
	}
}

func BenchmarkUltimateGoalController_FilterContent(b *testing.B) {
	controller := NewUltimateGoalController("criar uma API REST simples apenas para CRUD de usuários")

	sampleContent := `{
		"structure": {
			"directories": ["src", "tests", "docs", "docker"],
			"files": {
				"main.go": {"content": "package main"},
				"handlers.go": {"content": "package handlers"},
				"docker-compose.yml": {"content": "version: '3'"},
				"swagger.yaml": {"content": "swagger: '2.0'"}
			}
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = controller.FilterGeneratedContent(sampleContent)
	}
}
