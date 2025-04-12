package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExtractAndCreateProject extrai diretamente os diretórios e arquivos do JSON
// e cria a estrutura do projeto sem tentar interpretar o conteúdo
func ExtractAndCreateProject(projectName string, jsonStr string) error {
	// Remover blocos de código markdown se presentes
	if strings.HasPrefix(jsonStr, "```json\n") && strings.HasSuffix(jsonStr, "\n```") {
		jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
		jsonStr = strings.TrimSuffix(jsonStr, "\n```")
	}

	var scaffoldResp ScaffoldResponse
	err := json.Unmarshal([]byte(jsonStr), &scaffoldResp)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório raiz '%s': %v", projectName, err)
	}
	fmt.Printf("\n📁 Criando diretório raiz: %s\n", projectName)

	// Criar diretórios
	if len(scaffoldResp.Structure.Directories) > 0 {
		fmt.Println("\n📂 Criando diretórios:")
		for _, dir := range scaffoldResp.Structure.Directories {
			dirPath := filepath.Join(projectName, dir)
			fmt.Printf("   ├── %s\n", dir)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório '%s': %v", dir, err)
			}
		}
	}

	// Criar arquivos
	if len(scaffoldResp.Structure.Files) > 0 {
		fmt.Println("\n📄 Criando arquivos:")
		for filePath, content := range scaffoldResp.Structure.Files {
			var contentStr string
			if strContent, ok := content.(string); ok {
				contentStr = strContent
				contentStr = ProcessEscapedChars(contentStr)
			} else {
				contentBytes, err := json.MarshalIndent(content, "", "  ")
				if err != nil {
					return fmt.Errorf("erro ao serializar conteúdo JSON para '%s': %v", filePath, err)
				}
				contentStr = string(contentBytes)
			}

			fullPath := filepath.Join(projectName, filePath)
			fmt.Printf("   ├── %s\n", filePath)

			// Garantir que o diretório pai exista
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório pai para '%s': %v", filePath, err)
			}

			if err := os.WriteFile(fullPath, []byte(contentStr), 0644); err != nil {
				return fmt.Errorf("erro ao criar arquivo '%s': %v", filePath, err)
			}
		}
	}

	// Exibir resumo
	fmt.Printf("\n📊 Resumo da estrutura criada:\n")
	fmt.Printf("   ├── %d diretórios\n", len(scaffoldResp.Structure.Directories))
	fmt.Printf("   └── %d arquivos\n", len(scaffoldResp.Structure.Files))

	return nil
}

// As funções auxiliares foram movidas para o arquivo file_utils.go
