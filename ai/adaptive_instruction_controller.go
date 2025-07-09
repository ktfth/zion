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

	// Detectar palavras-chave de escopo
	if containsAny(desc, []string{"apenas", "somente", "só", "minimo", "básico", "simples"}) {
		aic.Scope = "minimal"
		aic.Constraints = append(aic.Constraints, "STRICT_MINIMAL_SCOPE")
	} else if containsAny(desc, []string{"completo", "completa", "abrangente", "extenso", "robusto", "detalhado"}) {
		aic.Scope = "comprehensive"
		aic.Requirements = append(aic.Requirements, "COMPREHENSIVE_IMPLEMENTATION")
	}

	// Detectar requisitos específicos
	if containsAny(desc, []string{"api", "rest", "endpoint", "servidor"}) {
		aic.Adaptations["include_api"] = true
		aic.Requirements = append(aic.Requirements, "API_ENDPOINTS")
	}

	if containsAny(desc, []string{"frontend", "interface", "ui", "web"}) {
		aic.Adaptations["include_frontend"] = true
		aic.Requirements = append(aic.Requirements, "FRONTEND_INTERFACE")
	}

	if containsAny(desc, []string{"teste", "test", "tdd", "bdd"}) {
		aic.Adaptations["include_tests"] = true
		aic.Requirements = append(aic.Requirements, "COMPREHENSIVE_TESTING")
	}

	if containsAny(desc, []string{"docker", "container", "deployment", "deploy"}) {
		aic.Adaptations["include_docker"] = true
		aic.Requirements = append(aic.Requirements, "CONTAINERIZATION")
	}

	if containsAny(desc, []string{"banco", "database", "db", "persistencia"}) {
		aic.Adaptations["include_database"] = true
		aic.Requirements = append(aic.Requirements, "DATABASE_INTEGRATION")
	}

	// Detectar exclusões explícitas
	if containsAny(desc, []string{"sem teste", "sem test", "não teste", "sem docker", "sem banco"}) {
		aic.Constraints = append(aic.Constraints, "EXPLICIT_EXCLUSIONS")
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
MINIMAL SCOPE MODE - CRITÉRIO DE CAMALEÃO:
- Inclua APENAS o que foi EXPLICITAMENTE solicitado
- NÃO adicione recursos extras, mesmo que sejam "boas práticas"
- Foque na funcionalidade ESSENCIAL para o propósito declarado
- Mantenha a estrutura mais simples possível
- Evite over-engineering ou features desnecessárias
- Cada arquivo deve ter JUSTIFICATIVA direta no propósito
`)
	case "comprehensive":
		instructions.WriteString(`
COMPREHENSIVE SCOPE MODE - CRITÉRIO DE CAMALEÃO:
- Implemente uma solução completa e robusta
- Inclua boas práticas e recursos avançados
- Adicione configurações de ambiente profissionais
- Implemente padrões de arquitetura escaláveis
- Inclua documentação abrangente
- Adicione testes, CI/CD e ferramentas de desenvolvimento
`)
	default:
		instructions.WriteString(`
STANDARD SCOPE MODE - CRITÉRIO DE CAMALEÃO:
- Equilibre funcionalidade essencial com boas práticas
- Inclua recursos padrão para a linguagem/framework
- Mantenha estrutura organizada mas não excessiva
- Adicione documentação básica
- Inclua configurações padrão de desenvolvimento
`)
	}

	return instructions.String()
}

// buildAdaptiveInstructions constrói instruções adaptativas
func (aic *AdaptiveInstructionController) buildAdaptiveInstructions(profile InstructionProfile) string {
	var instructions strings.Builder

	instructions.WriteString("\nADAPTIVE BEHAVIOR RULES - CRITÉRIO DE CAMALEÃO:\n")

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

	rules.WriteString("\nVALIDATION RULES - CRITÉRIO DE CAMALEÃO:\n")
	rules.WriteString(fmt.Sprintf("QUALITY THRESHOLD: %.1f%%\n", profile.QualityThreshold))
	rules.WriteString(fmt.Sprintf("STRICTNESS LEVEL: %d/10\n", profile.StrictnessLevel))

	rules.WriteString("\nCOMPLIANCE CHECKS:\n")
	rules.WriteString("1. Verificar se TODOS os requisitos explícitos foram atendidos\n")
	rules.WriteString("2. Validar que NENHUM constraint foi violado\n")
	rules.WriteString("3. Confirmar alinhamento com o propósito declarado\n")
	rules.WriteString("4. Verificar consistência entre camadas\n")
	rules.WriteString("5. Validar qualidade do código gerado\n")

	if profile.StrictnessLevel >= 8 {
		rules.WriteString("\nSTRICT MODE ACTIVE:\n")
		rules.WriteString("- Rejeitar qualquer desvio do escopo definido\n")
		rules.WriteString("- Aplicar validação rigorosa de cada componente\n")
		rules.WriteString("- Priorizar precisão sobre abrangência\n")
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
	// Implementar lógica para detectar recursos extras
	return false
}

// hasAPIEndpoints verifica se há endpoints de API
func (aic *AdaptiveInstructionController) hasAPIEndpoints(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar endpoints de API
	return false
}

// hasFrontendInterface verifica se há interface frontend
func (aic *AdaptiveInstructionController) hasFrontendInterface(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar interface frontend
	return false
}

// hasComprehensiveTesting verifica se há testes abrangentes
func (aic *AdaptiveInstructionController) hasComprehensiveTesting(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar testes abrangentes
	return false
}

// hasDatabaseIntegration verifica se há integração com banco de dados
func (aic *AdaptiveInstructionController) hasDatabaseIntegration(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar integração com banco de dados
	return false
}

// hasContainerization verifica se há containerização
func (aic *AdaptiveInstructionController) hasContainerization(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar containerização
	return false
}

// hasUnnecessaryComponents verifica se há componentes desnecessários
func (aic *AdaptiveInstructionController) hasUnnecessaryComponents(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar componentes desnecessários
	return false
}

// hasExplicitlyExcludedComponents verifica se há componentes explicitamente excluídos
func (aic *AdaptiveInstructionController) hasExplicitlyExcludedComponents(jsonData map[string]interface{}) bool {
	// Implementar lógica para detectar componentes explicitamente excluídos
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
