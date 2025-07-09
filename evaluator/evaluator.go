package evaluator

import (
	"fmt"
	"strings"
)

// ProjectEvaluator representa o sistema de avaliação de projetos
type ProjectEvaluator struct {
	rules    []EvaluationRule
	enableAI bool
}

// EvaluationResult representa o resultado da avaliação
type EvaluationResult struct {
	Valid       bool                   `json:"valid"`
	Score       float64                `json:"score"`
	MaxScore    float64                `json:"maxScore"`
	Issues      []Issue                `json:"issues"`
	Suggestions []string               `json:"suggestions"`
	Quality     QualityLevel           `json:"quality"`
	Details     map[string]interface{} `json:"details"`
}

// Issue representa um problema encontrado na avaliação
type Issue struct {
	Type        IssueType `json:"type"`
	Severity    Severity  `json:"severity"`
	Category    Category  `json:"category"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Suggestion  string    `json:"suggestion"`
}

// Enums para classificação
type IssueType string
type Severity string
type Category string
type QualityLevel string

const (
	// Tipos de issues
	IssueStructure     IssueType = "structure"
	IssueNaming        IssueType = "naming"
	IssueDependency    IssueType = "dependency"
	IssueConfiguration IssueType = "configuration"
	IssueBestPractice  IssueType = "best_practice"
	IssueSecurity      IssueType = "security"

	// Severidades
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"

	// Categorias
	CategoryStructure       Category = "structure"
	CategoryNaming          Category = "naming"
	CategoryDependencies    Category = "dependencies"
	CategoryConfiguration   Category = "configuration"
	CategorySecurity        Category = "security"
	CategoryPerformance     Category = "performance"
	CategoryMaintainability Category = "maintainability"
	CategoryBestPractice    Category = "best_practice"

	// Níveis de qualidade
	QualityExcellent QualityLevel = "excellent"
	QualityGood      QualityLevel = "good"
	QualityFair      QualityLevel = "fair"
	QualityPoor      QualityLevel = "poor"
	QualityCritical  QualityLevel = "critical"
)

// EvaluationRule define uma regra de avaliação
type EvaluationRule interface {
	Name() string
	Description() string
	Category() Category
	Weight() float64
	Evaluate(structure *ProjectStructure) (float64, []Issue)
}

// NewProjectEvaluator cria uma nova instância do avaliador
func NewProjectEvaluator() *ProjectEvaluator {
	evaluator := &ProjectEvaluator{
		rules:    make([]EvaluationRule, 0),
		enableAI: false, // Padrão: IA desabilitada
	}

	// Registrar regras padrão
	evaluator.RegisterDefaultRules()

	return evaluator
}

// EnableAIEvaluation habilita ou desabilita a avaliação por IA
func (pe *ProjectEvaluator) EnableAIEvaluation(enable bool) {
	pe.enableAI = enable
}

// RegisterRule registra uma nova regra de avaliação
func (pe *ProjectEvaluator) RegisterRule(rule EvaluationRule) {
	pe.rules = append(pe.rules, rule)
}

// RegisterDefaultRules registra as regras padrão de avaliação
func (pe *ProjectEvaluator) RegisterDefaultRules() {
	// Regras de estrutura
	pe.RegisterRule(&DirectoryStructureRule{})
	pe.RegisterRule(&FileNamingRule{})
	pe.RegisterRule(&RequiredFilesRule{})

	// Regras de dependências
	pe.RegisterRule(&DependencyConsistencyRule{})
	pe.RegisterRule(&SecurityVulnerabilityRule{})

	// Regras de configuração
	pe.RegisterRule(&ConfigurationValidityRule{})
	pe.RegisterRule(&BuildConfigurationRule{})

	// Regras de melhores práticas
	pe.RegisterRule(&BestPracticesRule{})
	pe.RegisterRule(&DocumentationRule{})
	pe.RegisterRule(&TestStructureRule{})

	// Avaliação por IA (condicional)
	if pe.enableAI {
		pe.RegisterRule(&AIEvaluatorRule{})
	}
}

// EvaluateProject avalia um projeto baseado na estrutura fornecida
func (pe *ProjectEvaluator) EvaluateProject(projectData string, language string) (*EvaluationResult, error) {
	// Extrair estrutura do projeto
	structure, err := ExtractProjectStructure(projectData, language)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair estrutura do projeto: %w", err)
	}

	result := &EvaluationResult{
		Valid:       true,
		Score:       0,
		MaxScore:    0,
		Issues:      make([]Issue, 0),
		Suggestions: make([]string, 0),
		Details:     make(map[string]interface{}),
	}

	// Executar todas as regras
	for _, rule := range pe.rules {
		score, issues := rule.Evaluate(structure)
		weight := rule.Weight()

		result.Score += score * weight
		result.MaxScore += weight
		result.Issues = append(result.Issues, issues...)

		// Armazenar detalhes por categoria
		category := string(rule.Category())
		if result.Details[category] == nil {
			result.Details[category] = make(map[string]interface{})
		}

		details := result.Details[category].(map[string]interface{})
		details[rule.Name()] = map[string]interface{}{
			"score":       score,
			"weight":      weight,
			"description": rule.Description(),
			"issues":      len(issues),
		}
	}

	// Calcular score final
	if result.MaxScore > 0 {
		result.Score = (result.Score / result.MaxScore) * 100
	}

	// Determinar nível de qualidade
	result.Quality = pe.determineQualityLevel(result.Score, result.Issues)

	// Verificar se o projeto é válido (sem issues críticos)
	result.Valid = pe.isProjectValid(result.Issues)

	// Gerar sugestões
	result.Suggestions = pe.generateSuggestions(result.Issues, structure)

	return result, nil
}

// determineQualityLevel determina o nível de qualidade baseado no score e issues
func (pe *ProjectEvaluator) determineQualityLevel(score float64, issues []Issue) QualityLevel {
	hasCritical := false
	hasHigh := false

	for _, issue := range issues {
		if issue.Severity == SeverityCritical {
			hasCritical = true
		} else if issue.Severity == SeverityHigh {
			hasHigh = true
		}
	}

	if hasCritical {
		return QualityCritical
	}

	if score >= 90 && !hasHigh {
		return QualityExcellent
	} else if score >= 75 {
		return QualityGood
	} else if score >= 60 {
		return QualityFair
	} else {
		return QualityPoor
	}
}

// isProjectValid verifica se o projeto é válido (sem issues críticos)
func (pe *ProjectEvaluator) isProjectValid(issues []Issue) bool {
	for _, issue := range issues {
		if issue.Severity == SeverityCritical {
			return false
		}
	}
	return true
}

// generateSuggestions gera sugestões baseadas nos issues encontrados
func (pe *ProjectEvaluator) generateSuggestions(issues []Issue, structure *ProjectStructure) []string {
	suggestions := make([]string, 0)
	categoryCount := make(map[Category]int)

	// Contar issues por categoria
	for _, issue := range issues {
		categoryCount[issue.Category]++
		if issue.Suggestion != "" {
			suggestions = append(suggestions, issue.Suggestion)
		}
	}

	// Sugestões gerais baseadas nas categorias mais problemáticas
	if categoryCount[CategoryStructure] > 3 {
		suggestions = append(suggestions, "Considere reorganizar a estrutura de diretórios seguindo as convenções da linguagem")
	}

	if categoryCount[CategorySecurity] > 0 {
		suggestions = append(suggestions, "Revise as configurações de segurança e dependências vulneráveis")
	}

	if categoryCount[CategoryDependencies] > 2 {
		suggestions = append(suggestions, "Verifique e atualize as dependências do projeto")
	}

	return suggestions
}

// GenerateReport gera um relatório detalhado da avaliação
func (pe *ProjectEvaluator) GenerateReport(result *EvaluationResult) string {
	var report strings.Builder

	report.WriteString("🔍 RELATÓRIO DE AVALIAÇÃO DO PROJETO\n")
	report.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Status geral
	if result.Valid {
		report.WriteString("✅ Status: VÁLIDO\n")
	} else {
		report.WriteString("❌ Status: INVÁLIDO\n")
	}

	report.WriteString(fmt.Sprintf("🎯 Score: %.1f/100\n", result.Score))
	report.WriteString(fmt.Sprintf("⭐ Qualidade: %s\n", pe.getQualityEmoji(result.Quality)))
	report.WriteString(fmt.Sprintf("⚠️  Issues encontrados: %d\n\n", len(result.Issues)))

	// Issues por severidade
	severityCount := make(map[Severity]int)
	for _, issue := range result.Issues {
		severityCount[issue.Severity]++
	}

	if len(severityCount) > 0 {
		report.WriteString("📊 DISTRIBUIÇÃO DE ISSUES\n")
		report.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		for severity, count := range severityCount {
			emoji := pe.getSeverityEmoji(severity)
			report.WriteString(fmt.Sprintf("%s %s: %d\n", emoji, strings.Title(string(severity)), count))
		}
		report.WriteString("\n")
	}

	// Issues detalhados
	if len(result.Issues) > 0 {
		report.WriteString("🔍 ISSUES DETALHADOS\n")
		report.WriteString("━━━━━━━━━━━━━━━━━━━━━━\n")

		for _, issue := range result.Issues {
			emoji := pe.getSeverityEmoji(issue.Severity)
			report.WriteString(fmt.Sprintf("%s [%s] %s\n", emoji, strings.ToUpper(string(issue.Category)), issue.Description))
			if issue.Location != "" {
				report.WriteString(fmt.Sprintf("   📍 Local: %s\n", issue.Location))
			}
			if issue.Suggestion != "" {
				report.WriteString(fmt.Sprintf("   💡 Sugestão: %s\n", issue.Suggestion))
			}
			report.WriteString("\n")
		}
	}

	// Sugestões
	if len(result.Suggestions) > 0 {
		report.WriteString("💡 SUGESTÕES DE MELHORIA\n")
		report.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		for i, suggestion := range result.Suggestions {
			report.WriteString(fmt.Sprintf("%d. %s\n", i+1, suggestion))
		}
		report.WriteString("\n")
	}

	// Detalhes por categoria
	report.WriteString("📈 ANÁLISE POR CATEGORIA\n")
	report.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	for category, details := range result.Details {
		report.WriteString(fmt.Sprintf("📂 %s:\n", strings.Title(category)))
		if categoryDetails, ok := details.(map[string]interface{}); ok {
			for ruleName, ruleDetails := range categoryDetails {
				if ruleInfo, ok := ruleDetails.(map[string]interface{}); ok {
					score := ruleInfo["score"].(float64)
					weight := ruleInfo["weight"].(float64)
					issues := ruleInfo["issues"].(int)

					finalScore := score * weight
					report.WriteString(fmt.Sprintf("  • %s: %.1f pontos", ruleName, finalScore))
					if issues > 0 {
						report.WriteString(fmt.Sprintf(" (%d issues)", issues))
					}
					report.WriteString("\n")
				}
			}
		}
		report.WriteString("\n")
	}

	return report.String()
}

// Helper functions
func (pe *ProjectEvaluator) getQualityEmoji(quality QualityLevel) string {
	switch quality {
	case QualityExcellent:
		return "🏆 EXCELENTE"
	case QualityGood:
		return "✅ BOM"
	case QualityFair:
		return "⚡ REGULAR"
	case QualityPoor:
		return "⚠️ RUIM"
	case QualityCritical:
		return "❌ CRÍTICO"
	default:
		return "❓ INDEFINIDO"
	}
}

func (pe *ProjectEvaluator) getSeverityEmoji(severity Severity) string {
	switch severity {
	case SeverityCritical:
		return "🚨"
	case SeverityHigh:
		return "❌"
	case SeverityMedium:
		return "⚠️"
	case SeverityLow:
		return "💡"
	case SeverityInfo:
		return "ℹ️"
	default:
		return "❓"
	}
}
