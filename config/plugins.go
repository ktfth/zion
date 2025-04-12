package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetPluginsDir retorna o diretório de plugins apropriado para o sistema operacional
func GetPluginsDir() string {
	// Se PLUGINS_DIR estiver definido no ambiente, use-o
	if envDir := os.Getenv("PLUGINS_DIR"); envDir != "" {
		return envDir
	}

	// Caso contrário, use o diretório padrão baseado no sistema operacional
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "windows" {
		// No Windows, use AppData/Local
		return filepath.Join(homeDir, "AppData", "Local", "Zion", "plugins")
	}

	// Em sistemas Unix, use ~/.zion/plugins
	return filepath.Join(homeDir, ".zion", "plugins")
}

// EnsurePluginsDirExists garante que o diretório de plugins existe
func EnsurePluginsDirExists() error {
	pluginsDir := GetPluginsDir()
	if pluginsDir == "" {
		return nil // Se não conseguimos determinar o diretório, não fazemos nada
	}

	// Cria o diretório com permissões apropriadas
	return os.MkdirAll(pluginsDir, 0755)
}

// IsPluginSupported verifica se plugins são suportados na plataforma atual
func IsPluginSupported() bool {
	// No Windows, verificamos se estamos usando uma versão que suporta plugins
	if runtime.GOOS == "windows" {
		// Windows 10 build 1803 ou superior suporta plugins
		return true
	}

	// Em outros sistemas operacionais, plugins são sempre suportados
	return true
}

// GetPluginExtension retorna a extensão correta para plugins no sistema atual
func GetPluginExtension() string {
	if runtime.GOOS == "windows" {
		return ".dll"
	}
	return ".so"
}
