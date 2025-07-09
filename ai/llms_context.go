package ai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LLMsContext representa o contexto extraído de arquivos llms.txt
type LLMsContext struct {
	ProjectRoot      string
	ContextFiles     []string
	HasLLMsFile      bool
	LLMsContent      string
	ParsedContext    map[string]string
	ProjectStructure []string
}

// ReadLLMsContext lê e processa arquivos llms.txt para obter contexto do projeto
func ReadLLMsContext(projectPath string) (*LLMsContext, error) {
	context := &LLMsContext{
		ProjectRoot:   projectPath,
		ParsedContext: make(map[string]string),
	}

	// Procurar por arquivo llms.txt
	llmsPath := filepath.Join(projectPath, "llms.txt")
	if _, err := os.Stat(llmsPath); err == nil {
		context.HasLLMsFile = true

		// Ler conteúdo do llms.txt
		content, err := ioutil.ReadFile(llmsPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler llms.txt: %v", err)
		}

		context.LLMsContent = string(content)

		// Processar conteúdo
		err = context.parseLLMsContent()
		if err != nil {
			return nil, fmt.Errorf("erro ao processar llms.txt: %v", err)
		}
	}

	// Análise da estrutura do projeto existente
	err := context.analyzeProjectStructure()
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar estrutura do projeto: %v", err)
	}

	return context, nil
}

// parseLLMsContent processa o conteúdo do arquivo llms.txt
func (ctx *LLMsContext) parseLLMsContent() error {
	lines := strings.Split(ctx.LLMsContent, "\n")

	var currentSection string
	var currentContent strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Detectar seções (marcadas com ### ou ## ou #)
		if strings.HasPrefix(line, "#") {
			// Salvar seção anterior se existir
			if currentSection != "" && currentContent.Len() > 0 {
				ctx.ParsedContext[currentSection] = strings.TrimSpace(currentContent.String())
				currentContent.Reset()
			}

			// Nova seção
			currentSection = strings.TrimSpace(strings.TrimLeft(line, "#"))
			continue
		}

		// Detectar referências a arquivos (iniciadas com ./ ou /)
		if strings.HasPrefix(line, "./") || strings.HasPrefix(line, "/") {
			ctx.ContextFiles = append(ctx.ContextFiles, line)
			continue
		}

		// Adicionar conteúdo à seção atual
		if currentSection != "" {
			currentContent.WriteString(line + "\n")
		} else {
			// Conteúdo sem seção específica
			if ctx.ParsedContext["general"] == "" {
				ctx.ParsedContext["general"] = ""
			}
			ctx.ParsedContext["general"] += line + "\n"
		}
	}

	// Salvar última seção
	if currentSection != "" && currentContent.Len() > 0 {
		ctx.ParsedContext[currentSection] = strings.TrimSpace(currentContent.String())
	}

	return nil
}

// analyzeProjectStructure analisa a estrutura do projeto existente
func (ctx *LLMsContext) analyzeProjectStructure() error {
	if _, err := os.Stat(ctx.ProjectRoot); os.IsNotExist(err) {
		// Projeto não existe ainda
		return nil
	}

	err := filepath.Walk(ctx.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar diretórios específicos
		if info.IsDir() {
			dirName := info.Name()
			if dirName == ".git" || dirName == "node_modules" || dirName == "vendor" ||
				dirName == "__pycache__" || dirName == ".next" || dirName == "dist" ||
				dirName == "build" || dirName == "target" {
				return filepath.SkipDir
			}
		}

		// Adicionar à estrutura
		relPath, err := filepath.Rel(ctx.ProjectRoot, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			ctx.ProjectStructure = append(ctx.ProjectStructure, relPath)
		}

		return nil
	})

	return err
}

// GetProjectDescription gera uma descrição do projeto baseada no contexto
func (ctx *LLMsContext) GetProjectDescription() string {
	if !ctx.HasLLMsFile {
		return ""
	}

	var description strings.Builder

	// Adicionar descrição geral se existir
	if general, exists := ctx.ParsedContext["general"]; exists && general != "" {
		description.WriteString(strings.TrimSpace(general))
		description.WriteString("\n\n")
	}

	// Adicionar informações específicas de seções importantes
	importantSections := []string{"purpose", "objetivo", "description", "descrição", "overview", "visão geral"}
	for _, section := range importantSections {
		if content, exists := ctx.ParsedContext[section]; exists && content != "" {
			description.WriteString(fmt.Sprintf("**%s**: %s\n", section, strings.TrimSpace(content)))
		}
	}

	return strings.TrimSpace(description.String())
}

// GetTechnicalRequirements extrai requisitos técnicos do contexto
func (ctx *LLMsContext) GetTechnicalRequirements() string {
	if !ctx.HasLLMsFile {
		return ""
	}

	var requirements strings.Builder

	// Procurar seções técnicas
	technicalSections := []string{"requirements", "requisitos", "tech stack", "tecnologias", "dependencies", "dependências", "architecture", "arquitetura"}
	for _, section := range technicalSections {
		if content, exists := ctx.ParsedContext[section]; exists && content != "" {
			requirements.WriteString(fmt.Sprintf("**%s**: %s\n", section, strings.TrimSpace(content)))
		}
	}

	return strings.TrimSpace(requirements.String())
}

// GetContextualFiles lê o conteúdo dos arquivos referenciados
func (ctx *LLMsContext) GetContextualFiles() (map[string]string, error) {
	contextFiles := make(map[string]string)

	for _, fileRef := range ctx.ContextFiles {
		// Resolver caminho relativo
		filePath := filepath.Join(ctx.ProjectRoot, fileRef)

		// Verificar se arquivo existe
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // Arquivo não existe, pular
		}

		// Ler conteúdo do arquivo
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler arquivo %s: %v", fileRef, err)
		}

		contextFiles[fileRef] = string(content)
	}

	return contextFiles, nil
}

// BuildContextualPrompt constrói um prompt enriquecido com o contexto do llms.txt
func (ctx *LLMsContext) BuildContextualPrompt(basePrompt string, userDescription string) (string, error) {
	if !ctx.HasLLMsFile && len(ctx.ProjectStructure) == 0 {
		return basePrompt, nil
	}

	var contextualPrompt strings.Builder

	// Adicionar prompt base
	contextualPrompt.WriteString(basePrompt)
	contextualPrompt.WriteString("\n\n")

	// NOVA FUNCIONALIDADE: Análise estrutural inteligente
	structuralAnalysis := ctx.GenerateStructuralRecommendations(userDescription)
	contextualPrompt.WriteString(structuralAnalysis)

	// Adicionar contexto do projeto (se llms.txt existir)
	if ctx.HasLLMsFile {
		contextualPrompt.WriteString("=== CONTEXTO DO PROJETO (de llms.txt) ===\n")

		// Descrição do projeto
		if desc := ctx.GetProjectDescription(); desc != "" {
			contextualPrompt.WriteString("DESCRIÇÃO DO PROJETO:\n")
			contextualPrompt.WriteString(desc)
			contextualPrompt.WriteString("\n\n")
		}

		// Requisitos técnicos
		if req := ctx.GetTechnicalRequirements(); req != "" {
			contextualPrompt.WriteString("REQUISITOS TÉCNICOS:\n")
			contextualPrompt.WriteString(req)
			contextualPrompt.WriteString("\n\n")
		}
	}

	// Estrutura existente do projeto
	if len(ctx.ProjectStructure) > 0 {
		contextualPrompt.WriteString("ESTRUTURA EXISTENTE DO PROJETO:\n")
		for _, item := range ctx.ProjectStructure {
			contextualPrompt.WriteString(fmt.Sprintf("- %s\n", item))
		}
		contextualPrompt.WriteString("\n")
	}

	// Arquivos de contexto
	contextFiles, err := ctx.GetContextualFiles()
	if err != nil {
		return "", err
	}

	if len(contextFiles) > 0 {
		contextualPrompt.WriteString("ARQUIVOS DE CONTEXTO:\n")
		for filePath, content := range contextFiles {
			contextualPrompt.WriteString(fmt.Sprintf("\n--- %s ---\n", filePath))
			// Limitar tamanho do conteúdo para evitar prompts muito grandes
			if len(content) > 2000 {
				content = content[:2000] + "\n... (truncado)"
			}
			contextualPrompt.WriteString(content)
			contextualPrompt.WriteString("\n")
		}
		contextualPrompt.WriteString("\n")
	}

	// Instruções específicas baseadas no contexto
	contextualPrompt.WriteString("INSTRUÇÕES ESPECÍFICAS:\n")
	contextualPrompt.WriteString("1. Use o contexto acima para gerar uma estrutura mais apropriada e específica\n")
	contextualPrompt.WriteString("2. Mantenha a consistência com arquivos e estruturas existentes\n")
	contextualPrompt.WriteString("3. Expanda ou melhore o projeto baseado nas informações do llms.txt\n")
	contextualPrompt.WriteString("4. Preserve configurações e padrões já estabelecidos\n")
	contextualPrompt.WriteString("5. Adicione apenas recursos complementares e melhorias\n\n")

	return contextualPrompt.String(), nil
}

// DetectProjectLanguage tenta detectar a linguagem do projeto baseada na estrutura
func (ctx *LLMsContext) DetectProjectLanguage() string {
	// Verificar por arquivos característicos na estrutura atual
	for _, item := range ctx.ProjectStructure {
		switch {
		case strings.HasSuffix(item, "package.json"):
			return "javascript"
		case strings.HasSuffix(item, "tsconfig.json"):
			return "typescript"
		case strings.HasSuffix(item, "go.mod") || strings.HasSuffix(item, "go.sum"):
			return "go"
		case strings.HasSuffix(item, "requirements.txt") || strings.HasSuffix(item, "pyproject.toml") || strings.HasSuffix(item, "setup.py"):
			return "python"
		case strings.HasSuffix(item, "Cargo.toml"):
			return "rust"
		case strings.HasSuffix(item, "pom.xml") || strings.HasSuffix(item, "build.gradle"):
			return "java"
		case strings.HasSuffix(item, "Gemfile"):
			return "ruby"
		case strings.HasSuffix(item, "composer.json"):
			return "php"
		case strings.HasSuffix(item, ".csproj") || strings.HasSuffix(item, ".sln"):
			return "csharp"
		}
	}

	// Verificar por extensões de arquivos se não encontrou arquivos de configuração
	languageCount := make(map[string]int)
	for _, item := range ctx.ProjectStructure {
		switch {
		case strings.HasSuffix(item, ".js") || strings.HasSuffix(item, ".mjs"):
			languageCount["javascript"]++
		case strings.HasSuffix(item, ".ts") || strings.HasSuffix(item, ".tsx"):
			languageCount["typescript"]++
		case strings.HasSuffix(item, ".go"):
			languageCount["go"]++
		case strings.HasSuffix(item, ".py"):
			languageCount["python"]++
		case strings.HasSuffix(item, ".rs"):
			languageCount["rust"]++
		case strings.HasSuffix(item, ".java"):
			languageCount["java"]++
		case strings.HasSuffix(item, ".rb"):
			languageCount["ruby"]++
		case strings.HasSuffix(item, ".php"):
			languageCount["php"]++
		case strings.HasSuffix(item, ".cs"):
			languageCount["csharp"]++
		}
	}

	// Retornar a linguagem com mais arquivos
	maxCount := 0
	detectedLang := ""
	for lang, count := range languageCount {
		if count > maxCount {
			maxCount = count
			detectedLang = lang
		}
	}

	if detectedLang != "" {
		return detectedLang
	}

	// Verificar no contexto do llms.txt se disponível
	if ctx.HasLLMsFile {
		content := strings.ToLower(ctx.LLMsContent)
		languages := map[string][]string{
			"javascript": {"javascript", "js", "node", "npm", "yarn"},
			"typescript": {"typescript", "ts", "tsx"},
			"python":     {"python", "py", "pip", "django", "flask"},
			"go":         {"golang", "go", "gin", "echo"},
			"rust":       {"rust", "cargo"},
			"java":       {"java", "maven", "gradle", "spring"},
			"php":        {"php", "composer", "laravel"},
			"ruby":       {"ruby", "rails", "gem"},
		}

		for lang, keywords := range languages {
			for _, keyword := range keywords {
				if strings.Contains(content, keyword) {
					return lang
				}
			}
		}
	}

	return ""
}

// CreateContextualProject cria arquivos e estruturas de forma inteligente em projetos existentes
func CreateContextualProject(projectName, response string, ctx *LLMsContext) error {
	// Detectar linguagem para validação
	language := ""
	if ctx != nil {
		language = ctx.DetectProjectLanguage()
	}

	// Validar estrutura antes de processar
	validation := ValidateProjectStructure(response, language)
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

	var scaffoldData ScaffoldResponse

	// Tentar fazer parse do JSON
	jsonErr := json.Unmarshal([]byte(response), &scaffoldData)
	if jsonErr != nil {
		// Se falhou, tentar método de fallback
		return createContextualProjectFallback(projectName, response, ctx)
	}

	// Se temos contexto existente, ser mais cuidadoso
	if ctx.HasLLMsFile && len(ctx.ProjectStructure) > 0 {
		return mergeWithExistingProject(scaffoldData, ctx)
	}

	// Caso contrário, usar método padrão
	return ExtractAndCreateProject(projectName, response)
}

// mergeWithExistingProject mescla a nova estrutura com o projeto existente
func mergeWithExistingProject(scaffoldData ScaffoldResponse, ctx *LLMsContext) error {
	projectRoot := ctx.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório atual: %v", err)
		}
	}

	fmt.Printf("🔧 Mesclando com projeto existente em: %s\n", projectRoot)

	// Criar diretórios que não existem
	for _, dir := range scaffoldData.Structure.Directories {
		dirPath := filepath.Join(projectRoot, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("📁 Criando diretório: %s\n", dir)
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
			}
		} else {
			fmt.Printf("📁 Diretório já existe: %s\n", dir)
		}
	}

	// Processar arquivos de forma inteligente
	for fileName, content := range scaffoldData.Structure.Files {
		filePath := filepath.Join(projectRoot, fileName)

		// Verificar se arquivo já existe
		if _, err := os.Stat(filePath); err == nil {
			// Arquivo existe - perguntar que ação tomar
			action := decideFileAction(fileName, filePath, content)
			switch action {
			case "skip":
				fmt.Printf("⏭️  Pulando arquivo existente: %s\n", fileName)
				continue
			case "backup":
				err = createBackupAndReplace(filePath, content)
				if err != nil {
					fmt.Printf("⚠️  Erro ao fazer backup de %s: %v\n", fileName, err)
					continue
				}
				fmt.Printf("💾 Backup criado e arquivo atualizado: %s\n", fileName)
			case "merge":
				err = mergeFileContent(filePath, content)
				if err != nil {
					fmt.Printf("⚠️  Erro ao mesclar %s: %v\n", fileName, err)
					continue
				}
				fmt.Printf("🔀 Arquivo mesclado: %s\n", fileName)
			case "replace":
				err = writeFileContent(filePath, content)
				if err != nil {
					fmt.Printf("⚠️  Erro ao substituir %s: %v\n", fileName, err)
					continue
				}
				fmt.Printf("🔄 Arquivo substituído: %s\n", fileName)
			}
		} else {
			// Arquivo não existe - criar normalmente
			fmt.Printf("📄 Criando novo arquivo: %s\n", fileName)

			// Criar diretório pai se necessário
			dir := filepath.Dir(filePath)
			if dir != "." {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return fmt.Errorf("erro ao criar diretório para %s: %v", fileName, err)
				}
			}

			err = writeFileContent(filePath, content)
			if err != nil {
				return fmt.Errorf("erro ao criar arquivo %s: %v", fileName, err)
			}
		}
	}

	return nil
}

// decideFileAction decide que ação tomar com um arquivo existente
func decideFileAction(fileName, filePath string, newContent interface{}) string {
	// Lógica baseada no tipo de arquivo
	switch {
	case strings.HasSuffix(fileName, ".md"):
		return "merge" // Arquivos markdown podem ser mesclados
	case strings.HasSuffix(fileName, "package.json"):
		return "merge" // package.json deve ser mesclado
	case strings.HasSuffix(fileName, ".json"):
		return "backup" // Outros JSONs - fazer backup
	case strings.HasSuffix(fileName, ".toml"):
		return "merge" // Arquivos TOML podem ser mesclados
	case strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml"):
		return "merge" // YAML podem ser mesclados
	case strings.HasSuffix(fileName, ".gitignore"):
		return "merge" // .gitignore deve ser mesclado
	case strings.HasSuffix(fileName, "README"):
		return "backup" // README importante - fazer backup
	case strings.Contains(fileName, "config"):
		return "backup" // Arquivos de configuração - backup
	default:
		return "skip" // Por padrão, pular arquivos existentes
	}
}

// createBackupAndReplace cria backup do arquivo existente e substitui pelo novo
func createBackupAndReplace(filePath string, newContent interface{}) error {
	// Criar backup
	backupPath := filePath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(backupPath, input, 0644)
	if err != nil {
		return err
	}

	// Substituir com novo conteúdo
	return writeFileContent(filePath, newContent)
}

// mergeFileContent mescla conteúdo novo com arquivo existente
func mergeFileContent(filePath string, newContent interface{}) error {
	// Implementação básica de merge baseada no tipo de arquivo
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		return mergeJSONFile(filePath, newContent)
	case ".md":
		return mergeMarkdownFile(filePath, newContent)
	case ".toml":
		return mergeTomlFile(filePath, newContent)
	case ".yaml", ".yml":
		return mergeYamlFile(filePath, newContent)
	default:
		return appendToFile(filePath, newContent)
	}
}

// mergeJSONFile mescla dois arquivos JSON
func mergeJSONFile(filePath string, newContent interface{}) error {
	// Ler arquivo existente
	existingData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var existing map[string]interface{}
	err = json.Unmarshal(existingData, &existing)
	if err != nil {
		return err
	}

	// Processar novo conteúdo
	var newData map[string]interface{}
	switch v := newContent.(type) {
	case string:
		err = json.Unmarshal([]byte(v), &newData)
		if err != nil {
			return err
		}
	case map[string]interface{}:
		newData = v
	default:
		return fmt.Errorf("tipo de conteúdo não suportado para JSON: %T", newContent)
	}

	// Mesclar dados
	for key, value := range newData {
		existing[key] = value
	}

	// Escrever arquivo mesclado
	mergedData, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, mergedData, 0644)
}

// mergeMarkdownFile mescla arquivos Markdown
func mergeMarkdownFile(filePath string, newContent interface{}) error {
	existingData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var newContentStr string
	switch v := newContent.(type) {
	case string:
		newContentStr = v
	default:
		return fmt.Errorf("tipo de conteúdo não suportado para Markdown: %T", newContent)
	}

	// Adicionar seção ao final
	mergedContent := string(existingData) + "\n\n" + newContentStr

	return ioutil.WriteFile(filePath, []byte(mergedContent), 0644)
}

// mergeTomlFile mescla arquivos TOML (implementação básica)
func mergeTomlFile(filePath string, newContent interface{}) error {
	// Para simplicidade, vamos fazer append.
	// Uma implementação completa usaria uma biblioteca TOML
	return appendToFile(filePath, newContent)
}

// mergeYamlFile mescla arquivos YAML (implementação básica)
func mergeYamlFile(filePath string, newContent interface{}) error {
	// Para simplicidade, vamos fazer append.
	// Uma implementação completa usaria uma biblioteca YAML
	return appendToFile(filePath, newContent)
}

// appendToFile adiciona conteúdo ao final do arquivo
func appendToFile(filePath string, newContent interface{}) error {
	var contentStr string
	switch v := newContent.(type) {
	case string:
		contentStr = v
	default:
		return fmt.Errorf("tipo de conteúdo não suportado: %T", newContent)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("\n" + contentStr)
	return err
}

// writeFileContent escreve conteúdo em um arquivo
func writeFileContent(filePath string, content interface{}) error {
	var data []byte
	var err error

	switch v := content.(type) {
	case string:
		data = []byte(v)
	case map[string]interface{}:
		data, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("tipo de conteúdo não suportado: %T", content)
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// createContextualProjectFallback método de fallback para criação contextual
func createContextualProjectFallback(projectName, response string, ctx *LLMsContext) error {
	// Se temos um projeto existente, salvar resposta em arquivo temporário
	if ctx.HasLLMsFile && len(ctx.ProjectStructure) > 0 {
		tempFile := "zion_scaffold_output.md"
		fmt.Printf("⚠️  Salvando resposta em %s para revisão manual\n", tempFile)
		return ioutil.WriteFile(tempFile, []byte(response), 0644)
	}

	// Caso contrário, usar método padrão
	return SaveRawResponse(projectName, response)
}

// AnalyzeProjectNeeds analisa a estrutura existente e identifica o que precisa ser criado
func (ctx *LLMsContext) AnalyzeProjectNeeds() *ProjectAnalysis {
	analysis := &ProjectAnalysis{
		Language:        ctx.DetectProjectLanguage(),
		ExistingFiles:   ctx.ProjectStructure,
		MissingPatterns: make(map[string][]string),
		Recommendations: []string{},
		ProjectType:     ctx.detectProjectType(),
	}

	// Análise baseada na linguagem detectada
	switch analysis.Language {
	case "javascript", "typescript":
		analysis = ctx.analyzeJSProject(analysis)
	case "go":
		analysis = ctx.analyzeGoProject(analysis)
	case "python":
		analysis = ctx.analyzePythonProject(analysis)
	case "java":
		analysis = ctx.analyzeJavaProject(analysis)
	case "rust":
		analysis = ctx.analyzeRustProject(analysis)
	default:
		analysis = ctx.analyzeGenericProject(analysis)
	}

	return analysis
}

// ProjectAnalysis contém informações sobre o que está faltando no projeto
type ProjectAnalysis struct {
	Language        string
	ProjectType     string
	ExistingFiles   []string
	MissingPatterns map[string][]string
	Recommendations []string
	StructuralGaps  []string
}

// detectProjectType identifica o tipo de projeto baseado na estrutura
func (ctx *LLMsContext) detectProjectType() string {
	// Analisar arquivos de configuração e estrutura
	hasAPI := false
	hasWeb := false
	hasCLI := false
	hasLib := false

	for _, item := range ctx.ProjectStructure {
		item = strings.ToLower(item)
		switch {
		case strings.Contains(item, "api") || strings.Contains(item, "server") || strings.Contains(item, "routes"):
			hasAPI = true
		case strings.Contains(item, "web") || strings.Contains(item, "frontend") || strings.Contains(item, "client"):
			hasWeb = true
		case strings.Contains(item, "cli") || strings.Contains(item, "cmd") || strings.Contains(item, "main"):
			hasCLI = true
		case strings.Contains(item, "lib") || strings.Contains(item, "package"):
			hasLib = true
		}
	}

	// Determinar tipo baseado nas características
	if hasAPI && hasWeb {
		return "fullstack"
	} else if hasAPI {
		return "api"
	} else if hasWeb {
		return "frontend"
	} else if hasCLI {
		return "cli"
	} else if hasLib {
		return "library"
	}

	return "generic"
}

// analyzeJSProject analisa projetos JavaScript/TypeScript
func (ctx *LLMsContext) analyzeJSProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	existingFiles := make(map[string]bool)
	for _, file := range analysis.ExistingFiles {
		existingFiles[file] = true
	}

	// Verificar estrutura típica de projetos JS
	expectedStructure := map[string][]string{
		"configuração": {
			"package.json", ".gitignore", ".eslintrc.json", ".prettierrc.json",
			"tsconfig.json", // para TypeScript
		},
		"código": {
			"src/index.js", "src/app.js", "src/main.js",
			"src/utils/", "src/lib/", "src/components/",
		},
		"testes": {
			"tests/", "test/", "__tests__/", "spec/",
			"jest.config.js", "vitest.config.js",
		},
		"documentação": {
			"README.md", "docs/", "CHANGELOG.md",
		},
		"build": {
			"webpack.config.js", "vite.config.js", "rollup.config.js",
			"build/", "dist/",
		},
		"deploy": {
			"Dockerfile", "docker-compose.yml", ".github/workflows/",
			"vercel.json", "netlify.toml",
		},
	}

	// Identificar o que está faltando
	for category, files := range expectedStructure {
		missing := []string{}
		for _, file := range files {
			found := false
			for existing := range existingFiles {
				if strings.Contains(existing, file) || strings.Contains(file, existing) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, file)
			}
		}
		if len(missing) > 0 {
			analysis.MissingPatterns[category] = missing
		}
	}

	// Gerar recomendações específicas
	if len(analysis.MissingPatterns["testes"]) > 0 {
		analysis.Recommendations = append(analysis.Recommendations, "Implementar testes automatizados com Jest ou Vitest")
	}
	if len(analysis.MissingPatterns["build"]) > 0 {
		analysis.Recommendations = append(analysis.Recommendations, "Configurar sistema de build (Webpack, Vite, ou Rollup)")
	}
	if len(analysis.MissingPatterns["deploy"]) > 0 {
		analysis.Recommendations = append(analysis.Recommendations, "Adicionar configuração de deploy (Docker, CI/CD)")
	}

	// Análise específica do tipo de projeto
	switch analysis.ProjectType {
	case "api":
		ctx.analyzeAPIStructure(analysis)
	case "frontend":
		ctx.analyzeFrontendStructure(analysis)
	case "fullstack":
		ctx.analyzeFullstackStructure(analysis)
	}

	return analysis
}

// analyzeAPIStructure analisa estrutura específica de APIs
func (ctx *LLMsContext) analyzeAPIStructure(analysis *ProjectAnalysis) {
	apiPatterns := []string{"routes/", "controllers/", "middleware/", "models/", "services/"}
	missing := []string{}

	for _, pattern := range apiPatterns {
		found := false
		for _, existing := range analysis.ExistingFiles {
			if strings.Contains(existing, pattern) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, pattern)
		}
	}

	if len(missing) > 0 {
		analysis.MissingPatterns["api"] = missing
		analysis.Recommendations = append(analysis.Recommendations,
			"Implementar estrutura de API com rotas, controllers e middleware")
	}

	// Verificar padrões de segurança
	securityFiles := []string{"auth/", "validation/", "security/"}
	securityMissing := []string{}
	for _, file := range securityFiles {
		found := false
		for _, existing := range analysis.ExistingFiles {
			if strings.Contains(existing, file) {
				found = true
				break
			}
		}
		if !found {
			securityMissing = append(securityMissing, file)
		}
	}

	if len(securityMissing) > 0 {
		analysis.MissingPatterns["segurança"] = securityMissing
		analysis.Recommendations = append(analysis.Recommendations,
			"Adicionar autenticação, validação e middleware de segurança")
	}
}

// analyzeFrontendStructure analisa estrutura específica de frontend
func (ctx *LLMsContext) analyzeFrontendStructure(analysis *ProjectAnalysis) {
	frontendPatterns := []string{"components/", "pages/", "hooks/", "styles/", "assets/"}
	missing := []string{}

	for _, pattern := range frontendPatterns {
		found := false
		for _, existing := range analysis.ExistingFiles {
			if strings.Contains(existing, pattern) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, pattern)
		}
	}

	if len(missing) > 0 {
		analysis.MissingPatterns["frontend"] = missing
		analysis.Recommendations = append(analysis.Recommendations,
			"Organizar estrutura de componentes, páginas e estilos")
	}
}

// analyzeFullstackStructure analisa estrutura de projetos fullstack
func (ctx *LLMsContext) analyzeFullstackStructure(analysis *ProjectAnalysis) {
	ctx.analyzeAPIStructure(analysis)
	ctx.analyzeFrontendStructure(analysis)

	// Verificar separação client/server
	hasClientServer := false
	for _, existing := range analysis.ExistingFiles {
		if (strings.Contains(existing, "client/") && strings.Contains(existing, "server/")) ||
			(strings.Contains(existing, "frontend/") && strings.Contains(existing, "backend/")) {
			hasClientServer = true
			break
		}
	}

	if !hasClientServer {
		analysis.Recommendations = append(analysis.Recommendations,
			"Separar código em diretórios client/ e server/ ou frontend/ e backend/")
	}
}

// analyzeGoProject analisa projetos Go
func (ctx *LLMsContext) analyzeGoProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	goStructure := map[string][]string{
		"configuração": {"go.mod", "go.sum", ".gitignore"},
		"código":       {"main.go", "cmd/", "internal/", "pkg/"},
		"testes":       {"*_test.go"},
		"documentação": {"README.md", "docs/"},
		"deploy":       {"Dockerfile", "Makefile"},
	}

	// Similar análise para Go...
	for category, patterns := range goStructure {
		missing := []string{}
		for _, pattern := range patterns {
			found := false
			for _, existing := range analysis.ExistingFiles {
				if strings.Contains(existing, pattern) ||
					(strings.Contains(pattern, "*") && strings.HasSuffix(existing, strings.TrimPrefix(pattern, "*"))) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, pattern)
			}
		}
		if len(missing) > 0 {
			analysis.MissingPatterns[category] = missing
		}
	}

	return analysis
}

// analyzePythonProject analisa projetos Python
func (ctx *LLMsContext) analyzePythonProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	pythonStructure := map[string][]string{
		"configuração": {"requirements.txt", "pyproject.toml", "setup.py", ".gitignore"},
		"código":       {"src/", "__init__.py", "main.py"},
		"testes":       {"tests/", "test_*.py"},
		"documentação": {"README.md", "docs/"},
		"ambiente":     {".env", "venv/", ".venv/"},
	}

	// Análise similar para Python...
	for category, patterns := range pythonStructure {
		missing := []string{}
		for _, pattern := range patterns {
			found := false
			for _, existing := range analysis.ExistingFiles {
				if strings.Contains(existing, pattern) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, pattern)
			}
		}
		if len(missing) > 0 {
			analysis.MissingPatterns[category] = missing
		}
	}

	return analysis
}

// analyzeJavaProject analisa projetos Java
func (ctx *LLMsContext) analyzeJavaProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	// Implementar análise Java...
	return analysis
}

// analyzeRustProject analisa projetos Rust
func (ctx *LLMsContext) analyzeRustProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	// Implementar análise Rust...
	return analysis
}

// analyzeGenericProject análise genérica para outras linguagens
func (ctx *LLMsContext) analyzeGenericProject(analysis *ProjectAnalysis) *ProjectAnalysis {
	// Análise genérica baseada em padrões comuns
	commonPatterns := map[string][]string{
		"código":       {"src/", "lib/"},
		"testes":       {"tests/", "test/"},
		"documentação": {"README.md", "docs/"},
		"configuração": {".gitignore"},
	}

	for category, patterns := range commonPatterns {
		missing := []string{}
		for _, pattern := range patterns {
			found := false
			for _, existing := range analysis.ExistingFiles {
				if strings.Contains(existing, pattern) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, pattern)
			}
		}
		if len(missing) > 0 {
			analysis.MissingPatterns[category] = missing
		}
	}

	return analysis
}

// GenerateStructuralRecommendations gera recomendações baseadas na análise
func (ctx *LLMsContext) GenerateStructuralRecommendations(userDescription string) string {
	analysis := ctx.AnalyzeProjectNeeds()

	var recommendations strings.Builder

	recommendations.WriteString("=== ANÁLISE ESTRUTURAL DO PROJETO ===\n\n")
	recommendations.WriteString(fmt.Sprintf("🔍 Linguagem: %s\n", analysis.Language))
	recommendations.WriteString(fmt.Sprintf("🏗️ Tipo de projeto: %s\n", analysis.ProjectType))
	recommendations.WriteString(fmt.Sprintf("📁 Arquivos existentes: %d\n\n", len(analysis.ExistingFiles)))

	// Mostrar estrutura atual
	if len(analysis.ExistingFiles) > 0 {
		recommendations.WriteString("📂 ESTRUTURA ATUAL:\n")
		for _, file := range analysis.ExistingFiles {
			if file != "." {
				recommendations.WriteString(fmt.Sprintf("   %s\n", file))
			}
		}
		recommendations.WriteString("\n")
	}

	// Mostrar gaps estruturais
	if len(analysis.MissingPatterns) > 0 {
		recommendations.WriteString("⚠️ GAPS ESTRUTURAIS IDENTIFICADOS:\n")
		for category, missing := range analysis.MissingPatterns {
			recommendations.WriteString(fmt.Sprintf("\n🔴 %s:\n", strings.ToUpper(category)))
			for _, item := range missing {
				recommendations.WriteString(fmt.Sprintf("   - %s\n", item))
			}
		}
		recommendations.WriteString("\n")
	}

	// Mostrar recomendações
	if len(analysis.Recommendations) > 0 {
		recommendations.WriteString("💡 RECOMENDAÇÕES AUTOMÁTICAS:\n")
		for i, rec := range analysis.Recommendations {
			recommendations.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec))
		}
		recommendations.WriteString("\n")
	}

	// Conectar com a solicitação do usuário
	if userDescription != "" {
		recommendations.WriteString("🎯 SOLICITAÇÃO DO USUÁRIO:\n")
		recommendations.WriteString(userDescription)
		recommendations.WriteString("\n\n")

		recommendations.WriteString("📋 ESTRATÉGIA DE IMPLEMENTAÇÃO:\n")
		recommendations.WriteString("Com base na análise estrutural e na solicitação, implemente:\n")
		recommendations.WriteString("1. Complete os gaps estruturais identificados\n")
		recommendations.WriteString("2. Adicione a funcionalidade solicitada\n")
		recommendations.WriteString("3. Mantenha consistência com a arquitetura existente\n")
		recommendations.WriteString("4. Siga as melhores práticas da linguagem\n\n")
	}

	return recommendations.String()
}
