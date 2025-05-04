package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GPTProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
}

type gptRequest struct {
	Model       string       `json:"model"`
	Messages    []gptMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
}

type gptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type gptResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewGPTProvider cria um novo provedor OpenAI GPT
func NewGPTProvider(config map[string]string) (Provider, error) {
	apiKey, ok := config["api_key"]
	if !ok {
		return nil, fmt.Errorf("api_key não encontrada na configuração do GPT")
	}

	// Valores padrão
	model := "gpt-3.5-turbo"
	maxTokens := 2048
	temperature := 0.7

	// Sobrescrever com valores da configuração se fornecidos
	if m, ok := config["model"]; ok {
		model = m
	}
	if t, ok := config["max_tokens"]; ok {
		fmt.Sscanf(t, "%d", &maxTokens)
	}
	if temp, ok := config["temperature"]; ok {
		fmt.Sscanf(temp, "%f", &temperature)
	}

	return &GPTProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}, nil
}

func (p *GPTProvider) Name() string {
	return "GPT"
}

func (p *GPTProvider) GenerateContent(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	request := gptRequest{
		Model: p.model,
		Messages: []gptMessage{
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

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

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

	var gptResp gptResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return "", fmt.Errorf("erro ao processar resposta: %v\nBody: %s", err, string(body))
	}

	if len(gptResp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada da API")
	}

	return gptResp.Choices[0].Message.Content, nil
}
