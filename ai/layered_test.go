package ai

import (
	"encoding/json"
	"fmt"
)

// TestLayeredGeneration testa o sistema de geração em camadas
func TestLayeredGeneration() {
	fmt.Println("🧪 Testando sistema de geração em camadas...")

	// Criar contexto de teste
	llmsContext := &LLMsContext{
		HasLLMsFile: false,
		ProjectStructure: []string{
			"main.go",
			"go.mod",
			"README.md",
		},
	}

	// Criar gerador
	generator, err := NewLayeredGenerator("go", "test-project", "API REST com autenticação JWT", llmsContext)
	if err != nil {
		fmt.Printf("❌ Erro ao criar gerador: %v\n", err)
		return
	}

	fmt.Println("✅ Gerador criado com sucesso")

	// Testar detecção de overflow
	smallPrompt := "Teste pequeno"
	largePrompt := fmt.Sprintf("Teste muito grande: %s", make([]byte, 1000000))

	fmt.Printf("📏 Prompt pequeno overflow: %v\n", DetectContextOverflow(smallPrompt, "openrouter"))
	fmt.Printf("📏 Prompt grande overflow: %v\n", DetectContextOverflow(largePrompt, "openrouter"))

	// Testar planejamento de camadas
	layers, err := generator.planLayers()
	if err != nil {
		fmt.Printf("❌ Erro ao planejar camadas: %v\n", err)
		return
	}

	fmt.Printf("📋 Camadas planejadas: %d\n", len(layers))
	for i, layer := range layers {
		fmt.Printf("   %d. %s - %s (foco: %v)\n", i+1, layer.Name, layer.Description, layer.Focus)
	}

	// Testar conversão para formato padrão
	testLayeredResponse := &LayeredResponse{
		ProjectInfo: struct {
			Name        string `json:"name"`
			Language    string `json:"language"`
			Description string `json:"description"`
		}{
			Name:        "test-project",
			Language:    "go",
			Description: "Projeto de teste",
		},
		Layers: []LayerResult{
			{
				LayerName:   "core",
				Description: "Estrutura básica",
				Directories: []string{"cmd", "pkg"},
				Files: map[string]interface{}{
					"main.go": map[string]interface{}{
						"content": "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}",
					},
					"go.mod": map[string]interface{}{
						"content": "module test-project\n\ngo 1.21",
					},
				},
				Dependencies: []string{"go 1.21"},
				NextSteps:    []string{"go mod tidy", "go run main.go"},
			},
		},
	}

	scaffoldResponse := generator.ConvertToScaffoldResponse(testLayeredResponse)

	// Serializar para verificar se funciona
	jsonData, err := json.MarshalIndent(scaffoldResponse, "", "  ")
	if err != nil {
		fmt.Printf("❌ Erro ao serializar resposta: %v\n", err)
		return
	}

	fmt.Printf("✅ Conversão bem-sucedida (%d bytes)\n", len(jsonData))
	fmt.Printf("📊 Resultado: %d diretórios, %d arquivos\n",
		len(scaffoldResponse.Structure.Directories),
		len(scaffoldResponse.Structure.Files))

	fmt.Println("🎉 Todos os testes passaram!")
}

// ValidateLayeredResponse valida se uma resposta em camadas está bem formada
func ValidateLayeredResponse(response *LayeredResponse) []string {
	var issues []string

	if response.ProjectInfo.Name == "" {
		issues = append(issues, "Nome do projeto não pode estar vazio")
	}

	if response.ProjectInfo.Language == "" {
		issues = append(issues, "Linguagem do projeto não pode estar vazia")
	}

	if len(response.Layers) == 0 {
		issues = append(issues, "Projeto deve ter pelo menos uma camada")
	}

	for i, layer := range response.Layers {
		if layer.LayerName == "" {
			issues = append(issues, fmt.Sprintf("Camada %d deve ter um nome", i+1))
		}

		if len(layer.Files) == 0 {
			issues = append(issues, fmt.Sprintf("Camada %s deve ter pelo menos um arquivo", layer.LayerName))
		}

		// Validar estrutura dos arquivos
		for fileName, content := range layer.Files {
			if fileName == "" {
				issues = append(issues, fmt.Sprintf("Nome de arquivo vazio na camada %s", layer.LayerName))
				continue
			}

			// Verificar se o conteúdo tem estrutura válida
			if contentMap, ok := content.(map[string]interface{}); ok {
				if _, hasContent := contentMap["content"]; !hasContent {
					issues = append(issues, fmt.Sprintf("Arquivo %s na camada %s deve ter campo 'content'", fileName, layer.LayerName))
				}
			}
		}
	}

	return issues
}
