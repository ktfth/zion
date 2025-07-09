package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// UltimateGoalController controla a geração baseada no objetivo final
type UltimateGoalController struct {
	Goal          string
	Intent        string
	Scope         string
	RequiredFiles []string
	RequiredDirs  []string
	ExcludedFiles []string
	ExcludedDirs  []string
	Keywords      []string
	Priority      int // 1-10 (10 = máxima prioridade)
	Adaptations   map[string]interface{}
}

// GoalAnalysis representa a análise do objetivo final
type GoalAnalysis struct {
	PrimaryGoal      string   `json:"primary_goal"`
	SecondaryGoals   []string `json:"secondary_goals"`
	KeyComponents    []string `json:"key_components"`
	RequiredFiles    []string `json:"required_files"`
	RequiredDirs     []string `json:"required_dirs"`
	UnnecessaryFiles []string `json:"unnecessary_files"`
	UnnecessaryDirs  []string `json:"unnecessary_dirs"`
	Confidence       float64  `json:"confidence"`
}

// NewUltimateGoalController cria um novo controlador baseado no objetivo final
func NewUltimateGoalController(description string) *UltimateGoalController {
	controller := &UltimateGoalController{
		Goal:          extractPrimaryGoal(description),
		Intent:        analyzeUserIntent(description),
		Scope:         determineOptimalScope(description),
		RequiredFiles: make([]string, 0),
		RequiredDirs:  make([]string, 0),
		ExcludedFiles: make([]string, 0),
		ExcludedDirs:  make([]string, 0),
		Keywords:      extractKeywords(description),
		Priority:      calculatePriority(description),
		Adaptations:   make(map[string]interface{}),
	}

	// Análise detalhada do objetivo
	controller.analyzeGoal(description)

	return controller
}

// extractPrimaryGoal extrai o objetivo principal da descrição
func extractPrimaryGoal(description string) string {
	desc := strings.ToLower(description)

	// Padrões para identificar o objetivo principal
	goalPatterns := []struct {
		pattern string
		weight  int
	}{
		{`criar\s+(?:um[a]?\s+)?([^\.]+)`, 5},
		{`desenvolver\s+(?:um[a]?\s+)?([^\.]+)`, 5},
		{`implementar\s+(?:um[a]?\s+)?([^\.]+)`, 5},
		{`construir\s+(?:um[a]?\s+)?([^\.]+)`, 5},
		{`fazer\s+(?:um[a]?\s+)?([^\.]+)`, 4},
		{`gerar\s+(?:um[a]?\s+)?([^\.]+)`, 4},
		{`build\s+(?:a\s+)?([^\.]+)`, 5},
		{`create\s+(?:a\s+)?([^\.]+)`, 5},
		{`develop\s+(?:a\s+)?([^\.]+)`, 5},
		{`implement\s+(?:a\s+)?([^\.]+)`, 5},
		{`make\s+(?:a\s+)?([^\.]+)`, 4},
	}

	bestMatch := ""
	bestWeight := 0

	for _, pattern := range goalPatterns {
		regex := regexp.MustCompile(pattern.pattern)
		matches := regex.FindStringSubmatch(desc)
		if len(matches) > 1 && pattern.weight > bestWeight {
			bestMatch = strings.TrimSpace(matches[1])
			bestWeight = pattern.weight
		}
	}

	if bestMatch != "" {
		return bestMatch
	}

	// Se não encontrou padrão específico, usar a primeira sentença
	sentences := strings.Split(description, ".")
	if len(sentences) > 0 {
		return strings.TrimSpace(sentences[0])
	}

	return description
}

// analyzeUserIntent analisa a intenção do usuário
func analyzeUserIntent(description string) string {
	desc := strings.ToLower(description)

	// Análise de intenções específicas
	if containsAny(desc, []string{"apenas", "somente", "só", "específico", "exato", "preciso"}) {
		return "specific_minimal"
	}

	if containsAny(desc, []string{"completo", "completa", "abrangente", "extenso", "detalhado"}) {
		return "comprehensive"
	}

	if containsAny(desc, []string{"rápido", "simples", "básico", "essencial", "direto"}) {
		return "quick_essential"
	}

	if containsAny(desc, []string{"profissional", "enterprise", "produção", "robusto"}) {
		return "professional"
	}

	if containsAny(desc, []string{"protótipo", "teste", "experimento", "prova de conceito"}) {
		return "prototype"
	}

	return "standard"
}

// determineOptimalScope determina o escopo ótimo baseado no objetivo
func determineOptimalScope(description string) string {
	desc := strings.ToLower(description)

	// Análise de escopo baseada em palavras-chave
	minimalScore := 0
	comprehensiveScore := 0

	minimalKeywords := []string{"apenas", "somente", "só", "mínimo", "básico", "simples", "essencial", "direto", "rápido", "específico"}
	comprehensiveKeywords := []string{"completo", "completa", "abrangente", "extenso", "detalhado", "robusto", "profissional", "enterprise", "produção", "avançado"}

	for _, keyword := range minimalKeywords {
		if strings.Contains(desc, keyword) {
			minimalScore++
		}
	}

	for _, keyword := range comprehensiveKeywords {
		if strings.Contains(desc, keyword) {
			comprehensiveScore++
		}
	}

	// Análise contextual adicional
	if containsAny(desc, []string{"api rest", "crud", "backend", "frontend", "database"}) {
		if minimalScore > comprehensiveScore {
			return "focused"
		} else if comprehensiveScore > minimalScore {
			return "extensive"
		}
	}

	if minimalScore > comprehensiveScore {
		return "minimal"
	} else if comprehensiveScore > minimalScore {
		return "comprehensive"
	}

	return "balanced"
}

// extractKeywords extrai palavras-chave importantes da descrição
func extractKeywords(description string) []string {
	desc := strings.ToLower(description)
	keywords := make([]string, 0)

	// Tecnologias e frameworks
	techKeywords := []string{
		"react", "vue", "angular", "svelte", "next", "nuxt",
		"express", "fastify", "nest", "koa", "hapi",
		"mongodb", "postgres", "mysql", "redis", "sqlite",
		"docker", "kubernetes", "aws", "azure", "gcp",
		"jwt", "oauth", "auth", "authentication", "authorization",
		"api", "rest", "graphql", "grpc", "websocket",
		"typescript", "javascript", "python", "go", "rust", "java",
		"frontend", "backend", "fullstack", "microservice",
		"cli", "web", "mobile", "desktop", "server",
		"test", "testing", "tdd", "bdd", "unit", "integration",
		"ci/cd", "pipeline", "deployment", "devops",
	}

	for _, keyword := range techKeywords {
		if strings.Contains(desc, keyword) {
			keywords = append(keywords, keyword)
		}
	}

	return keywords
}

// calculatePriority calcula a prioridade baseada na urgência e especificidade
func calculatePriority(description string) int {
	desc := strings.ToLower(description)
	priority := 5 // Base

	// Aumentar prioridade para palavras de urgência
	if containsAny(desc, []string{"urgente", "rápido", "imediato", "agora", "urgent", "asap", "quickly"}) {
		priority += 3
	}

	// Aumentar prioridade para especificidade
	if containsAny(desc, []string{"apenas", "somente", "só", "específico", "exato", "preciso"}) {
		priority += 2
	}

	// Diminuir prioridade para escopo amplo
	if containsAny(desc, []string{"completo", "abrangente", "extenso", "detalhado", "robusto"}) {
		priority -= 1
	}

	if priority > 10 {
		priority = 10
	}
	if priority < 1 {
		priority = 1
	}

	return priority
}

// analyzeGoal análisa o objetivo em detalhes
func (ugc *UltimateGoalController) analyzeGoal(description string) {
	desc := strings.ToLower(description)

	// Análise de componentes necessários baseada no objetivo
	if containsAny(desc, []string{"api", "rest", "endpoint", "servidor", "service"}) {
		ugc.RequiredFiles = append(ugc.RequiredFiles, "main.go", "handlers.go", "routes.go", "server.go")
		ugc.RequiredDirs = append(ugc.RequiredDirs, "api", "handlers", "routes")
		ugc.Adaptations["api_focus"] = true
	}

	if containsAny(desc, []string{"database", "banco", "persistencia", "storage", "db"}) {
		ugc.RequiredFiles = append(ugc.RequiredFiles, "database.go", "models.go", "migrations.go")
		ugc.RequiredDirs = append(ugc.RequiredDirs, "database", "models", "migrations")
		ugc.Adaptations["database_focus"] = true
	}

	if containsAny(desc, []string{"frontend", "web", "interface", "ui", "react", "vue"}) {
		ugc.RequiredFiles = append(ugc.RequiredFiles, "index.html", "main.js", "App.js", "package.json")
		ugc.RequiredDirs = append(ugc.RequiredDirs, "src", "components", "pages", "assets")
		ugc.Adaptations["frontend_focus"] = true
	}

	if containsAny(desc, []string{"cli", "command", "terminal", "script", "tool"}) {
		ugc.RequiredFiles = append(ugc.RequiredFiles, "main.go", "cmd.go", "cli.go")
		ugc.RequiredDirs = append(ugc.RequiredDirs, "cmd", "cli")
		ugc.Adaptations["cli_focus"] = true
	}

	// Análise de exclusões baseada no escopo
	if ugc.Scope == "minimal" || ugc.Priority >= 8 {
		// Excluir arquivos não essenciais para escopo mínimo
		ugc.ExcludedFiles = append(ugc.ExcludedFiles,
			"docker-compose.yml", "Dockerfile",
			"test_helpers.go", "integration_test.go",
			"docs.md", "examples.md", "CHANGELOG.md",
			"benchmark_test.go", "performance_test.go",
			"docker-compose.prod.yml", "docker-compose.dev.yml",
			"helm-chart.yaml", "kubernetes.yaml",
			"swagger.yaml", "openapi.yaml",
		)
		ugc.ExcludedDirs = append(ugc.ExcludedDirs,
			"examples", "benchmarks", "performance",
			"kubernetes", "helm", "docker", "deployment",
			"migrations", "seeds", "fixtures",
		)
	}

	// Exclusões específicas baseadas no objetivo
	if dbFocus, ok := ugc.Adaptations["database_focus"]; !ok || !dbFocus.(bool) {
		if ugc.Scope == "minimal" {
			ugc.ExcludedFiles = append(ugc.ExcludedFiles, "database.go", "models.go", "migrations.go")
			ugc.ExcludedDirs = append(ugc.ExcludedDirs, "database", "models", "migrations")
		}
	}

	if frontendFocus, ok := ugc.Adaptations["frontend_focus"]; !ok || !frontendFocus.(bool) {
		if ugc.Scope == "minimal" {
			ugc.ExcludedFiles = append(ugc.ExcludedFiles, "index.html", "style.css", "app.js")
			ugc.ExcludedDirs = append(ugc.ExcludedDirs, "frontend", "web", "static", "assets")
		}
	}
}

// BuildGoalFocusedPrompt constrói um prompt focado no objetivo final
func (ugc *UltimateGoalController) BuildGoalFocusedPrompt(basePrompt string) string {
	var prompt strings.Builder

	prompt.WriteString(basePrompt)
	prompt.WriteString("\n\n")

	// Instruções de foco no objetivo
	prompt.WriteString("🎯 ULTIMATE GOAL FOCUS SYSTEM - SISTEMA DE FOCO NO OBJETIVO FINAL\n")
	prompt.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	prompt.WriteString(fmt.Sprintf("OBJETIVO PRINCIPAL: %s\n", ugc.Goal))
	prompt.WriteString(fmt.Sprintf("INTENÇÃO DO USUÁRIO: %s\n", ugc.Intent))
	prompt.WriteString(fmt.Sprintf("ESCOPO DETERMINADO: %s\n", ugc.Scope))
	prompt.WriteString(fmt.Sprintf("PRIORIDADE: %d/10\n\n", ugc.Priority))

	// Regras específicas baseadas no objetivo
	prompt.WriteString("🔥 REGRAS CRÍTICAS DE FOCO:\n")
	prompt.WriteString("1. IMPLEMENTE APENAS o que é DIRETAMENTE necessário para o OBJETIVO PRINCIPAL\n")
	prompt.WriteString("2. CADA ARQUIVO deve ter JUSTIFICATIVA CLARA no contexto do objetivo\n")
	prompt.WriteString("3. ELIMINE recursos que não contribuem para o objetivo final\n")
	prompt.WriteString("4. PRIORIZE funcionalidade CORE sobre estrutura elaborada\n")
	prompt.WriteString("5. MANTENHA laser focus no propósito específico declarado\n\n")

	// Componentes obrigatórios
	if len(ugc.RequiredFiles) > 0 {
		prompt.WriteString("📋 ARQUIVOS OBRIGATÓRIOS (relacionados ao objetivo):\n")
		for _, file := range ugc.RequiredFiles {
			prompt.WriteString(fmt.Sprintf("   ✓ %s\n", file))
		}
		prompt.WriteString("\n")
	}

	if len(ugc.RequiredDirs) > 0 {
		prompt.WriteString("📁 DIRETÓRIOS OBRIGATÓRIOS (relacionados ao objetivo):\n")
		for _, dir := range ugc.RequiredDirs {
			prompt.WriteString(fmt.Sprintf("   ✓ %s\n", dir))
		}
		prompt.WriteString("\n")
	}

	// Exclusões específicas
	if len(ugc.ExcludedFiles) > 0 {
		prompt.WriteString("🚫 ARQUIVOS PROIBIDOS (não relacionados ao objetivo):\n")
		for _, file := range ugc.ExcludedFiles {
			prompt.WriteString(fmt.Sprintf("   ✗ %s\n", file))
		}
		prompt.WriteString("\n")
	}

	if len(ugc.ExcludedDirs) > 0 {
		prompt.WriteString("🚫 DIRETÓRIOS PROIBIDOS (não relacionados ao objetivo):\n")
		for _, dir := range ugc.ExcludedDirs {
			prompt.WriteString(fmt.Sprintf("   ✗ %s\n", dir))
		}
		prompt.WriteString("\n")
	}

	// Palavras-chave para foco
	if len(ugc.Keywords) > 0 {
		prompt.WriteString("🔍 PALAVRAS-CHAVE DE FOCO:\n")
		prompt.WriteString(fmt.Sprintf("   %s\n\n", strings.Join(ugc.Keywords, ", ")))
	}

	// Instruções específicas do escopo
	prompt.WriteString("📐 INSTRUÇÕES DE ESCOPO:\n")
	switch ugc.Scope {
	case "minimal":
		prompt.WriteString("   • FOQUE EXCLUSIVAMENTE no objetivo declarado\n")
		prompt.WriteString("   • ELIMINE qualquer componente não essencial\n")
		prompt.WriteString("   • MANTENHA arquitetura mais simples possível\n")
		prompt.WriteString("   • EVITE over-engineering e padrões complexos\n")
	case "focused":
		prompt.WriteString("   • CONCENTRE-SE no objetivo mas inclua suporte básico\n")
		prompt.WriteString("   • ADICIONE apenas recursos diretamente relacionados\n")
		prompt.WriteString("   • MANTENHA estrutura organizada mas não complexa\n")
	case "balanced":
		prompt.WriteString("   • EQUILIBRE funcionalidade essencial com boas práticas\n")
		prompt.WriteString("   • INCLUA recursos padrão relevantes\n")
		prompt.WriteString("   • MANTENHA equilíbrio entre simplicidade e robustez\n")
	case "comprehensive":
		prompt.WriteString("   • IMPLEMENTE solução completa baseada no objetivo\n")
		prompt.WriteString("   • INCLUA todas as boas práticas relevantes\n")
		prompt.WriteString("   • ADICIONE configurações profissionais\n")
	}

	prompt.WriteString("\n")

	// Validação de conformidade
	prompt.WriteString("✅ CRITÉRIOS DE VALIDAÇÃO:\n")
	prompt.WriteString("   • Todo arquivo deve ter propósito claro no contexto do objetivo\n")
	prompt.WriteString("   • Estrutura deve ser mínima mas funcional\n")
	prompt.WriteString("   • Não deve haver recursos desnecessários\n")
	prompt.WriteString("   • Código deve estar diretamente relacionado ao objetivo\n")
	prompt.WriteString("   • Arquitetura deve ser apropriada para o escopo\n\n")

	// Instrução final
	prompt.WriteString("🎯 LEMBRE-SE: O objetivo é entregar EXATAMENTE o que foi solicitado,\n")
	prompt.WriteString("    sem adições desnecessárias, mantendo foco laser no propósito final.\n")
	prompt.WriteString("    Seja um camaleão que se adapta precisamente ao objetivo!\n\n")

	return prompt.String()
}

// FilterGeneratedContent filtra o conteúdo gerado baseado no objetivo
func (ugc *UltimateGoalController) FilterGeneratedContent(content string) (string, error) {
	// Parse do JSON gerado
	var structure map[string]interface{}
	if err := json.Unmarshal([]byte(content), &structure); err != nil {
		return "", fmt.Errorf("erro ao fazer parse do conteúdo: %v", err)
	}

	// Filtrar arquivos e diretórios
	if projectStructure, ok := structure["structure"].(map[string]interface{}); ok {
		// Filtrar arquivos
		if files, ok := projectStructure["files"].(map[string]interface{}); ok {
			filteredFiles := make(map[string]interface{})

			for filename, fileContent := range files {
				if ugc.shouldIncludeFile(filename) {
					filteredFiles[filename] = fileContent
				}
			}

			projectStructure["files"] = filteredFiles
		}

		// Filtrar diretórios
		if dirs, ok := projectStructure["directories"].([]interface{}); ok {
			filteredDirs := make([]interface{}, 0)

			for _, dir := range dirs {
				if dirName, ok := dir.(string); ok {
					if ugc.shouldIncludeDir(dirName) {
						filteredDirs = append(filteredDirs, dir)
					}
				}
			}

			projectStructure["directories"] = filteredDirs
		}
	}

	// Converter de volta para JSON
	filteredContent, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao gerar JSON filtrado: %v", err)
	}

	return string(filteredContent), nil
}

// shouldIncludeFile determina se um arquivo deve ser incluído baseado no objetivo
func (ugc *UltimateGoalController) shouldIncludeFile(filename string) bool {
	// Verificar se está na lista de arquivos obrigatórios
	for _, required := range ugc.RequiredFiles {
		if filename == required || strings.Contains(filename, required) {
			return true
		}
	}

	// Verificar se está na lista de arquivos excluídos
	for _, excluded := range ugc.ExcludedFiles {
		if filename == excluded || strings.Contains(filename, excluded) {
			return false
		}
	}

	// Análise baseada no objetivo e escopo
	if ugc.Scope == "minimal" {
		// Para escopo mínimo, ser mais restritivo
		return ugc.isEssentialFile(filename)
	}

	// Para outros escopos, incluir se não estiver explicitamente excluído
	return true
}

// shouldIncludeDir determina se um diretório deve ser incluído baseado no objetivo
func (ugc *UltimateGoalController) shouldIncludeDir(dirname string) bool {
	// Verificar se está na lista de diretórios obrigatórios
	for _, required := range ugc.RequiredDirs {
		if dirname == required || strings.Contains(dirname, required) {
			return true
		}
	}

	// Verificar se está na lista de diretórios excluídos
	for _, excluded := range ugc.ExcludedDirs {
		if dirname == excluded || strings.Contains(dirname, excluded) {
			return false
		}
	}

	// Análise baseada no objetivo e escopo
	if ugc.Scope == "minimal" {
		// Para escopo mínimo, ser mais restritivo
		return ugc.isEssentialDir(dirname)
	}

	// Para outros escopos, incluir se não estiver explicitamente excluído
	return true
}

// isEssentialFile determina se um arquivo é essencial para o objetivo
func (ugc *UltimateGoalController) isEssentialFile(filename string) bool {
	// Arquivos sempre essenciais
	essentialFiles := []string{
		"main.go", "main.js", "main.py", "main.ts", "main.rs",
		"index.js", "index.ts", "index.html", "index.py",
		"app.go", "app.js", "app.py", "app.ts",
		"package.json", "go.mod", "requirements.txt", "Cargo.toml",
		"README.md", "LICENSE", "Makefile",
	}

	for _, essential := range essentialFiles {
		if filename == essential {
			return true
		}
	}

	// Análise baseada nas palavras-chave do objetivo
	lowerFilename := strings.ToLower(filename)
	for _, keyword := range ugc.Keywords {
		if strings.Contains(lowerFilename, keyword) {
			return true
		}
	}

	return false
}

// isEssentialDir determina se um diretório é essencial para o objetivo
func (ugc *UltimateGoalController) isEssentialDir(dirname string) bool {
	// Diretórios sempre essenciais
	essentialDirs := []string{
		"src", "lib", "pkg", "cmd", "api", "main",
	}

	for _, essential := range essentialDirs {
		if dirname == essential {
			return true
		}
	}

	// Análise baseada nas palavras-chave do objetivo
	lowerDirname := strings.ToLower(dirname)
	for _, keyword := range ugc.Keywords {
		if strings.Contains(lowerDirname, keyword) {
			return true
		}
	}

	return false
}

// GetGoalSummary retorna um resumo do objetivo analisado
func (ugc *UltimateGoalController) GetGoalSummary() GoalAnalysis {
	return GoalAnalysis{
		PrimaryGoal:      ugc.Goal,
		SecondaryGoals:   []string{}, // TODO: implementar análise de objetivos secundários
		KeyComponents:    ugc.Keywords,
		RequiredFiles:    ugc.RequiredFiles,
		RequiredDirs:     ugc.RequiredDirs,
		UnnecessaryFiles: ugc.ExcludedFiles,
		UnnecessaryDirs:  ugc.ExcludedDirs,
		Confidence:       ugc.calculateConfidence(),
	}
}

// calculateConfidence calcula a confiança na análise do objetivo
func (ugc *UltimateGoalController) calculateConfidence() float64 {
	confidence := 0.5 // Base

	// Aumentar confiança para objetivos específicos
	if ugc.Priority >= 8 {
		confidence += 0.3
	}

	// Aumentar confiança para palavras-chave claras
	if len(ugc.Keywords) > 0 {
		confidence += 0.2
	}

	// Aumentar confiança para escopo bem definido
	if ugc.Scope == "minimal" || ugc.Scope == "focused" {
		confidence += 0.2
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}
