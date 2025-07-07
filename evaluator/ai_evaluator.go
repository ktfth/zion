package evaluator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ktfth/zion/ai/providers"
)

// AIEvaluatorRule usa IA para avaliar a qualidade do projeto
type AIEvaluatorRule struct{}

func (r *AIEvaluatorRule) Name() string {
	return "AIEvaluation"
}

func (r *AIEvaluatorRule) Description() string {
	return "Avaliação inteligente da qualidade do projeto usando IA"
}

func (r *AIEvaluatorRule) Category() Category {
	return CategoryStructure
}

func (r *AIEvaluatorRule) Weight() float64 {
	return 20.0 // Peso alto para a avaliação de IA
}

func (r *AIEvaluatorRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Tentar avaliação por IA
	aiIssues, aiScore, err := r.evaluateWithAI(structure)
	if err != nil {
		// Se a IA falhar, retorna score neutro e adiciona aviso
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: fmt.Sprintf("Avaliação por IA não disponível: %v", err),
			Suggestion:  "Avaliação baseada apenas em regras locais",
		})
		return score, issues
	}

	return aiScore, aiIssues
}

func (r *AIEvaluatorRule) evaluateWithAI(structure *ProjectStructure) ([]Issue, float64, error) {
	// Serializar estrutura do projeto para envio à IA
	projectJSON, err := r.serializeProjectStructure(structure)
	if err != nil {
		return nil, 1.0, fmt.Errorf("erro ao serializar projeto: %v", err)
	}

	// Criar prompt para IA
	prompt := r.createEvaluationPrompt(structure.Language, projectJSON)

	// Fazer chamada para IA
	response, err := r.callAI(prompt)
	if err != nil {
		return nil, 1.0, fmt.Errorf("erro na chamada para IA: %v", err)
	}

	// Parsear resposta da IA
	return r.parseAIResponse(response)
}

func (r *AIEvaluatorRule) serializeProjectStructure(structure *ProjectStructure) (string, error) {
	// Criar estrutura simplificada para envio
	simplified := map[string]interface{}{
		"language":     structure.Language,
		"projectName":  structure.ProjectName,
		"directories":  structure.Directories,
		"files":        r.getFileList(structure),
		"dependencies": structure.Dependencies,
		"metadata":     structure.Metadata,
	}

	data, err := json.MarshalIndent(simplified, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *AIEvaluatorRule) getFileList(structure *ProjectStructure) map[string]interface{} {
	files := make(map[string]interface{})

	for fileName, fileInfo := range structure.Files {
		files[fileName] = map[string]interface{}{
			"path":      fileInfo.Path,
			"extension": fileInfo.Extension,
			"isConfig":  fileInfo.IsConfig,
			"size":      fileInfo.Size,
			// Incluir apenas primeiros 200 caracteres do conteúdo para economizar tokens
			"contentPreview": r.truncateContent(fileInfo.Content, 200),
		}
	}

	return files
}

func (r *AIEvaluatorRule) truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

func (r *AIEvaluatorRule) createEvaluationPrompt(language, projectJSON string) string {
	return fmt.Sprintf(`Você é um especialista em avaliação de qualidade de código e estrutura de projetos de software.

Analise a estrutura do projeto %s abaixo e forneça uma avaliação detalhada:

%s

Por favor, avalie o projeto considerando:

1. **Estrutura de Diretórios**: Organização e convenções
2. **Arquivos Essenciais**: Presença de arquivos obrigatórios
3. **Segurança**: Detecção de vulnerabilidades e dados sensíveis
4. **Boas Práticas**: Aderência às convenções da linguagem
5. **Manutenibilidade**: Facilidade de manutenção e documentação

Retorne sua análise no seguinte formato JSON:

{
  "score": <número entre 0-100>,
  "issues": [
    {
      "severity": "critical|high|medium|low",
      "category": "structure|security|naming|dependencies|maintainability",
      "description": "Descrição do problema",
      "suggestion": "Como corrigir"
    }
  ],
  "summary": "Resumo da avaliação",
  "hasCriticalIssues": <true|false>
}

Seja rigoroso na avaliação. Issues críticos devem bloquear a criação do projeto.
Foque em problemas que realmente impactam a qualidade, segurança e manutenibilidade.`, language, projectJSON)
}

func (r *AIEvaluatorRule) callAI(prompt string) (string, error) {
	// Obter provedor de IA configurado
	providerManager := providers.NewProviderManager()

	// Usar configuração padrão ou obter do ambiente
	providerName := "openrouter" // Pode ser configurável
	config := make(map[string]string)

	provider, err := providerManager.GetProvider(providerName, config)
	if err != nil {
		return "", fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Fazer a chamada
	response, err := provider.GenerateContent(prompt)
	if err != nil {
		return "", fmt.Errorf("erro na geração de conteúdo: %v", err)
	}

	return response, nil
}

func (r *AIEvaluatorRule) parseAIResponse(response string) ([]Issue, float64, error) {
	// Extrair JSON da resposta
	jsonContent := r.extractJSONFromResponse(response)
	if jsonContent == "" {
		return nil, 1.0, fmt.Errorf("resposta da IA não contém JSON válido")
	}

	// Parsear JSON
	var aiEvaluation struct {
		Score             float64 `json:"score"`
		HasCriticalIssues bool    `json:"hasCriticalIssues"`
		Summary           string  `json:"summary"`
		Issues            []struct {
			Severity    string `json:"severity"`
			Category    string `json:"category"`
			Description string `json:"description"`
			Suggestion  string `json:"suggestion"`
		} `json:"issues"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &aiEvaluation); err != nil {
		return nil, 1.0, fmt.Errorf("erro ao parsear resposta da IA: %v", err)
	}

	// Converter para nosso formato
	var issues []Issue
	for _, aiIssue := range aiEvaluation.Issues {
		severity := r.convertSeverity(aiIssue.Severity)
		category := r.convertCategory(aiIssue.Category)

		issues = append(issues, Issue{
			Type:        IssueStructure, // Padrão para issues da IA
			Severity:    severity,
			Category:    category,
			Description: fmt.Sprintf("[IA] %s", aiIssue.Description),
			Suggestion:  aiIssue.Suggestion,
		})
	}

	// Normalizar score (0-100 para 0-1)
	score := aiEvaluation.Score / 100.0
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return issues, score, nil
}

func (r *AIEvaluatorRule) extractJSONFromResponse(response string) string {
	// Procurar por blocos JSON na resposta
	start := strings.Index(response, "{")
	if start == -1 {
		return ""
	}

	// Encontrar JSON balanceado
	braceCount := 0
	var result bytes.Buffer

	for i := start; i < len(response); i++ {
		char := response[i]
		result.WriteByte(char)

		if char == '{' {
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 {
				break
			}
		}
	}

	jsonStr := result.String()

	// Validar se é JSON válido
	var test interface{}
	if json.Unmarshal([]byte(jsonStr), &test) == nil {
		return jsonStr
	}

	return ""
}

func (r *AIEvaluatorRule) convertSeverity(aiSeverity string) Severity {
	switch strings.ToLower(aiSeverity) {
	case "critical":
		return SeverityCritical
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	default:
		return SeverityMedium
	}
}

func (r *AIEvaluatorRule) convertCategory(aiCategory string) Category {
	switch strings.ToLower(aiCategory) {
	case "structure":
		return CategoryStructure
	case "security":
		return CategorySecurity
	case "naming":
		return CategoryNaming
	case "dependencies":
		return CategoryDependencies
	case "maintainability":
		return CategoryMaintainability
	case "configuration":
		return CategoryConfiguration
	default:
		return CategoryStructure
	}
}
