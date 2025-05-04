package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey string
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
	Role  string       `json:"role"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// NewGeminiProvider cria um novo provedor Gemini
func NewGeminiProvider(config map[string]string) (Provider, error) {
	apiKey, ok := config["api_key"]
	if !ok {
		return nil, fmt.Errorf("api_key não encontrada na configuração do Gemini")
	}
	return &GeminiProvider{apiKey: apiKey}, nil
}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

func (p *GeminiProvider) GenerateContent(prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", p.apiKey)

	request := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
				Role: "user",
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("erro ao processar resposta: %v\nBody: %s", err, string(body))
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada da API")
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("resposta sem conteúdo")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
