package evaluator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ktfth/zion/ai/providers"
	"github.com/ktfth/zion/config"
)

// AIEvaluatorRule usa IA para avaliar a qualidade do projeto
type AIEvaluatorRule struct {
	provider providers.Provider
	cache    map[string]*EvaluationCache
}

// EvaluationCache armazena resultados de avaliação
type EvaluationCache struct {
	Result    EvaluationResult `json:"result"`
	Timestamp time.Time        `json:"timestamp"`
	Hash      string           `json:"hash"`
}

// AIEvaluationResponse representa resposta da IA
type AIEvaluationResponse struct {
	OverallScore      float64               `json:"overall_score"`
	CategoryScores    map[string]float64    `json:"category_scores"`
	Issues            []AIIssue             `json:"issues"`
	Suggestions       []string              `json:"suggestions"`
	BestPractices     []string              `json:"best_practices"`
	SecurityConcerns  []string              `json:"security_concerns"`
	PerformanceIssues []string              `json:"performance_issues"`
	Maintainability   float64               `json:"maintainability"`
	Scalability       float64               `json:"scalability"`
	CodeQuality       float64               `json:"code_quality"`
	Architecture      ArchitectureAnalysis  `json:"architecture"`
	Dependencies      DependencyAnalysis    `json:"dependencies"`
	Testing           TestingAnalysis       `json:"testing"`
	Documentation     DocumentationAnalysis `json:"documentation"`
	Innovation        float64               `json:"innovation"`
	Complexity        float64               `json:"complexity"`
	Readability       float64               `json:"readability"`
	Reusability       float64               `json:"reusability"`
}

// AIIssue representa um problema identificado pela IA
type AIIssue struct {
	Type          string  `json:"type"`
	Severity      string  `json:"severity"`
	Category      string  `json:"category"`
	Description   string  `json:"description"`
	Location      string  `json:"location"`
	Suggestion    string  `json:"suggestion"`
	Impact        float64 `json:"impact"`
	Priority      int     `json:"priority"`
	FixComplexity string  `json:"fix_complexity"`
	AutoFixable   bool    `json:"auto_fixable"`
}

// ArchitectureAnalysis análise arquitetural
type ArchitectureAnalysis struct {
	Pattern         string   `json:"pattern"`
	Clarity         float64  `json:"clarity"`
	Consistency     float64  `json:"consistency"`
	Separation      float64  `json:"separation"`
	Coupling        float64  `json:"coupling"`
	Cohesion        float64  `json:"cohesion"`
	Modularity      float64  `json:"modularity"`
	Extensibility   float64  `json:"extensibility"`
	Violations      []string `json:"violations"`
	Recommendations []string `json:"recommendations"`
}

// DependencyAnalysis análise de dependências
type DependencyAnalysis struct {
	TotalCount             int      `json:"total_count"`
	OutdatedCount          int      `json:"outdated_count"`
	SecurityVulns          int      `json:"security_vulns"`
	UnusedCount            int      `json:"unused_count"`
	RedundantCount         int      `json:"redundant_count"`
	HealthScore            float64  `json:"health_score"`
	Recommendations        []string `json:"recommendations"`
	CriticalDeps           []string `json:"critical_deps"`
	AlternativeSuggestions []string `json:"alternative_suggestions"`
}

// TestingAnalysis análise de testes
type TestingAnalysis struct {
	Coverage         float64  `json:"coverage"`
	Quality          float64  `json:"quality"`
	TestTypes        []string `json:"test_types"`
	MissingTests     []string `json:"missing_tests"`
	TestSmells       []string `json:"test_smells"`
	Recommendations  []string `json:"recommendations"`
	UnitTests        bool     `json:"unit_tests"`
	IntegrationTests bool     `json:"integration_tests"`
	E2ETests         bool     `json:"e2e_tests"`
	MockUsage        float64  `json:"mock_usage"`
}

// DocumentationAnalysis análise de documentação
type DocumentationAnalysis struct {
	Coverage         float64  `json:"coverage"`
	Quality          float64  `json:"quality"`
	Completeness     float64  `json:"completeness"`
	Clarity          float64  `json:"clarity"`
	Examples         bool     `json:"examples"`
	APIDocumentation bool     `json:"api_documentation"`
	UserGuide        bool     `json:"user_guide"`
	DevGuide         bool     `json:"dev_guide"`
	Changelog        bool     `json:"changelog"`
	MissingDocs      []string `json:"missing_docs"`
	Recommendations  []string `json:"recommendations"`
}

func (r *AIEvaluatorRule) Name() string {
	return "AIEvaluation"
}

func (r *AIEvaluatorRule) Description() string {
	return "Avaliação inteligente da qualidade do projeto usando IA avançada"
}

func (r *AIEvaluatorRule) Category() Category {
	return CategoryStructure
}

func (r *AIEvaluatorRule) Weight() float64 {
	return 25.0 // Peso ainda maior para avaliação de IA melhorada
}

func (r *AIEvaluatorRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Inicializar cache se não existir
	if r.cache == nil {
		r.cache = make(map[string]*EvaluationCache)
	}

	// Verificar cache
	cacheKey := r.generateCacheKey(structure)
	if cached, exists := r.cache[cacheKey]; exists {
		if time.Since(cached.Timestamp) < 30*time.Minute {
			// Converter resultado cached para formato atual
			for _, issue := range cached.Result.Issues {
				issues = append(issues, issue)
			}
			return cached.Result.Score, issues
		}
	}

	// Tentar avaliação por IA
	aiResponse, err := r.evaluateWithAdvancedAI(structure)
	if err != nil {
		// Se a IA falhar, usar avaliação básica
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: fmt.Sprintf("Avaliação por IA não disponível: %v", err),
			Suggestion:  "Avaliação baseada apenas em regras locais",
		})
		return score, issues
	}

	// Converter resposta da IA para formato interno
	score = aiResponse.OverallScore
	issues = r.convertAIIssues(aiResponse.Issues)

	// Adicionar sugestões como issues informativos
	for _, suggestion := range aiResponse.Suggestions {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityLow,
			Category:    CategoryBestPractice,
			Description: suggestion,
			Suggestion:  "Implementar melhoria sugerida",
		})
	}

	// Adicionar problemas de segurança
	for _, security := range aiResponse.SecurityConcerns {
		issues = append(issues, Issue{
			Type:        IssueSecurity,
			Severity:    SeverityHigh,
			Category:    CategorySecurity,
			Description: security,
			Suggestion:  "Corrigir vulnerabilidade de segurança",
		})
	}

	// Adicionar problemas de performance
	for _, perf := range aiResponse.PerformanceIssues {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryPerformance,
			Description: perf,
			Suggestion:  "Otimizar performance",
		})
	}

	// Armazenar no cache
	r.cache[cacheKey] = &EvaluationCache{
		Result: EvaluationResult{
			Score:  score,
			Issues: issues,
		},
		Timestamp: time.Now(),
		Hash:      cacheKey,
	}

	return score, issues
}

func (r *AIEvaluatorRule) evaluateWithAdvancedAI(structure *ProjectStructure) (*AIEvaluationResponse, error) {
	// Obter provider de IA
	if r.provider == nil {
		cfg := config.LoadConfig()
		provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, cfg.GetAIConfig())
		if err != nil {
			return nil, fmt.Errorf("erro ao obter provedor de IA: %v", err)
		}
		r.provider = provider
	}

	// Serializar estrutura do projeto
	projectJSON, err := r.serializeAdvancedStructure(structure)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar projeto: %v", err)
	}

	// Criar prompt avançado
	prompt := r.createAdvancedEvaluationPrompt(structure, projectJSON)

	// Fazer chamada para IA
	response, err := r.provider.GenerateContent(prompt)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada para IA: %v", err)
	}

	// Parsear resposta da IA
	return r.parseAdvancedAIResponse(response)
}

func (r *AIEvaluatorRule) serializeAdvancedStructure(structure *ProjectStructure) (string, error) {
	// Análise mais detalhada da estrutura
	analysis := map[string]interface{}{
		"basic_info": map[string]interface{}{
			"language":     structure.Language,
			"project_name": structure.ProjectName,
			"total_files":  len(structure.Files),
			"total_dirs":   len(structure.Directories),
			"total_deps":   len(structure.Dependencies),
		},
		"file_analysis":            r.analyzeFiles(structure),
		"directory_analysis":       r.analyzeDirectories(structure),
		"dependency_analysis":      r.analyzeDependencies(structure),
		"architecture_analysis":    r.analyzeArchitecture(structure),
		"quality_metrics":          r.calculateQualityMetrics(structure),
		"complexity_metrics":       r.calculateComplexityMetrics(structure),
		"security_analysis":        r.performSecurityAnalysis(structure),
		"performance_analysis":     r.performPerformanceAnalysis(structure),
		"maintainability_analysis": r.performMaintainabilityAnalysis(structure),
		"documentation_analysis":   r.performDocumentationAnalysis(structure),
		"testing_analysis":         r.performTestingAnalysis(structure),
	}

	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *AIEvaluatorRule) analyzeFiles(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"by_type":             make(map[string]int),
		"by_size":             make(map[string]int),
		"largest_files":       make([]string, 0),
		"generated_files":     make([]string, 0),
		"config_files":        make([]string, 0),
		"test_files":          make([]string, 0),
		"documentation_files": make([]string, 0),
		"code_files":          make([]string, 0),
	}

	byType := analysis["by_type"].(map[string]int)
	bySize := analysis["by_size"].(map[string]int)
	largestFiles := analysis["largest_files"].([]string)
	generatedFiles := analysis["generated_files"].([]string)
	configFiles := analysis["config_files"].([]string)
	testFiles := analysis["test_files"].([]string)
	docFiles := analysis["documentation_files"].([]string)
	codeFiles := analysis["code_files"].([]string)

	for fileName, fileInfo := range structure.Files {
		// Analisar por tipo
		ext := r.getFileExtension(fileName)
		byType[ext]++

		// Analisar por tamanho
		sizeCategory := r.getSizeCategory(int64(fileInfo.Size))
		bySize[sizeCategory]++

		// Arquivos maiores
		if fileInfo.Size > 10000 {
			largestFiles = append(largestFiles, fileName)
		}

		// Arquivos gerados
		if r.isGeneratedFile(fileName, fileInfo.Content) {
			generatedFiles = append(generatedFiles, fileName)
		}

		// Arquivos de configuração
		if fileInfo.IsConfig {
			configFiles = append(configFiles, fileName)
		}

		// Arquivos de teste
		if r.isTestFile(fileName) {
			testFiles = append(testFiles, fileName)
		}

		// Arquivos de documentação
		if r.isDocumentationFile(fileName) {
			docFiles = append(docFiles, fileName)
		}

		// Arquivos de código
		if r.isCodeFile(fileName) {
			codeFiles = append(codeFiles, fileName)
		}
	}

	analysis["by_type"] = byType
	analysis["by_size"] = bySize
	analysis["largest_files"] = largestFiles
	analysis["generated_files"] = generatedFiles
	analysis["config_files"] = configFiles
	analysis["test_files"] = testFiles
	analysis["documentation_files"] = docFiles
	analysis["code_files"] = codeFiles

	return analysis
}

func (r *AIEvaluatorRule) analyzeDirectories(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"depth":                 0,
		"structure_type":        "unknown",
		"organization":          "unknown",
		"common_patterns":       make([]string, 0),
		"missing_directories":   make([]string, 0),
		"redundant_directories": make([]string, 0),
	}

	// Calcular profundidade
	maxDepth := 0
	for _, dir := range structure.Directories {
		depth := strings.Count(dir, "/") + strings.Count(dir, "\\")
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	analysis["depth"] = maxDepth

	// Detectar padrões comuns
	patterns := []string{}
	if r.hasDirectory(structure.Directories, "src") {
		patterns = append(patterns, "src_pattern")
	}
	if r.hasDirectory(structure.Directories, "lib") {
		patterns = append(patterns, "lib_pattern")
	}
	if r.hasDirectory(structure.Directories, "test") || r.hasDirectory(structure.Directories, "tests") {
		patterns = append(patterns, "test_pattern")
	}
	if r.hasDirectory(structure.Directories, "docs") {
		patterns = append(patterns, "docs_pattern")
	}

	analysis["common_patterns"] = patterns

	return analysis
}

func (r *AIEvaluatorRule) analyzeDependencies(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"total_count":         len(structure.Dependencies),
		"by_type":             make(map[string]int),
		"popular_deps":        make([]string, 0),
		"potential_conflicts": make([]string, 0),
		"security_concerns":   make([]string, 0),
		"outdated_deps":       make([]string, 0),
	}

	byType := analysis["by_type"].(map[string]int)
	popularDeps := analysis["popular_deps"].([]string)

	for depName, _ := range structure.Dependencies {
		// Categorizar por tipo
		depType := r.categorizeDependency(depName)
		byType[depType]++

		// Dependências populares
		if r.isPopularDependency(depName) {
			popularDeps = append(popularDeps, depName)
		}
	}

	analysis["by_type"] = byType
	analysis["popular_deps"] = popularDeps

	return analysis
}

func (r *AIEvaluatorRule) analyzeArchitecture(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"pattern":                "unknown",
		"layered":                false,
		"modular":                false,
		"separation_of_concerns": 0.0,
		"coupling_level":         "unknown",
		"cohesion_level":         "unknown",
	}

	// Detectar padrões arquiteturais
	if r.hasLayeredStructure(structure) {
		analysis["pattern"] = "layered"
		analysis["layered"] = true
	}

	if r.hasModularStructure(structure) {
		analysis["modular"] = true
	}

	// Calcular separação de responsabilidades
	analysis["separation_of_concerns"] = r.calculateSeparationOfConcerns(structure)

	return analysis
}

func (r *AIEvaluatorRule) calculateQualityMetrics(structure *ProjectStructure) map[string]interface{} {
	metrics := map[string]interface{}{
		"code_to_comment_ratio": 0.0,
		"duplication_index":     0.0,
		"naming_consistency":    0.0,
		"file_organization":     0.0,
		"configuration_quality": 0.0,
	}

	// Calcular métricas básicas
	totalLines := 0
	commentLines := 0

	for _, fileInfo := range structure.Files {
		if r.isCodeFile(fileInfo.Path) {
			lines := strings.Count(fileInfo.Content, "\n")
			totalLines += lines
			commentLines += r.countCommentLines(fileInfo.Content)
		}
	}

	if totalLines > 0 {
		metrics["code_to_comment_ratio"] = float64(commentLines) / float64(totalLines)
	}

	return metrics
}

func (r *AIEvaluatorRule) calculateComplexityMetrics(structure *ProjectStructure) map[string]interface{} {
	metrics := map[string]interface{}{
		"cyclomatic_complexity": 0.0,
		"nesting_depth":         0.0,
		"function_length":       0.0,
		"class_complexity":      0.0,
	}

	// Calcular complexidade ciclomática simplificada
	totalComplexity := 0
	fileCount := 0

	for _, fileInfo := range structure.Files {
		if r.isCodeFile(fileInfo.Path) {
			complexity := r.calculateCyclomaticComplexity(fileInfo.Content)
			totalComplexity += complexity
			fileCount++
		}
	}

	if fileCount > 0 {
		metrics["cyclomatic_complexity"] = float64(totalComplexity) / float64(fileCount)
	}

	return metrics
}

func (r *AIEvaluatorRule) performSecurityAnalysis(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"potential_vulnerabilities": make([]string, 0),
		"security_patterns":         make([]string, 0),
		"sensitive_data_exposure":   make([]string, 0),
		"authentication_issues":     make([]string, 0),
		"authorization_issues":      make([]string, 0),
	}

	vulnerabilities := make([]string, 0)

	// Verificar padrões de segurança
	for fileName, fileInfo := range structure.Files {
		content := strings.ToLower(fileInfo.Content)

		// Verificar exposição de dados sensíveis
		if strings.Contains(content, "password") && strings.Contains(content, "=") {
			vulnerabilities = append(vulnerabilities, fmt.Sprintf("Possível exposição de senha em %s", fileName))
		}

		// Verificar injeção SQL
		if strings.Contains(content, "select") && strings.Contains(content, "+") {
			vulnerabilities = append(vulnerabilities, fmt.Sprintf("Possível injeção SQL em %s", fileName))
		}

		// Verificar API keys hardcoded
		if strings.Contains(content, "api_key") || strings.Contains(content, "apikey") {
			vulnerabilities = append(vulnerabilities, fmt.Sprintf("Possível API key hardcoded em %s", fileName))
		}
	}

	analysis["potential_vulnerabilities"] = vulnerabilities

	return analysis
}

func (r *AIEvaluatorRule) performPerformanceAnalysis(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"potential_bottlenecks":      make([]string, 0),
		"resource_usage":             make(map[string]interface{}),
		"optimization_opportunities": make([]string, 0),
		"async_patterns":             make([]string, 0),
	}

	bottlenecks := analysis["potential_bottlenecks"].([]string)
	optimizations := analysis["optimization_opportunities"].([]string)

	// Analisar arquivos em busca de problemas de performance
	for fileName, fileInfo := range structure.Files {
		content := strings.ToLower(fileInfo.Content)

		// Verificar loops aninhados
		if strings.Count(content, "for") > 2 {
			bottlenecks = append(bottlenecks, fmt.Sprintf("Múltiplos loops em %s", fileName))
		}

		// Verificar operações síncronas
		if strings.Contains(content, "sync") && !strings.Contains(content, "async") {
			optimizations = append(optimizations, fmt.Sprintf("Considerar operações assíncronas em %s", fileName))
		}
	}

	analysis["potential_bottlenecks"] = bottlenecks
	analysis["optimization_opportunities"] = optimizations

	return analysis
}

func (r *AIEvaluatorRule) performMaintainabilityAnalysis(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"code_smells":               make([]string, 0),
		"refactoring_opportunities": make([]string, 0),
		"documentation_coverage":    0.0,
		"test_coverage":             0.0,
		"code_reusability":          0.0,
	}

	codeSmells := analysis["code_smells"].([]string)
	refactoring := analysis["refactoring_opportunities"].([]string)

	// Detectar code smells
	for fileName, fileInfo := range structure.Files {
		if r.isCodeFile(fileInfo.Path) {
			content := fileInfo.Content

			// Funções muito longas
			if strings.Count(content, "\n") > 100 {
				codeSmells = append(codeSmells, fmt.Sprintf("Arquivo muito longo: %s", fileName))
			}

			// Muitas responsabilidades
			if strings.Count(content, "class") > 5 {
				refactoring = append(refactoring, fmt.Sprintf("Muitas classes em %s", fileName))
			}
		}
	}

	analysis["code_smells"] = codeSmells
	analysis["refactoring_opportunities"] = refactoring

	return analysis
}

func (r *AIEvaluatorRule) performDocumentationAnalysis(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"readme_present":   false,
		"api_docs_present": false,
		"inline_comments":  0.0,
		"doc_quality":      0.0,
		"missing_docs":     make([]string, 0),
	}

	missingDocs := analysis["missing_docs"].([]string)

	// Verificar presença de README
	hasReadme := false
	for fileName := range structure.Files {
		if strings.ToLower(fileName) == "readme.md" || strings.ToLower(fileName) == "readme.txt" {
			hasReadme = true
			break
		}
	}
	analysis["readme_present"] = hasReadme

	if !hasReadme {
		missingDocs = append(missingDocs, "README.md")
	}

	// Verificar documentação de API
	hasAPIDocs := false
	for fileName := range structure.Files {
		if strings.Contains(strings.ToLower(fileName), "api") && strings.Contains(strings.ToLower(fileName), "doc") {
			hasAPIDocs = true
			break
		}
	}
	analysis["api_docs_present"] = hasAPIDocs

	analysis["missing_docs"] = missingDocs

	return analysis
}

func (r *AIEvaluatorRule) performTestingAnalysis(structure *ProjectStructure) map[string]interface{} {
	analysis := map[string]interface{}{
		"test_files_present":     false,
		"test_coverage_estimate": 0.0,
		"test_types":             make([]string, 0),
		"missing_tests":          make([]string, 0),
	}

	testTypes := analysis["test_types"].([]string)
	missingTests := analysis["missing_tests"].([]string)

	// Contar arquivos de teste
	testFiles := 0
	codeFiles := 0

	for fileName := range structure.Files {
		if r.isTestFile(fileName) {
			testFiles++

			// Identificar tipos de teste
			if strings.Contains(strings.ToLower(fileName), "unit") {
				testTypes = append(testTypes, "unit")
			}
			if strings.Contains(strings.ToLower(fileName), "integration") {
				testTypes = append(testTypes, "integration")
			}
			if strings.Contains(strings.ToLower(fileName), "e2e") {
				testTypes = append(testTypes, "e2e")
			}
		} else if r.isCodeFile(fileName) {
			codeFiles++
		}
	}

	analysis["test_files_present"] = testFiles > 0

	if codeFiles > 0 {
		analysis["test_coverage_estimate"] = float64(testFiles) / float64(codeFiles) * 100
	}

	if testFiles == 0 {
		missingTests = append(missingTests, "Nenhum arquivo de teste encontrado")
	}

	analysis["test_types"] = testTypes
	analysis["missing_tests"] = missingTests

	return analysis
}

func (r *AIEvaluatorRule) createAdvancedEvaluationPrompt(structure *ProjectStructure, analysisJSON string) string {
	prompt := fmt.Sprintf(`🤖 AVALIAÇÃO AVANÇADA DE QUALIDADE DE CÓDIGO

Você é um especialista em qualidade de código com conhecimento profundo em arquitetura de software, boas práticas e análise estática.

📋 ANÁLISE DO PROJETO:
%s

🎯 TAREFAS:
1. Avalie a qualidade geral do projeto (0.0 a 1.0)
2. Identifique problemas específicos e sua severidade
3. Forneça sugestões concretas de melhoria
4. Analise aspectos de segurança, performance e manutenibilidade
5. Avalie a arquitetura e organização do código

📊 CRITÉRIOS DE AVALIAÇÃO:
- Estrutura e organização (20%%)
- Qualidade do código (20%%)
- Arquitetura e design (15%%)
- Segurança (15%%)
- Performance (10%%)
- Manutenibilidade (10%%)
- Documentação (5%%)
- Testes (5%%)

🔍 RESPOSTA ESPERADA (JSON):
{
  "overall_score": 0.85,
  "category_scores": {
    "structure": 0.9,
    "code_quality": 0.8,
    "architecture": 0.85,
    "security": 0.7,
    "performance": 0.8,
    "maintainability": 0.9,
    "documentation": 0.6,
    "testing": 0.7
  },
  "issues": [
    {
      "type": "security",
      "severity": "high",
      "category": "security",
      "description": "Possível exposição de credenciais",
      "location": "config/database.js",
      "suggestion": "Usar variáveis de ambiente",
      "impact": 0.8,
      "priority": 1,
      "fix_complexity": "medium",
      "auto_fixable": false
    }
  ],
  "suggestions": [
    "Implementar testes unitários para cobertura completa",
    "Adicionar documentação da API",
    "Configurar linting e formatação automática"
  ],
  "best_practices": [
    "Uso correto de padrões de design",
    "Separação clara de responsabilidades"
  ],
  "security_concerns": [
    "Validação de entrada insuficiente",
    "Falta de sanitização de dados"
  ],
  "performance_issues": [
    "Queries N+1 potenciais",
    "Falta de cache em operações custosas"
  ],
  "maintainability": 0.85,
  "scalability": 0.80,
  "code_quality": 0.82,
  "architecture": {
    "pattern": "layered",
    "clarity": 0.85,
    "consistency": 0.90,
    "separation": 0.80,
    "coupling": 0.75,
    "cohesion": 0.85,
    "modularity": 0.80,
    "extensibility": 0.75,
    "violations": ["Tight coupling between layers"],
    "recommendations": ["Implement dependency injection"]
  },
  "dependencies": {
    "total_count": 25,
    "outdated_count": 3,
    "security_vulns": 1,
    "unused_count": 2,
    "redundant_count": 1,
    "health_score": 0.85,
    "recommendations": ["Update outdated dependencies", "Remove unused packages"],
    "critical_deps": ["express", "react"],
    "alternative_suggestions": ["Consider fastify instead of express for better performance"]
  },
  "testing": {
    "coverage": 0.65,
    "quality": 0.70,
    "test_types": ["unit", "integration"],
    "missing_tests": ["E2E tests", "Performance tests"],
    "test_smells": ["Long test methods", "Too many assertions"],
    "recommendations": ["Add more edge case tests", "Improve test isolation"],
    "unit_tests": true,
    "integration_tests": true,
    "e2e_tests": false,
    "mock_usage": 0.60
  },
  "documentation": {
    "coverage": 0.60,
    "quality": 0.65,
    "completeness": 0.70,
    "clarity": 0.75,
    "examples": true,
    "api_documentation": false,
    "user_guide": true,
    "dev_guide": false,
    "changelog": false,
    "missing_docs": ["API documentation", "Developer guide"],
    "recommendations": ["Add inline code comments", "Create comprehensive API docs"]
  },
  "innovation": 0.75,
  "complexity": 0.60,
  "readability": 0.85,
  "reusability": 0.70
}

Forneça uma análise detalhada e precisa do projeto.`, analysisJSON)

	return prompt
}

func (r *AIEvaluatorRule) parseAdvancedAIResponse(response string) (*AIEvaluationResponse, error) {
	// Extrair JSON da resposta
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("JSON não encontrado na resposta")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var aiResponse AIEvaluationResponse
	if err := json.Unmarshal([]byte(jsonStr), &aiResponse); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON: %v", err)
	}

	return &aiResponse, nil
}

func (r *AIEvaluatorRule) convertAIIssues(aiIssues []AIIssue) []Issue {
	issues := make([]Issue, 0)

	for _, aiIssue := range aiIssues {
		issue := Issue{
			Type:        r.convertIssueType(aiIssue.Type),
			Severity:    r.convertSeverity(aiIssue.Severity),
			Category:    r.convertCategory(aiIssue.Category),
			Description: aiIssue.Description,
			Location:    aiIssue.Location,
			Suggestion:  aiIssue.Suggestion,
		}
		issues = append(issues, issue)
	}

	return issues
}

func (r *AIEvaluatorRule) convertIssueType(aiType string) IssueType {
	switch aiType {
	case "structure":
		return IssueStructure
	case "naming":
		return IssueNaming
	case "dependency":
		return IssueDependency
	case "configuration":
		return IssueConfiguration
	case "best_practice":
		return IssueBestPractice
	case "security":
		return IssueSecurity
	default:
		return IssueStructure
	}
}

func (r *AIEvaluatorRule) convertSeverity(aiSeverity string) Severity {
	switch aiSeverity {
	case "low":
		return SeverityLow
	case "medium":
		return SeverityMedium
	case "high":
		return SeverityHigh
	case "critical":
		return SeverityCritical
	default:
		return SeverityMedium
	}
}

func (r *AIEvaluatorRule) convertCategory(aiCategory string) Category {
	switch aiCategory {
	case "structure":
		return CategoryStructure
	case "naming":
		return CategoryNaming
	case "dependency":
		return CategoryDependencies
	case "configuration":
		return CategoryConfiguration
	case "best_practice":
		return CategoryBestPractice
	case "security":
		return CategorySecurity
	case "performance":
		return CategoryPerformance
	default:
		return CategoryStructure
	}
}

// Métodos auxiliares para análise
func (r *AIEvaluatorRule) generateCacheKey(structure *ProjectStructure) string {
	// Gerar chave baseada no hash do projeto
	data := fmt.Sprintf("%s|%s|%d|%d",
		structure.Language,
		structure.ProjectName,
		len(structure.Files),
		len(structure.Dependencies))

	return fmt.Sprintf("%x", data)
}

func (r *AIEvaluatorRule) getFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

func (r *AIEvaluatorRule) getSizeCategory(size int64) string {
	if size < 1024 {
		return "small"
	} else if size < 10240 {
		return "medium"
	} else if size < 102400 {
		return "large"
	} else {
		return "xlarge"
	}
}

func (r *AIEvaluatorRule) isGeneratedFile(fileName, content string) bool {
	// Verificar marcadores de arquivos gerados
	if strings.Contains(content, "@generated") ||
		strings.Contains(content, "AUTO-GENERATED") ||
		strings.Contains(content, "automatically generated") {
		return true
	}

	// Verificar nomes de arquivos típicos
	generatedPatterns := []string{
		".generated.", "_generated.", "generated_",
		".build.", "_build.", "build_",
		".dist.", "_dist.", "dist_",
	}

	for _, pattern := range generatedPatterns {
		if strings.Contains(fileName, pattern) {
			return true
		}
	}

	return false
}

func (r *AIEvaluatorRule) isTestFile(fileName string) bool {
	testPatterns := []string{
		"test", "spec", "_test.", ".test.",
		"_spec.", ".spec.", "Test", "Spec",
	}

	for _, pattern := range testPatterns {
		if strings.Contains(fileName, pattern) {
			return true
		}
	}

	return false
}

func (r *AIEvaluatorRule) isDocumentationFile(fileName string) bool {
	docPatterns := []string{
		"README", "readme", "CHANGELOG", "changelog",
		"CONTRIBUTING", "contributing", "LICENSE", "license",
		"INSTALL", "install", "USAGE", "usage",
		".md", ".txt", ".rst", ".adoc",
	}

	for _, pattern := range docPatterns {
		if strings.Contains(fileName, pattern) {
			return true
		}
	}

	return false
}

func (r *AIEvaluatorRule) isCodeFile(fileName string) bool {
	codeExtensions := []string{
		".js", ".ts", ".jsx", ".tsx", ".py", ".go", ".java",
		".cpp", ".c", ".cs", ".rb", ".php", ".rs", ".kt",
		".swift", ".dart", ".scala", ".clj", ".hs", ".elm",
	}

	for _, ext := range codeExtensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}

func (r *AIEvaluatorRule) hasDirectory(directories []string, dirName string) bool {
	for _, dir := range directories {
		if strings.Contains(dir, dirName) {
			return true
		}
	}
	return false
}

func (r *AIEvaluatorRule) categorizeDependency(dep string) string {
	// Categorizar dependências por tipo
	if strings.Contains(dep, "test") || strings.Contains(dep, "jest") || strings.Contains(dep, "mocha") {
		return "testing"
	}
	if strings.Contains(dep, "webpack") || strings.Contains(dep, "babel") || strings.Contains(dep, "rollup") {
		return "build"
	}
	if strings.Contains(dep, "express") || strings.Contains(dep, "fastify") || strings.Contains(dep, "gin") {
		return "framework"
	}
	if strings.Contains(dep, "react") || strings.Contains(dep, "vue") || strings.Contains(dep, "angular") {
		return "frontend"
	}
	if strings.Contains(dep, "database") || strings.Contains(dep, "sql") || strings.Contains(dep, "mongo") {
		return "database"
	}

	return "other"
}

func (r *AIEvaluatorRule) isPopularDependency(dep string) bool {
	popularDeps := []string{
		"react", "vue", "angular", "express", "fastify",
		"lodash", "moment", "axios", "jquery", "bootstrap",
		"webpack", "babel", "typescript", "eslint", "prettier",
	}

	for _, popular := range popularDeps {
		if strings.Contains(dep, popular) {
			return true
		}
	}

	return false
}

func (r *AIEvaluatorRule) hasLayeredStructure(structure *ProjectStructure) bool {
	layers := []string{"controller", "service", "repository", "model", "view"}
	foundLayers := 0

	for _, dir := range structure.Directories {
		for _, layer := range layers {
			if strings.Contains(strings.ToLower(dir), layer) {
				foundLayers++
				break
			}
		}
	}

	return foundLayers >= 3
}

func (r *AIEvaluatorRule) hasModularStructure(structure *ProjectStructure) bool {
	modules := []string{"module", "component", "feature", "domain"}

	for _, dir := range structure.Directories {
		for _, module := range modules {
			if strings.Contains(strings.ToLower(dir), module) {
				return true
			}
		}
	}

	return false
}

func (r *AIEvaluatorRule) calculateSeparationOfConcerns(structure *ProjectStructure) float64 {
	// Cálculo simplificado baseado na organização de diretórios
	totalFiles := len(structure.Files)
	organizedFiles := 0

	organizationPatterns := []string{
		"controller", "service", "model", "view", "component",
		"util", "helper", "middleware", "config", "test",
	}

	for fileName := range structure.Files {
		for _, pattern := range organizationPatterns {
			if strings.Contains(strings.ToLower(fileName), pattern) {
				organizedFiles++
				break
			}
		}
	}

	if totalFiles == 0 {
		return 0.0
	}

	return float64(organizedFiles) / float64(totalFiles)
}

func (r *AIEvaluatorRule) countCommentLines(content string) int {
	lines := strings.Split(content, "\n")
	commentCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "/*") ||
			strings.HasPrefix(trimmed, "*") {
			commentCount++
		}
	}

	return commentCount
}

func (r *AIEvaluatorRule) calculateCyclomaticComplexity(content string) int {
	// Cálculo simplificado de complexidade ciclomática
	complexity := 1 // Complexidade base

	// Contadores para construções que aumentam complexidade
	keywords := []string{
		"if", "else", "elif", "for", "while", "switch", "case",
		"try", "catch", "finally", "&&", "||", "?", ":",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range keywords {
		complexity += strings.Count(contentLower, keyword)
	}

	return complexity
}
