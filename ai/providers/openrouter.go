package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenRouterProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	baseURL     string
}

type openRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []openRouterMessage `json:"messages"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenRouterProvider cria um novo provedor OpenRouter
func NewOpenRouterProvider(config map[string]string) (Provider, error) {
	apiKey, ok := config["api_key"]
	if !ok {
		return nil, fmt.Errorf("api_key não encontrada na configuração do OpenRouter")
	}

	// Valores padrão
	model := "microsoft/wizardlm-2-8x22b"
	maxTokens := 2048
	temperature := 0.7
	baseURL := "https://openrouter.ai/api/v1"

	// Sobrescrever com valores da configuração se fornecidos
	if m, ok := config["model"]; ok && m != "" {
		model = m
	}
	if t, ok := config["max_tokens"]; ok && t != "" {
		fmt.Sscanf(t, "%d", &maxTokens)
	}
	if temp, ok := config["temperature"]; ok && temp != "" {
		fmt.Sscanf(temp, "%f", &temperature)
	}
	if url, ok := config["base_url"]; ok && url != "" {
		baseURL = url
	}

	return &OpenRouterProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		baseURL:     baseURL,
	}, nil
}

func (p *OpenRouterProvider) Name() string {
	return "OpenRouter"
}

func (p *OpenRouterProvider) GenerateContent(prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)

	request := openRouterRequest{
		Model: p.model,
		Messages: []openRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("HTTP-Referer", "https://github.com/zion-ai")
	req.Header.Set("X-Title", "Zion AI")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na chamada API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp openRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("erro ao processar resposta: %v\nBody: %s", err, string(body))
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada da API")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}
