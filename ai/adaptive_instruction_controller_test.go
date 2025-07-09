package ai

import (
	"testing"
)

func TestAdaptiveInstructionController(t *testing.T) {
	testCases := []struct {
		name                 string
		description          string
		expectedScope        string
		expectedStrictness   int
		expectedRequirements []string
		expectedConstraints  []string
	}{
		{
			name:                 "Minimal API",
			description:          "uma API REST simples apenas para CRUD de usuários",
			expectedScope:        "minimal",
			expectedStrictness:   8,
			expectedRequirements: []string{"API_ENDPOINTS"},
			expectedConstraints:  []string{"STRICT_MINIMAL_SCOPE"},
		},
		{
			name:                 "Complete Application",
			description:          "uma aplicação completa com frontend, backend, testes e docker",
			expectedScope:        "comprehensive",
			expectedStrictness:   6,
			expectedRequirements: []string{"API_ENDPOINTS", "FRONTEND_INTERFACE", "COMPREHENSIVE_TESTING", "CONTAINERIZATION"},
			expectedConstraints:  []string{},
		},
		{
			name:                 "Standard Project",
			description:          "um projeto para gerenciamento de tarefas",
			expectedScope:        "standard",
			expectedStrictness:   5,
			expectedRequirements: []string{},
			expectedConstraints:  []string{},
		},
		{
			name:                 "Minimal CLI",
			description:          "um comando CLI básico somente para converter arquivos",
			expectedScope:        "minimal",
			expectedStrictness:   8,
			expectedRequirements: []string{},
			expectedConstraints:  []string{"STRICT_MINIMAL_SCOPE"},
		},
		{
			name:                 "Frontend Only",
			description:          "uma interface web apenas para exibir dados",
			expectedScope:        "minimal",
			expectedStrictness:   8,
			expectedRequirements: []string{"FRONTEND_INTERFACE"},
			expectedConstraints:  []string{"STRICT_MINIMAL_SCOPE"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewAdaptiveInstructionController("general_application", "javascript", tc.description)
			profile := controller.GetInstructionProfile()

			// Verificar escopo
			if controller.Scope != tc.expectedScope {
				t.Errorf("Expected scope %s, got %s", tc.expectedScope, controller.Scope)
			}

			// Verificar nível de rigidez
			if profile.StrictnessLevel != tc.expectedStrictness {
				t.Errorf("Expected strictness %d, got %d", tc.expectedStrictness, profile.StrictnessLevel)
			}

			// Verificar requisitos
			if len(tc.expectedRequirements) > 0 {
				for _, req := range tc.expectedRequirements {
					if !contains(controller.Requirements, req) {
						t.Errorf("Expected requirement %s not found", req)
					}
				}
			}

			// Verificar constraints
			if len(tc.expectedConstraints) > 0 {
				for _, constraint := range tc.expectedConstraints {
					if !contains(controller.Constraints, constraint) {
						t.Errorf("Expected constraint %s not found", constraint)
					}
				}
			}
		})
	}
}

func TestAdaptivePromptBuilding(t *testing.T) {
	controller := NewAdaptiveInstructionController("backend_api", "javascript", "uma API REST simples apenas para CRUD de usuários")

	basePrompt := "Crie uma API REST"
	adaptivePrompt := controller.BuildAdaptivePrompt(basePrompt)

	// Verificar se o prompt contém elementos esperados
	expectedElements := []string{
		"ADAPTIVE INSTRUCTION CONTROL",
		"STRICTNESS LEVEL:",
		"SCOPE CONTROL:",
		"MINIMAL SCOPE MODE",
		"CRITÉRIO DE CAMALEÃO",
		"VALIDATION RULES",
	}

	for _, element := range expectedElements {
		if !contains([]string{adaptivePrompt}, element) {
			t.Errorf("Expected element '%s' not found in adaptive prompt", element)
		}
	}
}

func TestInstructionCompliance(t *testing.T) {
	controller := NewAdaptiveInstructionController("backend_api", "javascript", "uma API REST simples apenas para CRUD de usuários")

	// Teste com resposta conforme
	compliantResponse := `{
		"structure": {
			"directories": ["src", "config"],
			"files": {
				"src/index.js": {
					"content": "const express = require('express');\nconst app = express();\napp.listen(3000);"
				},
				"package.json": {
					"content": {
						"name": "crud-api",
						"version": "1.0.0",
						"main": "src/index.js"
					}
				}
			}
		}
	}`

	compliance, err := controller.ValidateInstructionCompliance(compliantResponse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !compliance.IsCompliant {
		t.Errorf("Expected compliant response, got non-compliant")
	}

	// Teste com resposta não conforme (JSON inválido)
	nonCompliantResponse := `{
		"structure": {
			"directories": ["src", "config"
			"files": {
				"missing": "closing bracket"
			}
		}
	}`

	compliance, err = controller.ValidateInstructionCompliance(nonCompliantResponse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if compliance.IsCompliant {
		t.Errorf("Expected non-compliant response, got compliant")
	}

	if !contains(compliance.ViolatedRules, "INVALID_JSON_STRUCTURE") {
		t.Errorf("Expected INVALID_JSON_STRUCTURE violation not found")
	}
}

func TestProjectTypeDetection(t *testing.T) {
	testCases := []struct {
		description  string
		expectedType string
	}{
		{
			description:  "uma API REST para gerenciamento de usuários",
			expectedType: "backend_api",
		},
		{
			description:  "uma interface web para exibir dados",
			expectedType: "frontend_web",
		},
		{
			description:  "um comando CLI para converter arquivos",
			expectedType: "cli_tool",
		},
		{
			description:  "uma biblioteca JavaScript para validação",
			expectedType: "library",
		},
		{
			description:  "um bot para automatizar tarefas",
			expectedType: "automation_bot",
		},
		{
			description:  "um jogo simples em Python",
			expectedType: "game",
		},
		{
			description:  "um app mobile para iOS",
			expectedType: "mobile_app",
		},
		{
			description:  "um dashboard administrativo",
			expectedType: "admin_dashboard",
		},
		{
			description:  "um sistema de gerenciamento",
			expectedType: "general_application",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			detectedType := detectProjectType(tc.description)
			if detectedType != tc.expectedType {
				t.Errorf("Expected project type %s, got %s", tc.expectedType, detectedType)
			}
		})
	}
}

func TestScopeAdaptation(t *testing.T) {
	testCases := []struct {
		description   string
		expectedScope string
	}{
		{
			description:   "apenas uma função simples",
			expectedScope: "minimal",
		},
		{
			description:   "somente o básico para funcionar",
			expectedScope: "minimal",
		},
		{
			description:   "uma aplicação completa e robusta",
			expectedScope: "comprehensive",
		},
		{
			description:   "um projeto abrangente com tudo",
			expectedScope: "comprehensive",
		},
		{
			description:   "um sistema de gerenciamento",
			expectedScope: "standard",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			controller := NewAdaptiveInstructionController("general_application", "javascript", tc.description)
			if controller.Scope != tc.expectedScope {
				t.Errorf("Expected scope %s, got %s", tc.expectedScope, controller.Scope)
			}
		})
	}
}

func TestAdaptationFlags(t *testing.T) {
	testCases := []struct {
		description   string
		expectedFlags map[string]bool
	}{
		{
			description: "uma API REST com frontend e testes",
			expectedFlags: map[string]bool{
				"include_api":      true,
				"include_frontend": true,
				"include_tests":    true,
			},
		},
		{
			description: "uma aplicação com banco de dados e docker",
			expectedFlags: map[string]bool{
				"include_database": true,
				"include_docker":   true,
			},
		},
		{
			description: "apenas uma interface simples",
			expectedFlags: map[string]bool{
				"include_frontend": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			controller := NewAdaptiveInstructionController("general_application", "javascript", tc.description)

			for flag, expected := range tc.expectedFlags {
				if controller.Adaptations[flag] != expected {
					t.Errorf("Expected adaptation flag %s to be %v, got %v", flag, expected, controller.Adaptations[flag])
				}
			}
		})
	}
}

// Helper function para verificar se um slice contém um elemento
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Benchmark para testar performance
func BenchmarkAdaptiveInstructionController(b *testing.B) {
	description := "uma API REST completa com frontend, backend, testes e docker"

	for i := 0; i < b.N; i++ {
		controller := NewAdaptiveInstructionController("backend_api", "javascript", description)
		profile := controller.GetInstructionProfile()
		_ = controller.BuildAdaptivePrompt("Crie um projeto")
		_ = profile
	}
}
