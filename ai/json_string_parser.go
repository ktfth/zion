package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectStructure representa a estrutura do projeto
type ProjectStructure struct {
	Structure struct {
		Directories []string            `json:"directories"`
		Files       map[string]string   `json:"files"`
	} `json:"structure"`
}

// ParseJSONAsStringContent analisa o JSON tratando o conteúdo dos arquivos como strings simples
// sem tentar interpretar o conteúdo como JSON
func ParseJSONAsStringContent(jsonStr string) (*ProjectStructure, error) {
	// Primeiro, vamos verificar se o JSON já é válido como está
	var projectStructure ProjectStructure
	err := json.Unmarshal([]byte(jsonStr), &projectStructure)
	if err == nil {
		return &projectStructure, nil
	}

	// Se não for válido, vamos tentar um parsing personalizado
	fmt.Println("JSON não é válido como está. Tentando parsing personalizado...")

	// Usar regex para extrair o conteúdo dos arquivos como strings brutas
	// Primeiro, vamos encontrar a seção "files" no JSON
	filesPattern := regexp.MustCompile(`"files"\s*:\s*\{([^}]*)\}`)
	filesMatch := filesPattern.FindStringSubmatch(jsonStr)
	
	if len(filesMatch) < 2 {
		return nil, fmt.Errorf("não foi possível encontrar a seção 'files' no JSON")
	}
	
	filesContent := filesMatch[1]
	
	// Agora vamos substituir o conteúdo dos arquivos por placeholders
	fileContentPattern := regexp.MustCompile(`"([^"]+)"\s*:\s*"((?:\\.|[^"\\])*)"|"([^"]+)"\s*:\s*\{([^}]*)\}`)
	
	// Mapa para armazenar os conteúdos originais
	fileContents := make(map[string]string)
	
	// Contador para gerar placeholders únicos
	counter := 0
	
	// Substituir conteúdos por placeholders
	modifiedFilesContent := fileContentPattern.ReplaceAllStringFunc(filesContent, func(match string) string {
		// Extrair o nome do arquivo e o conteúdo
		parts := fileContentPattern.FindStringSubmatch(match)
		
		var fileName, content string
		if parts[1] != "" {
			// Caso simples: "arquivo": "conteúdo"
			fileName = parts[1]
			content = parts[2]
		} else if parts[3] != "" {
			// Caso complexo: "arquivo": {...}
			fileName = parts[3]
			content = "{" + parts[4] + "}"
		} else {
			// Não deveria chegar aqui
			return match
		}
		
		// Gerar um placeholder único
		placeholder := fmt.Sprintf("__PLACEHOLDER_%d__", counter)
		counter++
		
		// Armazenar o conteúdo original
		fileContents[placeholder] = content
		
		// Retornar a string com o placeholder
		return fmt.Sprintf("\"%s\": \"%s\"", fileName, placeholder)
	})
	
	// Reconstruir o JSON com os placeholders
	modifiedJSON := strings.Replace(jsonStr, filesMatch[1], modifiedFilesContent, 1)
	
	// Agora o JSON deve ser válido
	err = json.Unmarshal([]byte(modifiedJSON), &projectStructure)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar JSON modificado: %v", err)
	}
	
	// Substituir os placeholders pelos conteúdos originais
	for fileName, content := range projectStructure.Structure.Files {
		for placeholder, originalContent := range fileContents {
			if content == placeholder {
				// Substituir o placeholder pelo conteúdo original
				projectStructure.Structure.Files[fileName] = originalContent
				break
			}
		}
	}
	
	return &projectStructure, nil
}

// CreateProjectFromStringParsedJSON cria a estrutura do projeto a partir do JSON analisado
func CreateProjectFromStringParsedJSON(projectName string, jsonStr string) error {
	// Analisar o JSON tratando o conteúdo dos arquivos como strings
	projectStructure, err := ParseJSONAsStringContent(jsonStr)
	if err != nil {
		return fmt.Errorf("erro ao analisar JSON: %v", err)
	}
	
	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}
	
	// Criar os diretórios
	for _, dir := range projectStructure.Structure.Directories {
		dirPath := filepath.Join(projectName, dir)
		fmt.Println("Criando diretório:", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}
	
	// Criar os arquivos
	for filePath, content := range projectStructure.Structure.Files {
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
	}
	
	fmt.Println("\nEstrutura do projeto criada com sucesso em:", projectName)
	return nil
}
