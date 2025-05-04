package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GeminiAPIKey string
	OpenAIAPIKey string
	AIProvider   string
	HomeDir      string
	PluginsDir   string
}

func LoadConfig() *Config {
	// Obter o diretório home do usuário
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// Definir o diretório .zion dentro do home
	zionDir := filepath.Join(homeDir, ".zion")

	// Garantir que o diretório .zion existe
	if _, err := os.Stat(zionDir); os.IsNotExist(err) {
		os.MkdirAll(zionDir, 0755)
	}

	// Definir o diretório de plugins dentro de .zion
	pluginsDir := filepath.Join(zionDir, "plugins")

	// Garantir que o diretório de plugins existe
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		os.MkdirAll(pluginsDir, 0755)
	}

	// Determinar o provedor de IA baseado na variável de ambiente ou usar o padrão
	aiProvider := os.Getenv("ZION_AI_PROVIDER")
	if aiProvider == "" {
		// Se GEMINI_API_KEY estiver definida, usa Gemini como padrão
		if os.Getenv("GEMINI_API_KEY") != "" {
			aiProvider = "gemini"
		} else if os.Getenv("OPENAI_API_KEY") != "" {
			// Se OPENAI_API_KEY estiver definida, usa GPT como padrão
			aiProvider = "gpt"
		} else {
			// Se nenhuma chave estiver definida, usa Gemini como padrão
			aiProvider = "gemini"
		}
	}

	return &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		AIProvider:   aiProvider,
		HomeDir:      zionDir,
		PluginsDir:   pluginsDir,
	}
}

// GetAIConfig retorna a configuração específica para o provedor de IA selecionado
func (c *Config) GetAIConfig() map[string]string {
	switch c.AIProvider {
	case "gpt":
		return map[string]string{
			"api_key":    c.OpenAIAPIKey,
			"model":      os.Getenv("OPENAI_MODEL"),      // opcional
			"max_tokens": os.Getenv("OPENAI_MAX_TOKENS"), // opcional
		}
	default: // gemini
		return map[string]string{
			"api_key": c.GeminiAPIKey,
		}
	}
}
