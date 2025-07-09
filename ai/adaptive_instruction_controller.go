package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AdaptiveInstructionController controla as instruções de forma adaptativa
type AdaptiveInstructionController struct {
	ProjectType   string
	Language      string
	Scope         string // "minimal", "standard", "comprehensive"
	TargetPurpose string
	Constraints   []string
	Requirements  []string
	Adaptations   map[string]interface{}
}

// InstructionProfile define perfis de instruções para diferentes contextos
type InstructionProfile struct {
	Name             string
	Description      string
	StrictnessLevel  int // 1-10 (1=flexível, 10=rigoroso)
	ScopeControl     string
	FocusAreas       []string
	ExclusionRules   []string
	QualityThreshold float64
	AdaptiveRules    map[string]string
}

// NewAdaptiveInstructionController cria um novo controlador adaptativo
func NewAdaptiveInstructionController(projectType, language, description string) *AdaptiveInstructionController {
	controller := &AdaptiveInstructionController{
		ProjectType:   projectType,
		Language:      language,
		Scope:         "standard",
		TargetPurpose: description,
		Constraints:   make([]string, 0),
		Requirements:  make([]string, 0),
		Adaptations:   make(map[string]interface{}),
	}

	// Analisar descrição para extrair intenções
	controller.analyzeIntent(description)

	return controller
}

// analyzeIntent analisa a descrição para extrair intenções e ajustar o controle
func (aic *AdaptiveInstructionController) analyzeIntent(description string) {
	desc := strings.ToLower(description)

	// Detectar palavras-chave de escopo com maior precisão
	minimalKeywords := []string{"apenas", "somente", "só", "minimo", "básico", "simples", "essencial", "necessário", "específico", "exato", "direto", "clean", "limpo"}
	comprehensiveKeywords := []string{"completo", "completa", "abrangente", "extenso", "robusto", "detalhado", "full", "total", "toda", "tudo", "máximo", "avançado", "profissional", "enterprise"}

	minimalCount := 0
	comprehensiveCount := 0

	for _, keyword := range minimalKeywords {
		if strings.Contains(desc, keyword) {
			minimalCount++
		}
	}

	for _, keyword := range comprehensiveKeywords {
		if strings.Contains(desc, keyword) {
			comprehensiveCount++
		}
	}

	// Determinar escopo baseado na frequência de palavras-chave
	if minimalCount > comprehensiveCount {
		aic.Scope = "minimal"
		aic.Constraints = append(aic.Constraints, "STRICT_MINIMAL_SCOPE", "CHAMELEON_PRECISION")
	} else if comprehensiveCount > minimalCount {
		aic.Scope = "comprehensive"
		aic.Requirements = append(aic.Requirements, "COMPREHENSIVE_IMPLEMENTATION", "CHAMELEON_COMPLETENESS")
	}

	// Detectar requisitos específicos com maior precisão
	if containsAny(desc, []string{"api", "rest", "endpoint", "servidor", "service", "webservice", "http"}) {
		aic.Adaptations["include_api"] = true
		aic.Requirements = append(aic.Requirements, "API_ENDPOINTS")
	}

	if containsAny(desc, []string{"frontend", "interface", "ui", "web", "página", "tela", "componente", "react", "vue", "angular"}) {
		aic.Adaptations["include_frontend"] = true
		aic.Requirements = append(aic.Requirements, "FRONTEND_INTERFACE")
	}

	if containsAny(desc, []string{"teste", "test", "tdd", "bdd", "unit", "integration", "spec", "validação"}) {
		aic.Adaptations["include_tests"] = true
		aic.Requirements = append(aic.Requirements, "COMPREHENSIVE_TESTING")
	}

	if containsAny(desc, []string{"docker", "container", "deployment", "deploy", "kubernetes", "k8s", "ci/cd", "pipeline"}) {
		aic.Adaptations["include_docker"] = true
		aic.Requirements = append(aic.Requirements, "CONTAINERIZATION")
	}

	if containsAny(desc, []string{"banco", "database", "db", "persistencia", "storage", "dados", "sql", "nosql", "mongodb", "mysql", "postgres"}) {
		aic.Adaptations["include_database"] = true
		aic.Requirements = append(aic.Requirements, "DATABASE_INTEGRATION")
	}

	if containsAny(desc, []string{"autenticação", "auth", "login", "segurança", "security", "jwt", "oauth", "session"}) {
		aic.Adaptations["include_auth"] = true
		aic.Requirements = append(aic.Requirements, "AUTHENTICATION_SYSTEM")
	}

	if containsAny(desc, []string{"cli", "command", "linha de comando", "terminal", "script", "tool", "ferramenta"}) {
		aic.Adaptations["include_cli"] = true
		aic.Requirements = append(aic.Requirements, "CLI_INTERFACE")
	}

	// Detectar exclusões explícitas com maior precisão
	exclusionPatterns := []string{
		"sem teste", "sem test", "não teste", "no test", "skip test",
		"sem docker", "sem container", "no docker", "skip docker",
		"sem banco", "sem db", "sem database", "no database", "skip database",
		"sem frontend", "sem ui", "no frontend", "skip frontend",
		"sem api", "no api", "skip api",
	}

	for _, pattern := range exclusionPatterns {
		if strings.Contains(desc, pattern) {
			aic.Constraints = append(aic.Constraints, "EXPLICIT_EXCLUSIONS")
			// Marcar o que deve ser excluído
			if strings.Contains(pattern, "teste") || strings.Contains(pattern, "test") {
				aic.Adaptations["exclude_tests"] = true
			}
			if strings.Contains(pattern, "docker") || strings.Contains(pattern, "container") {
				aic.Adaptations["exclude_docker"] = true
			}
			if strings.Contains(pattern, "banco") || strings.Contains(pattern, "database") || strings.Contains(pattern, "db") {
				aic.Adaptations["exclude_database"] = true
			}
			if strings.Contains(pattern, "frontend") || strings.Contains(pattern, "ui") {
				aic.Adaptations["exclude_frontend"] = true
			}
			if strings.Contains(pattern, "api") {
				aic.Adaptations["exclude_api"] = true
			}
		}
	}

	// Detectar foco específico no objetivo final
	if containsAny(desc, []string{"foco", "focus", "objetivo", "goal", "propósito", "purpose", "específico", "specific"}) {
		aic.Constraints = append(aic.Constraints, "ULTIMATE_GOAL_FOCUS")
		aic.Adaptations["chameleon_focus"] = true
	}

	// Detectar necessidade de adaptabilidade
	if containsAny(desc, []string{"adaptável", "adaptable", "flexível", "flexible", "camaleão", "chameleon", "dinâmico", "dynamic"}) {
		aic.Adaptations["adaptive_behavior"] = true
		aic.Requirements = append(aic.Requirements, "ADAPTIVE_STRUCTURE")
	}
}

// GetInstructionProfile retorna o perfil de instruções baseado no contexto
func (aic *AdaptiveInstructionController) GetInstructionProfile() InstructionProfile {
	profile := InstructionProfile{
		Name:             fmt.Sprintf("%s_%s_profile", aic.ProjectType, aic.Scope),
		Description:      fmt.Sprintf("Perfil adaptativo para %s em %s", aic.ProjectType, aic.Language),
		StrictnessLevel:  aic.calculateStrictnessLevel(),
		ScopeControl:     aic.Scope,
		FocusAreas:       aic.getFocusAreas(),
		ExclusionRules:   aic.getExclusionRules(),
		QualityThreshold: aic.calculateQualityThreshold(),
		AdaptiveRules:    aic.getAdaptiveRules(),
	}

	return profile
}

// calculateStrictnessLevel calcula o nível de rigidez das instruções
func (aic *AdaptiveInstructionController) calculateStrictnessLevel() int {
	base := 5 // Nível médio

	if aic.Scope == "minimal" {
		base = 8 // Muito rígido para escopo mínimo
	} else if aic.Scope == "comprehensive" {
		base = 6 // Moderadamente rígido para escopo abrangente
	}

	// Aumentar rigidez se há constraints explícitas
	if len(aic.Constraints) > 0 {
		base += 2
	}

	if base > 10 {
		base = 10
	}

	return base
}

// getFocusAreas retorna as áreas de foco baseadas nas adaptações
func (aic *AdaptiveInstructionController) getFocusAreas() []string {
	areas := []string{"core", "structure"}

	if aic.Adaptations["include_api"] == true {
		areas = append(areas, "api_layer")
	}

	if aic.Adaptations["include_frontend"] == true {
		areas = append(areas, "frontend_layer")
	}

	if aic.Adaptations["include_tests"] == true {
		areas = append(areas, "testing_layer")
	}

	if aic.Adaptations["include_database"] == true {
		areas = append(areas, "data_layer")
	}

	if aic.Adaptations["include_docker"] == true {
		areas = append(areas, "deployment_layer")
	}

	return areas
}

// getExclusionRules retorna regras de exclusão baseadas nos constraints
func (aic *AdaptiveInstructionController) getExclusionRules() []string {
	rules := make([]string, 0)

	if aic.Scope == "minimal" {
		rules = append(rules, "NO_EXTRA_FEATURES", "NO_OPTIONAL_COMPONENTS", "ESSENTIAL_ONLY")
	}

	if contains(aic.Constraints, "STRICT_MINIMAL_SCOPE") {
		rules = append(rules, "STRICT_SCOPE_ENFORCEMENT", "NO_FEATURE_CREEP")
	}

	if contains(aic.Constraints, "EXPLICIT_EXCLUSIONS") {
		rules = append(rules, "RESPECT_EXPLICIT_EXCLUSIONS")
	}

	return rules
}

// calculateQualityThreshold calcula o limiar de qualidade baseado no contexto
func (aic *AdaptiveInstructionController) calculateQualityThreshold() float64 {
	base := 70.0

	if aic.Scope == "minimal" {
		base = 80.0 // Maior qualidade para escopo mínimo (precisão)
	} else if aic.Scope == "comprehensive" {
		base = 60.0 // Menor limiar para escopo abrangente (flexibilidade)
	}

	return base
}

// getAdaptiveRules retorna regras adaptativas específicas
func (aic *AdaptiveInstructionController) getAdaptiveRules() map[string]string {
	rules := make(map[string]string)

	rules["scope_control"] = fmt.Sprintf("ENFORCE_%s_SCOPE", strings.ToUpper(aic.Scope))
	rules["language_adaptation"] = fmt.Sprintf("ADAPT_TO_%s_PATTERNS", strings.ToUpper(aic.Language))
	rules["purpose_alignment"] = "ALIGN_WITH_STATED_PURPOSE"

	if len(aic.Requirements) > 0 {
		rules["requirement_enforcement"] = "ENFORCE_SPECIFIC_REQUIREMENTS"
	}

	if len(aic.Constraints) > 0 {
		rules["constraint_compliance"] = "STRICT_CONSTRAINT_COMPLIANCE"
	}

	return rules
}

// BuildAdaptivePrompt constrói um prompt adaptativo baseado no perfil
func (aic *AdaptiveInstructionController) BuildAdaptivePrompt(basePrompt string) string {
	profile := aic.GetInstructionProfile()

	var prompt strings.Builder

	// Cabeçalho adaptativo
	prompt.WriteString(fmt.Sprintf(`ADAPTIVE INSTRUCTION CONTROL SYSTEM - PROFILE: %s

STRICTNESS LEVEL: %d/10
SCOPE CONTROL: %s
QUALITY THRESHOLD: %.1f%%

`, profile.Name, profile.StrictnessLevel, profile.ScopeControl, profile.QualityThreshold))

	// Instruções de escopo
	prompt.WriteString(aic.buildScopeInstructions(profile))

	// Prompt base
	prompt.WriteString(basePrompt)

	// Instruções adaptativas específicas
	prompt.WriteString(aic.buildAdaptiveInstructions(profile))

	// Regras de validação
	prompt.WriteString(aic.buildValidationRules(profile))

	return prompt.String()
}

// buildScopeInstructions constrói instruções específicas de escopo
func (aic *AdaptiveInstructionController) buildScopeInstructions(profile InstructionProfile) string {
	var instructions strings.Builder

	instructions.WriteString("SCOPE CONTROL INSTRUCTIONS:\n")

	switch profile.ScopeControl {
	case "minimal":
		instructions.WriteString(`
MINIMAL SCOPE MODE - CRITÉRIO DE CAMALEÃO ADAPTATIVO:
- FOQUE EXCLUSIVAMENTE no objetivo final declarado
- Implemente APENAS o que é ESSENCIAL para atingir o propósito
- NÃO adicione recursos extras, features opcionais ou "nice to have"
- ELIMINE qualquer componente que não seja DIRETAMENTE necessário
- Mantenha arquitetura mais simples e direta possível
- Cada arquivo deve ter JUSTIFICATIVA clara e específica no contexto
- ADAPTE-SE precisamente ao escopo solicitado sem expansões
- PRIORIZE funcionalidade sobre estrutura elaborada
- EVITE over-engineering, padrões complexos desnecessários
- FOQUE na entrega do valor core do projeto
`)
	case "comprehensive":
		instructions.WriteString(`
COMPREHENSIVE SCOPE MODE - CRITÉRIO DE CAMALEÃO ADAPTATIVO:
- Implemente uma solução COMPLETA e ROBUSTA
- Inclua TODAS as boas práticas relevantes para o domínio
- Adicione configurações profissionais de ambiente e deployment
- Implemente padrões de arquitetura escaláveis e maintíveis
- Inclua documentação técnica abrangente e exemplos
- Adicione testes unitários, integração e CI/CD
- Implemente logging, monitoring e error handling completos
- Adicione ferramentas de desenvolvimento e debugging
- ADAPTE-SE ao contexto para fornecer solução enterprise-grade
- EXPANDA funcionalidades com base no propósito declarado
`)
	default:
		instructions.WriteString(`
STANDARD SCOPE MODE - CRITÉRIO DE CAMALEÃO ADAPTATIVO:
- EQUILIBRE funcionalidade essencial com práticas recomendadas
- Inclua recursos padrão relevantes para linguagem/framework
- Mantenha estrutura organizada mas não excessivamente complexa
- Adicione documentação básica e configurações padrão
- Implemente error handling e logging básicos
- ADAPTE-SE ao contexto mantendo equilíbrio entre simplicidade e robustez
- FOQUE no propósito principal com suporte adequado
- Evite tanto minimalismo excessivo quanto over-engineering
`)
	}

	return instructions.String()
}

// buildAdaptiveInstructions constrói instruções adaptativas
func (aic *AdaptiveInstructionController) buildAdaptiveInstructions(profile InstructionProfile) string {
	var instructions strings.Builder

	instructions.WriteString("\nADAPTIVE BEHAVIOR RULES - SISTEMA CAMALEÃO:\n")

	// Adicionar regras de adaptabilidade específicas
	instructions.WriteString("PRINCÍPIO DO CAMALEÃO:\n")
	instructions.WriteString("- ADAPTE-SE precisamente ao contexto e propósito declarado\n")
	instructions.WriteString("- MUDE a estrutura e complexidade baseado no objetivo final\n")
	instructions.WriteString("- ELIMINE componentes desnecessários para o propósito específico\n")
	instructions.WriteString("- MANTENHA consistência e coerência em toda a solução\n")
	instructions.WriteString("- PRIORIZE o valor entregue sobre padrões genéricos\n")
	instructions.WriteString("- SEJA PRECISO no escopo sem adicionar elementos não solicitados\n")

	// Regras baseadas em adaptações específicas
	if aic.Adaptations["chameleon_focus"] == true {
		instructions.WriteString("\nFOCO CAMALEÃO ATIVADO:\n")
		instructions.WriteString("- CONCENTRE-SE exclusivamente no objetivo final declarado\n")
		instructions.WriteString("- REJEITE qualquer tentativa de expandir o escopo\n")
		instructions.WriteString("- MANTENHA laser focus no propósito específico\n")
	}

	if aic.Adaptations["adaptive_behavior"] == true {
		instructions.WriteString("\nCOMPORTAMENTO ADAPTATIVO ATIVADO:\n")
		instructions.WriteString("- AJUSTE a arquitetura baseada no contexto específico\n")
		instructions.WriteString("- MUDE padrões e estruturas conforme necessário\n")
		instructions.WriteString("- SEJA FLEXÍVEL na implementação mantendo coerência\n")
	}

	// Regras de foco
	if len(profile.FocusAreas) > 0 {
		instructions.WriteString(fmt.Sprintf("FOCUS AREAS: %s\n", strings.Join(profile.FocusAreas, ", ")))
	}

	// Regras de exclusão
	if len(profile.ExclusionRules) > 0 {
		instructions.WriteString("EXCLUSION RULES:\n")
		for _, rule := range profile.ExclusionRules {
			instructions.WriteString(fmt.Sprintf("- %s\n", rule))
		}
	}

	// Requisitos específicos
	if len(aic.Requirements) > 0 {
		instructions.WriteString("SPECIFIC REQUIREMENTS:\n")
		for _, req := range aic.Requirements {
			instructions.WriteString(fmt.Sprintf("- %s\n", req))
		}
	}

	// Constraints
	if len(aic.Constraints) > 0 {
		instructions.WriteString("STRICT CONSTRAINTS:\n")
		for _, constraint := range aic.Constraints {
			instructions.WriteString(fmt.Sprintf("- %s\n", constraint))
		}
	}

	// Adaptações específicas
	if len(aic.Adaptations) > 0 {
		instructions.WriteString("ADAPTIVE FEATURES:\n")
		for key, value := range aic.Adaptations {
			instructions.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
		}
	}

	return instructions.String()
}

// buildValidationRules constrói regras de validação específicas
func (aic *AdaptiveInstructionController) buildValidationRules(profile InstructionProfile) string {
	var rules strings.Builder

	rules.WriteString("\nVALIDATION RULES - SISTEMA CAMALEÃO:\n")
	rules.WriteString(fmt.Sprintf("QUALITY THRESHOLD: %.1f%%\n", profile.QualityThreshold))
	rules.WriteString(fmt.Sprintf("STRICTNESS LEVEL: %d/10\n", profile.StrictnessLevel))

	rules.WriteString("\nCOMPLIANCE CHECKS - VALIDAÇÃO CAMALEÃO:\n")
	rules.WriteString("1. VERIFICAR alinhamento preciso com o objetivo final\n")
	rules.WriteString("2. VALIDAR que TODOS os requisitos explícitos foram atendidos\n")
	rules.WriteString("3. CONFIRMAR que NENHUM constraint foi violado\n")
	rules.WriteString("4. GARANTIR que não há componentes desnecessários\n")
	rules.WriteString("5. VERIFICAR consistência e coerência entre todos os elementos\n")
	rules.WriteString("6. VALIDAR qualidade e funcionalidade do código gerado\n")
	rules.WriteString("7. ASSEGURAR que o escopo não foi expandido sem justificativa\n")
	rules.WriteString("8. CONFIRMAR que a solução é direta e focada no propósito\n")

	if profile.StrictnessLevel >= 8 {
		rules.WriteString("\nSTRICT CHAMELEON MODE ACTIVE:\n")
		rules.WriteString("- REJEITAR qualquer desvio do escopo rigorosamente definido\n")
		rules.WriteString("- APLICAR validação ultra-rigorosa de cada componente\n")
		rules.WriteString("- PRIORIZAR precisão absoluta sobre abrangência\n")
		rules.WriteString("- ELIMINAR qualquer elemento que não tenha justificativa direta\n")
		rules.WriteString("- MANTER foco laser no objetivo final declarado\n")
	}

	// Adicionar regras específicas baseadas no escopo
	if profile.ScopeControl == "minimal" {
		rules.WriteString("\nMINIMAL SCOPE VALIDATION:\n")
		rules.WriteString("- GARANTIR que apenas o essencial foi implementado\n")
		rules.WriteString("- REJEITAR qualquer feature extra ou 'nice to have'\n")
		rules.WriteString("- VALIDAR que cada arquivo tem justificativa clara\n")
	}

	return rules.String()
}

// ValidateInstructionCompliance valida se a resposta está em conformidade com as instruções
func (aic *AdaptiveInstructionController) ValidateInstructionCompliance(response string) (*InstructionComplianceResult, error) {
	profile := aic.GetInstructionProfile()

	result := &InstructionComplianceResult{
		IsCompliant:         true,
		ComplianceScore:     100.0,
		ViolatedRules:       make([]string, 0),
		MissingRequirements: make([]string, 0),
		ScopeDeviations:     make([]string, 0),
		QualityIssues:       make([]string, 0),
	}

	// Verificar estrutura JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
		result.IsCompliant = false
		result.ViolatedRules = append(result.ViolatedRules, "INVALID_JSON_STRUCTURE")
		result.ComplianceScore -= 50.0
		return result, nil
	}

	// Verificar requisitos específicos
	result = aic.checkRequirements(result, jsonData)

	// Verificar constraints
	result = aic.checkConstraints(result, jsonData)

	// Verificar escopo
	result = aic.checkScopeCompliance(result, jsonData)

	// Verificar qualidade
	result = aic.checkQualityStandards(result, jsonData, profile)

	// Determinar conformidade geral
	result.IsCompliant = result.ComplianceScore >= profile.QualityThreshold

	return result, nil
}

// InstructionComplianceResult representa o resultado da validação de conformidade
type InstructionComplianceResult struct {
	IsCompliant         bool     `json:"is_compliant"`
	ComplianceScore     float64  `json:"compliance_score"`
	ViolatedRules       []string `json:"violated_rules"`
	MissingRequirements []string `json:"missing_requirements"`
	ScopeDeviations     []string `json:"scope_deviations"`
	QualityIssues       []string `json:"quality_issues"`
}

// checkRequirements verifica se todos os requisitos foram atendidos
func (aic *AdaptiveInstructionController) checkRequirements(result *InstructionComplianceResult, jsonData map[string]interface{}) *InstructionComplianceResult {
	for _, req := range aic.Requirements {
		if !aic.requirementMet(req, jsonData) {
			result.MissingRequirements = append(result.MissingRequirements, req)
			result.ComplianceScore -= 10.0
		}
	}
	return result
}

// checkConstraints verifica se todos os constraints foram respeitados
func (aic *AdaptiveInstructionController) checkConstraints(result *InstructionComplianceResult, jsonData map[string]interface{}) *InstructionComplianceResult {
	for _, constraint := range aic.Constraints {
		if aic.constraintViolated(constraint, jsonData) {
			result.ViolatedRules = append(result.ViolatedRules, constraint)
			result.ComplianceScore -= 15.0
		}
	}
	return result
}

// checkScopeCompliance verifica conformidade com o escopo definido
func (aic *AdaptiveInstructionController) checkScopeCompliance(result *InstructionComplianceResult, jsonData map[string]interface{}) *InstructionComplianceResult {
	if aic.Scope == "minimal" {
		// Verificar se não há recursos extras
		if aic.hasExtraFeatures(jsonData) {
			result.ScopeDeviations = append(result.ScopeDeviations, "EXTRA_FEATURES_DETECTED")
			result.ComplianceScore -= 20.0
		}
	}
	return result
}

// checkQualityStandards verifica padrões de qualidade
func (aic *AdaptiveInstructionController) checkQualityStandards(result *InstructionComplianceResult, jsonData map[string]interface{}, profile InstructionProfile) *InstructionComplianceResult {
	// Implementar verificações de qualidade específicas
	return result
}

// requirementMet verifica se um requisito específico foi atendido
func (aic *AdaptiveInstructionController) requirementMet(requirement string, jsonData map[string]interface{}) bool {
	// Implementar lógica específica para cada tipo de requisito
	switch requirement {
	case "API_ENDPOINTS":
		return aic.hasAPIEndpoints(jsonData)
	case "FRONTEND_INTERFACE":
		return aic.hasFrontendInterface(jsonData)
	case "COMPREHENSIVE_TESTING":
		return aic.hasComprehensiveTesting(jsonData)
	case "DATABASE_INTEGRATION":
		return aic.hasDatabaseIntegration(jsonData)
	case "CONTAINERIZATION":
		return aic.hasContainerization(jsonData)
	default:
		return true
	}
}

// constraintViolated verifica se um constraint foi violado
func (aic *AdaptiveInstructionController) constraintViolated(constraint string, jsonData map[string]interface{}) bool {
	switch constraint {
	case "STRICT_MINIMAL_SCOPE":
		return aic.hasUnnecessaryComponents(jsonData)
	case "EXPLICIT_EXCLUSIONS":
		return aic.hasExplicitlyExcludedComponents(jsonData)
	default:
		return false
	}
}

// hasExtraFeatures verifica se há recursos extras não solicitados
func (aic *AdaptiveInstructionController) hasExtraFeatures(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		// Verificar se há arquivos de teste quando não solicitados
		if aic.Adaptations["exclude_tests"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "test") ||
					strings.Contains(strings.ToLower(filename), "spec") ||
					strings.Contains(strings.ToLower(filename), ".test.") ||
					strings.Contains(strings.ToLower(filename), ".spec.") {
					return true
				}
			}
		}

		// Verificar se há arquivos Docker quando não solicitados
		if aic.Adaptations["exclude_docker"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "docker") ||
					strings.Contains(strings.ToLower(filename), "compose") {
					return true
				}
			}
		}

		// Verificar se há muitos arquivos para escopo mínimo
		if aic.Scope == "minimal" && len(files) > 10 {
			return true
		}
	}
	return false
}

// hasAPIEndpoints verifica se há endpoints de API
func (aic *AdaptiveInstructionController) hasAPIEndpoints(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		for filename, content := range files {
			if strings.Contains(strings.ToLower(filename), "route") ||
				strings.Contains(strings.ToLower(filename), "controller") ||
				strings.Contains(strings.ToLower(filename), "handler") ||
				strings.Contains(strings.ToLower(filename), "endpoint") {
				return true
			}

			// Verificar conteúdo do arquivo
			if contentStr, ok := content.(string); ok {
				contentLower := strings.ToLower(contentStr)
				if strings.Contains(contentLower, "router") ||
					strings.Contains(contentLower, "endpoint") ||
					strings.Contains(contentLower, "http") ||
					strings.Contains(contentLower, "rest") ||
					strings.Contains(contentLower, "api") {
					return true
				}
			}
		}
	}
	return false
}

// hasFrontendInterface verifica se há interface frontend
func (aic *AdaptiveInstructionController) hasFrontendInterface(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		for filename, content := range files {
			if strings.Contains(strings.ToLower(filename), "component") ||
				strings.Contains(strings.ToLower(filename), "page") ||
				strings.Contains(strings.ToLower(filename), "view") ||
				strings.Contains(strings.ToLower(filename), ".html") ||
				strings.Contains(strings.ToLower(filename), ".css") ||
				strings.Contains(strings.ToLower(filename), ".scss") ||
				strings.Contains(strings.ToLower(filename), ".jsx") ||
				strings.Contains(strings.ToLower(filename), ".tsx") ||
				strings.Contains(strings.ToLower(filename), ".vue") {
				return true
			}

			// Verificar conteúdo do arquivo
			if contentStr, ok := content.(string); ok {
				contentLower := strings.ToLower(contentStr)
				if strings.Contains(contentLower, "component") ||
					strings.Contains(contentLower, "render") ||
					strings.Contains(contentLower, "jsx") ||
					strings.Contains(contentLower, "html") ||
					strings.Contains(contentLower, "css") {
					return true
				}
			}
		}
	}
	return false
}

// hasComprehensiveTesting verifica se há testes abrangentes
func (aic *AdaptiveInstructionController) hasComprehensiveTesting(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		testFileCount := 0
		for filename := range files {
			if strings.Contains(strings.ToLower(filename), "test") ||
				strings.Contains(strings.ToLower(filename), "spec") ||
				strings.Contains(strings.ToLower(filename), ".test.") ||
				strings.Contains(strings.ToLower(filename), ".spec.") {
				testFileCount++
			}
		}
		return testFileCount >= 2 // Pelo menos 2 arquivos de teste para ser considerado "abrangente"
	}
	return false
}

// hasDatabaseIntegration verifica se há integração com banco de dados
func (aic *AdaptiveInstructionController) hasDatabaseIntegration(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		for filename, content := range files {
			if strings.Contains(strings.ToLower(filename), "database") ||
				strings.Contains(strings.ToLower(filename), "db") ||
				strings.Contains(strings.ToLower(filename), "model") ||
				strings.Contains(strings.ToLower(filename), "schema") ||
				strings.Contains(strings.ToLower(filename), "migration") {
				return true
			}

			// Verificar conteúdo do arquivo
			if contentStr, ok := content.(string); ok {
				contentLower := strings.ToLower(contentStr)
				if strings.Contains(contentLower, "database") ||
					strings.Contains(contentLower, "mongoose") ||
					strings.Contains(contentLower, "sequelize") ||
					strings.Contains(contentLower, "prisma") ||
					strings.Contains(contentLower, "sql") ||
					strings.Contains(contentLower, "mongodb") ||
					strings.Contains(contentLower, "mysql") ||
					strings.Contains(contentLower, "postgres") {
					return true
				}
			}
		}
	}
	return false
}

// hasContainerization verifica se há containerização
func (aic *AdaptiveInstructionController) hasContainerization(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		for filename := range files {
			if strings.Contains(strings.ToLower(filename), "docker") ||
				strings.Contains(strings.ToLower(filename), "compose") ||
				strings.Contains(strings.ToLower(filename), "k8s") ||
				strings.Contains(strings.ToLower(filename), "kubernetes") {
				return true
			}
		}
	}
	return false
}

// hasUnnecessaryComponents verifica se há componentes desnecessários
func (aic *AdaptiveInstructionController) hasUnnecessaryComponents(jsonData map[string]interface{}) bool {
	if aic.Scope != "minimal" {
		return false
	}

	// Para escopo mínimo, verificar se há componentes desnecessários
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		unnecessaryPatterns := []string{
			"example", "sample", "demo", "template", "boilerplate",
			"readme", "license", "changelog", "contributing",
			"eslint", "prettier", "husky", "commitlint",
		}

		for filename := range files {
			for _, pattern := range unnecessaryPatterns {
				if strings.Contains(strings.ToLower(filename), pattern) {
					return true
				}
			}
		}

		// Verificar se há muitos arquivos de configuração
		configFileCount := 0
		for filename := range files {
			if strings.Contains(strings.ToLower(filename), "config") ||
				strings.Contains(strings.ToLower(filename), ".env") ||
				strings.Contains(strings.ToLower(filename), "settings") {
				configFileCount++
			}
		}

		if configFileCount > 3 {
			return true
		}
	}
	return false
}

// hasExplicitlyExcludedComponents verifica se há componentes explicitamente excluídos
func (aic *AdaptiveInstructionController) hasExplicitlyExcludedComponents(jsonData map[string]interface{}) bool {
	if files, ok := jsonData["files"].(map[string]interface{}); ok {
		// Verificar exclusões específicas
		if aic.Adaptations["exclude_tests"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "test") ||
					strings.Contains(strings.ToLower(filename), "spec") {
					return true
				}
			}
		}

		if aic.Adaptations["exclude_docker"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "docker") {
					return true
				}
			}
		}

		if aic.Adaptations["exclude_database"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "database") ||
					strings.Contains(strings.ToLower(filename), "db") ||
					strings.Contains(strings.ToLower(filename), "model") {
					return true
				}
			}
		}

		if aic.Adaptations["exclude_frontend"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "component") ||
					strings.Contains(strings.ToLower(filename), "page") ||
					strings.Contains(strings.ToLower(filename), ".html") ||
					strings.Contains(strings.ToLower(filename), ".css") {
					return true
				}
			}
		}

		if aic.Adaptations["exclude_api"] == true {
			for filename := range files {
				if strings.Contains(strings.ToLower(filename), "route") ||
					strings.Contains(strings.ToLower(filename), "controller") ||
					strings.Contains(strings.ToLower(filename), "handler") {
					return true
				}
			}
		}
	}
	return false
}

// containsAny verifica se a string contém qualquer uma das substrings
func containsAny(str string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}

// contains verifica se um slice contém um elemento
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
