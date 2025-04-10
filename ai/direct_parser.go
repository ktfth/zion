package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DirectParseAndCreateProject é uma função que faz o parsing direto da resposta da API
// e cria a estrutura do projeto sem depender do parsing JSON padrão
func DirectParseAndCreateProject(projectName, response string) error {
	fmt.Println("Usando parser direto para criar a estrutura do projeto...")
	
	// Verificar se a resposta contém package.json com @types
	if strings.Contains(response, "@types/") {
		fmt.Println("Detectado package.json com @types, usando parser especializado...")
		// Usar o parser especializado para package.json
		return ParseAndCreateTypeScriptProject(projectName, response)
	}
	
	// Extrair diretórios e arquivos diretamente da resposta
	directories, files := ExtractProjectStructure(response)
	
	if len(directories) == 0 && len(files) == 0 {
		return fmt.Errorf("não foi possível extrair diretórios ou arquivos da resposta")
	}
	
	fmt.Printf("Estrutura extraída: %d diretórios, %d arquivos\n", 
		len(directories), len(files))
	
	// Criar o diretório raiz do projeto
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}
	
	fmt.Printf("Criando estrutura do projeto em: %s\n", projectDir)
	
	// Criar diretórios
	for _, dir := range directories {
		dirPath := filepath.Join(projectDir, dir)
		fmt.Printf("Criando diretório: %s\n", dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
	}
	
	// Criar arquivos
	for filePath, content := range files {
		fullPath := filepath.Join(projectDir, filePath)
		
		// Garantir que o diretório pai exista
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório pai para %s: %v", filePath, err)
		}
		
		fmt.Printf("Criando arquivo: %s (%d bytes)\n", filePath, len(content))
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
		}
	}
	
	fmt.Printf("\nEstrutura do projeto criada com sucesso em: %s\n", projectDir)
	return nil
}

// ExtractDirectoriesAndFiles extrai diretórios e arquivos diretamente da resposta da API
// sem depender do parsing JSON padrão
func ExtractDirectoriesAndFiles(response string) ([]string, map[string]string) {
	// Remover blocos de código markdown
	cleanResponse := response
	if strings.Contains(response, "```") {
		// Extrair apenas o conteúdo entre ```json e ```
		jsonBlockRegex := regexp.MustCompile("```json\\s*\\n([\\s\\S]*?)\\n\\s*```")
		matches := jsonBlockRegex.FindStringSubmatch(response)
		if len(matches) > 1 {
			cleanResponse = matches[1]
		} else {
			// Se não encontrar ```json, tenta encontrar qualquer bloco de código
			blockRegex := regexp.MustCompile("```\\s*\\n([\\s\\S]*?)\\n\\s*```")
			matches := blockRegex.FindStringSubmatch(response)
			if len(matches) > 1 {
				cleanResponse = matches[1]
			}
		}
	}
	
	// Extrair diretórios
	directories := []string{}
	dirRegex := regexp.MustCompile(`"directories":\s*\[\s*([\s\S]*?)\s*\]`)
	dirMatches := dirRegex.FindStringSubmatch(cleanResponse)
	if len(dirMatches) > 1 {
		dirList := dirMatches[1]
		// Extrair cada diretório
		dirItemRegex := regexp.MustCompile(`"([^"]+)"`)
		dirItemMatches := dirItemRegex.FindAllStringSubmatch(dirList, -1)
		for _, match := range dirItemMatches {
			if len(match) > 1 {
				directories = append(directories, match[1])
			}
		}
	}
	
	// Extrair arquivos
	files := make(map[string]string)
	fileBlockRegex := regexp.MustCompile(`"files":\s*\{\s*([\s\S]*?)\s*\}[\s\n]*\}`)
	fileBlockMatches := fileBlockRegex.FindStringSubmatch(cleanResponse)
	if len(fileBlockMatches) > 1 {
		fileBlock := fileBlockMatches[1]
		
		// Dividir o bloco de arquivos em linhas
		lines := strings.Split(fileBlock, "\n")
		
		// Variáveis para controlar o parsing
		currentFile := ""
		currentContent := ""
		inFileContent := false
		
		// Processar cada linha
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			
			// Se não estamos dentro do conteúdo de um arquivo
			if !inFileContent {
				// Procurar por um novo arquivo
				fileStartRegex := regexp.MustCompile(`"([^"]+)":\s*"(.*)`)
				fileStartMatches := fileStartRegex.FindStringSubmatch(trimmedLine)
				if len(fileStartMatches) > 2 {
					currentFile = fileStartMatches[1]
					currentContent = fileStartMatches[2]
					inFileContent = true
					
					// Verificar se o conteúdo termina na mesma linha
					if strings.HasSuffix(trimmedLine, "\",") || strings.HasSuffix(trimmedLine, "\"") {
						// Remover a vírgula ou aspas finais
						if strings.HasSuffix(trimmedLine, "\",") {
							currentContent = strings.TrimSuffix(currentContent, "\",")
						} else {
							currentContent = strings.TrimSuffix(currentContent, "\"")
						}
						
						// Desescapar o conteúdo
						currentContent = strings.ReplaceAll(currentContent, "\\\"", "\"")
						currentContent = strings.ReplaceAll(currentContent, "\\n", "\n")
						currentContent = strings.ReplaceAll(currentContent, "\\t", "\t")
						
						// Adicionar o arquivo ao mapa
						files[currentFile] = currentContent
						
						// Resetar variáveis
						currentFile = ""
						currentContent = ""
						inFileContent = false
					}
				}
			} else {
				// Estamos dentro do conteúdo de um arquivo
				// Verificar se a linha termina o conteúdo
				if strings.HasSuffix(trimmedLine, "\",") || strings.HasSuffix(trimmedLine, "\"") {
					// Adicionar a linha ao conteúdo
					if strings.HasSuffix(trimmedLine, "\",") {
						currentContent += strings.TrimSuffix(trimmedLine, "\",")
					} else {
						currentContent += strings.TrimSuffix(trimmedLine, "\"")
					}
					
					// Desescapar o conteúdo
					currentContent = strings.ReplaceAll(currentContent, "\\\"", "\"")
					currentContent = strings.ReplaceAll(currentContent, "\\n", "\n")
					currentContent = strings.ReplaceAll(currentContent, "\\t", "\t")
					
					// Adicionar o arquivo ao mapa
					files[currentFile] = currentContent
					
					// Resetar variáveis
					currentFile = ""
					currentContent = ""
					inFileContent = false
				} else {
					// Adicionar a linha ao conteúdo
					currentContent += trimmedLine
				}
			}
		}
	}
	
	return directories, files
}
