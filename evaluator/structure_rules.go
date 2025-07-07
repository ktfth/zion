package evaluator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// DirectoryStructureRule verifica se a estrutura de diretórios segue boas práticas
type DirectoryStructureRule struct{}

func (r *DirectoryStructureRule) Name() string {
	return "DirectoryStructure"
}

func (r *DirectoryStructureRule) Description() string {
	return "Verifica se a estrutura de diretórios segue convenções e boas práticas"
}

func (r *DirectoryStructureRule) Category() Category {
	return CategoryStructure
}

func (r *DirectoryStructureRule) Weight() float64 {
	return 15.0
}

func (r *DirectoryStructureRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar estrutura por linguagem
	switch strings.ToLower(structure.Language) {
	case "go", "golang":
		score, issues = r.evaluateGoStructure(structure)
	case "javascript", "js", "typescript", "ts":
		score, issues = r.evaluateJSStructure(structure)
	case "python", "py":
		score, issues = r.evaluatePythonStructure(structure)
	case "java":
		score, issues = r.evaluateJavaStructure(structure)
	case "rust":
		score, issues = r.evaluateRustStructure(structure)
	default:
		score, issues = r.evaluateGenericStructure(structure)
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluateGoStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	recommendedDirs := []string{"cmd", "pkg", "internal"}

	hasRecommended := 0
	for _, dir := range recommendedDirs {
		if r.hasDirectory(structure, dir) {
			hasRecommended++
		}
	}

	if hasRecommended == 0 {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto Go não segue estrutura padrão (cmd/, pkg/, internal/)",
			Suggestion:  "Considere usar a estrutura padrão: cmd/ para executáveis, pkg/ para bibliotecas, internal/ para código privado",
		})
		score -= 0.3
	}

	// Verificar se tem diretórios desnecessários na raiz
	rootFiles := r.getRootFiles(structure)
	if len(rootFiles) > 10 {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: "Muitos arquivos na raiz do projeto",
			Suggestion:  "Organize arquivos em diretórios apropriados",
		})
		score -= 0.1
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluateJSStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar estrutura comum para projetos JS/TS
	if !r.hasDirectory(structure, "src") && !r.hasDirectory(structure, "lib") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto JavaScript/TypeScript sem diretório src/ ou lib/",
			Suggestion:  "Crie um diretório src/ para organizar o código fonte",
		})
		score -= 0.3
	}

	// Verificar se tem node_modules na estrutura (não deveria)
	if r.hasDirectory(structure, "node_modules") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: "Diretório node_modules incluído na estrutura do projeto",
			Suggestion:  "Adicione node_modules ao .gitignore",
		})
		score -= 0.1
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluatePythonStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar se tem estrutura de pacote
	if !r.hasFile(structure, "__init__.py") && !r.hasDirectory(structure, "src") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto Python sem estrutura de pacote clara",
			Suggestion:  "Crie arquivos __init__.py ou organize código em diretório src/",
		})
		score -= 0.3
	}

	// Verificar se tem __pycache__ na estrutura
	if r.hasDirectory(structure, "__pycache__") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: "Diretório __pycache__ incluído na estrutura do projeto",
			Suggestion:  "Adicione __pycache__ ao .gitignore",
		})
		score -= 0.1
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluateJavaStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar estrutura Maven/Gradle padrão
	hasMavenStructure := r.hasDirectory(structure, "src/main/java") || r.hasDirectory(structure, "src\\main\\java")

	if !hasMavenStructure {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto Java não segue estrutura Maven/Gradle padrão",
			Suggestion:  "Use a estrutura padrão: src/main/java, src/test/java",
		})
		score -= 0.3
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluateRustStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar estrutura padrão do Cargo
	if !r.hasDirectory(structure, "src") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto Rust sem diretório src/",
			Suggestion:  "Crie um diretório src/ para o código fonte",
		})
		score -= 0.3
	}

	if !r.hasFile(structure, "src/main.rs") && !r.hasFile(structure, "src/lib.rs") {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto Rust sem main.rs ou lib.rs",
			Suggestion:  "Crie src/main.rs para executáveis ou src/lib.rs para bibliotecas",
		})
		score -= 0.2
	}

	return score, issues
}

func (r *DirectoryStructureRule) evaluateGenericStructure(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificações genéricas
	if len(structure.Directories) == 0 && len(structure.Files) <= 1 {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityCritical,
			Category:    CategoryStructure,
			Description: "Projeto com estrutura insuficiente - muito poucos arquivos e diretórios",
			Suggestion:  "Organize o código em uma estrutura apropriada com múltiplos arquivos e diretórios",
		})
		score -= 0.8
	} else if len(structure.Directories) == 0 {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityHigh,
			Category:    CategoryStructure,
			Description: "Projeto sem estrutura de diretórios definida",
			Suggestion:  "Organize o código em diretórios apropriados",
		})
		score -= 0.5
	}

	return score, issues
}

// FileNamingRule verifica convenções de nomenclatura de arquivos
type FileNamingRule struct{}

func (r *FileNamingRule) Name() string {
	return "FileNaming"
}

func (r *FileNamingRule) Description() string {
	return "Verifica se os nomes dos arquivos seguem convenções apropriadas"
}

func (r *FileNamingRule) Category() Category {
	return CategoryNaming
}

func (r *FileNamingRule) Weight() float64 {
	return 8.0
}

func (r *FileNamingRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar nomes de arquivos
	for fileName := range structure.Files {
		fileIssues := r.validateFileName(fileName, structure.Language)
		issues = append(issues, fileIssues...)

		if len(fileIssues) > 0 {
			score -= 0.1
		}
	}

	// Garantir que o score não fique negativo
	if score < 0 {
		score = 0
	}

	return score, issues
}

func (r *FileNamingRule) validateFileName(fileName, language string) []Issue {
	var issues []Issue

	// Verificar caracteres problemáticos
	if strings.Contains(fileName, " ") {
		issues = append(issues, Issue{
			Type:        IssueNaming,
			Severity:    SeverityMedium,
			Category:    CategoryNaming,
			Description: fmt.Sprintf("Arquivo '%s' contém espaços no nome", fileName),
			Location:    fileName,
			Suggestion:  "Use underscore (_) ou hífen (-) em vez de espaços",
		})
	}

	// Verificar caracteres especiais problemáticos
	problematicChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "=", "+", "[", "]", "{", "}", "|", "\\", ":", ";", "\"", "'", "<", ">", "?", ","}
	for _, char := range problematicChars {
		if strings.Contains(fileName, char) {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityLow,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo '%s' contém caractere problemático: %s", fileName, char),
				Location:    fileName,
				Suggestion:  "Evite caracteres especiais em nomes de arquivos",
			})
			break
		}
	}

	// Verificações específicas por linguagem
	switch strings.ToLower(language) {
	case "go", "golang":
		issues = append(issues, r.validateGoFileName(fileName)...)
	case "python", "py":
		issues = append(issues, r.validatePythonFileName(fileName)...)
	case "java":
		issues = append(issues, r.validateJavaFileName(fileName)...)
	}

	return issues
}

func (r *FileNamingRule) validateGoFileName(fileName string) []Issue {
	var issues []Issue

	if strings.HasSuffix(fileName, ".go") {
		baseName := strings.TrimSuffix(filepath.Base(fileName), ".go")

		// Go usa snake_case para arquivos
		if strings.Contains(baseName, "-") {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityLow,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo Go '%s' usa hífen em vez de underscore", fileName),
				Location:    fileName,
				Suggestion:  "Use underscore (_) em nomes de arquivos Go",
			})
		}

		// Verificar se nome está em camelCase (deveria ser snake_case)
		camelCaseRegex := regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)+$`)
		if camelCaseRegex.MatchString(baseName) {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityLow,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo Go '%s' usa camelCase em vez de snake_case", fileName),
				Location:    fileName,
				Suggestion:  "Use snake_case para nomes de arquivos Go",
			})
		}
	}

	return issues
}

func (r *FileNamingRule) validatePythonFileName(fileName string) []Issue {
	var issues []Issue

	if strings.HasSuffix(fileName, ".py") {
		baseName := strings.TrimSuffix(filepath.Base(fileName), ".py")

		// Python usa snake_case
		if strings.Contains(baseName, "-") {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityLow,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo Python '%s' usa hífen em vez de underscore", fileName),
				Location:    fileName,
				Suggestion:  "Use underscore (_) em nomes de arquivos Python",
			})
		}

		// Verificar camelCase
		camelCaseRegex := regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)+$`)
		if camelCaseRegex.MatchString(baseName) {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityLow,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo Python '%s' usa camelCase em vez de snake_case", fileName),
				Location:    fileName,
				Suggestion:  "Use snake_case para nomes de arquivos Python",
			})
		}
	}

	return issues
}

func (r *FileNamingRule) validateJavaFileName(fileName string) []Issue {
	var issues []Issue

	if strings.HasSuffix(fileName, ".java") {
		baseName := strings.TrimSuffix(filepath.Base(fileName), ".java")

		// Java usa PascalCase para classes
		pascalCaseRegex := regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
		if !pascalCaseRegex.MatchString(baseName) {
			issues = append(issues, Issue{
				Type:        IssueNaming,
				Severity:    SeverityMedium,
				Category:    CategoryNaming,
				Description: fmt.Sprintf("Arquivo Java '%s' não usa PascalCase", fileName),
				Location:    fileName,
				Suggestion:  "Use PascalCase para nomes de classes Java",
			})
		}
	}

	return issues
}

// RequiredFilesRule verifica se arquivos obrigatórios estão presentes
type RequiredFilesRule struct{}

func (r *RequiredFilesRule) Name() string {
	return "RequiredFiles"
}

func (r *RequiredFilesRule) Description() string {
	return "Verifica se arquivos essenciais estão presentes no projeto"
}

func (r *RequiredFilesRule) Category() Category {
	return CategoryStructure
}

func (r *RequiredFilesRule) Weight() float64 {
	return 12.0
}

func (r *RequiredFilesRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar README
	if !r.hasReadme(structure) {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityMedium,
			Category:    CategoryStructure,
			Description: "Projeto não possui arquivo README",
			Suggestion:  "Crie um arquivo README.md com documentação do projeto",
		})
		score -= 0.3
	}

	// Verificar .gitignore
	if !r.hasGitignore(structure) {
		issues = append(issues, Issue{
			Type:        IssueStructure,
			Severity:    SeverityLow,
			Category:    CategoryStructure,
			Description: "Projeto não possui arquivo .gitignore",
			Suggestion:  "Crie um arquivo .gitignore apropriado para a linguagem",
		})
		score -= 0.1
	}

	// Verificações específicas por linguagem
	languageIssues, languageScore := r.checkLanguageSpecificFiles(structure)
	issues = append(issues, languageIssues...)
	score += languageScore - 1.0 // Ajustar score

	if score < 0 {
		score = 0
	}

	return score, issues
}

func (r *RequiredFilesRule) checkLanguageSpecificFiles(structure *ProjectStructure) ([]Issue, float64) {
	var issues []Issue
	score := 1.0

	switch strings.ToLower(structure.Language) {
	case "go", "golang":
		if !r.hasFile(structure, "go.mod") {
			issues = append(issues, Issue{
				Type:        IssueStructure,
				Severity:    SeverityCritical,
				Category:    CategoryStructure,
				Description: "Projeto Go sem go.mod",
				Suggestion:  "Execute 'go mod init <module-name>' para criar go.mod",
			})
			score -= 0.4
		}

	case "javascript", "js", "typescript", "ts":
		if !r.hasFile(structure, "package.json") {
			issues = append(issues, Issue{
				Type:        IssueStructure,
				Severity:    SeverityCritical,
				Category:    CategoryStructure,
				Description: "Projeto JavaScript/TypeScript sem package.json",
				Suggestion:  "Execute 'npm init' para criar package.json",
			})
			score -= 0.4
		}

	case "python", "py":
		hasRequirements := r.hasFile(structure, "requirements.txt")
		hasPyproject := r.hasFile(structure, "pyproject.toml")

		if !hasRequirements && !hasPyproject {
			issues = append(issues, Issue{
				Type:        IssueStructure,
				Severity:    SeverityMedium,
				Category:    CategoryStructure,
				Description: "Projeto Python sem gerenciamento de dependências",
				Suggestion:  "Crie requirements.txt ou pyproject.toml para gerenciar dependências",
			})
			score -= 0.3
		}

	case "java":
		hasPom := r.hasFile(structure, "pom.xml")
		hasGradle := r.hasFile(structure, "build.gradle")

		if !hasPom && !hasGradle {
			issues = append(issues, Issue{
				Type:        IssueStructure,
				Severity:    SeverityCritical,
				Category:    CategoryStructure,
				Description: "Projeto Java sem arquivo de build (pom.xml ou build.gradle)",
				Suggestion:  "Configure Maven (pom.xml) ou Gradle (build.gradle)",
			})
			score -= 0.4
		}

	case "rust":
		if !r.hasFile(structure, "Cargo.toml") {
			issues = append(issues, Issue{
				Type:        IssueStructure,
				Severity:    SeverityCritical,
				Category:    CategoryStructure,
				Description: "Projeto Rust sem Cargo.toml",
				Suggestion:  "Execute 'cargo init' para criar projeto Rust",
			})
			score -= 0.4
		}
	}

	return issues, score
}

// Helper functions

func (r *DirectoryStructureRule) hasDirectory(structure *ProjectStructure, dirName string) bool {
	for _, dir := range structure.Directories {
		if strings.Contains(strings.ToLower(dir), strings.ToLower(dirName)) {
			return true
		}
	}
	return false
}

func (r *DirectoryStructureRule) hasFile(structure *ProjectStructure, fileName string) bool {
	for file := range structure.Files {
		if strings.EqualFold(filepath.Base(file), fileName) {
			return true
		}
	}
	return false
}

func (r *DirectoryStructureRule) getRootFiles(structure *ProjectStructure) []string {
	var rootFiles []string
	for file := range structure.Files {
		if !strings.Contains(file, "/") && !strings.Contains(file, "\\") {
			rootFiles = append(rootFiles, file)
		}
	}
	return rootFiles
}

func (r *RequiredFilesRule) hasReadme(structure *ProjectStructure) bool {
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "readme.md", "readme.txt", "README"}
	for _, readme := range readmeFiles {
		if r.hasFile(structure, readme) {
			return true
		}
	}
	return false
}

func (r *RequiredFilesRule) hasGitignore(structure *ProjectStructure) bool {
	return r.hasFile(structure, ".gitignore")
}

func (r *RequiredFilesRule) hasFile(structure *ProjectStructure, fileName string) bool {
	for file := range structure.Files {
		if strings.EqualFold(filepath.Base(file), fileName) {
			return true
		}
	}
	return false
}
