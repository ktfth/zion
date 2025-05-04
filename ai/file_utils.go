package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindUnescapedQuote encontra a próxima aspas não escapada
func FindUnescapedQuote(str string, startPos int) int {
	escaped := false
	for i := startPos; i < len(str); i++ {
		if str[i] == '\\' {
			escaped = !escaped
		} else if str[i] == '"' && !escaped {
			return i
		} else {
			escaped = false
		}
	}
	return -1
}

// FindMatchingBrace encontra a chave de fechamento correspondente
func FindMatchingBrace(str string, startPos int) int {
	count := 1
	for i := startPos; i < len(str); i++ {
		if str[i] == '{' {
			count++
		} else if str[i] == '}' {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}

// ProcessEscapedChars processes escaped characters in content
func ProcessEscapedChars(content string) string {
	// Replace escaped characters
	replacements := map[string]string{
		"\\\"": "\"",
		"\\n":  "\n",
		"\\t":  "\t",
		"\\r":  "\r",
		"\\\\": "\\",
	}

	result := content
	for escaped, unescaped := range replacements {
		result = strings.ReplaceAll(result, escaped, unescaped)
	}

	return result
}

// ProcessPackageJsonContent processa o conteúdo do package.json recursivamente
func ProcessPackageJsonContent(contentObj map[string]interface{}) {
	processPackageJsonContentWithDepth(contentObj, 0)
}

// processPackageJsonContentWithDepth handles the recursive processing with depth tracking
func processPackageJsonContentWithDepth(contentObj map[string]interface{}, depth int) {
	// Prevent infinite recursion
	if depth > 10 {
		return
	}

	// Garantir campos obrigatórios apenas no nível raiz
	if depth == 0 {
		if _, ok := contentObj["name"]; !ok {
			contentObj["name"] = "unnamed-project"
		}
		if _, ok := contentObj["version"]; !ok {
			contentObj["version"] = "1.0.0"
		}
		if _, ok := contentObj["scripts"]; !ok {
			contentObj["scripts"] = map[string]interface{}{
				"start": "node dist/index.js",
				"build": "tsc",
				"dev":   "nodemon src/index.ts",
			}
		}
	}

	// Process dependencies and devDependencies only at root level
	if depth == 0 {
		for _, field := range []string{"dependencies", "devDependencies"} {
			if deps, ok := contentObj[field].(map[string]interface{}); ok {
				processedDeps := make(map[string]interface{})
				for key, value := range deps {
					// Remove pkg: prefix from any key
					newKey := strings.TrimPrefix(key, "pkg:")
					// Ensure versions are strings
					if ver, ok := value.(string); ok {
						processedDeps[newKey] = ver
					} else {
						processedDeps[newKey] = "latest"
					}
				}
				contentObj[field] = processedDeps
			} else {
				// If field doesn't exist or isn't a map, create an empty one
				contentObj[field] = make(map[string]interface{})
			}
		}

		// Ensure minimum TypeScript dependencies
		if _, hasTS := contentObj["devDependencies"].(map[string]interface{})["typescript"]; hasTS {
			devDeps := contentObj["devDependencies"].(map[string]interface{})
			if _, ok := devDeps["@types/node"]; !ok {
				devDeps["@types/node"] = "^20.4.5"
			}
			if _, ok := devDeps["nodemon"]; !ok {
				devDeps["nodemon"] = "^3.0.1"
			}
		}
	}

	// Process other fields that might have pkg: prefix
	for key, value := range contentObj {
		if strings.HasPrefix(key, "pkg:") {
			newKey := strings.TrimPrefix(key, "pkg:")
			contentObj[newKey] = value
			delete(contentObj, key)
		}

		// Recursively process nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok && key != "dependencies" && key != "devDependencies" {
			processPackageJsonContentWithDepth(nestedMap, depth+1)
		}
	}
}

// CreateFile cria um arquivo com o conteúdo especificado
func CreateFile(projectName, filePath, content string) error {
	// Skip empty files
	if strings.TrimSpace(content) == "" {
		fmt.Printf("⚠️  Pulando arquivo vazio: %s\n", filePath)
		return nil
	}

	fullPath := filepath.Join(projectName, filePath)

	// Garantir que o diretório pai exista
	parentDir := filepath.Dir(fullPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório pai para %s: %v", filePath, err)
	}

	fmt.Println("Criando arquivo:", filePath)

	// Escrever o conteúdo no arquivo
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
	}

	return nil
}

// IsJSONFile verifica se o arquivo é um arquivo JSON ou similar que não deve ter escapes adicionados
func IsJSONFile(filePath string) bool {
	jsonFiles := []string{
		"package.json",
		"tsconfig.json",
		"angular.json",
		"next.config.js",
		"webpack.config.js",
		".eslintrc.json",
		"composer.json",
		"manifest.json",
		"app.json",
		"project.json",
	}

	fileName := filepath.Base(filePath)

	for _, jsonFile := range jsonFiles {
		if fileName == jsonFile {
			return true
		}
	}

	// Verificar pela extensão
	ext := filepath.Ext(filePath)
	return ext == ".json" || ext == ".jsonc"
}

// PreserveJSONFormat preserva o formato original de arquivos JSON, removendo apenas escapes desnecessários
func PreserveJSONFormat(content string) string {
	// Remover escapes desnecessários que podem afetar o formato JSON
	unescaped := content

	// Remover escapes de caracteres especiais que não devem ser escapados em JSON
	unescaped = strings.ReplaceAll(unescaped, "\\@", "@")
	unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")

	// Preservar quebras de linha reais, mas remover escapes de quebra de linha
	unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
	unescaped = strings.ReplaceAll(unescaped, "\\r", "\r")
	unescaped = strings.ReplaceAll(unescaped, "\\t", "\t")

	return unescaped
}
