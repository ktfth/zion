package evaluator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// DependencyConsistencyRule verifica a consistência e qualidade das dependências
type DependencyConsistencyRule struct{}

func (r *DependencyConsistencyRule) Name() string {
	return "DependencyConsistency"
}

func (r *DependencyConsistencyRule) Description() string {
	return "Verifica a consistência, atualidade e qualidade das dependências"
}

func (r *DependencyConsistencyRule) Category() Category {
	return CategoryDependencies
}

func (r *DependencyConsistencyRule) Weight() float64 {
	return 10.0
}

func (r *DependencyConsistencyRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	if len(structure.Dependencies) == 0 {
		// Sem dependências pode ser normal para projetos simples
		return score, issues
	}

	// Verificar versões das dependências
	versionIssues := r.checkVersions(structure.Dependencies)
	issues = append(issues, versionIssues...)

	// Verificar dependências duplicadas ou conflitantes
	duplicateIssues := r.checkDuplicates(structure.Dependencies)
	issues = append(issues, duplicateIssues...)

	// Verificar dependências desnecessárias ou suspeitas
	suspiciousIssues := r.checkSuspiciousDependencies(structure.Dependencies, structure.Language)
	issues = append(issues, suspiciousIssues...)

	// Calcular score baseado no número de issues
	if len(issues) > 0 {
		score -= float64(len(issues)) * 0.1
		if score < 0 {
			score = 0
		}
	}

	return score, issues
}

func (r *DependencyConsistencyRule) checkVersions(dependencies map[string]interface{}) []Issue {
	var issues []Issue

	for name, version := range dependencies {
		versionStr := fmt.Sprintf("%v", version)

		// Verificar se usa "latest" ou "*"
		if versionStr == "latest" || versionStr == "*" {
			issues = append(issues, Issue{
				Type:        IssueDependency,
				Severity:    SeverityMedium,
				Category:    CategoryDependencies,
				Description: fmt.Sprintf("Dependência '%s' usa versão '%s'", name, versionStr),
				Location:    name,
				Suggestion:  "Especifique uma versão específica para maior estabilidade",
			})
		}

		// Verificar versões very old ou problemáticas
		if r.isOldVersion(name, versionStr) {
			issues = append(issues, Issue{
				Type:        IssueDependency,
				Severity:    SeverityLow,
				Category:    CategoryDependencies,
				Description: fmt.Sprintf("Dependência '%s' pode estar desatualizada (v%s)", name, versionStr),
				Location:    name,
				Suggestion:  "Verifique se há versões mais recentes disponíveis",
			})
		}
	}

	return issues
}

func (r *DependencyConsistencyRule) checkDuplicates(dependencies map[string]interface{}) []Issue {
	var issues []Issue

	// Mapear dependências similares
	similarDeps := make(map[string][]string)

	for name := range dependencies {
		baseName := r.getBaseDependencyName(name)
		similarDeps[baseName] = append(similarDeps[baseName], name)
	}

	// Verificar duplicatas
	for baseName, deps := range similarDeps {
		if len(deps) > 1 {
			issues = append(issues, Issue{
				Type:        IssueDependency,
				Severity:    SeverityMedium,
				Category:    CategoryDependencies,
				Description: fmt.Sprintf("Possíveis dependências duplicadas: %v", deps),
				Location:    baseName,
				Suggestion:  "Verifique se realmente precisa de todas essas dependências similares",
			})
		}
	}

	return issues
}

func (r *DependencyConsistencyRule) checkSuspiciousDependencies(dependencies map[string]interface{}, language string) []Issue {
	var issues []Issue

	// Lista de dependências potencialmente problemáticas
	suspicious := map[string]string{
		"lodash":  "Considere usar funções nativas do JavaScript quando possível",
		"moment":  "Considere usar date-fns ou dayjs como alternativas mais leves",
		"request": "Biblioteca depreciada, use axios ou node-fetch",
		"jquery":  "Considere usar JavaScript vanilla ou frameworks modernos",
	}

	for name := range dependencies {
		lowerName := strings.ToLower(name)
		for suspName, suggestion := range suspicious {
			if strings.Contains(lowerName, suspName) {
				issues = append(issues, Issue{
					Type:        IssueDependency,
					Severity:    SeverityLow,
					Category:    CategoryDependencies,
					Description: fmt.Sprintf("Dependência '%s' pode ter alternativas melhores", name),
					Location:    name,
					Suggestion:  suggestion,
				})
			}
		}
	}

	return issues
}

func (r *DependencyConsistencyRule) getBaseDependencyName(name string) string {
	// Remove prefixos comuns e variações
	name = strings.ToLower(name)
	name = strings.TrimPrefix(name, "@")
	name = strings.Split(name, "/")[0]
	return name
}

func (r *DependencyConsistencyRule) isOldVersion(name, version string) bool {
	// Verificações básicas para versões antigas
	versionRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)
	matches := versionRegex.FindStringSubmatch(version)

	if len(matches) >= 2 {
		// Se versão major é 0, pode estar em desenvolvimento
		if matches[1] == "0" {
			return false
		}
	}

	// Lista de dependências com versões conhecidamente antigas
	oldVersions := map[string][]string{
		"react":     {"15", "16"},
		"angular":   {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		"vue":       {"1", "2"},
		"jquery":    {"1", "2"},
		"bootstrap": {"3", "4"},
	}

	for dep, oldVers := range oldVersions {
		if strings.Contains(strings.ToLower(name), dep) {
			for _, oldVer := range oldVers {
				if strings.HasPrefix(version, oldVer+".") {
					return true
				}
			}
		}
	}

	return false
}

// SecurityVulnerabilityRule verifica vulnerabilidades de segurança conhecidas
type SecurityVulnerabilityRule struct{}

func (r *SecurityVulnerabilityRule) Name() string {
	return "SecurityVulnerability"
}

func (r *SecurityVulnerabilityRule) Description() string {
	return "Verifica dependências com vulnerabilidades de segurança conhecidas"
}

func (r *SecurityVulnerabilityRule) Category() Category {
	return CategorySecurity
}

func (r *SecurityVulnerabilityRule) Weight() float64 {
	return 15.0
}

func (r *SecurityVulnerabilityRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar dependências com vulnerabilidades conhecidas
	vulnerableIssues := r.checkVulnerableDependencies(structure.Dependencies)
	issues = append(issues, vulnerableIssues...)

	// Verificar configurações de segurança
	securityIssues := r.checkSecurityConfigurations(structure)
	issues = append(issues, securityIssues...)

	// Calcular score baseado na severidade dos issues
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

func (r *SecurityVulnerabilityRule) checkVulnerableDependencies(dependencies map[string]interface{}) []Issue {
	var issues []Issue

	// Lista de dependências com vulnerabilidades conhecidas
	vulnerable := map[string]map[string]string{
		"lodash": {
			"severity": "high",
			"reason":   "Versões antigas têm vulnerabilidades de prototype pollution",
		},
		"moment": {
			"severity": "medium",
			"reason":   "Biblioteca grande com possíveis vulnerabilidades",
		},
		"request": {
			"severity": "critical",
			"reason":   "Biblioteca depreciada com vulnerabilidades não corrigidas",
		},
		"tar": {
			"severity": "critical",
			"reason":   "Versões antigas têm vulnerabilidades críticas de path traversal",
		},
		"minimist": {
			"severity": "critical",
			"reason":   "Prototype pollution crítica em versões antigas",
		},
		"node-tar": {
			"severity": "critical",
			"reason":   "Vulnerabilidades críticas de segurança",
		},
		"serialize-javascript": {
			"severity": "high",
			"reason":   "Vulnerabilidades de XSS em versões antigas",
		},
	}

	for depName, version := range dependencies {
		lowerName := strings.ToLower(depName)
		versionStr := fmt.Sprintf("%v", version)

		for vulnName, info := range vulnerable {
			if strings.Contains(lowerName, vulnName) {
				severity := SeverityMedium
				switch info["severity"] {
				case "critical":
					severity = SeverityCritical
				case "high":
					severity = SeverityHigh
				case "low":
					severity = SeverityLow
				}

				issues = append(issues, Issue{
					Type:        IssueSecurity,
					Severity:    severity,
					Category:    CategorySecurity,
					Description: fmt.Sprintf("Dependência '%s' v%s tem vulnerabilidades %s", depName, versionStr, info["severity"]),
					Location:    depName,
					Suggestion:  info["reason"] + ". Atualize para versão mais recente ou use alternativa",
				})
			}
		}
	}

	return issues
}

func (r *SecurityVulnerabilityRule) checkSecurityConfigurations(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se há arquivos de configuração expostos
	sensitiveFiles := []string{".env", "config.json", "secrets.json", "private.key", ".pem", "password", "secret"}

	for fileName, fileInfo := range structure.Files {
		lowerFileName := strings.ToLower(fileName)

		for _, sensitive := range sensitiveFiles {
			if strings.Contains(lowerFileName, sensitive) {
				severity := SeverityHigh
				if strings.Contains(lowerFileName, ".env") || strings.Contains(lowerFileName, "secret") {
					severity = SeverityCritical
				}

				issues = append(issues, Issue{
					Type:        IssueSecurity,
					Severity:    severity,
					Category:    CategorySecurity,
					Description: fmt.Sprintf("Arquivo sensível '%s' presente na estrutura", fileName),
					Location:    fileName,
					Suggestion:  "Adicione arquivos sensíveis ao .gitignore e use variáveis de ambiente",
				})

				// Verificar conteúdo do arquivo se disponível
				if fileInfo.Content != "" && (strings.Contains(lowerFileName, ".env") || strings.Contains(lowerFileName, "config")) {
					contentIssues := r.checkSensitiveContent(fileName, fileInfo.Content)
					issues = append(issues, contentIssues...)
				}
			}
		}
	}

	// Verificar configurações específicas por linguagem
	switch strings.ToLower(structure.Language) {
	case "javascript", "js", "typescript", "ts":
		issues = append(issues, r.checkJSSecurityConfig(structure)...)
	case "python", "py":
		issues = append(issues, r.checkPythonSecurityConfig(structure)...)
	}

	return issues
}

func (r *SecurityVulnerabilityRule) checkSensitiveContent(fileName, content string) []Issue {
	var issues []Issue

	// Padrões de dados sensíveis
	sensitivePatterns := map[string]string{
		"password": "Senha hardcoded detectada",
		"secret":   "Chave secreta hardcoded detectada",
		"api_key":  "API key hardcoded detectada",
		"token":    "Token hardcoded detectado",
		"private":  "Chave privada hardcoded detectada",
	}

	lowerContent := strings.ToLower(content)

	for pattern, description := range sensitivePatterns {
		if strings.Contains(lowerContent, pattern) {
			issues = append(issues, Issue{
				Type:        IssueSecurity,
				Severity:    SeverityCritical,
				Category:    CategorySecurity,
				Description: fmt.Sprintf("%s em %s", description, fileName),
				Location:    fileName,
				Suggestion:  "Remova dados sensíveis e use variáveis de ambiente",
			})
		}
	}

	return issues
}

func (r *SecurityVulnerabilityRule) checkJSSecurityConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar se package.json tem configurações inseguras
	if packageFile, exists := structure.Files["package.json"]; exists {
		content := strings.ToLower(packageFile.Content)

		// Verificar se permite HTTP em produção
		if strings.Contains(content, "\"allow-http\"") {
			issues = append(issues, Issue{
				Type:        IssueSecurity,
				Severity:    SeverityMedium,
				Category:    CategorySecurity,
				Description: "Configuração permite HTTP em potencial ambiente de produção",
				Location:    "package.json",
				Suggestion:  "Use HTTPS em produção",
			})
		}
	}

	return issues
}

func (r *SecurityVulnerabilityRule) checkPythonSecurityConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar requirements.txt para dependências inseguras
	if reqFile, exists := structure.Files["requirements.txt"]; exists {
		content := strings.ToLower(reqFile.Content)

		insecurePythonPackages := []string{"pickle", "yaml.load", "eval", "exec"}
		for _, pkg := range insecurePythonPackages {
			if strings.Contains(content, pkg) {
				issues = append(issues, Issue{
					Type:        IssueSecurity,
					Severity:    SeverityHigh,
					Category:    CategorySecurity,
					Description: fmt.Sprintf("Uso potencialmente inseguro de '%s'", pkg),
					Location:    "requirements.txt",
					Suggestion:  "Revise o uso dessas funções potencialmente perigosas",
				})
			}
		}
	}

	return issues
}

// ConfigurationValidityRule verifica a validade das configurações
type ConfigurationValidityRule struct{}

func (r *ConfigurationValidityRule) Name() string {
	return "ConfigurationValidity"
}

func (r *ConfigurationValidityRule) Description() string {
	return "Verifica se as configurações do projeto são válidas e consistentes"
}

func (r *ConfigurationValidityRule) Category() Category {
	return CategoryConfiguration
}

func (r *ConfigurationValidityRule) Weight() float64 {
	return 8.0
}

func (r *ConfigurationValidityRule) Evaluate(structure *ProjectStructure) (float64, []Issue) {
	var issues []Issue
	score := 1.0

	// Verificar arquivos de configuração por linguagem
	switch strings.ToLower(structure.Language) {
	case "javascript", "js", "typescript", "ts":
		configIssues := r.validateJSConfig(structure)
		issues = append(issues, configIssues...)
	case "python", "py":
		configIssues := r.validatePythonConfig(structure)
		issues = append(issues, configIssues...)
	case "go", "golang":
		configIssues := r.validateGoConfig(structure)
		issues = append(issues, configIssues...)
	case "java":
		configIssues := r.validateJavaConfig(structure)
		issues = append(issues, configIssues...)
	}

	// Calcular score baseado nos issues
	if len(issues) > 0 {
		score -= float64(len(issues)) * 0.15
		if score < 0 {
			score = 0
		}
	}

	return score, issues
}

func (r *ConfigurationValidityRule) validateJSConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar package.json
	if packageFile, exists := structure.Files["package.json"]; exists {
		packageIssues := r.validatePackageJson(packageFile.Content)
		issues = append(issues, packageIssues...)
	}

	// Verificar tsconfig.json se for TypeScript
	if tsconfigFile, exists := structure.Files["tsconfig.json"]; exists {
		tsconfigIssues := r.validateTsConfig(tsconfigFile.Content)
		issues = append(issues, tsconfigIssues...)
	}

	return issues
}

func (r *ConfigurationValidityRule) validatePackageJson(content string) []Issue {
	var issues []Issue

	// Tentar parsear JSON
	var packageData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &packageData); err != nil {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityHigh,
			Category:    CategoryConfiguration,
			Description: "package.json contém JSON inválido",
			Location:    "package.json",
			Suggestion:  "Corrija a sintaxe JSON no package.json",
		})
		return issues
	}

	// Verificar campos obrigatórios
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		if _, exists := packageData[field]; !exists {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityMedium,
				Category:    CategoryConfiguration,
				Description: fmt.Sprintf("package.json não tem campo obrigatório '%s'", field),
				Location:    "package.json",
				Suggestion:  fmt.Sprintf("Adicione o campo '%s' ao package.json", field),
			})
		}
	}

	return issues
}

func (r *ConfigurationValidityRule) validateTsConfig(content string) []Issue {
	var issues []Issue

	// Tentar parsear JSON
	var tsconfigData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &tsconfigData); err != nil {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityHigh,
			Category:    CategoryConfiguration,
			Description: "tsconfig.json contém JSON inválido",
			Location:    "tsconfig.json",
			Suggestion:  "Corrija a sintaxe JSON no tsconfig.json",
		})
		return issues
	}

	// Verificar configurações recomendadas
	if compilerOptions, exists := tsconfigData["compilerOptions"]; exists {
		if options, ok := compilerOptions.(map[string]interface{}); ok {
			if strict, exists := options["strict"]; !exists || strict != true {
				issues = append(issues, Issue{
					Type:        IssueConfiguration,
					Severity:    SeverityLow,
					Category:    CategoryConfiguration,
					Description: "TypeScript não está em modo strict",
					Location:    "tsconfig.json",
					Suggestion:  "Ative o modo strict para melhor verificação de tipos",
				})
			}
		}
	}

	return issues
}

func (r *ConfigurationValidityRule) validatePythonConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar requirements.txt se existir
	if reqFile, exists := structure.Files["requirements.txt"]; exists {
		reqIssues := r.validateRequirementsTxt(reqFile.Content)
		issues = append(issues, reqIssues...)
	}

	// Verificar pyproject.toml se existir
	if pyprojectFile, exists := structure.Files["pyproject.toml"]; exists {
		pyprojectIssues := r.validatePyprojectToml(pyprojectFile.Content)
		issues = append(issues, pyprojectIssues...)
	}

	return issues
}

func (r *ConfigurationValidityRule) validateRequirementsTxt(content string) []Issue {
	var issues []Issue

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Verificar formato básico
		if !strings.Contains(line, "==") && !strings.Contains(line, ">=") && !strings.Contains(line, "<=") {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityLow,
				Category:    CategoryConfiguration,
				Description: fmt.Sprintf("Linha %d em requirements.txt pode estar mal formatada", i+1),
				Location:    "requirements.txt",
				Suggestion:  "Use formato 'package==version' ou 'package>=version'",
			})
		}
	}

	return issues
}

func (r *ConfigurationValidityRule) validatePyprojectToml(content string) []Issue {
	var issues []Issue

	// Verificação básica se contém seções obrigatórias
	if !strings.Contains(content, "[build-system]") && !strings.Contains(content, "[project]") {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityMedium,
			Category:    CategoryConfiguration,
			Description: "pyproject.toml não contém seções obrigatórias",
			Location:    "pyproject.toml",
			Suggestion:  "Adicione seções [build-system] e/ou [project]",
		})
	}

	return issues
}

func (r *ConfigurationValidityRule) validateGoConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar go.mod se existir
	if gomodFile, exists := structure.Files["go.mod"]; exists {
		gomodIssues := r.validateGoMod(gomodFile.Content)
		issues = append(issues, gomodIssues...)
	}

	return issues
}

func (r *ConfigurationValidityRule) validateGoMod(content string) []Issue {
	var issues []Issue

	// Verificar se tem declaração de módulo
	if !strings.Contains(content, "module ") {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityHigh,
			Category:    CategoryConfiguration,
			Description: "go.mod não contém declaração de módulo",
			Location:    "go.mod",
			Suggestion:  "Adicione 'module <nome-do-modulo>' ao go.mod",
		})
	}

	// Verificar versão do Go
	if !strings.Contains(content, "go ") {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityMedium,
			Category:    CategoryConfiguration,
			Description: "go.mod não especifica versão do Go",
			Location:    "go.mod",
			Suggestion:  "Adicione 'go <versao>' ao go.mod",
		})
	}

	return issues
}

func (r *ConfigurationValidityRule) validateJavaConfig(structure *ProjectStructure) []Issue {
	var issues []Issue

	// Verificar pom.xml se for Maven
	if pomFile, exists := structure.Files["pom.xml"]; exists {
		pomIssues := r.validatePomXml(pomFile.Content)
		issues = append(issues, pomIssues...)
	}

	// Verificar build.gradle se for Gradle
	if gradleFile, exists := structure.Files["build.gradle"]; exists {
		gradleIssues := r.validateBuildGradle(gradleFile.Content)
		issues = append(issues, gradleIssues...)
	}

	return issues
}

func (r *ConfigurationValidityRule) validatePomXml(content string) []Issue {
	var issues []Issue

	// Verificações básicas de XML válido e estrutura Maven
	if !strings.Contains(content, "<project>") {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityHigh,
			Category:    CategoryConfiguration,
			Description: "pom.xml não contém elemento <project> raiz",
			Location:    "pom.xml",
			Suggestion:  "Corrija a estrutura XML do Maven",
		})
	}

	requiredElements := []string{"<groupId>", "<artifactId>", "<version>"}
	for _, element := range requiredElements {
		if !strings.Contains(content, element) {
			issues = append(issues, Issue{
				Type:        IssueConfiguration,
				Severity:    SeverityMedium,
				Category:    CategoryConfiguration,
				Description: fmt.Sprintf("pom.xml não contém elemento obrigatório %s", element),
				Location:    "pom.xml",
				Suggestion:  fmt.Sprintf("Adicione o elemento %s ao pom.xml", element),
			})
		}
	}

	return issues
}

func (r *ConfigurationValidityRule) validateBuildGradle(content string) []Issue {
	var issues []Issue

	// Verificações básicas de Gradle
	if !strings.Contains(content, "plugins") && !strings.Contains(content, "apply plugin") {
		issues = append(issues, Issue{
			Type:        IssueConfiguration,
			Severity:    SeverityMedium,
			Category:    CategoryConfiguration,
			Description: "build.gradle não define plugins",
			Location:    "build.gradle",
			Suggestion:  "Adicione plugins necessários ao build.gradle",
		})
	}

	return issues
}
