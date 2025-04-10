package ai

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractProjectStructure extrai diretórios e arquivos da resposta da API usando regex
// Isso é usado como último recurso quando o parsing do JSON falha
func ExtractProjectStructure(response string) ([]string, map[string]string) {
	directories := []string{}
	files := map[string]string{}

	// Extrai diretórios
	fmt.Println("Extraindo diretórios...")
	dirRegex := regexp.MustCompile(`"directories"\s*:\s*\[(.*?)\]`)
	dirMatches := dirRegex.FindStringSubmatch(response)
	
	if len(dirMatches) > 1 {
		dirList := dirMatches[1]
		// Remove espaços em branco e aspas extras
		dirList = strings.TrimSpace(dirList)
		// Divide a lista por vírgulas
		dirItems := strings.Split(dirList, ",")
		
		for _, item := range dirItems {
			// Remove aspas e espaços
			dir := strings.Trim(strings.TrimSpace(item), "\"'")
			if dir != "" {
				directories = append(directories, dir)
				fmt.Printf("Diretório encontrado: %s\n", dir)
			}
		}
	}

	// Extrai arquivos
	fmt.Println("Extraindo arquivos...")
	fileRegex := regexp.MustCompile(`"files"\s*:\s*\{(.*?)\}`)
	fileMatches := fileRegex.FindStringSubmatch(response)
	
	if len(fileMatches) > 1 {
		fileList := fileMatches[1]
		
		// Extrai pares de chave-valor (nome do arquivo e conteúdo)
		// Padrão: "arquivo.ext": "conteúdo"
		fileItemRegex := regexp.MustCompile(`"([^"]+)"\s*:\s*"((?:\\"|[^"])*)"`)
		fileItemMatches := fileItemRegex.FindAllStringSubmatch(fileList, -1)
		
		for _, match := range fileItemMatches {
			if len(match) > 2 {
				fileName := match[1]
				fileContent := match[2]
				
				// Desescapa aspas no conteúdo
				fileContent = strings.ReplaceAll(fileContent, "\\\"", "\"")
				fileContent = strings.ReplaceAll(fileContent, "\\n", "\n")
				fileContent = strings.ReplaceAll(fileContent, "\\t", "\t")
				
				files[fileName] = fileContent
				fmt.Printf("Arquivo encontrado: %s (%d bytes)\n", fileName, len(fileContent))
			}
		}
	}

	// Caso especial: tenta extrair arquivos em formato alternativo
	if len(files) == 0 {
		fmt.Println("Tentando formato alternativo para arquivos...")
		altFileRegex := regexp.MustCompile(`"([^"]+\.(?:js|ts|json|md|html|css|jsx|tsx))"\s*:\s*"((?:\\"|[^"])*)"`)
		altFileMatches := altFileRegex.FindAllStringSubmatch(response, -1)
		
		for _, match := range altFileMatches {
			if len(match) > 2 {
				fileName := match[1]
				fileContent := match[2]
				
				// Desescapa aspas no conteúdo
				fileContent = strings.ReplaceAll(fileContent, "\\\"", "\"")
				fileContent = strings.ReplaceAll(fileContent, "\\n", "\n")
				fileContent = strings.ReplaceAll(fileContent, "\\t", "\t")
				
				files[fileName] = fileContent
				fmt.Printf("Arquivo encontrado (formato alt): %s (%d bytes)\n", fileName, len(fileContent))
			}
		}
	}

	return directories, files
}
