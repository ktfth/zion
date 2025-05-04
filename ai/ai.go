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

	// Criar o contexto de scaffold para os plugins
	ctx := &plugins.ScaffoldContext{
		ProjectName: projectName,
		Language:    language,
		Description: description,
	}

	// Executar o hook BeforeGeneration para todos os plugins
	ctx = plugins.ExecuteHook(plugins.BeforeGeneration, ctx)

	// Construir a descrição do projeto com mais detalhes e boas práticas
	projectDesc := fmt.Sprintf(`Você é um especialista em desenvolvimento de software com vasta experiência em %s.
Crie uma estrutura moderna e profissional para um projeto chamado '%s'.

O projeto deve seguir:
1. Arquitetura limpa e modular
2. Padrões de projeto adequados à linguagem %s
3. Estrutura de diretórios organizada e escalável
4. Configuração de ambiente flexível
5. Documentação clara e objetiva`, language, projectName, language)

	// Adicionar descrição específica se fornecida
	if description != "" {
		projectDesc += fmt.Sprintf(`\n\nRequisitos específicos:\n%s`, description)
	}

	// Adicionar requisitos específicos por linguagem
	projectDesc += getLanguageSpecificRequirements(language)

	prompt := fmt.Sprintf(`%s

IMPORTANTE: Para garantir um JSON válido, siga estas regras:

1. Use apenas aspas duplas (") para strings
2. Para valores de arquivos JSON (como package.json), use a seguinte sintaxe:
   "arquivo.json": {
     "content": {
       // conteúdo do JSON aqui
     }
   }
3. Para outros arquivos de texto, use a seguinte sintaxe:
   "arquivo.txt": {
     "content": "conteúdo do arquivo"
   }
4. Para nomes de pacotes npm que começam com @, use a seguinte sintaxe:
   "dependencies": {
     "pkg:@types/node": "^20.4.8",
     "pkg:@typescript-eslint/parser": "^6.7.5"
   }

Retorne um JSON com esta estrutura exata:
{
  "structure": {
    "directories": ["dir1", "dir2"],
    "files": {
      "arquivo.json": {
        "content": {
          // conteúdo JSON aqui
        }
      },
      "arquivo.txt": {
        "content": "conteúdo texto aqui"
      }
    }
  }
}`, projectDesc)

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

	// Gerar conteúdo usando o provedor
	fmt.Printf("🤖 Usando provedor de IA: %s\n", provider.Name())
	response, err := provider.GenerateContent(prompt)
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
		processedResponse, err := processScaffoldResponse(response)
		if err != nil {
			return "", fmt.Errorf("erro ao processar resposta: %v", err)
		}
		response = processedResponse
	}

	return response, nil
}

// getLanguageSpecificRequirements retorna requisitos específicos para cada linguagem
func getLanguageSpecificRequirements(language string) string {
	switch strings.ToLower(language) {
	case "js", "javascript":
		return `\n\nRequisitos específicos para JavaScript:
1. Estrutura moderna com ES6+
2. Sistema de módulos ES
3. Configuração de linting (ESLint)
4. Configuração de formatação (Prettier)
5. Scripts NPM úteis
6. Testes unitários configurados
7. Documentação com JSDoc`

	case "ts", "typescript":
		return `\n\nRequisitos específicos para TypeScript:
1. Configuração do TSConfig otimizada
2. Tipos bem definidos
3. Estrutura de módulos organizada
4. Configuração de linting (ESLint)
5. Configuração de formatação (Prettier)
6. Scripts NPM úteis
7. Testes unitários com Jest/Vitest
8. Documentação com TSDoc`

	case "go", "golang":
		return `\n\nRequisitos específicos para Go:
1. Estrutura de módulos Go
2. Padrões idiomáticos Go
3. Configuração de linting (golangci-lint)
4. Makefile com comandos úteis
5. Testes unitários
6. Documentação no estilo Go
7. Gerenciamento de dependências com go.mod`

	case "rs", "rust":
		return `\n\nRequisitos específicos para Rust:
1. Estrutura de workspace Cargo
2. Módulos bem organizados
3. Tratamento de erros robusto
4. Configuração de linting (clippy)
5. Testes unitários e de integração
6. Documentação com rustdoc
7. CI/CD com cargo`

	case "cs", "csharp":
		return `\n\nRequisitos específicos para C#:
1. Estrutura de solução .NET moderna
2. Organização em camadas (DDD/Clean Architecture)
3. Configuração de linting
4. Testes com xUnit/NUnit
5. Documentação XML
6. Scripts de build
7. Gerenciamento de dependências com NuGet`

	default:
		return ""
	}
}

// processScaffoldResponse processa a resposta do scaffold para garantir JSON válido
func processScaffoldResponse(response string) (string, error) {
	// Remove blocos de código markdown se presentes
	if strings.HasPrefix(response, "```json\n") && strings.HasSuffix(response, "\n```") {
		response = strings.TrimPrefix(response, "```json\n")
		response = strings.TrimSuffix(response, "\n```")
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
