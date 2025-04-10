package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GeminiAPIKey string
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

	return &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		HomeDir:      zionDir,
		PluginsDir:   pluginsDir,
	}
}
