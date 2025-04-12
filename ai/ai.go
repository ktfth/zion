package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"zion/config"
	"zion/plugins"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type ScaffoldResponse struct {
	Structure struct {
		Directories []string               `json:"directories"`
		Files       map[string]interface{} `json:"files"` // Changed to interface{}
	} `json:"structure"`
}

func callGeminiAPI(prompt string) (string, error) {
	cfg := config.LoadConfig()
	if cfg.GeminiAPIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY não configurada")
	}

	fmt.Println("🔑 Usando chave da API configurada")

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", cfg.GeminiAPIKey)

	request := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
				"role": "user",
			},
		},
		"safetySettings": []map[string]interface{}{
			{
				"category":  "HARM_CATEGORY_DANGEROUS_CONTENT",
				"threshold": "BLOCK_NONE",
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	fmt.Println("📡 Enviando requisição para a API Gemini...")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro na chamada API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	fmt.Println("📥 Resposta recebida da API")

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("erro ao processar resposta: %v\nBody: %s", err, string(body))
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada da API")
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("resposta sem conteúdo")
	}

	fmt.Println("🔍 Processando resposta...")

	// Extrai o JSON da resposta que pode estar dentro de um bloco markdown
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Remove blocos de código markdown se presentes
	if strings.HasPrefix(responseText, "```json\n") && strings.HasSuffix(responseText, "\n```") {
		fmt.Println("📝 Removendo blocos de código markdown")
		responseText = strings.TrimPrefix(responseText, "```json\n")
		responseText = strings.TrimSuffix(responseText, "\n```")
	}

	fmt.Println("🧹 Limpando e corrigindo JSON...")

	// Pré-processamento do JSON para lidar com caracteres especiais em nomes de pacotes
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &jsonMap); err != nil {
		fmt.Printf("⚠️  JSON inválido, tentando limpar: %v\n", err)

		// Se falhar, tenta limpar o JSON
		cleanedResponse := cleanJSONString(responseText)

		fmt.Println("🔄 Tentando parse do JSON limpo...")

		if err := json.Unmarshal([]byte(cleanedResponse), &jsonMap); err != nil {
			return "", fmt.Errorf("resposta não é um JSON válido mesmo após limpeza: %v\nResposta original:\n%s\n\nResposta limpa:\n%s", err, responseText, cleanedResponse)
		}

		fmt.Println("✅ JSON limpo e válido")
		responseText = cleanedResponse
	} else {
		fmt.Println("✅ JSON válido")
	}

	return responseText, nil
}

// cleanJSONString limpa e corrige problemas comuns em strings JSON
func cleanJSONString(input string) string {
	fmt.Println("🧰 Iniciando limpeza do JSON")

	// Remove caracteres invisíveis e espaços em branco extras
	input = strings.TrimSpace(input)
	fmt.Println("✂️  Removidos espaços em branco extras")

	// Pré-processa as dependências
	input = preprocessDependencies(input)
	fmt.Println("📦 Pré-processadas as dependências")

	// Corrige aspas dentro de strings
	input = fixQuotesInJSON(input)
	fmt.Println("🔧 Corrigidas aspas em strings")

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
	var result strings.Builder
	inString := false
	escaped := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		if c == '\\' && !escaped {
			escaped = true
			result.WriteByte(c)
			continue
		}

		if c == '"' && !escaped {
			inString = !inString
		}

		result.WriteByte(c)
		escaped = false
	}

	return result.String()
}

// fixNpmPackageNames corrige problemas com nomes de pacotes npm que contêm @
func fixNpmPackageNames(input string) string {
	// Regex para encontrar nomes de pacotes npm com @
	npmPackageRegex := regexp.MustCompile(`"(@[^"]+)":\s*"([^"]+)"`)

	// Substitui mantendo a estrutura correta
	fixed := npmPackageRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extrai o nome do pacote e a versão
		parts := npmPackageRegex.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		// Reconstrói com escape adequado
		packageName := strings.ReplaceAll(parts[1], "@", "\\u0040")
		return fmt.Sprintf(`"%s": "%s"`, packageName, parts[2])
	})

	// Se houve substituições, registra no log
	if fixed != input {
		fmt.Println("📝 Corrigidos pacotes npm com @")
	}

	return fixed
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
	switch strings.ToLower(language) {
	case "js", "javascript":
		projectDesc += `\n\nRequisitos específicos para JavaScript:
1. Estrutura moderna com ES6+
2. Sistema de módulos ES
3. Configuração de linting (ESLint)
4. Configuração de formatação (Prettier)
5. Scripts NPM úteis
6. Testes unitários configurados
7. Documentação com JSDoc`

	case "ts", "typescript":
		projectDesc += `\n\nRequisitos específicos para TypeScript:
1. Configuração do TSConfig otimizada
2. Tipos bem definidos
3. Estrutura de módulos organizada
4. Configuração de linting (ESLint)
5. Configuração de formatação (Prettier)
6. Scripts NPM úteis
7. Testes unitários com Jest/Vitest
8. Documentação com TSDoc`

	case "go", "golang":
		projectDesc += `\n\nRequisitos específicos para Go:
1. Estrutura de módulos Go
2. Padrões idiomáticos Go
3. Configuração de linting (golangci-lint)
4. Makefile com comandos úteis
5. Testes unitários
6. Documentação no estilo Go
7. Gerenciamento de dependências com go.mod`

	case "rs", "rust":
		projectDesc += `\n\nRequisitos específicos para Rust:
1. Estrutura de workspace Cargo
2. Módulos bem organizados
3. Tratamento de erros robusto
4. Configuração de linting (clippy)
5. Testes unitários e de integração
6. Documentação com rustdoc
7. CI/CD com cargo`

	case "cs", "csharp":
		projectDesc += `\n\nRequisitos específicos para C#:
1. Estrutura de solução .NET moderna
2. Organização em camadas (DDD/Clean Architecture)
3. Configuração de linting
4. Testes com xUnit/NUnit
5. Documentação XML
6. Scripts de build
7. Gerenciamento de dependências com NuGet`
	}

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

	response, err := callGeminiAPI(prompt)
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

// processScaffoldResponse processa a resposta do scaffold para garantir JSON válido
func processScaffoldResponse(response string) (string, error) {
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
		if contentMap, ok := fileContent.(map[string]interface{}); ok {
			if content, exists := contentMap["content"]; exists {
				// Se o conteúdo for um objeto JSON
				if contentObj, isObj := content.(map[string]interface{}); isObj {
					// Processar pacotes npm se for package.json
					if filename == "package.json" {
						if deps, hasDeps := contentObj["dependencies"].(map[string]interface{}); hasDeps {
							processedDeps := make(map[string]interface{})
							for key, value := range deps {
								if strings.HasPrefix(key, "pkg:@") {
									newKey := "@" + strings.TrimPrefix(key, "pkg:@")
									processedDeps[newKey] = value
								} else {
									processedDeps[key] = value
								}
							}
							contentObj["dependencies"] = processedDeps
						}
						if devDeps, hasDevDeps := contentObj["devDependencies"].(map[string]interface{}); hasDevDeps {
							processedDevDeps := make(map[string]interface{})
							for key, value := range devDeps {
								if strings.HasPrefix(key, "pkg:@") {
									newKey := "@" + strings.TrimPrefix(key, "pkg:@")
									processedDevDeps[newKey] = value
								} else {
									processedDevDeps[key] = value
								}
							}
							contentObj["devDependencies"] = processedDevDeps
						}
					}
					baseStruct.Structure.Files[filename] = contentObj
				} else {
					// Se for conteúdo de texto simples
					baseStruct.Structure.Files[filename] = content
				}
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
