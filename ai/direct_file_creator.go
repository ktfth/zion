package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateReadmeContent gera um README.md com instruções reais baseadas no projeto
func generateReadmeContent(projectName string, language string, scaffoldResp *ScaffoldResponse) string {
	// Detectar tipo de projeto baseado nos arquivos
	projectType := detectProjectTypeFromFiles(scaffoldResp)

	// Gerar conteúdo do README baseado no tipo e linguagem
	readme := fmt.Sprintf(`# %s

## Descrição

Este projeto foi gerado usando o Zion CLI e contém uma estrutura %s para %s.

## Estrutura do Projeto

`, projectName, projectType, language)

	// Adicionar informações sobre diretórios
	if len(scaffoldResp.Structure.Directories) > 0 {
		readme += "### Diretórios\n\n"
		for _, dir := range scaffoldResp.Structure.Directories {
			readme += fmt.Sprintf("- `%s/` - %s\n", dir, getDirectoryDescription(dir, language))
		}
		readme += "\n"
	}

	// Adicionar informações sobre arquivos importantes
	readme += "### Arquivos Principais\n\n"
	for filePath := range scaffoldResp.Structure.Files {
		if isImportantFile(filePath) {
			readme += fmt.Sprintf("- `%s` - %s\n", filePath, getFileDescription(filePath, language))
		}
	}

	// Adicionar instruções de instalação/configuração
	readme += generateInstallationInstructions(language, scaffoldResp)

	// Adicionar instruções de execução
	readme += generateRunInstructions(language, scaffoldResp)

	// Adicionar próximos passos
	readme += generateNextSteps(language, scaffoldResp)

	return readme
}

// detectProjectTypeFromFiles detecta o tipo de projeto baseado nos arquivos
func detectProjectTypeFromFiles(scaffoldResp *ScaffoldResponse) string {
	files := scaffoldResp.Structure.Files

	if _, hasPackageJson := files["package.json"]; hasPackageJson {
		if _, hasReact := files["src/App.tsx"]; hasReact {
			return "aplicação React"
		}
		if _, hasNext := files["next.config.js"]; hasNext {
			return "aplicação Next.js"
		}
		if _, hasExpress := files["src/server.js"]; hasExpress {
			return "servidor Express"
		}
		return "aplicação Node.js"
	}

	if _, hasGoMod := files["go.mod"]; hasGoMod {
		if _, hasMain := files["main.go"]; hasMain {
			return "aplicação Go"
		}
		return "módulo Go"
	}

	if _, hasRequirements := files["requirements.txt"]; hasRequirements {
		return "aplicação Python"
	}

	if _, hasCargoToml := files["Cargo.toml"]; hasCargoToml {
		return "aplicação Rust"
	}

	return "aplicação"
}

// getDirectoryDescription retorna uma descrição do diretório
func getDirectoryDescription(dir string, language string) string {
	descriptions := map[string]string{
		"src":         "Código fonte principal",
		"tests":       "Testes automatizados",
		"test":        "Testes automatizados",
		"docs":        "Documentação do projeto",
		"config":      "Arquivos de configuração",
		"scripts":     "Scripts utilitários",
		"dist":        "Código compilado/distribuível",
		"build":       "Arquivos de build",
		"public":      "Arquivos públicos (web)",
		"assets":      "Recursos estáticos",
		"components":  "Componentes reutilizáveis",
		"pages":       "Páginas da aplicação",
		"api":         "Endpoints da API",
		"lib":         "Bibliotecas e utilitários",
		"utils":       "Funções utilitárias",
		"middleware":  "Middleware da aplicação",
		"models":      "Modelos de dados",
		"services":    "Serviços da aplicação",
		"controllers": "Controladores da API",
		"routes":      "Definições de rotas",
		"database":    "Arquivos relacionados ao banco de dados",
		"migrations":  "Migrações do banco de dados",
		"internal":    "Código interno (Go)",
		"cmd":         "Comandos executáveis (Go)",
		"pkg":         "Pacotes reutilizáveis (Go)",
	}

	if desc, found := descriptions[dir]; found {
		return desc
	}

	return "Diretório do projeto"
}

// getFileDescription retorna uma descrição do arquivo
func getFileDescription(filePath string, language string) string {
	fileName := filepath.Base(filePath)

	descriptions := map[string]string{
		"package.json":       "Configuração do projeto Node.js e dependências",
		"go.mod":             "Definição do módulo Go e dependências",
		"Cargo.toml":         "Configuração do projeto Rust",
		"requirements.txt":   "Dependências Python",
		"main.go":            "Ponto de entrada da aplicação Go",
		"main.py":            "Ponto de entrada da aplicação Python",
		"index.js":           "Ponto de entrada da aplicação JavaScript",
		"index.ts":           "Ponto de entrada da aplicação TypeScript",
		"tsconfig.json":      "Configuração do TypeScript",
		".gitignore":         "Arquivos ignorados pelo Git",
		"Dockerfile":         "Configuração do container Docker",
		"docker-compose.yml": "Configuração do Docker Compose",
		"Makefile":           "Comandos de build e automação",
		".env":               "Variáveis de ambiente",
		".env.example":       "Exemplo de variáveis de ambiente",
	}

	if desc, found := descriptions[fileName]; found {
		return desc
	}

	// Verificar por extensão
	ext := filepath.Ext(fileName)
	extDescriptions := map[string]string{
		".js":   "Arquivo JavaScript",
		".ts":   "Arquivo TypeScript",
		".go":   "Arquivo Go",
		".py":   "Arquivo Python",
		".rs":   "Arquivo Rust",
		".json": "Arquivo de configuração JSON",
		".yaml": "Arquivo de configuração YAML",
		".yml":  "Arquivo de configuração YAML",
		".toml": "Arquivo de configuração TOML",
		".md":   "Documentação Markdown",
	}

	if desc, found := extDescriptions[ext]; found {
		return desc
	}

	return "Arquivo do projeto"
}

// isImportantFile determina se um arquivo é importante o suficiente para ser mencionado no README
func isImportantFile(filePath string) bool {
	importantFiles := []string{
		"package.json", "go.mod", "Cargo.toml", "requirements.txt",
		"main.go", "main.py", "index.js", "index.ts", "src/index.js", "src/index.ts",
		"tsconfig.json", "Dockerfile", "docker-compose.yml", "Makefile",
		".env.example", "config.yaml", "config.json",
	}

	for _, important := range importantFiles {
		if filePath == important {
			return true
		}
	}

	return false
}

// generateInstallationInstructions gera instruções de instalação
func generateInstallationInstructions(language string, scaffoldResp *ScaffoldResponse) string {
	var instructions strings.Builder

	instructions.WriteString("## Instalação e Configuração\n\n")

	switch strings.ToLower(language) {
	case "javascript", "js", "typescript", "ts":
		if _, hasPackageJson := scaffoldResp.Structure.Files["package.json"]; hasPackageJson {
			instructions.WriteString("### Pré-requisitos\n\n")
			instructions.WriteString("- Node.js (versão 14 ou superior)\n")
			instructions.WriteString("- npm ou yarn\n\n")

			instructions.WriteString("### Instalação de Dependências\n\n")
			instructions.WriteString("```bash\n")
			instructions.WriteString("# Usando npm\n")
			instructions.WriteString("npm install\n\n")
			instructions.WriteString("# Ou usando yarn\n")
			instructions.WriteString("yarn install\n")
			instructions.WriteString("```\n\n")
		}
	case "go", "golang":
		if _, hasGoMod := scaffoldResp.Structure.Files["go.mod"]; hasGoMod {
			instructions.WriteString("### Pré-requisitos\n\n")
			instructions.WriteString("- Go (versão 1.19 ou superior)\n\n")

			instructions.WriteString("### Instalação de Dependências\n\n")
			instructions.WriteString("```bash\n")
			instructions.WriteString("go mod tidy\n")
			instructions.WriteString("```\n\n")
		}
	case "python", "py":
		if _, hasRequirements := scaffoldResp.Structure.Files["requirements.txt"]; hasRequirements {
			instructions.WriteString("### Pré-requisitos\n\n")
			instructions.WriteString("- Python (versão 3.8 ou superior)\n")
			instructions.WriteString("- pip\n\n")

			instructions.WriteString("### Ambiente Virtual (Recomendado)\n\n")
			instructions.WriteString("```bash\n")
			instructions.WriteString("python -m venv venv\n")
			instructions.WriteString("source venv/bin/activate  # No Windows: venv\\Scripts\\activate\n")
			instructions.WriteString("```\n\n")

			instructions.WriteString("### Instalação de Dependências\n\n")
			instructions.WriteString("```bash\n")
			instructions.WriteString("pip install -r requirements.txt\n")
			instructions.WriteString("```\n\n")
		}
	case "rust", "rs":
		if _, hasCargoToml := scaffoldResp.Structure.Files["Cargo.toml"]; hasCargoToml {
			instructions.WriteString("### Pré-requisitos\n\n")
			instructions.WriteString("- Rust (versão 1.70 ou superior)\n")
			instructions.WriteString("- Cargo\n\n")

			instructions.WriteString("### Instalação de Dependências\n\n")
			instructions.WriteString("```bash\n")
			instructions.WriteString("cargo build\n")
			instructions.WriteString("```\n\n")
		}
	}

	return instructions.String()
}

// generateRunInstructions gera instruções de execução
func generateRunInstructions(language string, scaffoldResp *ScaffoldResponse) string {
	var instructions strings.Builder

	instructions.WriteString("## Como Executar\n\n")

	switch strings.ToLower(language) {
	case "javascript", "js", "typescript", "ts":
		if packageJsonContent, hasPackageJson := scaffoldResp.Structure.Files["package.json"]; hasPackageJson {
			instructions.WriteString("### Desenvolvimento\n\n")
			instructions.WriteString("```bash\n")

			// Tentar extrair scripts do package.json
			if contentMap, ok := packageJsonContent.(map[string]interface{}); ok {
				if scripts, hasScripts := contentMap["scripts"].(map[string]interface{}); hasScripts {
					if _, hasDev := scripts["dev"]; hasDev {
						instructions.WriteString("npm run dev\n")
					} else if _, hasStart := scripts["start"]; hasStart {
						instructions.WriteString("npm start\n")
					}
				}
			}

			instructions.WriteString("```\n\n")

			if strings.Contains(strings.ToLower(language), "typescript") {
				instructions.WriteString("### Build\n\n")
				instructions.WriteString("```bash\n")
				instructions.WriteString("npm run build\n")
				instructions.WriteString("```\n\n")
			}
		}
	case "go", "golang":
		instructions.WriteString("### Executar diretamente\n\n")
		instructions.WriteString("```bash\n")
		instructions.WriteString("go run main.go\n")
		instructions.WriteString("```\n\n")

		instructions.WriteString("### Compilar e executar\n\n")
		instructions.WriteString("```bash\n")
		instructions.WriteString("go build -o app\n")
		instructions.WriteString("./app\n")
		instructions.WriteString("```\n\n")
	case "python", "py":
		instructions.WriteString("### Executar aplicação\n\n")
		instructions.WriteString("```bash\n")
		if _, hasMainPy := scaffoldResp.Structure.Files["main.py"]; hasMainPy {
			instructions.WriteString("python main.py\n")
		} else if _, hasAppPy := scaffoldResp.Structure.Files["app.py"]; hasAppPy {
			instructions.WriteString("python app.py\n")
		} else {
			instructions.WriteString("python src/main.py\n")
		}
		instructions.WriteString("```\n\n")
	case "rust", "rs":
		instructions.WriteString("### Executar aplicação\n\n")
		instructions.WriteString("```bash\n")
		instructions.WriteString("cargo run\n")
		instructions.WriteString("```\n\n")

		instructions.WriteString("### Compilar para produção\n\n")
		instructions.WriteString("```bash\n")
		instructions.WriteString("cargo build --release\n")
		instructions.WriteString("```\n\n")
	}

	return instructions.String()
}

// generateNextSteps gera próximos passos
func generateNextSteps(language string, scaffoldResp *ScaffoldResponse) string {
	var steps strings.Builder

	steps.WriteString("## Próximos Passos\n\n")

	// Verificar se há arquivos de exemplo de ambiente
	if _, hasEnvExample := scaffoldResp.Structure.Files[".env.example"]; hasEnvExample {
		steps.WriteString("1. Configure as variáveis de ambiente:\n")
		steps.WriteString("   ```bash\n")
		steps.WriteString("   cp .env.example .env\n")
		steps.WriteString("   # Edite o arquivo .env com suas configurações\n")
		steps.WriteString("   ```\n\n")
	}

	// Verificar se há testes
	hasTests := false
	for filePath := range scaffoldResp.Structure.Files {
		if strings.Contains(filePath, "test") || strings.Contains(filePath, "spec") {
			hasTests = true
			break
		}
	}

	if hasTests {
		steps.WriteString("2. Execute os testes:\n")
		steps.WriteString("   ```bash\n")
		switch strings.ToLower(language) {
		case "javascript", "js", "typescript", "ts":
			steps.WriteString("   npm test\n")
		case "go", "golang":
			steps.WriteString("   go test ./...\n")
		case "python", "py":
			steps.WriteString("   pytest\n")
		case "rust", "rs":
			steps.WriteString("   cargo test\n")
		}
		steps.WriteString("   ```\n\n")
	}

	// Verificar se há Docker
	if _, hasDockerfile := scaffoldResp.Structure.Files["Dockerfile"]; hasDockerfile {
		steps.WriteString("3. Executar com Docker:\n")
		steps.WriteString("   ```bash\n")

		// Tentar extrair o nome do projeto
		projectName := "app"
		if packageJsonContent, hasPackageJson := scaffoldResp.Structure.Files["package.json"]; hasPackageJson {
			if contentMap, ok := packageJsonContent.(map[string]interface{}); ok {
				if name, hasName := contentMap["name"].(string); hasName {
					projectName = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
				}
			}
		}

		steps.WriteString("   docker build -t " + projectName + " .\n")
		steps.WriteString("   docker run -p 3000:3000 " + projectName + "\n")
		steps.WriteString("   ```\n\n")
	}

	steps.WriteString("4. Personalize o código conforme suas necessidades\n")
	steps.WriteString("5. Adicione mais funcionalidades ao projeto\n")
	steps.WriteString("6. Configure CI/CD se necessário\n\n")

	steps.WriteString("## Recursos Adicionais\n\n")
	switch strings.ToLower(language) {
	case "javascript", "js":
		steps.WriteString("- [Documentação do Node.js](https://nodejs.org/docs/)\n")
		steps.WriteString("- [Guia do JavaScript](https://developer.mozilla.org/pt-BR/docs/Web/JavaScript)\n")
	case "typescript", "ts":
		steps.WriteString("- [Documentação do TypeScript](https://www.typescriptlang.org/docs/)\n")
		steps.WriteString("- [Documentação do Node.js](https://nodejs.org/docs/)\n")
	case "go", "golang":
		steps.WriteString("- [Documentação do Go](https://golang.org/doc/)\n")
		steps.WriteString("- [Go by Example](https://gobyexample.com/)\n")
	case "python", "py":
		steps.WriteString("- [Documentação do Python](https://docs.python.org/3/)\n")
		steps.WriteString("- [Guia de Estilo PEP 8](https://www.python.org/dev/peps/pep-0008/)\n")
	case "rust", "rs":
		steps.WriteString("- [Documentação do Rust](https://doc.rust-lang.org/)\n")
		steps.WriteString("- [Rust by Example](https://doc.rust-lang.org/rust-by-example/)\n")
	}

	steps.WriteString("\n---\n")
	steps.WriteString("*Este projeto foi gerado usando [Zion CLI](https://github.com/ktfth/zion)*\n")

	return steps.String()
}

// detectLanguageFromProject detecta a linguagem do projeto baseado nos arquivos
func detectLanguageFromProject(scaffoldResp *ScaffoldResponse) string {
	files := scaffoldResp.Structure.Files

	// Verificar por arquivos específicos de linguagem
	if _, hasPackageJson := files["package.json"]; hasPackageJson {
		// Verificar se há arquivos TypeScript
		for filePath := range files {
			if strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx") {
				return "typescript"
			}
		}
		return "javascript"
	}

	if _, hasGoMod := files["go.mod"]; hasGoMod {
		return "go"
	}

	if _, hasRequirements := files["requirements.txt"]; hasRequirements {
		return "python"
	}

	if _, hasCargoToml := files["Cargo.toml"]; hasCargoToml {
		return "rust"
	}

	// Verificar por extensões de arquivo
	for filePath := range files {
		ext := filepath.Ext(filePath)
		switch ext {
		case ".js", ".jsx":
			return "javascript"
		case ".ts", ".tsx":
			return "typescript"
		case ".go":
			return "go"
		case ".py":
			return "python"
		case ".rs":
			return "rust"
		case ".java":
			return "java"
		case ".cs":
			return "csharp"
		case ".cpp", ".cc", ".cxx":
			return "cpp"
		case ".c":
			return "c"
		}
	}

	return "unknown"
}

// ExtractAndCreateProject extrai diretamente os diretórios e arquivos do JSON
// e cria a estrutura do projeto sem tentar interpretar o conteúdo
func ExtractAndCreateProject(projectName string, jsonStr string) error {
	// Remover blocos de código markdown se presentes
	if strings.HasPrefix(jsonStr, "```json\n") && strings.HasSuffix(jsonStr, "\n```") {
		jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
		jsonStr = strings.TrimSuffix(jsonStr, "\n```")
	}

	// Validar estrutura antes de processar
	validation := ValidateProjectStructure(jsonStr, "")
	if !validation.IsValid {
		fmt.Printf("⚠️  Estrutura do projeto apresenta problemas:\n")
		for _, issue := range validation.Issues {
			fmt.Printf("   • %s\n", issue)
		}
		fmt.Printf("📊 Pontuação de qualidade: %.1f/100\n", validation.Score)

		// Se o score for muito baixo, falhar
		if validation.Score < 50 {
			return fmt.Errorf("projeto não passou na validação (score: %.1f/100)", validation.Score)
		}

		// Caso contrário, mostrar avisos mas continuar
		fmt.Printf("⚠️  Continuando com avisos...\n")
	} else {
		fmt.Printf("✅ Estrutura validada com sucesso (score: %.1f/100)\n", validation.Score)
		if len(validation.Suggestions) > 0 {
			fmt.Printf("💡 Sugestões de melhoria:\n")
			for _, suggestion := range validation.Suggestions {
				fmt.Printf("   • %s\n", suggestion)
			}
		}
	}

	var scaffoldResp ScaffoldResponse
	err := json.Unmarshal([]byte(jsonStr), &scaffoldResp)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório raiz '%s': %v", projectName, err)
	}
	fmt.Printf("\n📁 Criando diretório raiz: %s\n", projectName)

	// Criar diretórios
	if len(scaffoldResp.Structure.Directories) > 0 {
		fmt.Println("\n📂 Criando diretórios:")
		for _, dir := range scaffoldResp.Structure.Directories {
			dirPath := filepath.Join(projectName, dir)
			fmt.Printf("   ├── %s\n", dir)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório '%s': %v", dir, err)
			}
		}
	}

	// Processar cada arquivo
	if len(scaffoldResp.Structure.Files) > 0 {
		fmt.Println("\n📄 Criando arquivos:")

		// Verificar se package.json existe
		hasPackageJson := false
		for filePath := range scaffoldResp.Structure.Files {
			if filePath == "package.json" {
				hasPackageJson = true
				break
			}
		}

		// Se não houver package.json e for um projeto TypeScript/JavaScript, criar um padrão
		if !hasPackageJson && (strings.Contains(strings.ToLower(projectName), "ts") ||
			strings.Contains(strings.ToLower(projectName), "js")) {
			scaffoldResp.Structure.Files["package.json"] = map[string]interface{}{
				"name":        projectName,
				"version":     "1.0.0",
				"description": "Generated by Zion",
				"main":        "dist/index.js",
				"scripts": map[string]string{
					"build": "tsc",
					"start": "node dist/index.js",
					"dev":   "nodemon src/index.ts",
				},
				"dependencies": map[string]string{
					"express": "^4.18.2",
				},
				"devDependencies": map[string]string{
					"@types/express": "^4.17.17",
					"@types/node":    "^20.4.5",
					"nodemon":        "^3.0.1",
					"typescript":     "^5.1.6",
				},
			}
			fmt.Println("   ℹ️  package.json não encontrado, criando versão padrão")
		}

		for filePath, content := range scaffoldResp.Structure.Files {
			var contentStr string

			switch content := content.(type) {
			case map[string]interface{}:
				if fileContent, ok := content["content"]; ok {
					switch fc := fileContent.(type) {
					case string:
						contentStr = ProcessEscapedChars(fc)
					case map[string]interface{}:
						// Se o conteúdo é um objeto JSON (ex: package.json)
						if filePath == "package.json" {
							ProcessPackageJsonContent(fc)
						}
						jsonBytes, err := json.MarshalIndent(fc, "", "  ")
						if err != nil {
							return fmt.Errorf("erro ao serializar conteúdo JSON para '%s': %v", filePath, err)
						}
						contentStr = string(jsonBytes)
					}
				} else {
					// Se não tem "content", assume que o map é o conteúdo direto (para package.json)
					if filePath == "package.json" {
						ProcessPackageJsonContent(content)
					}
					jsonBytes, err := json.MarshalIndent(content, "", "  ")
					if err != nil {
						return fmt.Errorf("erro ao serializar conteúdo direto para '%s': %v", filePath, err)
					}
					contentStr = string(jsonBytes)
				}
			case string:
				if filePath == "package.json" {
					// Se package.json vier como string, tentar fazer parse e processar
					var jsonContent map[string]interface{}
					if err := json.Unmarshal([]byte(content), &jsonContent); err == nil {
						ProcessPackageJsonContent(jsonContent)
						jsonBytes, err := json.MarshalIndent(jsonContent, "", "  ")
						if err != nil {
							return fmt.Errorf("erro ao serializar package.json processado: %v", err)
						}
						contentStr = string(jsonBytes)
					} else {
						contentStr = ProcessEscapedChars(content)
					}
				} else {
					contentStr = ProcessEscapedChars(content)
				}
			}

			// Skip empty files
			if strings.TrimSpace(contentStr) == "" {
				fmt.Printf("   ⚠️  Pulando arquivo vazio: %s\n", filePath)
				continue
			}

			fullPath := filepath.Join(projectName, filePath)
			fmt.Printf("   ├── %s\n", filePath)

			// Garantir que o diretório pai exista
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório pai para '%s': %v", filePath, err)
			}

			if err := os.WriteFile(fullPath, []byte(contentStr), 0644); err != nil {
				return fmt.Errorf("erro ao criar arquivo '%s': %v", filePath, err)
			}
		}
	}

	// Sempre gerar ou melhorar o README.md com instruções reais
	fmt.Println("\n📋 Gerando README.md com instruções reais...")
	readmeContent := generateReadmeContent(projectName, detectLanguageFromProject(&scaffoldResp), &scaffoldResp)

	// Verificar se já existe um README.md
	readmePath := filepath.Join(projectName, "README.md")
	if _, existsInFiles := scaffoldResp.Structure.Files["README.md"]; !existsInFiles {
		// Criar novo README.md
		if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
			fmt.Printf("   ⚠️  Erro ao criar README.md: %v\n", err)
		} else {
			fmt.Printf("   ✅ README.md criado com instruções detalhadas\n")
		}
	} else {
		// Substituir o README.md existente se for muito básico
		if existingContent, ok := scaffoldResp.Structure.Files["README.md"]; ok {
			var existingReadme string
			switch content := existingContent.(type) {
			case map[string]interface{}:
				if fileContent, ok := content["content"]; ok {
					if strContent, ok := fileContent.(string); ok {
						existingReadme = ProcessEscapedChars(strContent)
					}
				}
			case string:
				existingReadme = ProcessEscapedChars(content)
			}

			// Verificar se o README existente é muito básico (menos de 200 caracteres ou sem seções)
			if len(strings.TrimSpace(existingReadme)) < 200 || (!strings.Contains(existingReadme, "##") && !strings.Contains(existingReadme, "Instalação") && !strings.Contains(existingReadme, "Como Executar")) {
				if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
					fmt.Printf("   ⚠️  Erro ao melhorar README.md: %v\n", err)
				} else {
					fmt.Printf("   ✅ README.md melhorado com instruções detalhadas\n")
				}
			} else {
				fmt.Printf("   ℹ️  README.md já contém instruções adequadas\n")
			}
		}
	}

	// Exibir resumo
	fmt.Printf("\n📊 Resumo da estrutura criada:\n")
	fmt.Printf("   ├── %d diretórios\n", len(scaffoldResp.Structure.Directories))
	fmt.Printf("   └── %d arquivos\n", len(scaffoldResp.Structure.Files))

	return nil
}
