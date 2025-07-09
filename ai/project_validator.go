package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ProjectValidationResult contém o resultado da validação
type ProjectValidationResult struct {
	IsValid     bool     `json:"is_valid"`
	Issues      []string `json:"issues"`
	Suggestions []string `json:"suggestions"`
	Score       float64  `json:"score"`
}

// ValidateProjectStructure valida a estrutura de um projeto gerado
func ValidateProjectStructure(response string, language string) *ProjectValidationResult {
	result := &ProjectValidationResult{
		IsValid:     true,
		Issues:      make([]string, 0),
		Suggestions: make([]string, 0),
		Score:       100.0,
	}

	// Validar se é JSON válido
	if !isValidJSON(response) {
		result.IsValid = false
		result.Issues = append(result.Issues, "Resposta não é um JSON válido")
		result.Score -= 50
		return result
	}

	// Validar estrutura específica do projeto
	if isLayeredResponse(response) {
		validateLayeredStructure(response, result)
	} else {
		validateTraditionalStructure(response, result)
	}

	// Validar elementos específicos da linguagem
	validateLanguageSpecificElements(response, language, result)

	return result
}

// isValidJSON verifica se a string é um JSON válido
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// isLayeredResponse verifica se é uma resposta em camadas
func isLayeredResponse(response string) bool {
	return strings.Contains(response, `"layers"`) && strings.Contains(response, `"project_info"`)
}

// validateLayeredStructure valida estrutura de resposta em camadas
func validateLayeredStructure(response string, result *ProjectValidationResult) {
	var layeredResponse LayeredResponse
	if err := json.Unmarshal([]byte(response), &layeredResponse); err != nil {
		result.IsValid = false
		result.Issues = append(result.Issues, "Erro ao parsear resposta em camadas")
		result.Score -= 30
		return
	}

	// Validar informações do projeto
	if layeredResponse.ProjectInfo.Name == "" {
		result.Issues = append(result.Issues, "Nome do projeto não pode estar vazio")
		result.Score -= 10
	}

	if layeredResponse.ProjectInfo.Language == "" {
		result.Issues = append(result.Issues, "Linguagem do projeto não pode estar vazia")
		result.Score -= 10
	}

	// Validar camadas
	if len(layeredResponse.Layers) == 0 {
		result.IsValid = false
		result.Issues = append(result.Issues, "Projeto deve ter pelo menos uma camada")
		result.Score -= 40
		return
	}

	totalFiles := 0
	for i, layer := range layeredResponse.Layers {
		if layer.LayerName == "" {
			result.Issues = append(result.Issues, fmt.Sprintf("Camada %d deve ter um nome", i+1))
			result.Score -= 5
		}

		if len(layer.Files) == 0 {
			result.Issues = append(result.Issues, fmt.Sprintf("Camada %s deve ter pelo menos um arquivo", layer.LayerName))
			result.Score -= 10
		}

		totalFiles += len(layer.Files)
	}

	if totalFiles < 3 {
		result.Issues = append(result.Issues, "Projeto deve ter pelo menos 3 arquivos")
		result.Score -= 20
	}
}

// validateTraditionalStructure valida estrutura tradicional
func validateTraditionalStructure(response string, result *ProjectValidationResult) {
	var scaffoldResponse ScaffoldResponse
	if err := json.Unmarshal([]byte(response), &scaffoldResponse); err != nil {
		result.IsValid = false
		result.Issues = append(result.Issues, "Erro ao parsear resposta tradicional")
		result.Score -= 30
		return
	}

	// Validar diretórios
	if len(scaffoldResponse.Structure.Directories) == 0 {
		result.Issues = append(result.Issues, "Projeto deve ter pelo menos um diretório")
		result.Score -= 15
	}

	// Validar arquivos
	if len(scaffoldResponse.Structure.Files) == 0 {
		result.IsValid = false
		result.Issues = append(result.Issues, "Projeto deve ter pelo menos um arquivo")
		result.Score -= 40
		return
	}

	if len(scaffoldResponse.Structure.Files) < 3 {
		result.Issues = append(result.Issues, "Projeto deve ter pelo menos 3 arquivos")
		result.Score -= 20
	}
}

// validateLanguageSpecificElements valida elementos específicos da linguagem
func validateLanguageSpecificElements(response string, language string, result *ProjectValidationResult) {
	switch strings.ToLower(language) {
	case "javascript", "typescript":
		if !strings.Contains(response, "package.json") {
			result.Issues = append(result.Issues, "Projeto JavaScript/TypeScript deve ter package.json")
			result.Score -= 15
		}
	case "python":
		if !strings.Contains(response, "requirements.txt") && !strings.Contains(response, "pyproject.toml") {
			result.Issues = append(result.Issues, "Projeto Python deve ter requirements.txt ou pyproject.toml")
			result.Score -= 15
		}
	case "go":
		if !strings.Contains(response, "go.mod") {
			result.Issues = append(result.Issues, "Projeto Go deve ter go.mod")
			result.Score -= 15
		}
	case "java":
		if !strings.Contains(response, "pom.xml") && !strings.Contains(response, "build.gradle") {
			result.Issues = append(result.Issues, "Projeto Java deve ter pom.xml ou build.gradle")
			result.Score -= 15
		}
	}

	// Validar arquivos essenciais
	if !strings.Contains(response, "README") {
		result.Suggestions = append(result.Suggestions, "Considere adicionar um arquivo README")
		result.Score -= 5
	}
}
