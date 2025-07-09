package main

import (
	"fmt"
	"log"

	"github.com/ktfth/zion/ai"
)

func main() {
	// Exemplo 1: API REST Mínima
	fmt.Println("=== Exemplo 1: API REST Mínima ===")
	testMinimalAPI()

	// Exemplo 2: Aplicação Completa
	fmt.Println("\n=== Exemplo 2: Aplicação Completa ===")
	testCompleteApplication()

	// Exemplo 3: CLI Tool Simples
	fmt.Println("\n=== Exemplo 3: CLI Tool Simples ===")
	testCLITool()

	// Exemplo 4: Frontend com Testes
	fmt.Println("\n=== Exemplo 4: Frontend com Testes ===")
	testFrontendWithTests()
}

func testMinimalAPI() {
	description := "uma API REST simples apenas para CRUD de usuários"
	controller := ai.NewAdaptiveInstructionController("backend_api", "javascript", description)

	fmt.Printf("Tipo de Projeto: %s\n", controller.ProjectType)
	fmt.Printf("Escopo: %s\n", controller.Scope)
	fmt.Printf("Nível de Rigidez: %d/10\n", controller.GetInstructionProfile().StrictnessLevel)
	fmt.Printf("Requisitos: %v\n", controller.Requirements)
	fmt.Printf("Restrições: %v\n", controller.Constraints)

	// Simular geração de prompt
	basePrompt := "Crie uma API REST para CRUD de usuários"
	adaptivePrompt := controller.BuildAdaptivePrompt(basePrompt)

	fmt.Printf("\nPrompt adaptativo gerado (%d caracteres)\n", len(adaptivePrompt))
	fmt.Printf("Contém instruções de escopo mínimo: %t\n",
		contains(adaptivePrompt, "MINIMAL SCOPE MODE"))
	fmt.Printf("Contém critério de camaleão: %t\n",
		contains(adaptivePrompt, "CRITÉRIO DE CAMALEÃO"))
}

func testCompleteApplication() {
	description := "uma aplicação completa com frontend, backend, testes e docker"
	controller := ai.NewAdaptiveInstructionController("general_application", "javascript", description)

	fmt.Printf("Tipo de Projeto: %s\n", controller.ProjectType)
	fmt.Printf("Escopo: %s\n", controller.Scope)
	fmt.Printf("Nível de Rigidez: %d/10\n", controller.GetInstructionProfile().StrictnessLevel)
	fmt.Printf("Requisitos: %v\n", controller.Requirements)
	fmt.Printf("Adaptações: %v\n", controller.Adaptations)

	// Verificar se todas as adaptações foram detectadas
	expectedAdaptations := []string{"include_api", "include_frontend", "include_tests", "include_docker"}
	for _, adaptation := range expectedAdaptations {
		if controller.Adaptations[adaptation] == true {
			fmt.Printf("✅ Adaptação %s detectada\n", adaptation)
		} else {
			fmt.Printf("❌ Adaptação %s não detectada\n", adaptation)
		}
	}
}

func testCLITool() {
	description := "um comando CLI básico somente para converter arquivos CSV para JSON"
	controller := ai.NewAdaptiveInstructionController("cli_tool", "go", description)

	fmt.Printf("Tipo de Projeto: %s\n", controller.ProjectType)
	fmt.Printf("Escopo: %s\n", controller.Scope)
	fmt.Printf("Nível de Rigidez: %d/10\n", controller.GetInstructionProfile().StrictnessLevel)

	profile := controller.GetInstructionProfile()
	fmt.Printf("Áreas de Foco: %v\n", profile.FocusAreas)
	fmt.Printf("Regras de Exclusão: %v\n", profile.ExclusionRules)
	fmt.Printf("Limiar de Qualidade: %.1f%%\n", profile.QualityThreshold)

	// Testar validação de conformidade
	testResponse := `{
		"structure": {
			"directories": ["cmd", "pkg"],
			"files": {
				"main.go": {
					"content": "package main\n\nfunc main() {\n\t// Convert CSV to JSON\n}"
				},
				"go.mod": {
					"content": "module csv-converter\n\ngo 1.21"
				}
			}
		}
	}`

	compliance, err := controller.ValidateInstructionCompliance(testResponse)
	if err != nil {
		log.Printf("Erro na validação: %v", err)
		return
	}

	fmt.Printf("Conformidade: %t\n", compliance.IsCompliant)
	fmt.Printf("Score: %.1f%%\n", compliance.ComplianceScore)

	if len(compliance.ViolatedRules) > 0 {
		fmt.Printf("Regras Violadas: %v\n", compliance.ViolatedRules)
	}
}

func testFrontendWithTests() {
	description := "uma interface web com testes unitários para exibir dashboard"
	controller := ai.NewAdaptiveInstructionController("frontend_web", "typescript", description)

	fmt.Printf("Tipo de Projeto: %s\n", controller.ProjectType)
	fmt.Printf("Escopo: %s\n", controller.Scope)
	fmt.Printf("Nível de Rigidez: %d/10\n", controller.GetInstructionProfile().StrictnessLevel)

	// Verificar adaptações específicas
	if controller.Adaptations["include_frontend"] == true {
		fmt.Println("✅ Frontend detectado")
	}
	if controller.Adaptations["include_tests"] == true {
		fmt.Println("✅ Testes detectados")
	}

	// Simular geração de camadas adaptativas
	if lg, err := ai.NewLayeredGenerator("typescript", "dashboard-app", description, nil); err == nil {
		fmt.Println("✅ Gerador em camadas criado com sucesso")

		// Aqui seria executado o planejamento de camadas
		// que usaria o controlador adaptativo
		fmt.Println("🔧 Planejamento de camadas seria executado com adaptações:")
		fmt.Printf("   - Frontend: %t\n", controller.Adaptations["include_frontend"] == true)
		fmt.Printf("   - Testes: %t\n", controller.Adaptations["include_tests"] == true)
		fmt.Printf("   - Escopo: %s\n", controller.Scope)
	} else {
		log.Printf("Erro ao criar gerador: %v", err)
	}
}

// demonstrateLayeredGeneration mostra como o sistema funciona com geração em camadas
func demonstrateLayeredGeneration() {
	fmt.Println("\n=== Demonstração de Geração em Camadas ===")

	// Exemplo com diferentes escopos
	examples := []struct {
		name        string
		description string
		language    string
	}{
		{
			name:        "API Mínima",
			description: "apenas uma API REST básica para CRUD",
			language:    "javascript",
		},
		{
			name:        "App Completo",
			description: "aplicação completa com frontend, backend, testes e docker",
			language:    "typescript",
		},
		{
			name:        "CLI Simple",
			description: "comando CLI somente para processar arquivos",
			language:    "go",
		},
	}

	for _, example := range examples {
		fmt.Printf("\n--- %s ---\n", example.name)

		controller := ai.NewAdaptiveInstructionController(
			ai.DetectProjectType(example.description),
			example.language,
			example.description,
		)

		profile := controller.GetInstructionProfile()

		fmt.Printf("Escopo: %s\n", controller.Scope)
		fmt.Printf("Rigidez: %d/10\n", profile.StrictnessLevel)
		fmt.Printf("Qualidade: %.1f%%\n", profile.QualityThreshold)

		// Simular camadas que seriam geradas
		if lg, err := ai.NewLayeredGenerator(example.language, "test-project", example.description, nil); err == nil {
			fmt.Println("✅ Gerador configurado com adaptações")
			// Aqui seria executado lg.GenerateLayeredProject()
		} else {
			fmt.Printf("❌ Erro: %v\n", err)
		}
	}
}

// Helper function
func contains(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 &&
		(str == substr || (len(str) > len(substr) &&
			findSubstring(str, substr) >= 0))
}

func findSubstring(str, substr string) int {
	if len(substr) > len(str) {
		return -1
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
