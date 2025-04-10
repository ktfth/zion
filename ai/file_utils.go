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

// ProcessEscapedChars processa caracteres escapados no conteúdo
func ProcessEscapedChars(content string) string {
	// Substituir caracteres escapados
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

// CreateFile cria um arquivo com o conteúdo especificado
func CreateFile(projectName, filePath, content string) error {
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
