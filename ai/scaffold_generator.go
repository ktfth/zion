package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateProjectStructure cria a estrutura de diretórios e arquivos com base na resposta JSON
func CreateProjectStructure(projectName, jsonResponse string) error {
	fmt.Printf("\nIniciando criação da estrutura do projeto '%s'...\n", projectName)

	// Extrair o conteúdo JSON da resposta
	jsonContent := extractJSONContent(jsonResponse)
	fmt.Printf("Conteúdo JSON extraído:\n%s\n", jsonContent)

	// Decodificar o JSON
	var scaffoldResponse ScaffoldResponse
	err := json.Unmarshal([]byte(jsonContent), &scaffoldResponse)
	if err != nil {
		fmt.Printf("Erro ao decodificar JSON: %v\nTentando processamento alternativo...\n", err)

		// Se falhar, tenta usar o processador específico para package.json com @types
		if strings.Contains(jsonContent, "@types/") {
			fmt.Println("Detectado pacote npm com @, aplicando correção específica...")
			scaffoldResponse = ProcessNpmPackageJson(jsonContent)
		} else {
			// Tenta extrair usando regex
			fmt.Println("Tentando extrair estrutura usando regex...")
			diretoriosExtraidos, arquivosExtraidos := ExtractProjectStructure(jsonContent) // Use jsonContent aqui

			if len(diretoriosExtraidos) > 0 || len(arquivosExtraidos) > 0 {
				fmt.Printf("Estrutura extraída via regex: %d diretórios, %d arquivos\n",
					len(diretoriosExtraidos), len(arquivosExtraidos))

				// Preenche a estrutura manualmente
				scaffoldResponse.Structure.Directories = diretoriosExtraidos
				scaffoldResponse.Structure.Files = make(map[string]interface{}) // Criar o mapa correto
				for k, v := range arquivosExtraidos {
					scaffoldResponse.Structure.Files[k] = v
				}
			} else {
				return fmt.Errorf("não foi possível extrair a estrutura do projeto: %v", err)
			}
		}
	}

	// Criar o diretório raiz do projeto
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	fmt.Printf("Criando estrutura do projeto em: %s\n", projectDir)

	// Criar diretórios
	for _, dir := range scaffoldResponse.Structure.Directories {
		dirPath := filepath.Join(projectDir, dir)
		fmt.Printf("Criando diretório: %s\n", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}

	// Criar arquivos
	for filePath, content := range scaffoldResponse.Structure.Files {
		fullPath := filepath.Join(projectDir, filePath)

		// Garantir que o diretório pai exista
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório pai para %s: %v", filePath, err)
		}

		fmt.Printf("Criando arquivo: %s\n", filePath)
		var contentBytes []byte
		if strContent, ok := content.(string); ok {
			contentBytes = []byte(strContent)
		} else {
			contentBytes, err = json.MarshalIndent(content, "", "   ") // Alteração aqui para três espaços
			if err != nil {
				return fmt.Errorf("erro ao serializar conteúdo para %s: %v", filePath, err)
			}
		}
		if err := os.WriteFile(fullPath, contentBytes, 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
		}
	}

	fmt.Printf("\nEstrutura do projeto criada com sucesso em: %s\n", projectDir)
	return nil
}

// extractJSONContent extrai o conteúdo JSON da resposta, mesmo que esteja dentro de blocos de código markdown
func extractJSONContent(response string) string {
	// Verificar se a resposta está em um bloco de código markdown
	if strings.Contains(response, "```json") {
		parts := strings.Split(response, "```json")
		if len(parts) > 1 {
			endParts := strings.Split(parts[1], "```")
			if len(endParts) > 0 {
				return strings.TrimSpace(endParts[0])
			}
		}
	} else if strings.Contains(response, "```") {
		parts := strings.Split(response, "```")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	// Se não encontrar blocos de código, assume que a resposta já é JSON
	return response
}