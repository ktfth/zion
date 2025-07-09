package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ContextAnalyzer analisa o contexto do projeto de forma mais inteligente
type ContextAnalyzer struct {
	projectPath      string
	detectedPatterns []DetectedPattern
	insights         []ProjectInsight
	workspaceData    *WorkspaceData
}

// DetectedPattern representa um padrão detectado no projeto
type DetectedPattern struct {
	Type        PatternType `json:"type"`
	Confidence  float64     `json:"confidence"`
	Evidence    []string    `json:"evidence"`
	Description string      `json:"description"`
	Suggestions []string    `json:"suggestions"`
}

// PatternType define os tipos de padrões detectáveis
type PatternType string

const (
	PatternMicroservice    PatternType = "microservice"
	PatternMonorepo        PatternType = "monorepo"
	PatternAPI             PatternType = "api"
	PatternCLI             PatternType = "cli"
	PatternWebapp          PatternType = "webapp"
	PatternLibrary         PatternType = "library"
	PatternPlugin          PatternType = "plugin"
	PatternDataPipeline    PatternType = "data_pipeline"
	PatternMachineLearning PatternType = "machine_learning"
	PatternBlockchain      PatternType = "blockchain"
	PatternGameDev         PatternType = "game_dev"
	PatternMobile          PatternType = "mobile"
	PatternDesktop         PatternType = "desktop"
	PatternIoT             PatternType = "iot"
	PatternDevOps          PatternType = "devops"
)

// ProjectInsight representa insights gerados sobre o projeto
type ProjectInsight struct {
	Type        InsightType `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Priority    int         `json:"priority"`
	Actionable  bool        `json:"actionable"`
	Actions     []string    `json:"actions"`
}

// InsightType define os tipos de insights
type InsightType string

const (
	InsightArchitecture  InsightType = "architecture"
	InsightPerformance   InsightType = "performance"
	InsightSecurity      InsightType = "security"
	InsightMaintenance   InsightType = "maintenance"
	InsightScalability   InsightType = "scalability"
	InsightBestPractices InsightType = "best_practices"
	InsightDependencies  InsightType = "dependencies"
	InsightTesting       InsightType = "testing"
	InsightDocumentation InsightType = "documentation"
	InsightDeployment    InsightType = "deployment"
)

// WorkspaceData contém dados do workspace para análise
type WorkspaceData struct {
	Files           map[string]FileInfo `json:"files"`
	Dependencies    []Dependency        `json:"dependencies"`
	Configuration   map[string]string   `json:"configuration"`
	GitHistory      []GitCommit         `json:"git_history"`
	Patterns        []DetectedPattern   `json:"patterns"`
	Metrics         WorkspaceMetrics    `json:"metrics"`
	ProjectType     string              `json:"project_type"`
	PrimaryLanguage string              `json:"primary_language"`
	TotalFiles      int                 `json:"total_files"`
	ComplexityScore float64             `json:"complexity_score"`
	Technologies    []string            `json:"technologies"`
}

// FileInfo contém informações sobre um arquivo
type FileInfo struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	Language     string    `json:"language"`
	Complexity   int       `json:"complexity"`
	LOC          int       `json:"loc"`
	IsGenerated  bool      `json:"is_generated"`
	IsTest       bool      `json:"is_test"`
	IsConfig     bool      `json:"is_config"`
	Dependencies []string  `json:"dependencies"`
}

// Dependency representa uma dependência do projeto
type Dependency struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Type        string  `json:"type"` // production, development, peer
	Popularity  float64 `json:"popularity"`
	Security    float64 `json:"security"`
	Maintenance float64 `json:"maintenance"`
}

// GitCommit representa um commit do git
type GitCommit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Date      time.Time `json:"date"`
	Files     []string  `json:"files"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
}

// WorkspaceMetrics contém métricas do workspace
type WorkspaceMetrics struct {
	TotalFiles      int            `json:"total_files"`
	TotalLOC        int            `json:"total_loc"`
	Languages       map[string]int `json:"languages"`
	Complexity      float64        `json:"complexity"`
	TestCoverage    float64        `json:"test_coverage"`
	DependencyCount int            `json:"dependency_count"`
	LastActivity    time.Time      `json:"last_activity"`
	FileTypes       map[string]int `json:"file_types"`
	DirectoryDepth  int            `json:"directory_depth"`
	LargestFile     string         `json:"largest_file"`
	MostComplexFile string         `json:"most_complex_file"`
	ConfigFiles     []string       `json:"config_files"`
	TestFiles       []string       `json:"test_files"`
	TechnicalDebt   float64        `json:"technical_debt"`
	Maintainability float64        `json:"maintainability"`
}

// NewContextAnalyzer cria um novo analisador de contexto
func NewContextAnalyzer(projectPath string) *ContextAnalyzer {
	return &ContextAnalyzer{
		projectPath:      projectPath,
		detectedPatterns: make([]DetectedPattern, 0),
		insights:         make([]ProjectInsight, 0),
		workspaceData:    &WorkspaceData{},
	}
}

// AnalyzeContext analisa o contexto completo do projeto
func (ca *ContextAnalyzer) AnalyzeContext() (*WorkspaceData, error) {
	fmt.Printf("🔍 Analisando contexto do projeto: %s\n", ca.projectPath)

	// 1. Analisar estrutura de arquivos
	if err := ca.analyzeFileStructure(); err != nil {
		return nil, fmt.Errorf("erro ao analisar estrutura de arquivos: %v", err)
	}

	// 2. Detectar padrões arquiteturais
	if err := ca.detectArchitecturalPatterns(); err != nil {
		return nil, fmt.Errorf("erro ao detectar padrões: %v", err)
	}

	// 3. Analisar dependências
	if err := ca.analyzeDependencies(); err != nil {
		return nil, fmt.Errorf("erro ao analisar dependências: %v", err)
	}

	// 4. Extrair insights
	if err := ca.extractInsights(); err != nil {
		return nil, fmt.Errorf("erro ao extrair insights: %v", err)
	}

	// 5. Calcular métricas
	if err := ca.calculateMetrics(); err != nil {
		return nil, fmt.Errorf("erro ao calcular métricas: %v", err)
	}

	// 6. Analisar histórico git (se disponível)
	if err := ca.analyzeGitHistory(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao analisar histórico git: %v\n", err)
	}

	ca.workspaceData.Patterns = ca.detectedPatterns

	// Populate derived fields for easy access
	ca.workspaceData.TotalFiles = ca.workspaceData.Metrics.TotalFiles
	ca.workspaceData.ComplexityScore = ca.workspaceData.Metrics.Complexity

	// Determine primary language
	var primaryLang string
	maxCount := 0
	for lang, count := range ca.workspaceData.Metrics.Languages {
		if count > maxCount {
			maxCount = count
			primaryLang = lang
		}
	}
	ca.workspaceData.PrimaryLanguage = primaryLang

	// Simple project type detection
	ca.workspaceData.ProjectType = ca.detectProjectType()

	// Extract technologies from dependencies and patterns
	ca.workspaceData.Technologies = ca.extractTechnologies()

	fmt.Printf("✅ Análise concluída: %d padrões detectados, %d insights gerados\n",
		len(ca.detectedPatterns), len(ca.insights))

	return ca.workspaceData, nil
}

// analyzeFileStructure analisa a estrutura de arquivos do projeto
func (ca *ContextAnalyzer) analyzeFileStructure() error {
	ca.workspaceData.Files = make(map[string]FileInfo)

	return filepath.Walk(ca.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignorar erros de acesso
		}

		// Ignorar diretórios e arquivos ocultos
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Ignorar arquivos muito grandes (>10MB)
		if info.Size() > 10*1024*1024 {
			return nil
		}

		relativePath, _ := filepath.Rel(ca.projectPath, path)
		relativePath = strings.ReplaceAll(relativePath, "\\", "/")

		fileInfo := FileInfo{
			Path:     relativePath,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			Language: ca.detectLanguage(relativePath),
			IsTest:   ca.isTestFile(relativePath),
			IsConfig: ca.isConfigFile(relativePath),
		}

		// Analisar complexidade e LOC para arquivos de código
		if ca.isCodeFile(relativePath) {
			content, err := os.ReadFile(path)
			if err == nil {
				fileInfo.LOC = ca.countLOC(string(content))
				fileInfo.Complexity = ca.calculateComplexity(string(content))
				fileInfo.IsGenerated = ca.isGeneratedFile(string(content))
			}
		}

		ca.workspaceData.Files[relativePath] = fileInfo
		return nil
	})
}

// detectArchitecturalPatterns detecta padrões arquiteturais no projeto
func (ca *ContextAnalyzer) detectArchitecturalPatterns() error {
	patterns := []struct {
		pattern PatternType
		checker func() (float64, []string)
	}{
		{PatternMicroservice, ca.checkMicroservicePattern},
		{PatternMonorepo, ca.checkMonorepoPattern},
		{PatternAPI, ca.checkAPIPattern},
		{PatternCLI, ca.checkCLIPattern},
		{PatternWebapp, ca.checkWebappPattern},
		{PatternLibrary, ca.checkLibraryPattern},
		{PatternPlugin, ca.checkPluginPattern},
		{PatternDataPipeline, ca.checkDataPipelinePattern},
		{PatternMachineLearning, ca.checkMLPattern},
		{PatternBlockchain, ca.checkBlockchainPattern},
		{PatternGameDev, ca.checkGameDevPattern},
		{PatternMobile, ca.checkMobilePattern},
		{PatternDesktop, ca.checkDesktopPattern},
		{PatternIoT, ca.checkIoTPattern},
		{PatternDevOps, ca.checkDevOpsPattern},
	}

	for _, p := range patterns {
		confidence, evidence := p.checker()
		if confidence > 0.3 { // Threshold para considerar padrão detectado
			ca.detectedPatterns = append(ca.detectedPatterns, DetectedPattern{
				Type:        p.pattern,
				Confidence:  confidence,
				Evidence:    evidence,
				Description: ca.getPatternDescription(p.pattern),
				Suggestions: ca.getPatternSuggestions(p.pattern),
			})
		}
	}

	// Ordenar por confiança
	sort.Slice(ca.detectedPatterns, func(i, j int) bool {
		return ca.detectedPatterns[i].Confidence > ca.detectedPatterns[j].Confidence
	})

	return nil
}

// checkMicroservicePattern verifica padrões de microserviços
func (ca *ContextAnalyzer) checkMicroservicePattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura de diretórios
	if ca.hasDirectoryPattern([]string{"services", "microservices", "apps"}) {
		confidence += 0.3
		evidence = append(evidence, "Estrutura de diretórios indica microserviços")
	}

	// Verificar arquivos de configuração
	if ca.hasFiles([]string{"docker-compose.yml", "docker-compose.yaml", "k8s", "kubernetes"}) {
		confidence += 0.2
		evidence = append(evidence, "Arquivos de orquestração encontrados")
	}

	// Verificar múltiplos pontos de entrada
	if ca.countMainFiles() > 1 {
		confidence += 0.2
		evidence = append(evidence, "Múltiplos pontos de entrada detectados")
	}

	// Verificar arquivos de API
	if ca.hasFiles([]string{"api", "gateway", "proxy"}) {
		confidence += 0.1
		evidence = append(evidence, "Componentes de API Gateway detectados")
	}

	// Verificar mensageria
	if ca.hasPatternInFiles([]string{"kafka", "rabbitmq", "redis", "nats", "grpc"}) {
		confidence += 0.1
		evidence = append(evidence, "Sistemas de mensageria detectados")
	}

	// Verificar monitoramento
	if ca.hasFiles([]string{"prometheus", "grafana", "jaeger", "zipkin"}) {
		confidence += 0.1
		evidence = append(evidence, "Ferramentas de monitoramento detectadas")
	}

	return confidence, evidence
}

// checkMonorepoPattern verifica padrões de monorepo
func (ca *ContextAnalyzer) checkMonorepoPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura típica de monorepo
	if ca.hasDirectoryPattern([]string{"packages", "apps", "libs", "modules"}) {
		confidence += 0.4
		evidence = append(evidence, "Estrutura típica de monorepo detectada")
	}

	// Verificar ferramentas de monorepo
	if ca.hasFiles([]string{"lerna.json", "nx.json", "rush.json", "workspace"}) {
		confidence += 0.3
		evidence = append(evidence, "Ferramentas de monorepo detectadas")
	}

	// Verificar múltiplos package.json ou go.mod
	packageFiles := ca.countPackageFiles()
	if packageFiles > 2 {
		confidence += 0.2
		evidence = append(evidence, fmt.Sprintf("%d arquivos de package encontrados", packageFiles))
	}

	// Verificar yarn workspaces ou npm workspaces
	if ca.hasPatternInFiles([]string{"workspaces", "workspace:"}) {
		confidence += 0.1
		evidence = append(evidence, "Workspaces detectados")
	}

	return confidence, evidence
}

// checkAPIPattern verifica padrões de API
func (ca *ContextAnalyzer) checkAPIPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura de API
	if ca.hasDirectoryPattern([]string{"api", "routes", "controllers", "handlers", "endpoints"}) {
		confidence += 0.3
		evidence = append(evidence, "Estrutura de API detectada")
	}

	// Verificar frameworks de API
	if ca.hasPatternInFiles([]string{"express", "fastify", "gin", "echo", "flask", "django", "spring"}) {
		confidence += 0.2
		evidence = append(evidence, "Frameworks de API detectados")
	}

	// Verificar arquivos de especificação
	if ca.hasFiles([]string{"swagger", "openapi", "api.yaml", "api.json"}) {
		confidence += 0.2
		evidence = append(evidence, "Especificações de API encontradas")
	}

	// Verificar middleware
	if ca.hasPatternInFiles([]string{"middleware", "cors", "auth", "rate-limit"}) {
		confidence += 0.1
		evidence = append(evidence, "Middleware de API detectado")
	}

	// Verificar testes de API
	if ca.hasPatternInFiles([]string{"supertest", "httptest", "rest-assured"}) {
		confidence += 0.1
		evidence = append(evidence, "Testes de API detectados")
	}

	// Verificar rotas RESTful
	if ca.hasPatternInFiles([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}) {
		confidence += 0.1
		evidence = append(evidence, "Rotas RESTful detectadas")
	}

	return confidence, evidence
}

// Implementar outros checkers de padrões...
func (ca *ContextAnalyzer) checkCLIPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura de CLI
	if ca.hasDirectoryPattern([]string{"cmd", "cli", "commands"}) {
		confidence += 0.3
		evidence = append(evidence, "Estrutura de CLI detectada")
	}

	// Verificar frameworks de CLI
	if ca.hasPatternInFiles([]string{"cobra", "cli", "click", "argparse", "commander"}) {
		confidence += 0.2
		evidence = append(evidence, "Frameworks de CLI detectados")
	}

	// Verificar arquivo main típico de CLI
	if ca.hasMainWithArgs() {
		confidence += 0.2
		evidence = append(evidence, "Função main com argumentos detectada")
	}

	// Verificar flags e argumentos
	if ca.hasPatternInFiles([]string{"flag", "arg", "option", "command"}) {
		confidence += 0.1
		evidence = append(evidence, "Sistema de flags/argumentos detectado")
	}

	return confidence, evidence
}

// checkWebappPattern verifica padrões de aplicação web
func (ca *ContextAnalyzer) checkWebappPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura de webapp
	if ca.hasDirectoryPattern([]string{"public", "static", "assets", "templates", "views"}) {
		confidence += 0.3
		evidence = append(evidence, "Estrutura de webapp detectada")
	}

	// Verificar frameworks web
	if ca.hasPatternInFiles([]string{"react", "vue", "angular", "svelte", "next", "nuxt"}) {
		confidence += 0.2
		evidence = append(evidence, "Frameworks frontend detectados")
	}

	// Verificar arquivos web
	if ca.hasFiles([]string{"index.html", "app.html", "main.html"}) {
		confidence += 0.2
		evidence = append(evidence, "Arquivos HTML principais encontrados")
	}

	// Verificar CSS/styling
	if ca.hasFiles([]string{".css", ".scss", ".sass", ".less", ".styl"}) {
		confidence += 0.1
		evidence = append(evidence, "Arquivos de estilo detectados")
	}

	// Verificar bundlers
	if ca.hasFiles([]string{"webpack", "vite", "rollup", "parcel"}) {
		confidence += 0.1
		evidence = append(evidence, "Bundlers detectados")
	}

	return confidence, evidence
}

// Implementar outros checkers...
func (ca *ContextAnalyzer) checkLibraryPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	// Verificar estrutura de biblioteca
	if ca.hasDirectoryPattern([]string{"lib", "src", "pkg"}) {
		confidence += 0.2
		evidence = append(evidence, "Estrutura de biblioteca detectada")
	}

	// Verificar arquivos de distribuição
	if ca.hasFiles([]string{"package.json", "setup.py", "go.mod", "Cargo.toml"}) {
		confidence += 0.2
		evidence = append(evidence, "Arquivos de distribuição encontrados")
	}

	// Verificar exports
	if ca.hasPatternInFiles([]string{"export", "module.exports", "pub fn", "public class"}) {
		confidence += 0.1
		evidence = append(evidence, "Exports públicos detectados")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkPluginPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasDirectoryPattern([]string{"plugins", "extensions", "addons"}) {
		confidence += 0.3
		evidence = append(evidence, "Estrutura de plugins detectada")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkDataPipelinePattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"airflow", "luigi", "prefect", "dagster", "pipeline"}) {
		confidence += 0.3
		evidence = append(evidence, "Ferramentas de pipeline detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkMLPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"tensorflow", "pytorch", "scikit", "pandas", "numpy"}) {
		confidence += 0.3
		evidence = append(evidence, "Bibliotecas de ML detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkBlockchainPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"solidity", "web3", "ethereum", "smart contract"}) {
		confidence += 0.3
		evidence = append(evidence, "Tecnologias blockchain detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkGameDevPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"unity", "unreal", "godot", "game", "sprite"}) {
		confidence += 0.3
		evidence = append(evidence, "Ferramentas de game dev detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkMobilePattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"android", "ios", "flutter", "react-native", "xamarin"}) {
		confidence += 0.3
		evidence = append(evidence, "Tecnologias mobile detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkDesktopPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"electron", "tauri", "qt", "gtk", "winforms"}) {
		confidence += 0.3
		evidence = append(evidence, "Tecnologias desktop detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkIoTPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasPatternInFiles([]string{"mqtt", "iot", "sensor", "arduino", "raspberry"}) {
		confidence += 0.3
		evidence = append(evidence, "Tecnologias IoT detectadas")
	}

	return confidence, evidence
}

func (ca *ContextAnalyzer) checkDevOpsPattern() (float64, []string) {
	confidence := 0.0
	evidence := []string{}

	if ca.hasFiles([]string{"Dockerfile", ".github", ".gitlab-ci", "jenkins", "terraform"}) {
		confidence += 0.3
		evidence = append(evidence, "Ferramentas DevOps detectadas")
	}

	return confidence, evidence
}

// Métodos auxiliares para detecção de padrões
func (ca *ContextAnalyzer) hasDirectoryPattern(patterns []string) bool {
	for filePath := range ca.workspaceData.Files {
		dir := filepath.Dir(filePath)
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(dir), pattern) {
				return true
			}
		}
	}
	return false
}

func (ca *ContextAnalyzer) hasFiles(patterns []string) bool {
	for filePath := range ca.workspaceData.Files {
		fileName := strings.ToLower(filepath.Base(filePath))
		for _, pattern := range patterns {
			if strings.Contains(fileName, pattern) {
				return true
			}
		}
	}
	return false
}

func (ca *ContextAnalyzer) hasPatternInFiles(patterns []string) bool {
	for filePath := range ca.workspaceData.Files {
		if !ca.isCodeFile(filePath) {
			continue
		}

		content, err := os.ReadFile(filepath.Join(ca.projectPath, filePath))
		if err != nil {
			continue
		}

		contentStr := strings.ToLower(string(content))
		for _, pattern := range patterns {
			if strings.Contains(contentStr, pattern) {
				return true
			}
		}
	}
	return false
}

func (ca *ContextAnalyzer) countMainFiles() int {
	count := 0
	for filePath := range ca.workspaceData.Files {
		fileName := strings.ToLower(filepath.Base(filePath))
		if fileName == "main.go" || fileName == "main.py" || fileName == "main.js" ||
			fileName == "index.js" || fileName == "app.js" || fileName == "server.js" {
			count++
		}
	}
	return count
}

func (ca *ContextAnalyzer) countPackageFiles() int {
	count := 0
	for filePath := range ca.workspaceData.Files {
		fileName := strings.ToLower(filepath.Base(filePath))
		if fileName == "package.json" || fileName == "go.mod" || fileName == "cargo.toml" ||
			fileName == "setup.py" || fileName == "pom.xml" {
			count++
		}
	}
	return count
}

func (ca *ContextAnalyzer) hasMainWithArgs() bool {
	for filePath := range ca.workspaceData.Files {
		if !ca.isCodeFile(filePath) {
			continue
		}

		content, err := os.ReadFile(filepath.Join(ca.projectPath, filePath))
		if err != nil {
			continue
		}

		contentStr := string(content)
		if strings.Contains(contentStr, "func main()") ||
			strings.Contains(contentStr, "def main(") ||
			strings.Contains(contentStr, "int main(") {
			return true
		}
	}
	return false
}

// Métodos auxiliares para análise de arquivos
func (ca *ContextAnalyzer) detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	languageMap := map[string]string{
		".go":         "go",
		".py":         "python",
		".js":         "javascript",
		".ts":         "typescript",
		".java":       "java",
		".cpp":        "cpp",
		".c":          "c",
		".cs":         "csharp",
		".rb":         "ruby",
		".php":        "php",
		".rs":         "rust",
		".kt":         "kotlin",
		".swift":      "swift",
		".dart":       "dart",
		".scala":      "scala",
		".r":          "r",
		".m":          "objective-c",
		".sql":        "sql",
		".html":       "html",
		".css":        "css",
		".scss":       "scss",
		".sass":       "sass",
		".less":       "less",
		".vue":        "vue",
		".jsx":        "jsx",
		".tsx":        "tsx",
		".json":       "json",
		".yaml":       "yaml",
		".yml":        "yaml",
		".xml":        "xml",
		".md":         "markdown",
		".sh":         "bash",
		".ps1":        "powershell",
		".bat":        "batch",
		".dockerfile": "dockerfile",
		".tf":         "terraform",
		".proto":      "protobuf",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	return "unknown"
}

func (ca *ContextAnalyzer) isTestFile(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))
	return strings.Contains(fileName, "test") ||
		strings.Contains(fileName, "spec") ||
		strings.Contains(filePath, "/test/") ||
		strings.Contains(filePath, "/tests/") ||
		strings.Contains(filePath, "/__tests__/")
}

func (ca *ContextAnalyzer) isConfigFile(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))
	configFiles := []string{
		"config", "configuration", "settings", "env", "dockerfile",
		"package.json", "go.mod", "cargo.toml", "setup.py", "pom.xml",
		"webpack", "vite", "rollup", "babel", "eslint", "prettier",
		"tsconfig", "jest", "karma", "mocha", "cypress",
	}

	for _, config := range configFiles {
		if strings.Contains(fileName, config) {
			return true
		}
	}

	return false
}

func (ca *ContextAnalyzer) isCodeFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	codeExtensions := []string{
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".cs",
		".rb", ".php", ".rs", ".kt", ".swift", ".dart", ".scala",
		".r", ".m", ".sql", ".html", ".css", ".scss", ".sass",
		".less", ".vue", ".jsx", ".tsx",
	}

	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}

	return false
}

func (ca *ContextAnalyzer) countLOC(content string) int {
	lines := strings.Split(content, "\n")
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			count++
		}
	}
	return count
}

func (ca *ContextAnalyzer) calculateComplexity(content string) int {
	// Análise simplificada de complexidade ciclomática
	complexity := 1 // Complexidade base

	// Padrões que aumentam complexidade
	patterns := []string{
		"if", "else", "for", "while", "switch", "case", "catch", "&&", "||",
		"?", ":", "break", "continue", "return", "throw", "try",
	}

	for _, pattern := range patterns {
		complexity += strings.Count(strings.ToLower(content), pattern)
	}

	return complexity
}

func (ca *ContextAnalyzer) isGeneratedFile(content string) bool {
	// Verificar marcadores de arquivos gerados
	markers := []string{
		"@generated", "AUTO-GENERATED", "automatically generated",
		"DO NOT EDIT", "Code generated", "This file was generated",
		"autogenerated", "auto-generated",
	}

	contentUpper := strings.ToUpper(content)
	for _, marker := range markers {
		if strings.Contains(contentUpper, strings.ToUpper(marker)) {
			return true
		}
	}

	return false
}

// analyzeDependencies analisa as dependências do projeto
func (ca *ContextAnalyzer) analyzeDependencies() error {
	ca.workspaceData.Dependencies = make([]Dependency, 0)

	// Analisar package.json
	if err := ca.analyzePackageJson(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao analisar package.json: %v\n", err)
	}

	// Analisar go.mod
	if err := ca.analyzeGoMod(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao analisar go.mod: %v\n", err)
	}

	// Analisar requirements.txt
	if err := ca.analyzeRequirements(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao analisar requirements.txt: %v\n", err)
	}

	return nil
}

func (ca *ContextAnalyzer) analyzePackageJson() error {
	packagePath := filepath.Join(ca.projectPath, "package.json")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(packagePath)
	if err != nil {
		return err
	}

	var packageData map[string]interface{}
	if err := json.Unmarshal(content, &packageData); err != nil {
		return err
	}

	// Analisar dependencies
	if deps, ok := packageData["dependencies"].(map[string]interface{}); ok {
		for name, version := range deps {
			ca.workspaceData.Dependencies = append(ca.workspaceData.Dependencies, Dependency{
				Name:    name,
				Version: version.(string),
				Type:    "production",
			})
		}
	}

	// Analisar devDependencies
	if devDeps, ok := packageData["devDependencies"].(map[string]interface{}); ok {
		for name, version := range devDeps {
			ca.workspaceData.Dependencies = append(ca.workspaceData.Dependencies, Dependency{
				Name:    name,
				Version: version.(string),
				Type:    "development",
			})
		}
	}

	return nil
}

func (ca *ContextAnalyzer) analyzeGoMod() error {
	goModPath := filepath.Join(ca.projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require") {
			// Parse require lines
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[1]
				version := parts[2]
				ca.workspaceData.Dependencies = append(ca.workspaceData.Dependencies, Dependency{
					Name:    name,
					Version: version,
					Type:    "production",
				})
			}
		}
	}

	return nil
}

func (ca *ContextAnalyzer) analyzeRequirements() error {
	reqPath := filepath.Join(ca.projectPath, "requirements.txt")
	if _, err := os.Stat(reqPath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(reqPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// Parse dependency line
			parts := strings.Split(line, "==")
			name := parts[0]
			version := ""
			if len(parts) > 1 {
				version = parts[1]
			}

			ca.workspaceData.Dependencies = append(ca.workspaceData.Dependencies, Dependency{
				Name:    name,
				Version: version,
				Type:    "production",
			})
		}
	}

	return nil
}

// extractInsights extrai insights do projeto
func (ca *ContextAnalyzer) extractInsights() error {
	// Insight de arquitetura
	if len(ca.detectedPatterns) > 0 {
		primary := ca.detectedPatterns[0]
		ca.insights = append(ca.insights, ProjectInsight{
			Type:        InsightArchitecture,
			Title:       fmt.Sprintf("Padrão Arquitetural: %s", primary.Type),
			Description: fmt.Sprintf("Projeto segue padrão %s com %.1f%% de confiança", primary.Type, primary.Confidence*100),
			Priority:    1,
			Actionable:  true,
			Actions:     primary.Suggestions,
		})
	}

	// Insight de dependências
	if len(ca.workspaceData.Dependencies) > 50 {
		ca.insights = append(ca.insights, ProjectInsight{
			Type:        InsightDependencies,
			Title:       "Muitas Dependências",
			Description: fmt.Sprintf("Projeto tem %d dependências, considere revisar", len(ca.workspaceData.Dependencies)),
			Priority:    2,
			Actionable:  true,
			Actions:     []string{"Revisar dependências desnecessárias", "Considerar alternativas mais leves"},
		})
	}

	// Insight de testes
	testFileCount := 0
	for _, file := range ca.workspaceData.Files {
		if file.IsTest {
			testFileCount++
		}
	}

	if testFileCount == 0 {
		ca.insights = append(ca.insights, ProjectInsight{
			Type:        InsightTesting,
			Title:       "Sem Testes",
			Description: "Projeto não possui testes automatizados",
			Priority:    1,
			Actionable:  true,
			Actions:     []string{"Implementar testes unitários", "Configurar CI/CD"},
		})
	}

	return nil
}

// calculateMetrics calcula métricas do workspace
func (ca *ContextAnalyzer) calculateMetrics() error {
	metrics := &WorkspaceMetrics{
		Languages:   make(map[string]int),
		FileTypes:   make(map[string]int),
		ConfigFiles: make([]string, 0),
		TestFiles:   make([]string, 0),
	}

	totalLOC := 0
	totalComplexity := 0
	codeFileCount := 0

	for path, file := range ca.workspaceData.Files {
		metrics.TotalFiles++

		// Contar linguagens
		if file.Language != "unknown" {
			metrics.Languages[file.Language]++
		}

		// Contar tipos de arquivo
		ext := filepath.Ext(path)
		if ext != "" {
			metrics.FileTypes[ext]++
		}

		// Arquivos de configuração
		if file.IsConfig {
			metrics.ConfigFiles = append(metrics.ConfigFiles, path)
		}

		// Arquivos de teste
		if file.IsTest {
			metrics.TestFiles = append(metrics.TestFiles, path)
		}

		// Métricas de código
		if ca.isCodeFile(path) {
			codeFileCount++
			totalLOC += file.LOC
			totalComplexity += file.Complexity
		}

		// Arquivo maior
		if metrics.LargestFile == "" || file.Size > ca.workspaceData.Files[metrics.LargestFile].Size {
			metrics.LargestFile = path
		}

		// Arquivo mais complexo
		if metrics.MostComplexFile == "" || file.Complexity > ca.workspaceData.Files[metrics.MostComplexFile].Complexity {
			metrics.MostComplexFile = path
		}

		// Última atividade
		if file.ModTime.After(metrics.LastActivity) {
			metrics.LastActivity = file.ModTime
		}
	}

	metrics.TotalLOC = totalLOC
	metrics.DependencyCount = len(ca.workspaceData.Dependencies)

	if codeFileCount > 0 {
		metrics.Complexity = float64(totalComplexity) / float64(codeFileCount)
	}

	// Calcular cobertura de testes (estimativa)
	if len(metrics.TestFiles) > 0 {
		metrics.TestCoverage = float64(len(metrics.TestFiles)) / float64(codeFileCount) * 100
	}

	// Calcular profundidade de diretórios
	maxDepth := 0
	for path := range ca.workspaceData.Files {
		depth := strings.Count(path, "/") + strings.Count(path, "\\")
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	metrics.DirectoryDepth = maxDepth

	// Calcular dívida técnica (estimativa baseada em complexidade)
	if metrics.Complexity > 20 {
		metrics.TechnicalDebt = (metrics.Complexity - 20) / 100
	}

	// Calcular manutenibilidade
	metrics.Maintainability = 100 - (metrics.TechnicalDebt * 10)
	if metrics.Maintainability < 0 {
		metrics.Maintainability = 0
	}

	ca.workspaceData.Metrics = *metrics
	return nil
}

// analyzeGitHistory analisa o histórico do git
func (ca *ContextAnalyzer) analyzeGitHistory() error {
	// Verificar se é um repositório git
	gitDir := filepath.Join(ca.projectPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil
	}

	// Implementar análise básica do git (simplificada)
	// Na prática, isso usaria a biblioteca git ou chamaria comandos git

	return nil
}

// getPatternDescription retorna descrição do padrão
func (ca *ContextAnalyzer) getPatternDescription(pattern PatternType) string {
	descriptions := map[PatternType]string{
		PatternMicroservice:    "Arquitetura de microserviços com serviços independentes",
		PatternMonorepo:        "Repositório único com múltiplos projetos",
		PatternAPI:             "API REST ou GraphQL",
		PatternCLI:             "Interface de linha de comando",
		PatternWebapp:          "Aplicação web com frontend",
		PatternLibrary:         "Biblioteca reutilizável",
		PatternPlugin:          "Sistema de plugins extensível",
		PatternDataPipeline:    "Pipeline de processamento de dados",
		PatternMachineLearning: "Projeto de machine learning",
		PatternBlockchain:      "Aplicação blockchain",
		PatternGameDev:         "Desenvolvimento de jogos",
		PatternMobile:          "Aplicação mobile",
		PatternDesktop:         "Aplicação desktop",
		PatternIoT:             "Projeto Internet of Things",
		PatternDevOps:          "Ferramentas e automação DevOps",
	}

	if desc, exists := descriptions[pattern]; exists {
		return desc
	}

	return "Padrão não identificado"
}

// getPatternSuggestions retorna sugestões para o padrão
func (ca *ContextAnalyzer) getPatternSuggestions(pattern PatternType) []string {
	suggestions := map[PatternType][]string{
		PatternMicroservice: {
			"Implementar service discovery",
			"Configurar load balancing",
			"Adicionar circuit breaker",
			"Implementar distributed tracing",
		},
		PatternMonorepo: {
			"Configurar workspace dependencies",
			"Implementar shared build tools",
			"Adicionar linting compartilhado",
			"Configurar CI/CD otimizado",
		},
		PatternAPI: {
			"Implementar rate limiting",
			"Adicionar documentação OpenAPI",
			"Configurar CORS",
			"Implementar autenticação JWT",
		},
		PatternCLI: {
			"Adicionar auto-completion",
			"Implementar progress bars",
			"Adicionar colored output",
			"Configurar man pages",
		},
		PatternWebapp: {
			"Implementar lazy loading",
			"Adicionar PWA features",
			"Configurar SEO",
			"Implementar error boundaries",
		},
		PatternLibrary: {
			"Adicionar TypeScript definitions",
			"Implementar tree shaking",
			"Configurar semantic versioning",
			"Adicionar usage examples",
		},
	}

	if sugg, exists := suggestions[pattern]; exists {
		return sugg
	}

	return []string{"Seguir best practices da linguagem"}
}

// GetPrimaryPattern retorna o padrão principal detectado
func (ca *ContextAnalyzer) GetPrimaryPattern() *DetectedPattern {
	if len(ca.detectedPatterns) > 0 {
		return &ca.detectedPatterns[0]
	}
	return nil
}

// GetInsights retorna os insights gerados
func (ca *ContextAnalyzer) GetInsights() []ProjectInsight {
	return ca.insights
}

// GenerateContextualPrompt gera um prompt contextual baseado na análise
func (ca *ContextAnalyzer) GenerateContextualPrompt(basePrompt string) string {
	if ca.workspaceData == nil {
		return basePrompt
	}

	contextInfo := fmt.Sprintf("\n\n🔍 CONTEXTO INTELIGENTE DETECTADO:\n")

	// Adicionar padrões detectados
	if len(ca.detectedPatterns) > 0 {
		contextInfo += fmt.Sprintf("📋 Padrões Arquiteturais:\n")
		for _, pattern := range ca.detectedPatterns {
			contextInfo += fmt.Sprintf("   • %s (%.1f%% confiança)\n", pattern.Type, pattern.Confidence*100)
		}
	}

	// Adicionar métricas relevantes
	if ca.workspaceData.Metrics.TotalFiles > 0 {
		contextInfo += fmt.Sprintf("\n📊 Métricas do Projeto:\n")
		contextInfo += fmt.Sprintf("   • %d arquivos, %d LOC\n", ca.workspaceData.Metrics.TotalFiles, ca.workspaceData.Metrics.TotalLOC)
		contextInfo += fmt.Sprintf("   • Complexidade: %.1f\n", ca.workspaceData.Metrics.Complexity)
		contextInfo += fmt.Sprintf("   • %d dependências\n", ca.workspaceData.Metrics.DependencyCount)
	}

	// Adicionar linguagens principais
	if len(ca.workspaceData.Metrics.Languages) > 0 {
		contextInfo += fmt.Sprintf("\n🗣️ Linguagens Principais:\n")
		for lang, count := range ca.workspaceData.Metrics.Languages {
			if count > 1 {
				contextInfo += fmt.Sprintf("   • %s (%d arquivos)\n", lang, count)
			}
		}
	}

	// Adicionar insights prioritários
	if len(ca.insights) > 0 {
		contextInfo += fmt.Sprintf("\n💡 Insights Relevantes:\n")
		for _, insight := range ca.insights {
			if insight.Priority <= 2 {
				contextInfo += fmt.Sprintf("   • %s\n", insight.Title)
			}
		}
	}

	return basePrompt + contextInfo
}

// detectProjectType tries to determine the project type based on files and structure
func (ca *ContextAnalyzer) detectProjectType() string {
	// Check for specific config files and patterns
	if _, exists := ca.workspaceData.Files["package.json"]; exists {
		return "nodejs"
	}
	if _, exists := ca.workspaceData.Files["go.mod"]; exists {
		return "go"
	}
	if _, exists := ca.workspaceData.Files["requirements.txt"]; exists {
		return "python"
	}
	if _, exists := ca.workspaceData.Files["setup.py"]; exists {
		return "python"
	}
	if _, exists := ca.workspaceData.Files["Cargo.toml"]; exists {
		return "rust"
	}
	if _, exists := ca.workspaceData.Files["pom.xml"]; exists {
		return "java"
	}
	if _, exists := ca.workspaceData.Files["build.gradle"]; exists {
		return "java"
	}

	// Check by primary language
	if ca.workspaceData.PrimaryLanguage != "" {
		return ca.workspaceData.PrimaryLanguage
	}

	return "unknown"
}

// extractTechnologies extracts technologies from dependencies and detected patterns
func (ca *ContextAnalyzer) extractTechnologies() []string {
	techSet := make(map[string]bool)
	technologies := make([]string, 0)

	// Extract from dependencies
	for _, dep := range ca.workspaceData.Dependencies {
		techSet[dep.Name] = true
	}

	// Extract from detected patterns
	for _, pattern := range ca.workspaceData.Patterns {
		if pattern.Type == "framework" || pattern.Type == "library" {
			// Use the description or type as technology name
			if pattern.Description != "" {
				techSet[pattern.Description] = true
			} else {
				techSet[string(pattern.Type)] = true
			}
		}
	}

	// Extract from file extensions
	for lang := range ca.workspaceData.Metrics.Languages {
		techSet[lang] = true
	}

	// Convert set to slice
	for tech := range techSet {
		technologies = append(technologies, tech)
	}

	return technologies
}
