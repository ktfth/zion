package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		Directories []string                 `json:"directories"`
		Files       map[string]interface{} `json:"files"` // Changed to interface{}
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

	// Substituir caracteres problemáticos como '@' por sua versão escapada '\\u0040'
	cleanedResponse := strings.ReplaceAll(responseText, "@", "\\u0040")

	// Verifica se é um JSON válido
	var testJson map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedResponse), &testJson); err != nil {
		// Se não for um JSON válido, tentar escapar caracteres especiais
		fmt.Println("Resposta não é um JSON válido. Tentando escapar caracteres especiais...")
		
		// Tentar novamente com a resposta limpa
		if err := json.Unmarshal([]byte(responseText), &testJson); err != nil {
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

IMPORTANTE:
1. Retorne APENAS o JSON válido, sem explicações ou texto adicional.
2. Para arquivos JSON (como package.json), NÃO adicione escapes ou quebras de linha no conteúdo. Mantenha o conteúdo exatamente como seria em um arquivo real.
3. NÃO escape caracteres especiais como @ em nomes de pacotes npm.
4. Mantenha o formato original dos arquivos, sem adicionar escapes desnecessários.`, projectDesc, language)

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

	fmt.Printf("\nResposta final da API (após plugins):\n%s\n", response)

	err = CreateProjectStructure(projectName, response)
	if err != nil {
		return "", fmt.Errorf("erro ao criar estrutura do projeto: %v", err)
	}

	return response, nil
}

