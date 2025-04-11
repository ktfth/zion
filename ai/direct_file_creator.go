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
		return fmt.Errorf("erro ao fazer parse do JSON para ScaffoldResponse: %v", err)
	}

	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	// Criar diretórios
	for _, dir := range scaffoldResp.Structure.Directories {
		dirPath := filepath.Join(projectName, dir)
		fmt.Println("Criando diretório:", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}

	// Criar arquivos
	for filePath, content := range scaffoldResp.Structure.Files {
		var contentStr string
		if strContent, ok := content.(string); ok {
			contentStr = strContent
			contentStr = ProcessEscapedChars(contentStr)
		} else {
			contentBytes, err := json.MarshalIndent(content, "", "  ")
			if err != nil {
				return fmt.Errorf("erro ao serializar conteúdo JSON para %s: %v", filePath, err)
			}
			contentStr = string(contentBytes)
		}
		err := CreateFile(projectName, filePath, contentStr)
		if err != nil {
			return err
		}
	}

	fmt.Println("\nEstrutura do projeto criada com sucesso em:", projectName)
	return nil
}

// As funções auxiliares foram movidas para o arquivo file_utils.go
