package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RawProjectStructure representa a estrutura do projeto para parsing bruto
type RawProjectStructure struct {
	Structure struct {
		Directories []string          `json:"directories"`
		Files       map[string]string `json:"files"`
	} `json:"structure"`
}

// CreateProjectFromRawJSON cria a estrutura do projeto a partir do JSON bruto
// sem tentar interpretar o conteúdo dos arquivos como JSON
func CreateProjectFromRawJSON(projectName string, jsonStr string) error {
	// Remover blocos de código markdown se presentes
	if strings.HasPrefix(jsonStr, "```json\n") && strings.HasSuffix(jsonStr, "\n```") {
		jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
		jsonStr = strings.TrimSuffix(jsonStr, "\n```")
	}

	// Extrair diretórios usando regex
	dirRegex := regexp.MustCompile(`"directories"\s*:\s*\[\s*((?:"[^"]+"\s*,?\s*)*)\]`)
	dirMatch := dirRegex.FindStringSubmatch(jsonStr)
	
	var directories []string
	if len(dirMatch) > 1 {
		// Extrair cada diretório
		dirItemRegex := regexp.MustCompile(`"([^"]+)"`)
		dirItems := dirItemRegex.FindAllStringSubmatch(dirMatch[1], -1)
		
		for _, item := range dirItems {
			if len(item) > 1 {
				directories = append(directories, item[1])
			}
		}
	}
	
	// Extrair arquivos usando regex
	filesRegex := regexp.MustCompile(`"files"\s*:\s*\{([\s\S]*?)\}(\s*)\}(\s*)$`)
	filesMatch := filesRegex.FindStringSubmatch(jsonStr)
	
	if len(filesMatch) < 2 {
		return fmt.Errorf("não foi possível encontrar a seção 'files' no JSON")
	}
	
	filesContent := filesMatch[1]
	
	// Extrair pares de arquivo/conteúdo
	fileMap := make(map[string]string)
	
	// Posição atual no conteúdo
	pos := 0
	for pos < len(filesContent) {
		// Encontrar o próximo nome de arquivo
		fileNameRegex := regexp.MustCompile(`\s*"([^"]+)"\s*:`)
		fileNameMatch := fileNameRegex.FindStringSubmatchIndex(filesContent[pos:])
		
		if len(fileNameMatch) < 4 {
			break // Não encontrou mais arquivos
		}
		
		// Extrair o nome do arquivo
		fileName := filesContent[pos+fileNameMatch[2]:pos+fileNameMatch[3]]
		
		// Avançar para depois do ":"
		pos += fileNameMatch[1]
		
		// Verificar se o próximo caractere é "
		contentStart := strings.Index(filesContent[pos:], "\"")
		if contentStart < 0 {
			return fmt.Errorf("formato inválido para o conteúdo do arquivo %s", fileName)
		}
		
		pos += contentStart + 1 // Avançar para depois da primeira "
		
		// Encontrar o final do conteúdo (a próxima " não escapada)
		contentEnd := -1
		escaped := false
		for i := pos; i < len(filesContent); i++ {
			if filesContent[i] == '\\' {
				escaped = !escaped
			} else if filesContent[i] == '"' && !escaped {
				contentEnd = i
				break
			} else {
				escaped = false
			}
		}
		
		if contentEnd < 0 {
			return fmt.Errorf("não foi possível encontrar o final do conteúdo para o arquivo %s", fileName)
		}
		
		// Extrair o conteúdo do arquivo
		content := filesContent[pos:contentEnd]
		
		// Decodificar caracteres escapados
		content = strings.ReplaceAll(content, "\\\"", "\"")
		content = strings.ReplaceAll(content, "\\n", "\n")
		content = strings.ReplaceAll(content, "\\t", "\t")
		content = strings.ReplaceAll(content, "\\\\", "\\")
		
		// Adicionar ao mapa
		fileMap[fileName] = content
		
		// Avançar para depois da última "
		pos = contentEnd + 1
		
		// Verificar se há uma vírgula
		commaPos := strings.Index(filesContent[pos:], ",")
		if commaPos >= 0 {
			pos += commaPos + 1
		} else {
			// Se não houver vírgula, deve ser o último arquivo
			break
		}
	}
	
	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}
	
	// Criar os diretórios
	for _, dir := range directories {
		dirPath := filepath.Join(projectName, dir)
		fmt.Println("Criando diretório:", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}
	
	// Criar os arquivos
	for filePath, content := range fileMap {
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
