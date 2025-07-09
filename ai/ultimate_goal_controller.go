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

	// Detectar indicadores de escopo mínimo/específico
	minimalIndicators := []string{
		"hello world", "hello", "hola mundo", "olá mundo",
		"apenas", "somente", "só", "simples", "básico", "mínimo",
		"específico", "exato", "direto", "clean", "limpo",
		"teste", "exemplo", "demo", "prova de conceito",
	}

	isMinimalProject := false
	for _, indicator := range minimalIndicators {
		if strings.Contains(desc, indicator) {
			isMinimalProject = true
			ugc.Scope = "minimal"
			ugc.Priority = 9
			break
		}
	}

	// Para projetos mínimos, ser extremamente restritivo
	if isMinimalProject {
		ugc.addMinimalProjectRules()
		return
	}

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
		ugc.addStrictExclusionRules()
	}

	// Exclusões específicas baseadas no objetivo
	ugc.addConditionalExclusions()
}

// addMinimalProjectRules adiciona regras para projetos mínimos
func (ugc *UltimateGoalController) addMinimalProjectRules() {
	// Para projetos mínimos, permitir apenas arquivos absolutamente essenciais
	ugc.RequiredFiles = []string{"main.go", "go.mod", "README.md"}
	ugc.RequiredDirs = []string{} // Nenhum diretório adicional necessário

	// Excluir praticamente tudo que não é essencial
	ugc.ExcludedFiles = []string{
		// Arquivos de configuração desnecessários
		"docker-compose.yml", "Dockerfile", ".dockerignore",
		"kubernetes.yaml", "helm-chart.yaml", "skaffold.yaml",
		"docker-compose.prod.yml", "docker-compose.dev.yml",

		// Arquivos de teste desnecessários
		"test_helpers.go", "integration_test.go", "benchmark_test.go",
		"performance_test.go", "e2e_test.go", "functional_test.go",
		"main_test.go", "handler_test.go", "service_test.go",

		// Arquivos de documentação desnecessários
		"docs.md", "examples.md", "CHANGELOG.md", "CONTRIBUTING.md",
		"CODE_OF_CONDUCT.md", "SECURITY.md", "API.md", "USAGE.md",

		// Arquivos de CI/CD desnecessários
		".github/workflows/ci.yml", ".github/workflows/cd.yml",
		".gitlab-ci.yml", ".travis.yml", "Jenkinsfile", "azure-pipelines.yml",

		// Arquivos de banco de dados desnecessários
		"database.go", "models.go", "migrations.go", "seeds.go",
		"schema.sql", "init.sql", "migrate.go", "connection.go",

		// Arquivos de configuração avançada
		"config.yaml", "config.json", "settings.toml", "env.example",
		".env", ".env.local", ".env.production", ".env.development",

		// Arquivos de API desnecessários
		"handlers.go", "routes.go", "middleware.go", "server.go",
		"swagger.yaml", "openapi.yaml", "api.go", "rest.go",

		// Arquivos de frontend desnecessários
		"index.html", "style.css", "app.js", "package.json",
		"webpack.config.js", "tsconfig.json", "tailwind.config.js",

		// Arquivos de utilitários desnecessários
		"utils.go", "helpers.go", "constants.go", "types.go",
		"errors.go", "logger.go", "validator.go", "auth.go",

		// Arquivos de monitoramento desnecessários
		"metrics.go", "monitoring.go", "health.go", "prometheus.go",
		"grafana.json", "alertmanager.yml", "logs.go",

		// Arquivos de cache desnecessários
		"cache.go", "redis.go", "memcached.go", "storage.go",

		// Arquivos de build desnecessários
		"Makefile", "build.sh", "deploy.sh", "install.sh",
		"goreleaser.yml", "release.yml", "version.go",
	}

	ugc.ExcludedDirs = []string{
		// Diretórios de desenvolvimento desnecessários
		"examples", "benchmarks", "performance", "scripts",
		"tools", "build", "dist", "vendor", "node_modules",

		// Diretórios de infraestrutura desnecessários
		"kubernetes", "helm", "docker", "deployment", "infra",
		"terraform", "ansible", "vagrant", "compose",

		// Diretórios de dados desnecessários
		"migrations", "seeds", "fixtures", "data", "sql",
		"database", "models", "entities", "repositories",

		// Diretórios de configuração desnecessários
		"config", "configs", "settings", "environments",
		"secrets", "certs", "keys", "ssl",

		// Diretórios de documentação desnecessários
		"docs", "documentation", "wiki", "guides",
		"tutorials", "examples", "samples", "demos",

		// Diretórios de teste desnecessários
		"tests", "test", "testing", "spec", "specs",
		"integration", "e2e", "functional", "unit",

		// Diretórios de frontend desnecessários
		"frontend", "web", "static", "assets", "public",
		"components", "pages", "views", "templates",

		// Diretórios de API desnecessários
		"api", "handlers", "routes", "middleware", "controllers",
		"services", "endpoints", "rest", "graphql",

		// Diretórios de utilitários desnecessários
		"utils", "helpers", "common", "shared", "libs",
		"packages", "modules", "plugins", "extensions",

		// Diretórios de monitoramento desnecessários
		"monitoring", "metrics", "logs", "observability",
		"telemetry", "tracing", "profiling", "debug",

		// Diretórios de cache desnecessários
		"cache", "redis", "memcached", "storage", "tmp",

		// Diretórios de CI/CD desnecessários
		".github", ".gitlab", "ci", "cd", "pipelines",
		"workflows", "actions", "jobs", "stages",
	}
}

// addStrictExclusionRules adiciona regras de exclusão rigorosas
func (ugc *UltimateGoalController) addStrictExclusionRules() {
	// Adicionar arquivos comuns que geralmente são desnecessários
	commonUnnecessaryFiles := []string{
		"docker-compose.yml", "Dockerfile", ".dockerignore",
		"test_helpers.go", "integration_test.go", "benchmark_test.go",
		"docs.md", "examples.md", "CHANGELOG.md", "CONTRIBUTING.md",
		"performance_test.go", "docker-compose.prod.yml", "docker-compose.dev.yml",
		"helm-chart.yaml", "kubernetes.yaml", "swagger.yaml", "openapi.yaml",
		".github/workflows/ci.yml", ".gitlab-ci.yml", ".travis.yml",
		"Jenkinsfile", "azure-pipelines.yml", "goreleaser.yml",
		"config.yaml", "config.json", "settings.toml", ".env",
		"Makefile", "build.sh", "deploy.sh", "install.sh",
		"monitoring.go", "metrics.go", "health.go", "prometheus.go",
		"cache.go", "redis.go", "memcached.go", "logger.go",
		"errors.go", "validator.go", "auth.go", "middleware.go",
	}

	commonUnnecessaryDirs := []string{
		"examples", "benchmarks", "performance", "scripts", "tools",
		"kubernetes", "helm", "docker", "deployment", "infra",
		"migrations", "seeds", "fixtures", "data", "sql",
		"docs", "documentation", "wiki", "guides", "tutorials",
		"tests", "test", "testing", "spec", "specs", "integration",
		"monitoring", "metrics", "logs", "observability", "telemetry",
		"cache", "redis", "memcached", "storage", "tmp",
		".github", ".gitlab", "ci", "cd", "pipelines", "workflows",
		"vendor", "node_modules", "build", "dist", "public",
	}

	ugc.ExcludedFiles = append(ugc.ExcludedFiles, commonUnnecessaryFiles...)
	ugc.ExcludedDirs = append(ugc.ExcludedDirs, commonUnnecessaryDirs...)
}

// addConditionalExclusions adiciona exclusões condicionais baseadas no foco
func (ugc *UltimateGoalController) addConditionalExclusions() {
	// Se não há foco em database, excluir arquivos relacionados
	if dbFocus, ok := ugc.Adaptations["database_focus"]; !ok || !dbFocus.(bool) {
		dbFiles := []string{"database.go", "models.go", "migrations.go", "schema.sql", "connection.go"}
		dbDirs := []string{"database", "models", "migrations", "sql", "data"}
		ugc.ExcludedFiles = append(ugc.ExcludedFiles, dbFiles...)
		ugc.ExcludedDirs = append(ugc.ExcludedDirs, dbDirs...)
	}

	// Se não há foco em frontend, excluir arquivos relacionados
	if frontendFocus, ok := ugc.Adaptations["frontend_focus"]; !ok || !frontendFocus.(bool) {
		frontendFiles := []string{"index.html", "style.css", "app.js", "package.json", "webpack.config.js"}
		frontendDirs := []string{"frontend", "web", "static", "assets", "public", "components"}
		ugc.ExcludedFiles = append(ugc.ExcludedFiles, frontendFiles...)
		ugc.ExcludedDirs = append(ugc.ExcludedDirs, frontendDirs...)
	}

	// Se não há foco em API, excluir arquivos relacionados
	if apiFocus, ok := ugc.Adaptations["api_focus"]; !ok || !apiFocus.(bool) {
		apiFiles := []string{"handlers.go", "routes.go", "middleware.go", "server.go", "swagger.yaml"}
		apiDirs := []string{"api", "handlers", "routes", "middleware", "controllers", "services"}
		ugc.ExcludedFiles = append(ugc.ExcludedFiles, apiFiles...)
		ugc.ExcludedDirs = append(ugc.ExcludedDirs, apiDirs...)
	}

	// Se não há foco em CLI, excluir arquivos relacionados
	if cliFocus, ok := ugc.Adaptations["cli_focus"]; !ok || !cliFocus.(bool) {
		cliFiles := []string{"cmd.go", "cli.go", "flags.go", "commands.go"}
		cliDirs := []string{"cmd", "cli", "commands"}
		ugc.ExcludedFiles = append(ugc.ExcludedFiles, cliFiles...)
		ugc.ExcludedDirs = append(ugc.ExcludedDirs, cliDirs...)
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
		prompt.WriteString("   🎯 MODO ULTRA MINIMALISTA ATIVADO\n")
		prompt.WriteString("   • FOQUE EXCLUSIVAMENTE no objetivo declarado\n")
		prompt.WriteString("   • ELIMINE qualquer componente não essencial\n")
		prompt.WriteString("   • MANTENHA arquitetura mais simples possível\n")
		prompt.WriteString("   • EVITE over-engineering e padrões complexos\n")
		prompt.WriteString("   • GERE APENAS arquivos que são ABSOLUTAMENTE necessários\n")
		prompt.WriteString("   • PREFIRA um único arquivo quando possível\n")
		prompt.WriteString("   • NÃO crie diretórios desnecessários\n")
		prompt.WriteString("   • NÃO inclua arquivos de configuração avançada\n")
		prompt.WriteString("   • NÃO inclua arquivos de teste a menos que explicitamente solicitado\n")
		prompt.WriteString("   • NÃO inclua arquivos de documentação além de README básico\n")
		prompt.WriteString("   • NÃO inclua arquivos de Docker/containerização\n")
		prompt.WriteString("   • NÃO inclua arquivos de CI/CD\n")
		prompt.WriteString("   • NÃO inclua arquivos de monitoramento/logging\n")
		prompt.WriteString("   • EXEMPLO: Hello World = apenas main.go + go.mod + README.md\n")
	case "focused":
		prompt.WriteString("   • CONCENTRE-SE no objetivo mas inclua suporte básico\n")
		prompt.WriteString("   • ADICIONE apenas recursos diretamente relacionados\n")
		prompt.WriteString("   • MANTENHA estrutura organizada mas não complexa\n")
		prompt.WriteString("   • EVITE recursos que não contribuem para o objetivo\n")
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
		// Filtrar arquivos com maior rigor
		if files, ok := projectStructure["files"].(map[string]interface{}); ok {
			filteredFiles := make(map[string]interface{})
			totalFiles := len(files)

			for filename, fileContent := range files {
				if ugc.shouldIncludeFile(filename) {
					filteredFiles[filename] = fileContent
				}
			}

			// Log da filtragem para projetos mínimos
			if ugc.Scope == "minimal" {
				filteredCount := len(filteredFiles)
				if filteredCount < totalFiles {
					fmt.Printf("🧹 Filtro rigoroso aplicado: %d arquivos removidos (%d → %d)\n",
						totalFiles-filteredCount, totalFiles, filteredCount)
				}
			}

			projectStructure["files"] = filteredFiles
		}

		// Filtrar diretórios com maior rigor
		if dirs, ok := projectStructure["directories"].([]interface{}); ok {
			filteredDirs := make([]interface{}, 0)
			totalDirs := len(dirs)

			for _, dir := range dirs {
				if dirName, ok := dir.(string); ok {
					if ugc.shouldIncludeDir(dirName) {
						filteredDirs = append(filteredDirs, dir)
					}
				}
			}

			// Log da filtragem para projetos mínimos
			if ugc.Scope == "minimal" {
				filteredCount := len(filteredDirs)
				if filteredCount < totalDirs {
					fmt.Printf("🧹 Filtro rigoroso aplicado: %d diretórios removidos (%d → %d)\n",
						totalDirs-filteredCount, totalDirs, filteredCount)
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
	// Verificar se está na lista de arquivos explicitamente excluídos
	for _, excluded := range ugc.ExcludedFiles {
		if filename == excluded || strings.Contains(filename, excluded) {
			return false
		}
	}

	// Verificar se está na lista de arquivos obrigatórios
	for _, required := range ugc.RequiredFiles {
		if filename == required || strings.Contains(filename, required) {
			return true
		}
	}

	// Para projetos mínimos, aplicar filtro mais rigoroso
	if ugc.Scope == "minimal" || ugc.Priority >= 8 {
		return ugc.isMinimalEssentialFile(filename)
	}

	// Para outros escopos, usar análise padrão
	return ugc.isEssentialFile(filename)
}

// shouldIncludeDir determina se um diretório deve ser incluído baseado no objetivo
func (ugc *UltimateGoalController) shouldIncludeDir(dirname string) bool {
	// Verificar se está na lista de diretórios explicitamente excluídos
	for _, excluded := range ugc.ExcludedDirs {
		if dirname == excluded || strings.Contains(dirname, excluded) {
			return false
		}
	}

	// Verificar se está na lista de diretórios obrigatórios
	for _, required := range ugc.RequiredDirs {
		if dirname == required || strings.Contains(dirname, required) {
			return true
		}
	}

	// Para projetos mínimos, aplicar filtro mais rigoroso
	if ugc.Scope == "minimal" || ugc.Priority >= 8 {
		return ugc.isMinimalEssentialDir(dirname)
	}

	// Para outros escopos, usar análise padrão
	return ugc.isEssentialDir(dirname)
}

// isEssentialFile determina se um arquivo é essencial para o objetivo
func (ugc *UltimateGoalController) isEssentialFile(filename string) bool {
	// Para projetos com escopo mínimo, ser extremamente seletivo
	if ugc.Scope == "minimal" {
		return ugc.isMinimalEssentialFile(filename)
	}

	// Arquivos sempre essenciais para projetos normais
	essentialFiles := []string{
		"main.go", "main.js", "main.py", "main.ts", "main.rs",
		"index.js", "index.ts", "index.html", "index.py",
		"app.go", "app.js", "app.py", "app.ts",
		"package.json", "go.mod", "requirements.txt", "Cargo.toml",
		"README.md", "LICENSE",
	}

	for _, essential := range essentialFiles {
		if filename == essential {
			return true
		}
	}

	// Análise baseada nas palavras-chave do objetivo
	return ugc.isKeywordRelatedFile(filename)
}

// isMinimalEssentialFile determina se um arquivo é essencial para projetos mínimos
func (ugc *UltimateGoalController) isMinimalEssentialFile(filename string) bool {
	// Para projetos mínimos, apenas arquivos absolutamente necessários
	absoluteEssentials := []string{
		"main.go", "main.js", "main.py", "main.ts", "main.rs",
		"index.js", "index.ts", "index.py",
		"app.go", "app.js", "app.py", "app.ts",
		"go.mod", "package.json", "requirements.txt", "Cargo.toml",
		"README.md",
	}

	for _, essential := range absoluteEssentials {
		if filename == essential {
			return true
		}
	}

	// Para projetos mínimos, não incluir arquivos adicionais baseados em palavras-chave
	// a menos que seja explicitamente necessário
	return ugc.isExplicitlyRequired(filename)
}

// isKeywordRelatedFile verifica se o arquivo está relacionado às palavras-chave do objetivo
func (ugc *UltimateGoalController) isKeywordRelatedFile(filename string) bool {
	if len(ugc.Keywords) == 0 {
		return false
	}

	lowerFilename := strings.ToLower(filename)
	for _, keyword := range ugc.Keywords {
		if strings.Contains(lowerFilename, keyword) {
			return true
		}
	}

	return false
}

// isExplicitlyRequired verifica se um arquivo é explicitamente necessário
func (ugc *UltimateGoalController) isExplicitlyRequired(filename string) bool {
	for _, required := range ugc.RequiredFiles {
		if filename == required || strings.Contains(filename, required) {
			return true
		}
	}
	return false
}

// isEssentialDir determina se um diretório é essencial para o objetivo
func (ugc *UltimateGoalController) isEssentialDir(dirname string) bool {
	// Para projetos com escopo mínimo, ser extremamente seletivo
	if ugc.Scope == "minimal" {
		return ugc.isMinimalEssentialDir(dirname)
	}

	// Diretórios sempre essenciais para projetos normais
	essentialDirs := []string{
		"src", "lib", "pkg", "cmd", "api", "main",
	}

	for _, essential := range essentialDirs {
		if dirname == essential {
			return true
		}
	}

	// Análise baseada nas palavras-chave do objetivo
	return ugc.isKeywordRelatedDir(dirname)
}

// isMinimalEssentialDir determina se um diretório é essencial para projetos mínimos
func (ugc *UltimateGoalController) isMinimalEssentialDir(dirname string) bool {
	// Para projetos mínimos, evitar diretórios desnecessários
	// A maioria dos projetos mínimos pode funcionar sem diretórios adicionais

	// Apenas diretórios absolutamente necessários
	absoluteEssentials := []string{
		"src", "main", // Apenas se realmente necessário
	}

	for _, essential := range absoluteEssentials {
		if dirname == essential {
			// Verificar se é explicitamente necessário
			return ugc.isExplicitlyRequiredDir(dirname)
		}
	}

	return false
}

// isKeywordRelatedDir verifica se o diretório está relacionado às palavras-chave do objetivo
func (ugc *UltimateGoalController) isKeywordRelatedDir(dirname string) bool {
	if len(ugc.Keywords) == 0 {
		return false
	}

	lowerDirname := strings.ToLower(dirname)
	for _, keyword := range ugc.Keywords {
		if strings.Contains(lowerDirname, keyword) {
			return true
		}
	}

	return false
}

// isExplicitlyRequiredDir verifica se um diretório é explicitamente necessário
func (ugc *UltimateGoalController) isExplicitlyRequiredDir(dirname string) bool {
	for _, required := range ugc.RequiredDirs {
		if dirname == required || strings.Contains(dirname, required) {
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
