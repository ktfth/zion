package ai

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CacheKeyGenerator gera chaves de cache para diferentes tipos de requisições
type CacheKeyGenerator struct{}

// GenerateCacheKey gera uma chave única para o cache baseada nos parâmetros
func GenerateCacheKey(language, projectName, description string) string {
	// Normalizar entradas
	normalizedLang := strings.ToLower(strings.TrimSpace(language))
	normalizedName := strings.ToLower(strings.TrimSpace(projectName))
	normalizedDesc := strings.ToLower(strings.TrimSpace(description))

	// Criar string de entrada
	input := fmt.Sprintf("%s:%s:%s", normalizedLang, normalizedName, normalizedDesc)

	// Gerar hash
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("scaffold:%x", hash)
}

// GenerateContextCacheKey gera chave para análise de contexto
func GenerateContextCacheKey(projectPath string) string {
	absPath, _ := filepath.Abs(projectPath)
	hash := sha256.Sum256([]byte(absPath))
	return fmt.Sprintf("context:%x", hash)
}

// DetectProjectTypeFromPath detecta o tipo de projeto baseado no path
func DetectProjectTypeFromPath(projectPath string) string {
	// Lista de arquivos indicadores e seus tipos
	indicators := map[string]string{
		"go.mod":             "Go",
		"package.json":       "JavaScript/TypeScript",
		"requirements.txt":   "Python",
		"pyproject.toml":     "Python",
		"Cargo.toml":         "Rust",
		"pom.xml":            "Java",
		"build.gradle":       "Java",
		"composer.json":      "PHP",
		"Gemfile":            "Ruby",
		"mix.exs":            "Elixir",
		"project.clj":        "Clojure",
		"stack.yaml":         "Haskell",
		"dub.json":           "D",
		"shard.yml":          "Crystal",
		"pubspec.yaml":       "Dart",
		"Package.swift":      "Swift",
		"Makefile":           "C/C++",
		"CMakeLists.txt":     "C/C++",
		"configure.ac":       "C/C++",
		"Dockerfile":         "Docker",
		"docker-compose.yml": "Docker Compose",
		"terraform.tf":       "Terraform",
		"main.tf":            "Terraform",
		"ansible.cfg":        "Ansible",
		"playbook.yml":       "Ansible",
	}

	// Verificar arquivos no diretório
	for file, projectType := range indicators {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return projectType
		}
	}

	// Verificar extensões de arquivos
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignorar erros
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go":
			return fmt.Errorf("Go")
		case ".js", ".ts", ".jsx", ".tsx":
			return fmt.Errorf("JavaScript/TypeScript")
		case ".py":
			return fmt.Errorf("Python")
		case ".rs":
			return fmt.Errorf("Rust")
		case ".java", ".kt", ".scala":
			return fmt.Errorf("Java/JVM")
		case ".php":
			return fmt.Errorf("PHP")
		case ".rb":
			return fmt.Errorf("Ruby")
		case ".ex", ".exs":
			return fmt.Errorf("Elixir")
		case ".clj", ".cljs":
			return fmt.Errorf("Clojure")
		case ".hs":
			return fmt.Errorf("Haskell")
		case ".c", ".cpp", ".cc", ".cxx", ".h", ".hpp":
			return fmt.Errorf("C/C++")
		case ".cs":
			return fmt.Errorf("C#")
		case ".swift":
			return fmt.Errorf("Swift")
		case ".dart":
			return fmt.Errorf("Dart")
		}

		return nil
	})

	if err != nil {
		return err.Error() // Retorna o tipo detectado
	}

	return "Unknown"
}

// GetCacheDirectory retorna o diretório padrão para cache
func GetCacheDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".zion", "cache")
	}
	return filepath.Join(homeDir, ".zion", "cache")
}

// GetLearningDirectory retorna o diretório padrão para aprendizado
func GetLearningDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".zion", "learning")
	}
	return filepath.Join(homeDir, ".zion", "learning")
}

// GetFeedbackDirectory retorna o diretório padrão para feedback
func GetFeedbackDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".zion", "feedback")
	}
	return filepath.Join(homeDir, ".zion", "feedback")
}

// EnsureDirectoryExists cria um diretório se ele não existir
func EnsureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
