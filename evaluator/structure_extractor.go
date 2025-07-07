package evaluator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectStructure representa a estrutura extraída de um projeto
type ProjectStructure struct {
	ProjectName   string                 `json:"projectName"`
	Language      string                 `json:"language"`
	Directories   []string               `json:"directories"`
	Files         map[string]FileInfo    `json:"files"`
	Dependencies  map[string]interface{} `json:"dependencies"`
	Configuration map[string]interface{} `json:"configuration"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// FileInfo representa informações sobre um arquivo
type FileInfo struct {
	Path      string                 `json:"path"`
	Name      string                 `json:"name"`
	Extension string                 `json:"extension"`
	Content   string                 `json:"content"`
	IsConfig  bool                   `json:"isConfig"`
	IsBinary  bool                   `json:"isBinary"`
	Size      int                    `json:"size"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ExtractProjectStructure extrai a estrutura do projeto a partir dos dados fornecidos
func ExtractProjectStructure(projectData string, language string) (*ProjectStructure, error) {
	structure := &ProjectStructure{
		Language:      language,
		Directories:   make([]string, 0),
		Files:         make(map[string]FileInfo),
		Dependencies:  make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	// Tentar extrair JSON estruturado primeiro
	if err := tryExtractJSONStructure(projectData, structure); err == nil {
		return structure, nil
	}

	// Fallback para parsing de texto
	return extractFromTextData(projectData, language)
}

// tryExtractJSONStructure tenta extrair estrutura de dados JSON
func tryExtractJSONStructure(data string, structure *ProjectStructure) error {
	// Limpar e extrair JSON do texto
	jsonData := extractJSONFromText(data)
	if jsonData == "" {
		return fmt.Errorf("nenhum JSON válido encontrado")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &parsed); err != nil {
		return err
	}

	// Extrair estrutura se presente
	if structureData, ok := parsed["structure"]; ok {
		if structMap, ok := structureData.(map[string]interface{}); ok {
			extractDirectories(structMap, structure)
			extractFiles(structMap, structure)
		}
	}

	// Extrair dependências
	extractDependencies(parsed, structure)

	// Extrair configurações
	extractConfiguration(parsed, structure)

	// Extrair nome do projeto se presente
	if name, ok := parsed["projectName"].(string); ok {
		structure.ProjectName = name
	}

	return nil
}

// extractFromTextData extrai estrutura de dados de texto não estruturado
func extractFromTextData(data string, language string) (*ProjectStructure, error) {
	structure := &ProjectStructure{
		Language:      language,
		Directories:   make([]string, 0),
		Files:         make(map[string]FileInfo),
		Dependencies:  make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	lines := strings.Split(data, "\n")

	// Padrões de regex para diferentes elementos
	dirPattern := regexp.MustCompile(`(?i)(?:directory|folder|dir):\s*([^\n]+)`)
	filePattern := regexp.MustCompile(`(?i)(?:file|arquivo):\s*([^\n]+)`)
	packagePattern := regexp.MustCompile(`(?i)(?:package|dependency|dependência):\s*([^\n]+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extrair diretórios
		if matches := dirPattern.FindStringSubmatch(line); len(matches) > 1 {
			structure.Directories = append(structure.Directories, strings.TrimSpace(matches[1]))
		}

		// Extrair arquivos
		if matches := filePattern.FindStringSubmatch(line); len(matches) > 1 {
			fileName := strings.TrimSpace(matches[1])
			structure.Files[fileName] = FileInfo{
				Path:      fileName,
				Name:      filepath.Base(fileName),
				Extension: filepath.Ext(fileName),
				Content:   "",
			}
		}

		// Extrair dependências
		if matches := packagePattern.FindStringSubmatch(line); len(matches) > 1 {
			dep := strings.TrimSpace(matches[1])
			structure.Dependencies[dep] = "latest"
		}
	}

	// Análise específica por linguagem
	analyzeLanguageSpecifics(data, language, structure)

	return structure, nil
}

// analyzeLanguageSpecifics analisa aspectos específicos da linguagem
func analyzeLanguageSpecifics(data string, language string, structure *ProjectStructure) {
	switch strings.ToLower(language) {
	case "go", "golang":
		analyzeGoProject(data, structure)
	case "javascript", "js", "typescript", "ts":
		analyzeJSProject(data, structure)
	case "python", "py":
		analyzePythonProject(data, structure)
	case "java":
		analyzeJavaProject(data, structure)
	case "rust":
		analyzeRustProject(data, structure)
	case "c#", "csharp":
		analyzeCSharpProject(data, structure)
	}
}

// analyzeGoProject analisa projetos Go
func analyzeGoProject(data string, structure *ProjectStructure) {
	// Procurar go.mod
	if strings.Contains(data, "go.mod") {
		structure.Configuration["hasGoMod"] = true
		// Extrair module name
		moduleRegex := regexp.MustCompile(`module\s+([^\s\n]+)`)
		if matches := moduleRegex.FindStringSubmatch(data); len(matches) > 1 {
			structure.Metadata["moduleName"] = matches[1]
		}
	}

	// Procurar estrutura padrão
	commonDirs := []string{"cmd", "pkg", "internal", "api", "web", "configs", "scripts", "build", "deployments", "test"}
	for _, dir := range commonDirs {
		if strings.Contains(data, dir+"/") || strings.Contains(data, dir+"\\") {
			structure.Metadata["hasStandardStructure"] = true
			break
		}
	}
}

// analyzeJSProject analisa projetos JavaScript/TypeScript
func analyzeJSProject(data string, structure *ProjectStructure) {
	// Procurar package.json
	if strings.Contains(data, "package.json") {
		structure.Configuration["hasPackageJson"] = true

		// Extrair dependências do package.json se presente
		packageJsonRegex := regexp.MustCompile(`"dependencies":\s*{([^}]+)}`)
		if matches := packageJsonRegex.FindStringSubmatch(data); len(matches) > 1 {
			deps := parseDependenciesFromJson(matches[1])
			for k, v := range deps {
				structure.Dependencies[k] = v
			}
		}
	}

	// Verificar presença de TypeScript
	if strings.Contains(data, "tsconfig.json") || strings.Contains(data, ".ts") {
		structure.Metadata["isTypeScript"] = true
	}

	// Verificar frameworks
	frameworks := map[string]string{
		"react":   "React",
		"vue":     "Vue.js",
		"angular": "Angular",
		"express": "Express.js",
		"next":    "Next.js",
		"nuxt":    "Nuxt.js",
	}

	for fw, name := range frameworks {
		if strings.Contains(strings.ToLower(data), fw) {
			structure.Metadata["framework"] = name
			break
		}
	}
}

// analyzePythonProject analisa projetos Python
func analyzePythonProject(data string, structure *ProjectStructure) {
	// Procurar requirements.txt ou pyproject.toml
	if strings.Contains(data, "requirements.txt") {
		structure.Configuration["hasRequirements"] = true
	}

	if strings.Contains(data, "pyproject.toml") {
		structure.Configuration["hasPyproject"] = true
	}

	// Procurar estrutura de pacote
	if strings.Contains(data, "__init__.py") {
		structure.Metadata["isPackage"] = true
	}

	// Verificar frameworks
	frameworks := map[string]string{
		"django":    "Django",
		"flask":     "Flask",
		"fastapi":   "FastAPI",
		"streamlit": "Streamlit",
		"pytest":    "PyTest",
	}

	for fw, name := range frameworks {
		if strings.Contains(strings.ToLower(data), fw) {
			structure.Metadata["framework"] = name
			break
		}
	}
}

// analyzeJavaProject analisa projetos Java
func analyzeJavaProject(data string, structure *ProjectStructure) {
	// Procurar pom.xml (Maven) ou build.gradle (Gradle)
	if strings.Contains(data, "pom.xml") {
		structure.Configuration["buildTool"] = "Maven"
	} else if strings.Contains(data, "build.gradle") {
		structure.Configuration["buildTool"] = "Gradle"
	}

	// Verificar estrutura Maven padrão
	if strings.Contains(data, "src/main/java") {
		structure.Metadata["hasMavenStructure"] = true
	}

	// Verificar frameworks
	frameworks := map[string]string{
		"spring":    "Spring",
		"hibernate": "Hibernate",
		"junit":     "JUnit",
	}

	for fw, name := range frameworks {
		if strings.Contains(strings.ToLower(data), fw) {
			structure.Metadata["framework"] = name
			break
		}
	}
}

// analyzeRustProject analisa projetos Rust
func analyzeRustProject(data string, structure *ProjectStructure) {
	// Procurar Cargo.toml
	if strings.Contains(data, "Cargo.toml") {
		structure.Configuration["hasCargo"] = true
	}

	// Verificar estrutura padrão
	if strings.Contains(data, "src/main.rs") || strings.Contains(data, "src/lib.rs") {
		structure.Metadata["hasStandardStructure"] = true
	}
}

// analyzeCSharpProject analisa projetos C#
func analyzeCSharpProject(data string, structure *ProjectStructure) {
	// Procurar arquivos de projeto
	if strings.Contains(data, ".csproj") || strings.Contains(data, ".sln") {
		structure.Configuration["hasProjectFile"] = true
	}

	// Verificar frameworks
	frameworks := map[string]string{
		"asp.net": "ASP.NET",
		"blazor":  "Blazor",
		"entity":  "Entity Framework",
	}

	for fw, name := range frameworks {
		if strings.Contains(strings.ToLower(data), fw) {
			structure.Metadata["framework"] = name
			break
		}
	}
}

// Helper functions

func extractJSONFromText(text string) string {
	// Primeiro, tentar fazer parse do texto inteiro como JSON
	var test interface{}
	if json.Unmarshal([]byte(text), &test) == nil {
		return text
	}

	// Se falhar, procurar blocos JSON no texto usando regex mais robusta
	// Esta regex captura JSON balanceado com múltiplos níveis de aninhamento
	var jsonContent strings.Builder
	braceCount := 0
	inJSON := false

	for _, char := range text {
		if char == '{' {
			if !inJSON {
				inJSON = true
				jsonContent.Reset()
			}
			braceCount++
			jsonContent.WriteRune(char)
		} else if char == '}' && inJSON {
			braceCount--
			jsonContent.WriteRune(char)
			if braceCount == 0 {
				// JSON completo encontrado
				jsonStr := jsonContent.String()
				if json.Unmarshal([]byte(jsonStr), &test) == nil {
					return jsonStr
				}
				inJSON = false
			}
		} else if inJSON {
			jsonContent.WriteRune(char)
		}
	}

	return ""
}

func extractDirectories(structMap map[string]interface{}, structure *ProjectStructure) {
	if dirs, ok := structMap["directories"]; ok {
		if dirList, ok := dirs.([]interface{}); ok {
			for _, dir := range dirList {
				if dirStr, ok := dir.(string); ok {
					structure.Directories = append(structure.Directories, dirStr)
				}
			}
		}
	}
}

func extractFiles(structMap map[string]interface{}, structure *ProjectStructure) {
	if files, ok := structMap["files"]; ok {
		if fileMap, ok := files.(map[string]interface{}); ok {
			for fileName, fileContent := range fileMap {
				content := ""

				// Verificar se é string direta
				if contentStr, ok := fileContent.(string); ok {
					content = contentStr
				} else if contentMap, ok := fileContent.(map[string]interface{}); ok {
					// Verificar se é objeto com propriedade content
					if contentStr, ok := contentMap["content"].(string); ok {
						content = contentStr
					}
				}

				structure.Files[fileName] = FileInfo{
					Path:      fileName,
					Name:      filepath.Base(fileName),
					Extension: filepath.Ext(fileName),
					Content:   content,
					IsConfig:  isConfigFile(fileName),
					IsBinary:  isBinaryFile(fileName),
					Size:      len(content),
				}
			}
		}
	}
}

func extractDependencies(parsed map[string]interface{}, structure *ProjectStructure) {
	// Procurar várias formas de dependências
	depKeys := []string{"dependencies", "deps", "packages", "requirements"}

	for _, key := range depKeys {
		if deps, ok := parsed[key]; ok {
			if depMap, ok := deps.(map[string]interface{}); ok {
				for k, v := range depMap {
					structure.Dependencies[k] = v
				}
			}
		}
	}
}

func extractConfiguration(parsed map[string]interface{}, structure *ProjectStructure) {
	// Procurar configurações
	configKeys := []string{"config", "configuration", "settings", "options"}

	for _, key := range configKeys {
		if config, ok := parsed[key]; ok {
			if configMap, ok := config.(map[string]interface{}); ok {
				for k, v := range configMap {
					structure.Configuration[k] = v
				}
			}
		}
	}
}

func parseDependenciesFromJson(depsStr string) map[string]interface{} {
	deps := make(map[string]interface{})

	// Regex simples para extrair dependências do formato "name": "version"
	depRegex := regexp.MustCompile(`"([^"]+)":\s*"([^"]+)"`)
	matches := depRegex.FindAllStringSubmatch(depsStr, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			deps[match[1]] = match[2]
		}
	}

	return deps
}

func isConfigFile(fileName string) bool {
	configFiles := []string{
		".json", ".yaml", ".yml", ".toml", ".ini", ".conf", ".config",
		"Dockerfile", "docker-compose", ".env", "Makefile",
	}

	name := strings.ToLower(fileName)
	for _, ext := range configFiles {
		if strings.HasSuffix(name, ext) || strings.Contains(name, ext) {
			return true
		}
	}

	return false
}

func isBinaryFile(fileName string) bool {
	binaryExts := []string{
		".exe", ".dll", ".so", ".dylib", ".bin", ".app",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico",
		".mp3", ".mp4", ".avi", ".mov", ".pdf",
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	for _, binaryExt := range binaryExts {
		if ext == binaryExt {
			return true
		}
	}

	return false
}
