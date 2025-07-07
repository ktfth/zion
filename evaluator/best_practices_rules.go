package evaluator

import (
	"fmt"
	"strings"
)

// BuildConfigurationRule verifica configurações de build
type BuildConfigurationRule struct{}

func (r *BuildConfigurationRule) Name() string {
	return "BuildConfiguration"
}

func (r *BuildConfigurationRule) Description() string {
	return "Verifica se as configurações de build estão adequadas"
}

func (r *BuildConfigurationRule) Category() Category {
	return CategoryConfiguration
}

func (r *BuildConfigurationRule) Weight() float64 {
	return 7.0
}

func (r *BuildConfigurationRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar configurações específicas por linguagem
	switch strings.ToLower(structure.Language) {
	case "javascript", "js", "typescript", "ts":
		buildIssues := r.evaluateJSBuild(structure)
		issues = append(issues, buildIssues...)
	case "python", "py":
		buildIssues := r.evaluatePythonBuild(structure)
		issues = append(issues, buildIssues...)
	case "go", "golang":
		buildIssues := r.evaluateGoBuild(structure)
		issues = append(issues, buildIssues...)
	case "java":
		buildIssues := r.evaluateJavaBuild(structure)
		issues = append(issues, buildIssues...)
	case "rust":
		buildIssues := r.evaluateRustBuild(structure)
		issues = append(issues, buildIssues...)
	}

	// Ajustar score baseado nos issues
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityCritical:
			score -= 0.5
		case SeverityHigh:
			score -= 0.3
		case SeverityMedium:
			score -= 0.2
		case SeverityLow:
			score -= 0.1
		}
	}

	if score < 0 {
		score = 0
	}

	return score, issues
}

func (r *BuildConfigurationRule) evaluateJSBuild(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se tem scripts de build no package.json
	if packageFile, exists := structure.Files["package.json"]; exists {
		if !strings.Contains(packageFile.Content, "\"scripts\"") {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityMedium,
				Category:    CategoryConfiguration,
				Description: "package.json não define scripts de build",
				Location:    "package.json",
				Suggestion:  "Adicione scripts como 'build', 'start', 'test' ao package.json",
			})
		} else {
			// Verificar scripts essenciais
			content := packageFile.Content
			essentialScripts := []string{"start", "build"}

			for _, script := range essentialScripts {
				if !strings.Contains(content, fmt.Sprintf("\"%s\":", script)) {
					issues = append(issues, Issue{
						Type:        IssueConfiguration,
						Severity:    SeverityLow,
						Category:    CategoryConfiguration,
						Description: fmt.Sprintf("Script '%s' não encontrado no package.json", script),
						Location:    "package.json",
						Suggestion:  fmt.Sprintf("Adicione script '%s' para facilitar o desenvolvimento", script),
					})
				}
			}
		}
	}

	// Verificar ferramentas de build modernas
	buildTools := []string{"webpack.config.js", "vite.config.js", "rollup.config.js", "esbuild.config.js"}
	hasBuildTool := false

	for _, tool := range buildTools {
		if _, exists := structure.Files[tool]; exists {
			hasBuildTool = true
			break
		}
	}

	if !hasBuildTool && len(structure.Dependencies) > 5 {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityLow,
			Category:    CategoryConfiguration,
			Description: "Projeto complexo sem ferramenta de build configurada",
			Location:    "root",
			Suggestion:  "Considere usar Webpack, Vite, ou outra ferramenta de build moderna",
		})
	}

	return issues
}

func (r *BuildConfigurationRule) evaluatePythonBuild(structure *ProjectStructure) []Issue {
	var issues []Issue

	hasSetupPy := false
	hasPyprojectToml := false

	if _, exists := structure.Files["setup.py"]; exists {
		hasSetupPy = true
	}
	if _, exists := structure.Files["pyproject.toml"]; exists {
		hasPyprojectToml = true
	}

	// Para projetos com dependências, deve ter configuração de build
	if len(structure.Dependencies) > 0 && !hasSetupPy && !hasPyprojectToml {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityMedium,
			Category:    CategoryConfiguration,
			Description: "Projeto Python sem configuração de build/packaging",
			Location:    "root",
			Suggestion:  "Adicione setup.py ou pyproject.toml para configurar o projeto",
		})
	}

	// Verificar se tem Makefile para automação
	if _, exists := structure.Files["Makefile"]; !exists {
		if len(structure.Dependencies) > 3 {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityLow,
				Category:    CategoryConfiguration,
				Description: "Projeto complexo sem Makefile para automação",
				Location:    "root",
				Suggestion:  "Considere criar um Makefile para automatizar tarefas comuns",
			})
		}
	}

	return issues
}

func (r *BuildConfigurationRule) evaluateGoBuild(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se tem Makefile
	if _, exists := structure.Files["Makefile"]; !exists {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityLow,
			Category:    CategoryConfiguration,
			Description: "Projeto Go sem Makefile para automação",
			Location:    "root",
			Suggestion:  "Considere criar um Makefile para tarefas de build, test, e deploy",
		})
	}

	// Verificar se tem arquivos de CI/CD
	ciFiles := []string{".github/workflows", "Dockerfile", "docker-compose.yml"}
	hasCICD := false

	for _, ciFile := range ciFiles {
		if _, exists := structure.Files[ciFile]; exists {
			hasCICD = true
			break
		}
	}

	if !hasCICD && len(structure.Dependencies) > 2 {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityLow,
			Category:    CategoryConfiguration,
			Description: "Projeto sem configuração de CI/CD",
			Location:    "root",
			Suggestion:  "Considere adicionar GitHub Actions, Dockerfile ou docker-compose.yml",
		})
	}

	return issues
}

func (r *BuildConfigurationRule) evaluateJavaBuild(structure *ProjectStructure) []Issue {
	var issues []Issue

	hasMaven := false
	hasGradle := false

	if _, exists := structure.Files["pom.xml"]; exists {
		hasMaven = true
	}
	if _, exists := structure.Files["build.gradle"]; exists {
		hasGradle = true
	}

	if !hasMaven && !hasGradle {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityHigh,
			Category:    CategoryConfiguration,
			Description: "Projeto Java sem ferramenta de build (Maven/Gradle)",
			Location:    "root",
			Suggestion:  "Configure Maven (pom.xml) ou Gradle (build.gradle)",
		})
	}

	// Se tem Maven, verificar se tem wrapper
	if hasMaven {
		if _, exists := structure.Files["mvnw"]; !exists {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityLow,
				Category:    CategoryConfiguration,
				Description: "Projeto Maven sem wrapper (mvnw)",
				Location:    "root",
				Suggestion:  "Adicione Maven wrapper para facilitar builds",
			})
		}
	}

	// Se tem Gradle, verificar se tem wrapper
	if hasGradle {
		if _, exists := structure.Files["gradlew"]; !exists {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityLow,
				Category:    CategoryConfiguration,
				Description: "Projeto Gradle sem wrapper (gradlew)",
				Location:    "root",
				Suggestion:  "Adicione Gradle wrapper para facilitar builds",
			})
		}
	}

	return issues
}

func (r *BuildConfigurationRule) evaluateRustBuild(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar Cargo.toml
	if cargoFile, exists := structure.Files["Cargo.toml"]; exists {
		// Verificar se tem configurações de release otimizadas
		if !strings.Contains(cargoFile.Content, "[profile.release]") {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityLow,
				Category:    CategoryConfiguration,
				Description: "Cargo.toml sem perfil de release otimizado",
				Location:    "Cargo.toml",
				Suggestion:  "Adicione seção [profile.release] com otimizações",
			})
		}
	}

	return issues
}

// BestPracticesRule verifica aderência a melhores práticas gerais
type BestPracticesRule struct{}

func (r *BestPracticesRule) Name() string {
	return "BestPractices"
}

func (r *BestPracticesRule) Description() string {
	return "Verifica aderência a melhores práticas de desenvolvimento"
}

func (r *BestPracticesRule) Category() Category {
	return CategoryMaintainability
}

func (r *BestPracticesRule) Weight() float64 {
	return 10.0
}

func (r *BestPracticesRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar presença de arquivos importantes
	importantFiles := []string{".gitignore", "README.md", "LICENSE"}
	for _, file := range importantFiles {
		if !r.hasFile(structure, file) {
			severity := SeverityMedium
			if file == ".gitignore" {
				severity = SeverityHigh
			} else if file == "LICENSE" {
				severity = SeverityLow
			}

			issues = append(issues, Issue{
				Type:        IssueBestPractice,
				Severity:    severity,
				Category:    CategoryMaintainability,
				Description: fmt.Sprintf("Arquivo '%s' não encontrado", file),
				Suggestion:  fmt.Sprintf("Adicione arquivo '%s' ao projeto", file),
			})
		}
	}

	// Verificar estrutura de testes
	if !r.hasTestDirectory(structure) {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityMedium,
			Category:    CategoryMaintainability,
			Description: "Projeto sem estrutura de testes aparente",
			Suggestion:  "Crie diretório/arquivos de teste apropriados para a linguagem",
		})
	}

	// Verificar se há documentação além do README
	if !r.hasDocumentation(structure) {
		if len(structure.Files) > 10 { // Projetos maiores precisam mais documentação
			issues = append(issues, Issue{
				Type:        IssueBestPractice,
				Severity:    SeverityLow,
				Category:    CategoryMaintainability,
				Description: "Projeto complexo sem documentação adicional",
				Suggestion:  "Considere adicionar documentação em diretório docs/",
			})
		}
	}

	// Verificar configuração de editor
	if !r.hasEditorConfig(structure) {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityLow,
			Category:    CategoryMaintainability,
			Description: "Projeto sem configuração de editor (.editorconfig)",
			Suggestion:  "Adicione .editorconfig para consistência de formatação",
		})
	}

	// Calcular score
	if len(issues) > 0 {
		score -= float64(len(issues)) * 0.1
		if score < 0 {
			score = 0
		}
	}

	return score, issues
}

// DocumentationRule verifica qualidade da documentação
type DocumentationRule struct{}

func (r *DocumentationRule) Name() string {
	return "Documentation"
}

func (r *DocumentationRule) Description() string {
	return "Verifica qualidade e completude da documentação"
}

func (r *DocumentationRule) Category() Category {
	return CategoryMaintainability
}

func (r *DocumentationRule) Weight() float64 {
	return 6.0
}

func (r *DocumentationRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar README
	readmeIssues := r.evaluateReadme(structure)
	issues = append(issues, readmeIssues...)

	// Verificar comentários de código
	codeDocIssues := r.evaluateCodeDocumentation(structure)
	issues = append(issues, codeDocIssues...)

	// Verificar documentação API se aplicável
	apiDocIssues := r.evaluateAPIDocumentation(structure)
	issues = append(issues, apiDocIssues...)

	// Calcular score
	for _, issue := range issues {
		switch issue.Severity {
		case SeverityHigh:
			score -= 0.3
		case SeverityMedium:
			score -= 0.2
		case SeverityLow:
			score -= 0.1
		}
	}

	if score < 0 {
		score = 0
	}

	return score, issues
}

func (r *DocumentationRule) evaluateReadme(structure *ProjectStructure) []Issue {
	var issues []Issue

	readmeFile := r.findReadmeFile(structure)
	if readmeFile == nil {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityHigh,
			Category:    CategoryMaintainability,
			Description: "Projeto sem arquivo README",
			Suggestion:  "Crie um README.md com descrição, instalação e uso",
		})
		return issues
	}

	content := strings.ToLower(readmeFile.Content)

	// Verificar seções importantes
	importantSections := map[string]string{
		"install":     "Seção de instalação",
		"usage":       "Seção de uso/exemplos",
		"description": "Descrição do projeto",
	}

	for keyword, description := range importantSections {
		if !strings.Contains(content, keyword) &&
			!strings.Contains(content, "como usar") &&
			!strings.Contains(content, "como instalar") {
			issues = append(issues, Issue{
				Type:        IssueBestPractice,
				Severity:    SeverityLow,
				Category:    CategoryMaintainability,
				Description: fmt.Sprintf("README pode estar faltando: %s", description),
				Location:    "README",
				Suggestion:  fmt.Sprintf("Adicione %s ao README", description),
			})
		}
	}

	// Verificar se README é muito curto
	if len(readmeFile.Content) < 100 {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityMedium,
			Category:    CategoryMaintainability,
			Description: "README muito curto ou incompleto",
			Location:    "README",
			Suggestion:  "Expanda o README com mais informações sobre o projeto",
		})
	}

	return issues
}

func (r *DocumentationRule) evaluateCodeDocumentation(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Esta é uma verificação básica - em um sistema real,
	// analisaríamos o conteúdo dos arquivos de código
	codeFiles := 0
	for fileName := range structure.Files {
		if r.isCodeFile(fileName) {
			codeFiles++
		}
	}

	if codeFiles > 5 {
		// Para projetos com muitos arquivos de código, assumir que precisa de mais documentação
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityLow,
			Category:    CategoryMaintainability,
			Description: "Projeto com muitos arquivos pode precisar de mais documentação",
			Suggestion:  "Considere adicionar comentários e documentação inline",
		})
	}

	return issues
}

func (r *DocumentationRule) evaluateAPIDocumentation(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se parece ser uma API e se tem documentação
	isAPI := r.looksLikeAPI(structure)

	if isAPI {
		hasAPIDoc := r.hasAPIDocumentation(structure)
		if !hasAPIDoc {
			issues = append(issues, Issue{
				Type:        IssueBestPractice,
				Severity:    SeverityMedium,
				Category:    CategoryMaintainability,
				Description: "Projeto de API sem documentação específica",
				Suggestion:  "Considere usar OpenAPI/Swagger ou similar para documentar a API",
			})
		}
	}

	return issues
}

// TestStructureRule verifica estrutura de testes
type TestStructureRule struct{}

func (r *TestStructureRule) Name() string {
	return "TestStructure"
}

func (r *TestStructureRule) Description() string {
	return "Verifica se o projeto tem estrutura adequada de testes"
}

func (r *TestStructureRule) Category() Category {
	return CategoryMaintainability
}

func (r *TestStructureRule) Weight() float64 {
	return 8.0
}

func (r *TestStructureRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar presença de testes
	hasTests := r.hasTestFiles(structure)

	if !hasTests {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityMedium,
			Category:    CategoryMaintainability,
			Description: "Projeto sem estrutura de testes",
			Suggestion:  "Adicione testes unitários apropriados para a linguagem",
		})
		score -= 0.5
	} else {
		// Verificar qualidade da estrutura de testes
		testIssues := r.evaluateTestQuality(structure)
		issues = append(issues, testIssues...)

		if len(testIssues) > 0 {
			score -= float64(len(testIssues)) * 0.1
		}
	}

	if score < 0 {
		score = 0
	}

	return score, issues
}

func (r *TestStructureRule) hasTestFiles(structure *ProjectStructure) bool {
	testPatterns := []string{
		"test", "spec", "_test", ".test", "_spec", ".spec",
		"__test__", "__tests__", "tests/", "test/",
	}

	for fileName := range structure.Files {
		lowerName := strings.ToLower(fileName)
		for _, pattern := range testPatterns {
			if strings.Contains(lowerName, pattern) {
				return true
			}
		}
	}

	for _, dir := range structure.Directories {
		lowerDir := strings.ToLower(dir)
		for _, pattern := range testPatterns {
			if strings.Contains(lowerDir, pattern) {
				return true
			}
		}
	}

	return false
}

func (r *TestStructureRule) evaluateTestQuality(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se há framework de teste configurado
	hasTestFramework := r.hasTestFramework(structure)

	if !hasTestFramework {
		issues = append(issues, Issue{
			Type:        IssueBestPractice,
			Severity:    SeverityLow,
			Category:    CategoryMaintainability,
			Description: "Nenhum framework de teste identificado",
			Suggestion:  "Configure um framework de teste apropriado para a linguagem",
		})
	}

	return issues
}

func (r *TestStructureRule) hasTestFramework(structure *ProjectStructure) bool {
	testFrameworks := map[string][]string{
		"javascript": {"jest", "mocha", "jasmine", "cypress", "playwright"},
		"python":     {"pytest", "unittest", "nose", "tox"},
		"go":         {"testing", "testify", "ginkgo"},
		"java":       {"junit", "testng", "mockito"},
		"rust":       {"cargo test"},
	}

	language := strings.ToLower(structure.Language)
	frameworks, exists := testFrameworks[language]

	if !exists {
		return false
	}

	// Verificar dependências
	for depName := range structure.Dependencies {
		lowerDepName := strings.ToLower(depName)
		for _, framework := range frameworks {
			if strings.Contains(lowerDepName, framework) {
				return true
			}
		}
	}

	// Verificar arquivos de configuração
	for fileName := range structure.Files {
		lowerFileName := strings.ToLower(fileName)
		for _, framework := range frameworks {
			if strings.Contains(lowerFileName, framework) {
				return true
			}
		}
	}

	return false
}

// Helper functions para as regras de melhores práticas

func (r *BestPracticesRule) hasFile(structure *ProjectStructure, fileName string) bool {
	for file := range structure.Files {
		if strings.EqualFold(file, fileName) {
			return true
		}
	}
	return false
}

func (r *BestPracticesRule) hasTestDirectory(structure *ProjectStructure) bool {
	testDirs := []string{"test", "tests", "__tests__", "spec"}

	for _, dir := range structure.Directories {
		lowerDir := strings.ToLower(dir)
		for _, testDir := range testDirs {
			if strings.Contains(lowerDir, testDir) {
				return true
			}
		}
	}

	// Verificar arquivos de teste
	for fileName := range structure.Files {
		lowerName := strings.ToLower(fileName)
		if strings.Contains(lowerName, "test") || strings.Contains(lowerName, "spec") {
			return true
		}
	}

	return false
}

func (r *BestPracticesRule) hasDocumentation(structure *ProjectStructure) bool {
	docDirs := []string{"docs", "doc", "documentation"}

	for _, dir := range structure.Directories {
		lowerDir := strings.ToLower(dir)
		for _, docDir := range docDirs {
			if strings.Contains(lowerDir, docDir) {
				return true
			}
		}
	}
	return false
}

func (r *BestPracticesRule) hasEditorConfig(structure *ProjectStructure) bool {
	return r.hasFile(structure, ".editorconfig")
}

func (r *DocumentationRule) findReadmeFile(structure *ProjectStructure) *FileInfo {
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "readme.md", "README"}

	for fileName, fileInfo := range structure.Files {
		for _, readme := range readmeFiles {
			if strings.EqualFold(fileName, readme) {
				return &fileInfo
			}
		}
	}
	return nil
}

func (r *DocumentationRule) isCodeFile(fileName string) bool {
	codeExtensions := []string{".go", ".js", ".ts", ".py", ".java", ".rs", ".cpp", ".c", ".cs"}

	lowerName := strings.ToLower(fileName)
	for _, ext := range codeExtensions {
		if strings.HasSuffix(lowerName, ext) {
			return true
		}
	}
	return false
}

func (r *DocumentationRule) looksLikeAPI(structure *ProjectStructure) bool {
	apiIndicators := []string{"api", "server", "service", "endpoint", "router", "handler"}

	// Verificar nomes de arquivos e diretórios
	for fileName := range structure.Files {
		lowerName := strings.ToLower(fileName)
		for _, indicator := range apiIndicators {
			if strings.Contains(lowerName, indicator) {
				return true
			}
		}
	}

	for _, dir := range structure.Directories {
		lowerDir := strings.ToLower(dir)
		for _, indicator := range apiIndicators {
			if strings.Contains(lowerDir, indicator) {
				return true
			}
		}
	}

	// Verificar dependências comuns de API
	apiDeps := []string{"express", "fastapi", "gin", "echo", "spring", "django"}
	for depName := range structure.Dependencies {
		lowerDepName := strings.ToLower(depName)
		for _, apiDep := range apiDeps {
			if strings.Contains(lowerDepName, apiDep) {
				return true
			}
		}
	}

	return false
}

func (r *DocumentationRule) hasAPIDocumentation(structure *ProjectStructure) bool {
	apiDocFiles := []string{"swagger", "openapi", "api.md", "api.yml", "api.yaml"}

	for fileName := range structure.Files {
		lowerName := strings.ToLower(fileName)
		for _, docFile := range apiDocFiles {
			if strings.Contains(lowerName, docFile) {
				return true
			}
		}
	}

	return false
}
