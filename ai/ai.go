package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		Directories []string          `json:"directories"`
		Files       map[string]string `json:"files"`
	} `json:"structure"`
}

func callGeminiAPI(prompt string) (string, error) {
	cfg := config.LoadConfig()
	if cfg.GeminiAPIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY não configurada")
	}

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
				"category": "HARM_CATEGORY_DANGEROUS_CONTENT",
				"threshold": "BLOCK_NONE",
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	fmt.Printf("Enviando request para Gemini: %s\n", string(jsonData))

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

	fmt.Printf("Resposta da API: %s\n", string(body))

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

	// Extrai o JSON da resposta que pode estar dentro de um bloco markdown
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	
	// Remove blocos de código markdown se presentes
	if strings.HasPrefix(responseText, "```json\n") && strings.HasSuffix(responseText, "\n```") {
		responseText = strings.TrimPrefix(responseText, "```json\n")
		responseText = strings.TrimSuffix(responseText, "\n```")
	}

	// Verifica se é um JSON válido
	var testJson map[string]interface{}
	if err := json.Unmarshal([]byte(responseText), &testJson); err != nil {
		// Se não for um JSON válido, tentar escapar caracteres especiais
		fmt.Println("Resposta não é um JSON válido. Tentando escapar caracteres especiais...")
		
		// Substituir caracteres problemáticos como '@' por sua versão escapada '\\u0040'
		cleanedResponse := strings.ReplaceAll(responseText, "@", "\\u0040")
		
		// Tentar novamente com a resposta limpa
		if err := json.Unmarshal([]byte(cleanedResponse), &testJson); err != nil {
			// Se ainda não for um JSON válido, retornar erro
			return "", fmt.Errorf("resposta não é um JSON válido: %v\nResposta: %s", err, responseText)
		}
		
		// Se conseguiu processar, usar a resposta limpa
		responseText = cleanedResponse
	}

	return responseText, nil
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
	// Construir a descrição do projeto
	projectDesc := fmt.Sprintf(`Você é um especialista em desenvolvimento de software. Crie uma estrutura de projeto em '%s' chamada '%s'.`, language, projectName)
	
	// Adicionar descrição específica se fornecida
	if description != "" {
		projectDesc += fmt.Sprintf(`\n\nEspecificação: %s`, description)
	}
	
	prompt := fmt.Sprintf(`%s

RETORNE APENAS JSON válido com esta estrutura, sem texto adicional:
{
    "structure": {
        "directories": ["src", "pkg", ...],
        "files": {
            "go.mod": "module hello\n\ngo 1.20\n",
            "main.go": "package main...."
        }
    }
}

Requisitos:
1. Siga padrões da linguagem %s
2. Estrutura modular
3. README.md com instruções incluindo como compilar e instalar
4. Mínimo de dependências

IMPORTANTE: Retorne APENAS o JSON válido, sem explicações ou texto adicional.`, projectDesc, language)

	// Adicionar sistema de plugins apenas se explicitamente solicitado na descrição
	if strings.Contains(strings.ToLower(description), "plugin") {
		prompt += "\n\nImplemente um sistema de plugins básico que permita a integração de funcionalidades externas."
		
		// Mencionar plugins registrados apenas se existirem e forem solicitados
		if len(registeredPlugins) > 0 {
			prompt += "\n\nO sistema deve ser capaz de carregar plugins externos, mas NÃO crie implementações específicas de plugins na estrutura gerada."
		}
	}
	
	// Se for uma aplicação web, adicionar requisitos específicos
	if strings.Contains(strings.ToLower(description), "web") || strings.Contains(strings.ToLower(description), "site") {
		prompt += "\n\nComo esta é uma aplicação web, inclua:\n1. Estrutura para rotas/endpoints\n2. Templates/views\n3. Arquivos estáticos (CSS, JavaScript)\n4. Configuração para servidor web"
	}
	
	// Importante: Não adicionar plugins específicos na estrutura gerada
	prompt += "\n\nIMPORTANTE: NÃO crie diretórios ou arquivos para plugins específicos na estrutura gerada. Mantenha a estrutura limpa e focada apenas no essencial para a aplicação."

	// Atualizar o prompt no contexto
	ctx.Prompt = prompt

	// Executar o hook ModifyPrompt para todos os plugins
	ctx = plugins.ExecuteHook(plugins.ModifyPrompt, ctx)

	// Obter o prompt modificado pelos plugins
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

	fmt.Printf("\nProcessando resposta JSON...\n")
	fmt.Printf("Resposta bruta recebida da API:\n%s\n", response)

	// Extrai o conteúdo JSON da resposta (pode estar dentro de blocos de código markdown)
	responseText := response
	// Verifica se a resposta está em um bloco de código markdown
	if strings.Contains(response, "```json") {
		fmt.Println("Encontrado bloco de código JSON...")
		parts := strings.Split(response, "```json")
		if len(parts) > 1 {
			endParts := strings.Split(parts[1], "```")
			if len(endParts) > 0 {
				responseText = strings.TrimSpace(endParts[0])
				fmt.Println("Conteúdo JSON extraído do bloco de código...")
			}
		}
	} else if strings.Contains(response, "```") {
		fmt.Println("Encontrado bloco de código genérico...")
		parts := strings.Split(response, "```")
		if len(parts) > 1 {
			endParts := strings.Split(parts[1], "```")
			if len(endParts) > 0 {
				responseText = strings.TrimSpace(endParts[0])
				fmt.Println("Conteúdo possivelmente JSON extraído do bloco de código genérico...")
			}
		}
	} else {
		fmt.Println("Nenhum bloco de código encontrado, usando resposta bruta...")
	}
	
	fmt.Printf("Conteúdo JSON extraído:\n%s\n", responseText)

	// Processa a resposta e cria a estrutura do projeto
	var scaffoldResp ScaffoldResponse
	
	// Verifica se o erro é específico para o caractere @
	if strings.Contains(responseText, "@types/") {
		fmt.Println("Detectado pacote npm com @, aplicando correção específica...")
		// Aplica correção específica para package.json com @types
		scaffoldResp = ProcessNpmPackageJson(responseText)
		
		// Verifica se a estrutura foi preenchida corretamente
		if len(scaffoldResp.Structure.Directories) > 0 || len(scaffoldResp.Structure.Files) > 0 {
			fmt.Printf("Estrutura processada com sucesso: %d diretórios, %d arquivos\n",
				len(scaffoldResp.Structure.Directories), len(scaffoldResp.Structure.Files))
			goto structureProcessed
		}
	}
	
	// Tenta fazer o parse do JSON normal
	err = json.Unmarshal([]byte(responseText), &scaffoldResp)
	if err != nil {
		// Se falhar, tenta limpar o JSON para lidar com caracteres especiais
		fmt.Printf("Erro ao fazer parse do JSON: %v\nTentando corrigir o JSON...\n", err)
		
		// Cria uma versão mais limpa da resposta
		cleanedResponse := CleanJSONString(responseText)
		
		// Tenta novamente com a resposta limpa
		err = json.Unmarshal([]byte(cleanedResponse), &scaffoldResp)
		if err != nil {
			// Se ainda falhar, tenta extrair manualmente a estrutura do JSON
			fmt.Printf("Ainda há erro no parsing do JSON: %v\nTentando extrair manualmente...\n", err)
			
			// Tenta uma abordagem alternativa: usar regex para extrair diretórios e arquivos
			fmt.Println("Tentando extrair estrutura usando regex...")
			diretoriosExtraidos, arquivosExtraidos := ExtractProjectStructure(responseText)
			
			if len(diretoriosExtraidos) > 0 || len(arquivosExtraidos) > 0 {
				fmt.Printf("Estrutura extraída via regex: %d diretórios, %d arquivos\n", 
					len(diretoriosExtraidos), len(arquivosExtraidos))
				
				// Preenche a estrutura manualmente
				scaffoldResp.Structure.Directories = diretoriosExtraidos
				scaffoldResp.Structure.Files = arquivosExtraidos
			} else {
				// Se não conseguiu extrair, retorna erro
				fmt.Printf("JSON recebido:\n%s\n", responseText)
				return response, fmt.Errorf("resposta não é um JSON válido e não foi possível extrair a estrutura: %v", err)
			}
		}
	}

	// Ponto para pular quando o processamento já foi feito
structureProcessed:

	fmt.Printf("JSON parseado com sucesso. Diretórios: %d, Arquivos: %d\n", 
		len(scaffoldResp.Structure.Directories), 
		len(scaffoldResp.Structure.Files))

	// Verifica se temos diretórios e arquivos para criar
	if len(scaffoldResp.Structure.Directories) == 0 && len(scaffoldResp.Structure.Files) == 0 {
		return "", fmt.Errorf("resposta não contém diretórios ou arquivos para criar")
	}

	// Cria a estrutura de diretórios e arquivos
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório atual: %v", err)
	}

	// Usar o diretório atual como base e criar o projeto como subdiretório
	baseDir := filepath.Join(wd, projectName)
	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("erro ao obter caminho absoluto: %v", err)
	}
	
	// Exibir informações detalhadas para debug
	fmt.Printf("Diretório atual: %s\n", wd)
	fmt.Printf("Diretório do projeto: %s\n", baseDir)
	
	// Adicionar instruções para compilar e instalar
	fmt.Printf("\nApós a criação do projeto, execute os seguintes comandos:\n")
	fmt.Printf("cd %s\n", projectName)
	fmt.Printf("go mod tidy\n")
	fmt.Printf("go build -o %s\n", projectName)
	fmt.Printf("go install .\n")

	fmt.Printf("\nCriando projeto em: %s\n", baseDir)

	// Verificar se o diretório já existe
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		fmt.Printf("Diretório %s já existe. Removendo...", baseDir)
		os.RemoveAll(baseDir)
	}

	// Criar o diretório base
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório base: %v", err)
	}

	// Verificar se o diretório foi criado
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return "", fmt.Errorf("falha ao criar diretório base %s: %v", baseDir, err)
	} else {
		fmt.Printf("Diretório base criado com sucesso: %s\n", baseDir)
	}

	// Primeiro, cria todos os diretórios
	for _, dir := range scaffoldResp.Structure.Directories {
		fullPath := filepath.Join(baseDir, dir)
		fullPath, err = filepath.Abs(fullPath)
		if err != nil {
			return "", fmt.Errorf("erro ao obter caminho absoluto para %s: %v", dir, err)
		}

		fmt.Printf("Criando diretório: %s\n", fullPath)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", fmt.Errorf("erro ao criar diretório %s: %v", fullPath, err)
		}

		// Verificar se o diretório foi criado
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return "", fmt.Errorf("falha ao criar diretório %s: %v", fullPath, err)
		} else {
			fmt.Printf("Diretório criado com sucesso: %s\n", fullPath)
		}
	}

	// Depois, cria todos os arquivos
	for path, content := range scaffoldResp.Structure.Files {
		fullPath := filepath.Join(baseDir, path)
		fullPath, err = filepath.Abs(fullPath)
		if err != nil {
			return "", fmt.Errorf("erro ao obter caminho absoluto para %s: %v", path, err)
		}

		dir := filepath.Dir(fullPath)
		
		// Garante que o diretório existe (pode ser um que não estava na lista)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}

		fmt.Printf("Criando arquivo: %s\n", fullPath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("erro ao criar arquivo %s: %v", fullPath, err)
		}

		// Verificar se o arquivo foi criado
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return "", fmt.Errorf("falha ao criar arquivo %s: %v", fullPath, err)
		} else {
			fmt.Printf("Arquivo criado com sucesso: %s\n", fullPath)
		}
	}

	fmt.Printf("\nProjeto '%s' criado com sucesso!\n", projectName)

	return response, nil
}

