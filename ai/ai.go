package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ktfth/zion/ai/providers"
	"github.com/ktfth/zion/config"
	"github.com/ktfth/zion/plugins"
)

type ScaffoldResponse struct {
	Structure struct {
		Directories []string               `json:"directories"`
		Files       map[string]interface{} `json:"files"`
	} `json:"structure"`
}

// GenerateProjectScaffolding gera uma estrutura de projeto com base na linguagem, nome e descrição fornecidos
func GenerateProjectScaffolding(language, projectName, description string, registeredPlugins []string) (string, error) {
	// Substituir SamplePlugin por CorePlugin em registeredPlugins
	for i, plugin := range registeredPlugins {
		if plugin == "SamplePlugin" {
			registeredPlugins[i] = "CorePlugin"
		}
	}

	// Ler contexto do llms.txt se existir
	llmsContext, err := ReadLLMsContext(".")
	if err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao ler contexto llms.txt: %v\n", err)
		llmsContext = &LLMsContext{} // Contexto vazio
	}

	// Se o llms.txt existe, enriquecer os parâmetros
	if llmsContext.HasLLMsFile {
		fmt.Println("📖 Contexto llms.txt detectado - enriquecendo scaffolding...")

		// Auto-detectar linguagem se não especificada
		if language == "" {
			detectedLang := llmsContext.DetectProjectLanguage()
			if detectedLang != "" {
				language = detectedLang
				fmt.Printf("🔍 Linguagem detectada: %s\n", language)
			}
		}

		// Enriquecer descrição com contexto
		if projectDesc := llmsContext.GetProjectDescription(); projectDesc != "" {
			if description == "" {
				description = projectDesc
			} else {
				description = description + "\n\nContexto adicional:\n" + projectDesc
			}
		}
	}

	// Criar o contexto de scaffold para os plugins
	ctx := &plugins.ScaffoldContext{
		ProjectName: projectName,
		Language:    language,
		Description: description,
	}

	// Executar o hook BeforeGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.BeforeGeneration, ctx)

	// Construir o prompt usando o novo sistema melhorado
	prompt := buildImprovedPrompt(language, projectName, description)

	// SEMPRE enriquecer prompt com análise estrutural (com ou sem llms.txt)
	prompt, err = llmsContext.BuildContextualPrompt(prompt, description)
	if err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao construir prompt contextual: %v\n", err)
	}

	// Executar o hook ModifyPrompt para todos os plugins
	ctx.Prompt = prompt
	ctx = plugins.ExecuteHook(plugins.ModifyPrompt, ctx)
	prompt = ctx.Prompt

	// Obter o provedor de IA configurado
	cfg := config.LoadConfig()
	provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, cfg.GetAIConfig())
	if err != nil {
		return "", fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Declarar variável de resposta
	var response string

	// Usar uma estratégia unificada de geração
	response, err = generateWithUnifiedStrategy(provider, prompt, language, projectName, description, llmsContext)
	if err != nil {
		return "", err
	}

	// Atualizar a resposta no contexto
	ctx.Response = response

	// Executar o hook AfterGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.AfterGeneration, ctx)

	// Obter a resposta possivelmente modificada pelos plugins
	response = ctx.Response

	// Processar a resposta antes de retornar
	if response != "" {
		// Verificar se é uma resposta em camadas (contém "layers" e "project_info")
		if strings.Contains(response, `"layers"`) && strings.Contains(response, `"project_info"`) {
			// É uma resposta em camadas, não processar como scaffold tradicional
			return response, nil
		} else {
			// É uma resposta tradicional, processar normalmente
			processedResponse, err := processScaffoldResponse(response)
			if err != nil {
				return "", fmt.Errorf("erro ao processar resposta: %v", err)
			}
			response = processedResponse
		}
	}

	return response, nil
}

// GenerateProjectScaffoldingWithProvider gera uma estrutura de projeto com configurações customizadas de provider
func GenerateProjectScaffoldingWithProvider(language, projectName, description string, registeredPlugins []string, providerName, apiKey, model string) (string, error) {
	// Substituir SamplePlugin por CorePlugin em registeredPlugins
	for i, plugin := range registeredPlugins {
		if plugin == "SamplePlugin" {
			registeredPlugins[i] = "CorePlugin"
		}
	}

	// Ler contexto do llms.txt se existir
	llmsContext, err := ReadLLMsContext(".")
	if err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao ler contexto llms.txt: %v\n", err)
		llmsContext = &LLMsContext{} // Contexto vazio
	}

	// Se o llms.txt existe, enriquecer os parâmetros
	if llmsContext.HasLLMsFile {
		fmt.Println("📖 Contexto llms.txt detectado - enriquecendo scaffolding...")

		// Auto-detectar linguagem se não especificada
		if language == "" {
			detectedLang := llmsContext.DetectProjectLanguage()
			if detectedLang != "" {
				language = detectedLang
				fmt.Printf("🔍 Linguagem detectada: %s\n", language)
			}
		}

		// Enriquecer descrição com contexto
		if projectDesc := llmsContext.GetProjectDescription(); projectDesc != "" {
			if description == "" {
				description = projectDesc
			} else {
				description = description + "\n\nContexto adicional:\n" + projectDesc
			}
		}
	}

	// Criar o contexto de scaffold para os plugins
	ctx := &plugins.ScaffoldContext{
		ProjectName: projectName,
		Language:    language,
		Description: description,
	}

	// Executar o hook BeforeGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.BeforeGeneration, ctx)

	// Construir o prompt usando o novo sistema melhorado
	prompt := buildImprovedPrompt(language, projectName, description)

	// SEMPRE enriquecer prompt com análise estrutural (com ou sem llms.txt)
	prompt, err = llmsContext.BuildContextualPrompt(prompt, description)
	if err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao construir prompt contextual: %v\n", err)
	}

	// Executar o hook ModifyPrompt para todos os plugins
	ctx.Prompt = prompt
	ctx = plugins.ExecuteHook(plugins.ModifyPrompt, ctx)
	prompt = ctx.Prompt

	// Obter configurações do provider
	cfg := config.LoadConfig()

	// Usar provider customizado se especificado
	if providerName != "" {
		cfg.AIProvider = providerName
	}

	// Obter configuração do AI
	aiConfig := cfg.GetAIConfig()

	// Override da API key se especificado
	if apiKey != "" {
		aiConfig["api_key"] = apiKey
	}

	// Override do modelo se especificado
	if model != "" {
		aiConfig["model"] = model
	}

	// Criar o provider
	provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, aiConfig)
	if err != nil {
		return "", fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Verificar se há risco de overflow de contexto antes de gerar
	var response string
	if DetectContextOverflow(prompt, provider.Name()) {
		fmt.Printf("⚠️  Contexto muito grande detectado - usando geração em camadas\n")

		// Usar gerador em camadas
		layeredGen, err := NewLayeredGenerator(language, projectName, description, llmsContext)
		if err != nil {
			return "", fmt.Errorf("erro ao criar gerador em camadas: %v", err)
		}

		layeredResponse, err := layeredGen.GenerateLayeredProject()
		if err != nil {
			return "", fmt.Errorf("erro na geração em camadas: %v", err)
		}

		// Serializar a resposta em camadas diretamente
		responseBytes, err := json.MarshalIndent(layeredResponse, "", "  ")
		if err != nil {
			return "", fmt.Errorf("erro ao serializar resposta em camadas: %v", err)
		}

		response = string(responseBytes)
	} else {
		// Gerar conteúdo usando o provedor normalmente
		fmt.Printf("🤖 Usando provedor de IA: %s\n", provider.Name())
		response, err = provider.GenerateContent(prompt)
		if err != nil {
			// Verificar se é erro de contexto e tentar geração em camadas
			if IsContextOverflowError(err) {
				fmt.Printf("❌ Erro de contexto detectado - tentando geração em camadas\n")

				layeredGen, layerErr := NewLayeredGenerator(language, projectName, description, llmsContext)
				if layerErr != nil {
					return "", fmt.Errorf("erro original: %v, erro ao criar gerador em camadas: %v", err, layerErr)
				}

				layeredResponse, layerErr := layeredGen.GenerateLayeredProject()
				if layerErr != nil {
					return "", fmt.Errorf("erro original: %v, erro na geração em camadas: %v", err, layerErr)
				}

				// Serializar a resposta em camadas diretamente
				responseBytes, marshalErr := json.MarshalIndent(layeredResponse, "", "  ")
				if marshalErr != nil {
					return "", fmt.Errorf("erro original: %v, erro ao serializar resposta em camadas: %v", err, marshalErr)
				}

				response = string(responseBytes)
			} else {
				return "", err
			}
		}
	}

	// Atualizar a resposta no contexto
	ctx.Response = response

	// Executar o hook AfterGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.AfterGeneration, ctx)

	// Obter a resposta possivelmente modificada pelos plugins
	response = ctx.Response

	// Processar a resposta antes de retornar
	if response != "" {
		// Verificar se é uma resposta em camadas (contém "layers" e "project_info")
		if strings.Contains(response, `"layers"`) && strings.Contains(response, `"project_info"`) {
			// É uma resposta em camadas, não processar como scaffold tradicional
			return response, nil
		} else {
			// É uma resposta tradicional, processar normalmente
			processedResponse, err := processScaffoldResponse(response)
			if err != nil {
				return "", fmt.Errorf("erro ao processar resposta: %v", err)
			}
			response = processedResponse
		}
	}

	return response, nil
}

// TestLayeredGeneration executa testes do sistema de geração em camadas
func TestLayeredGeneration() {
	fmt.Printf("1️⃣ Testando detecção de overflow...\n")

	// Criar um prompt muito grande
	largePrompt := "Este é um prompt de teste "
	for i := 0; i < 1000; i++ {
		largePrompt += "muito grande com muitas palavras repetitivas para simular um contexto extenso "
	}

	isOverflow := DetectContextOverflow(largePrompt, "gpt")
	if isOverflow {
		fmt.Printf("✅ Detecção de overflow funcionando (prompt: %d chars, estimativa: %d tokens)\n",
			len(largePrompt), estimateTokens(largePrompt))
	} else {
		fmt.Printf("❌ Falha na detecção de overflow (prompt: %d chars, estimativa: %d tokens)\n",
			len(largePrompt), estimateTokens(largePrompt))
	}

	// Teste 2: Detecção de erro de contexto
	fmt.Printf("\n2️⃣ Testando detecção de erros de contexto...\n")

	testErrors := []string{
		"This endpoint's maximum context length is 200000 tokens",
		"API retornou status 400: token limit exceeded",
		"context too long",
		"input too long",
		"reduce the length of either one",
	}

	for _, errMsg := range testErrors {
		testErr := fmt.Errorf(errMsg)
		if IsContextOverflowError(testErr) {
			fmt.Printf("✅ Detectou erro de contexto: %s\n", errMsg[:min(50, len(errMsg))]+"}")
		} else {
			fmt.Printf("❌ Falha na detecção: %s\n", errMsg[:min(50, len(errMsg))]+"}")
		}
	}

	// Teste 3: Criação do gerador em camadas
	fmt.Printf("\n3️⃣ Testando criação do gerador em camadas...\n")

	llmsContext := &LLMsContext{}
	layeredGen, err := NewLayeredGenerator("go", "test-project", "Test project description", llmsContext)
	if err != nil {
		fmt.Printf("❌ Erro ao criar gerador: %v\n", err)
	} else {
		fmt.Printf("✅ Gerador em camadas criado com sucesso\n")
		fmt.Printf("   Provider: %s\n", layeredGen.provider.Name())
		fmt.Printf("   Max tokens: %d\n", layeredGen.maxTokens)
	}
}

// min retorna o menor dos dois valores int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Rest of the existing functions remain the same...

// getLanguageSpecificRequirements retorna requisitos específicos para cada linguagem
func getLanguageSpecificRequirements(language string) string {
	switch strings.ToLower(language) {
	case "js", "javascript":
		return `

JAVASCRIPT REQUIREMENTS:
1. Modern ES6+ structure with modules
2. Package.json with comprehensive scripts
3. ESLint configuration for code quality
4. Prettier for code formatting
5. Jest/Vitest for testing
6. JSDoc for documentation
7. Node.js best practices
8. Express.js framework setup (if web app)
9. Environment configuration (.env)
10. Build and deployment scripts`

	case "ts", "typescript":
		return `

TYPESCRIPT REQUIREMENTS:
1. Optimized TSConfig configuration
2. Well-defined types and interfaces
3. Organized module structure
4. ESLint with TypeScript rules
5. Prettier configuration
6. Jest/Vitest with TypeScript support
7. TSDoc for documentation
8. Proper import/export patterns
9. Build pipeline configuration
10. Type declaration files`

	case "go", "golang":
		return `

GO REQUIREMENTS:
1. Proper Go module structure (go.mod)
2. Idiomatic Go patterns and conventions
3. Golangci-lint configuration
4. Makefile with useful commands
5. Unit tests with table-driven tests
6. Go-style documentation comments
7. Proper error handling patterns
8. Dependency management with go.mod
9. Docker configuration
10. CI/CD pipeline setup`

	case "rs", "rust":
		return `

RUST REQUIREMENTS:
1. Cargo workspace structure
2. Well-organized modules and crates
3. Robust error handling with Result/Option
4. Clippy configuration for linting
5. Unit and integration tests
6. Rustdoc documentation
7. CI/CD with cargo commands
8. Proper use of traits and generics
9. Memory safety best practices
10. Performance optimization patterns`

	case "cs", "csharp":
		return `

C# REQUIREMENTS:
1. Modern .NET solution structure
2. Clean Architecture layers (DDD patterns)
3. Code analysis and linting rules
4. xUnit/NUnit testing framework
5. XML documentation comments
6. Build and deployment scripts
7. NuGet package management
8. Dependency injection setup
9. Configuration management
10. Logging and monitoring setup`

	case "python":
		return `

PYTHON REQUIREMENTS:
1. Modern Python packaging (pyproject.toml)
2. Virtual environment setup
3. Requirements.txt with pinned versions
4. Black formatter configuration
5. Flake8/Pylint for linting
6. Pytest for testing
7. Sphinx for documentation
8. Type hints (mypy support)
9. Docker configuration
10. CI/CD pipeline setup`

	default:
		return `

GENERAL REQUIREMENTS:
1. Clean and organized project structure
2. Comprehensive documentation
3. Testing framework setup
4. Code quality tools
5. Build and deployment scripts
6. Environment configuration
7. Version control setup
8. CI/CD pipeline basics
9. Dependency management
10. Error handling and logging`
	}
}

// processScaffoldResponse processa a resposta do scaffold para garantir JSON válido
func processScaffoldResponse(response string) (string, error) {
	// Remove blocos de código markdown se presentes
	if strings.HasPrefix(response, "```json\n") && strings.HasSuffix(response, "\n```") {
		response = strings.TrimPrefix(response, "```json\n")
		response = strings.TrimSuffix(response, "\n```")
	} else if strings.HasPrefix(response, "```json") && strings.HasSuffix(response, "```") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") && strings.HasSuffix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}

	// Limpar whitespaces extras
	response = strings.TrimSpace(response)

	// Tentar extrair JSON se há texto antes/depois
	jsonStart := -1
	jsonEnd := -1

	// Procurar pelo início do JSON (procurar por { que não esteja dentro de uma string)
	for i, char := range response {
		if char == '{' {
			jsonStart = i
			break
		}
	}

	// Procurar pelo fim do JSON (procurar pelo } balanceado)
	if jsonStart >= 0 {
		braceCount := 0
		inString := false
		escaped := false

		for i := jsonStart; i < len(response); i++ {
			char := response[i]

			if escaped {
				escaped = false
				continue
			}

			if char == '\\' {
				escaped = true
				continue
			}

			if char == '"' {
				inString = !inString
				continue
			}

			if !inString {
				if char == '{' {
					braceCount++
				} else if char == '}' {
					braceCount--
					if braceCount == 0 {
						jsonEnd = i
						break
					}
				}
			}
		}
	}

	// Se encontrou início e fim válidos, extrair o JSON
	if jsonStart >= 0 && jsonEnd >= 0 {
		response = response[jsonStart : jsonEnd+1]
	} else if jsonStart >= 0 {
		// Se só encontrou o início, tentar do início até o fim
		response = response[jsonStart:]
	}

	// Primeiro, vamos tentar fazer parse do JSON base
	var baseStruct struct {
		Structure struct {
			Directories []string               `json:"directories"`
			Files       map[string]interface{} `json:"files"`
		} `json:"structure"`
	}

	if err := json.Unmarshal([]byte(response), &baseStruct); err != nil {
		return "", fmt.Errorf("erro no parse inicial: %v", err)
	}

	// Processar cada arquivo
	for filename, fileContent := range baseStruct.Structure.Files {
		switch content := fileContent.(type) {
		case map[string]interface{}:
			if filename == "package.json" {
				if jsonContent, ok := content["content"].(map[string]interface{}); ok {
					ProcessPackageJsonContent(jsonContent)
					baseStruct.Structure.Files[filename] = jsonContent
				}
			} else {
				baseStruct.Structure.Files[filename] = content
			}
		case string:
			if filename == "package.json" {
				// Se o conteúdo é uma string, tentar fazer parse para JSON
				var jsonContent map[string]interface{}
				if err := json.Unmarshal([]byte(content), &jsonContent); err == nil {
					ProcessPackageJsonContent(jsonContent)
					baseStruct.Structure.Files[filename] = jsonContent
				}
			} else {
				baseStruct.Structure.Files[filename] = content
			}
		}
	}

	// Converter de volta para JSON
	result, err := json.MarshalIndent(baseStruct, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao gerar JSON final: %v", err)
	}

	// Validar a estrutura resultante
	validation := ValidateProjectStructure(string(result), "")
	if !validation.IsValid {
		fmt.Printf("⚠️  Estrutura do projeto apresenta problemas:\n")
		for _, issue := range validation.Issues {
			fmt.Printf("   • %s\n", issue)
		}
		fmt.Printf("📊 Pontuação de qualidade: %.1f/100\n", validation.Score)

		// Se o score for muito baixo, falhar
		if validation.Score < 30 {
			return "", fmt.Errorf("projeto não passou na validação após processamento (score: %.1f/100)", validation.Score)
		}

		// Caso contrário, mostrar avisos mas continuar
		fmt.Printf("⚠️  Continuando com avisos...\n")
	} else {
		fmt.Printf("✅ Estrutura processada e validada com sucesso (score: %.1f/100)\n", validation.Score)
		if len(validation.Suggestions) > 0 {
			fmt.Printf("💡 Sugestões de melhoria:\n")
			for _, suggestion := range validation.Suggestions {
				fmt.Printf("   • %s\n", suggestion)
			}
		}
	}

	return string(result), nil
}

// cleanJSONString limpa e corrige problemas comuns em strings JSON
func cleanJSONString(input string) string {
	// Remove caracteres invisíveis e espaços em branco extras
	input = strings.TrimSpace(input)

	// Pré-processa as dependências
	input = preprocessDependencies(input)

	// Corrige aspas dentro de strings
	input = fixQuotesInJSON(input)

	return input
}

// preprocessDependencies faz um pré-processamento específico nas seções de dependências
func preprocessDependencies(input string) string {
	// Regex para encontrar blocos de dependencies e devDependencies
	depsRegex := regexp.MustCompile(`"(dev)?dependencies"\s*:\s*{([^}]+)}`)

	return depsRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Processa cada pacote dentro do bloco de dependências
		packageRegex := regexp.MustCompile(`"(@[^"]+)"\s*:\s*"([^"]+)"`)
		processed := packageRegex.ReplaceAllString(match, `"\\u0040$1": "$2"`)
		return processed
	})
}

// fixQuotesInJSON corrige problemas com aspas em strings JSON
func fixQuotesInJSON(input string) string {
	// Regex para encontrar aspas simples em valores
	valueRegex := regexp.MustCompile(`:\s*'([^']*)'`)
	return valueRegex.ReplaceAllString(input, `: "$1"`)
}

// buildBasePrompt constrói o prompt base com instruções específicas
func buildBasePrompt(language, projectName, description string) string {
	basePrompt := fmt.Sprintf(`ROLE: Você é um arquiteto de software sênior especializado em %s com 15+ anos de experiência.

TASK: Crie uma estrutura de projeto moderna, profissional e escalável para '%s'.

QUALITY STANDARDS:
1. Arquitetura limpa e modular seguindo princípios SOLID
2. Padrões de projeto adequados para %s
3. Estrutura de diretórios organizada e escalável
4. Configuração de ambiente flexível e robusta
5. Documentação completa e clara
6. Testes automatizados configurados
7. Ferramentas de desenvolvimento (linting, formatação)
8. Scripts de build e deployment
9. Configuração de CI/CD básica
10. Tratamento de erros e logging`, language, projectName, language)

	// Adicionar descrição específica se fornecida
	if description != "" {
		basePrompt += fmt.Sprintf(`

PROJECT REQUIREMENTS:
%s`, description)
	}

	// Adicionar requisitos específicos por linguagem
	basePrompt += getLanguageSpecificRequirements(language)

	return basePrompt
}

// buildJSONInstructions constrói instruções específicas para formato JSON
func buildJSONInstructions(language string) string {
	instructions := `JSON FORMATTING RULES:

1. STRUCTURE: Use exactly this structure without modifications
2. ENCODING: UTF-8 encoding for all content
3. QUOTES: Use double quotes (") for all strings
4. ESCAPING: Properly escape special characters (\n, \t, \", \\)
5. OBJECTS: For JSON files (package.json, etc.), use nested object in "content"
6. STRINGS: For text files, use string in "content"
7. ARRAYS: Use proper array syntax with square brackets
8. BOOLEANS: Use true/false (lowercase, no quotes)
9. NUMBERS: Use numeric values without quotes
10. NULL: Use null (lowercase, no quotes) for null values`

	// Adicionar instruções específicas por linguagem
	switch strings.ToLower(language) {
	case "js", "javascript", "ts", "typescript":
		instructions += `

JAVASCRIPT/TYPESCRIPT SPECIFIC:
- For package.json: Use proper semver versions (^1.0.0)
- For npm packages starting with @: Use normal syntax "@package/name"
- For scripts: Use cross-platform commands when possible
- For dependencies: Include both dependencies and devDependencies
- For TypeScript: Include proper type definitions`

	case "go", "golang":
		instructions += `

GO SPECIFIC:
- For go.mod: Use proper module syntax
- For main.go: Include proper package declaration
- For imports: Use proper import grouping
- For comments: Use Go-style comments (// and /* */)
- For configuration: Use proper Go idioms`

	case "python":
		instructions += `

PYTHON SPECIFIC:
- For requirements.txt: Use proper version specifications
- For setup.py/pyproject.toml: Use modern Python packaging
- For __init__.py: Include proper package initialization
- For imports: Use proper Python import style
- For configuration: Use proper Python configuration patterns`
	}

	return instructions
}

// buildLanguageExamples constrói exemplos específicos da linguagem
func buildLanguageExamples(language string) string {
	switch strings.ToLower(language) {
	case "js", "javascript":
		return `JAVASCRIPT EXAMPLE STRUCTURE:
{
  "structure": {
    "directories": ["src", "tests", "docs", "config"],
    "files": {
      "package.json": {
        "content": {
          "name": "project-name",
          "version": "1.0.0",
          "description": "Project description",
          "main": "src/index.js",
          "scripts": {
            "start": "node src/index.js",
            "dev": "nodemon src/index.js",
            "test": "jest",
            "lint": "eslint src/",
            "format": "prettier --write src/"
          },
          "dependencies": {
            "express": "^4.18.0"
          },
          "devDependencies": {
            "jest": "^29.0.0",
            "eslint": "^8.0.0",
            "prettier": "^3.0.0",
            "nodemon": "^3.0.0"
          }
        }
      },
      "src/index.js": {
        "content": "// Main application entry point\nconsole.log('Hello, World!');"
      },
      ".eslintrc.json": {
        "content": {
          "env": {
            "node": true,
            "es6": true
          },
          "extends": ["eslint:recommended"],
          "parserOptions": {
            "ecmaVersion": 2022,
            "sourceType": "module"
          }
        }
      }
    }
  }
}`

	case "ts", "typescript":
		return `TYPESCRIPT EXAMPLE STRUCTURE:
{
  "structure": {
    "directories": ["src", "tests", "dist", "docs"],
    "files": {
      "package.json": {
        "content": {
          "name": "project-name",
          "version": "1.0.0",
          "description": "Project description",
          "main": "dist/index.js",
          "scripts": {
            "build": "tsc",
            "start": "node dist/index.js",
            "dev": "ts-node src/index.ts",
            "test": "jest"
          },
          "dependencies": {},
          "devDependencies": {
            "typescript": "^5.0.0",
            "@types/node": "^20.0.0",
            "ts-node": "^10.0.0",
            "jest": "^29.0.0"
          }
        }
      },
      "tsconfig.json": {
        "content": {
          "compilerOptions": {
            "target": "ES2022",
            "module": "commonjs",
            "outDir": "./dist",
            "rootDir": "./src",
            "strict": true,
            "esModuleInterop": true
          },
          "include": ["src/**/*"],
          "exclude": ["node_modules", "dist"]
        }
      },
      "src/index.ts": {
        "content": "// Main application entry point\nconsole.log('Hello, TypeScript!');"
      }
    }
  }
}`

	case "go", "golang":
		return `GO EXAMPLE STRUCTURE:
{
  "structure": {
    "directories": ["cmd", "internal", "pkg", "api", "docs"],
    "files": {
      "go.mod": {
        "content": "module project-name\n\ngo 1.21\n\nrequire (\n\t// Add dependencies here\n)"
      },
      "main.go": {
        "content": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, Go!\")\n}"
      },
      "Makefile": {
        "content": ".PHONY: build test clean\n\nbuild:\n\tgo build -o bin/app ./cmd/app\n\ntest:\n\tgo test ./...\n\nclean:\n\trm -rf bin/"
      }
    }
  }
}`

	case "python":
		return `PYTHON EXAMPLE STRUCTURE:
{
  "structure": {
    "directories": ["src", "tests", "docs", "scripts"],
    "files": {
      "requirements.txt": {
        "content": "# Production dependencies\nrequests>=2.31.0\n\n# Development dependencies\npytest>=7.0.0\nblack>=23.0.0\nflake8>=6.0.0"
      },
      "setup.py": {
        "content": "from setuptools import setup, find_packages\n\nsetup(\n    name=\"project-name\",\n    version=\"0.1.0\",\n    packages=find_packages(where=\"src\"),\n    package_dir={\"\": \"src\"},\n    install_requires=[\n        \"requests>=2.31.0\",\n    ],\n)"
      },
      "src/main.py": {
        "content": "#!/usr/bin/env python3\n\"\"\"Main application entry point.\"\"\"\n\ndef main():\n    \"\"\"Main function.\"\"\"\n    print(\"Hello, Python!\")\n\nif __name__ == \"__main__\":\n    main()"
      }
    }
  }
}`

	default:
		return fmt.Sprintf(`GENERIC EXAMPLE STRUCTURE:
{
  "structure": {
    "directories": ["src", "tests", "docs"],
    "files": {
      "README.md": {
        "content": "# Project Name\n\nProject description here."
      },
      "src/main.%s": {
        "content": "// Main application file"
      }
    }
  }
}`, getFileExtension(language))
	}
}

// getFileExtension retorna a extensão de arquivo apropriada para a linguagem
func getFileExtension(language string) string {
	switch strings.ToLower(language) {
	case "js", "javascript":
		return "js"
	case "ts", "typescript":
		return "ts"
	case "go", "golang":
		return "go"
	case "python", "py":
		return "py"
	case "rust", "rs":
		return "rs"
	case "java":
		return "java"
	case "csharp", "cs":
		return "cs"
	case "cpp", "c++":
		return "cpp"
	case "c":
		return "c"
	default:
		return "txt"
	}
}

// BuildImprovedPrompt constrói um prompt melhorado e mais consistente (função exportada)
func BuildImprovedPrompt(language, projectName, description string) string {
	return buildImprovedPrompt(language, projectName, description)
}

// buildImprovedPrompt constrói um prompt melhorado e mais consistente
func buildImprovedPrompt(language, projectName, description string) string {
	// Criar controlador adaptativo para geração normal
	instructionController := NewAdaptiveInstructionController(detectProjectType(description), language, description)

	// Construir o prompt base com instruções específicas
	basePrompt := buildBasePrompt(language, projectName, description)

	// Aplicar controle adaptativo de instruções
	adaptivePrompt := instructionController.BuildAdaptivePrompt(basePrompt)

	// Construir instruções de formato JSON específicas
	jsonInstructions := buildJSONInstructions(language)

	// Construir exemplos específicos da linguagem
	languageExamples := buildLanguageExamples(language)

	// Construir o prompt final
	return fmt.Sprintf(`%s

%s

%s

FORMATO DE SAÍDA OBRIGATÓRIO:
Retorne APENAS um JSON válido seguindo exatamente esta estrutura:

{
  "structure": {
    "directories": ["lista", "de", "diretórios"],
    "files": {
      "nome-do-arquivo.ext": {
        "content": "conteúdo do arquivo para arquivos de texto"
      },
      "package.json": {
        "content": {
          "name": "nome-do-projeto",
          "version": "1.0.0"
        }
      }
    }
  }
}

REGRAS CRÍTICAS:
1. Responda APENAS com JSON válido
2. Não adicione texto explicativo antes ou depois do JSON
3. Use aspas duplas para todas as strings
4. Garanta que todos os colchetes e chaves estão balanceados
5. Não use comentários dentro do JSON
6. Para arquivos JSON, use objeto aninhado em "content"
7. Para arquivos de texto, use string em "content"
8. Escape caracteres especiais adequadamente (\n, \t, \", \\)
9. Use valores booleanos sem aspas (true, false)
10. Use valores numéricos sem aspas

VALIDAÇÃO:
- O JSON deve passar em JSON.parse() sem erros
- Todos os arquivos devem ter conteúdo válido e realista
- A estrutura deve ser coerente com a linguagem %s
- Inclua pelo menos 5-8 arquivos essenciais
- Inclua pelo menos 3-5 diretórios organizados
- Conteúdo dos arquivos deve ser funcional e não apenas placeholder
- OBEDEÇA RIGOROSAMENTE às instruções de escopo e restrições especificadas`, adaptivePrompt, jsonInstructions, languageExamples, language)
}

// generateWithUnifiedStrategy implementa uma estratégia unificada de geração
func generateWithUnifiedStrategy(provider providers.Provider, prompt, language, projectName, description string, llmsContext *LLMsContext) (string, error) {
	// Estratégia 1: Tentar geração normal primeiro
	fmt.Printf("🤖 Usando provedor de IA: %s\n", provider.Name())

	// Verificar se o prompt é muito grande
	if EstimateTokens(prompt) > determineMaxTokens(provider.Name()) {
		fmt.Printf("⚠️  Contexto muito grande detectado - usando geração em camadas\n")
		return generateWithLayeredStrategy(provider, language, projectName, description, llmsContext)
	}

	// Tentar geração normal
	response, err := provider.GenerateContent(prompt)
	if err != nil {
		// Se é erro de contexto, tentar camadas
		if IsContextOverflowError(err) {
			fmt.Printf("❌ Erro de contexto detectado - tentando geração em camadas\n")
			return generateWithLayeredStrategy(provider, language, projectName, description, llmsContext)
		}
		return "", err
	}

	return response, nil
}

// generateWithLayeredStrategy gera usando o sistema de camadas
func generateWithLayeredStrategy(provider providers.Provider, language, projectName, description string, llmsContext *LLMsContext) (string, error) {
	layeredGen, err := NewLayeredGenerator(language, projectName, description, llmsContext)
	if err != nil {
		return "", fmt.Errorf("erro ao criar gerador em camadas: %v", err)
	}

	layeredResponse, err := layeredGen.GenerateLayeredProject()
	if err != nil {
		return "", fmt.Errorf("erro na geração em camadas: %v", err)
	}

	// Serializar a resposta em camadas
	responseBytes, err := json.MarshalIndent(layeredResponse, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta em camadas: %v", err)
	}

	return string(responseBytes), nil
}

// CreateProjectWithUnifiedStrategy usa uma estratégia unificada para criar projetos
func CreateProjectWithUnifiedStrategy(projectName, response string, llmsContext *LLMsContext, isContextualMode bool) error {
	// Detectar linguagem para validação
	language := ""
	if llmsContext != nil {
		language = llmsContext.DetectProjectLanguage()
	}

	// Validar estrutura antes de criar o projeto
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

	// Determinar se é resposta em camadas analisando a estrutura JSON
	var layeredResponse LayeredResponse
	if json.Unmarshal([]byte(response), &layeredResponse) == nil && len(layeredResponse.Layers) > 0 {
		// É resposta em camadas - usar criação em camadas
		return CreateLayeredProject(projectName, &layeredResponse)
	}

	// Resposta tradicional
	if isContextualMode && llmsContext != nil && llmsContext.HasLLMsFile {
		// Modo contextual com llms.txt - usar criação contextual
		return CreateContextualProject(projectName, response, llmsContext)
	}

	// Modo normal - usar criação padrão
	return ExtractAndCreateProject(projectName, response)
}
