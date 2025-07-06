package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ClaudeProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	baseURL     string
}

type claudeRequest struct {
	Model       string          `json:"model"`
	Messages    []claudeMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewClaudeProvider cria um novo provedor Claude (via OpenRouter)
func NewClaudeProvider(config map[string]string) (Provider, error) {
	apiKey, ok := config["api_key"]
	if !ok {
		return nil, fmt.Errorf("api_key não encontrada na configuração do Claude")
	}

	// Valores padrão para Claude
	model := "anthropic/claude-3-haiku"
	maxTokens := 4096
	temperature := 0.7
	baseURL := "https://openrouter.ai/api/v1"

	// Overrides da configuração
	if configModel, ok := config["model"]; ok && configModel != "" {
		model = configModel
	}
	if configMaxTokens, ok := config["max_tokens"]; ok && configMaxTokens != "" {
		if tokens, err := strconv.Atoi(configMaxTokens); err == nil {
			maxTokens = tokens
		}
	}
	if configTemperature, ok := config["temperature"]; ok && configTemperature != "" {
		if temp, err := strconv.ParseFloat(configTemperature, 64); err == nil {
			temperature = temp
		}
	}
	if configBaseURL, ok := config["base_url"]; ok && configBaseURL != "" {
		baseURL = configBaseURL
	}

	return &ClaudeProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		baseURL:     baseURL,
	}, nil
}

func (p *ClaudeProvider) Name() string {
	return "Claude"
}

func (p *ClaudeProvider) GenerateContent(prompt string) (string, error) {
	// Preparar requisição
	req := claudeRequest{
		Model: p.model,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	// Converter para JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar requisição: %v", err)
	}

	// Criar requisição HTTP
	httpReq, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição HTTP: %v", err)
	}

	// Configurar headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/ktfth/zion")
	httpReq.Header.Set("X-Title", "Zion CLI")

	// Enviar requisição
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	// Verificar status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro na API (status %d): %s", resp.StatusCode, string(body))
	}

	// Parsear resposta
	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta: %v", err)
	}

	// Extrair conteúdo
	if len(claudeResp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta recebida")
	}

	return claudeResp.Choices[0].Message.Content, nil
}
