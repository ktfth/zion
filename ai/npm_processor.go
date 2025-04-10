package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ProcessNpmPackageJson processa especificamente respostas JSON que contêm package.json com @types
// Esta função é uma solução específica para o problema de parsing do JSON quando há pacotes npm com @
func ProcessNpmPackageJson(responseText string) ScaffoldResponse {
	fmt.Println("Processando package.json com @types...")
	
	// Inicializa a resposta
	var scaffoldResp ScaffoldResponse
	scaffoldResp.Structure.Directories = []string{}
	scaffoldResp.Structure.Files = make(map[string]string)
	
	// Extrai diretórios
	dirRegex := regexp.MustCompile(`"directories":\s*\[(.*?)\]`)
	dirMatches := dirRegex.FindStringSubmatch(responseText)
	
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
				scaffoldResp.Structure.Directories = append(scaffoldResp.Structure.Directories, dir)
				fmt.Printf("Diretório encontrado: %s\n", dir)
			}
		}
	}
	
	// Extrai arquivos
	// Primeiro, vamos extrair o bloco "files"
	filesRegex := regexp.MustCompile(`"files":\s*\{([\s\S]*?)\}[\s]*\}[\s]*\}`)
	filesMatches := filesRegex.FindStringSubmatch(responseText)
	
	if len(filesMatches) > 1 {
		filesBlock := filesMatches[1]
		
		// Agora vamos extrair cada arquivo individualmente
		// Padrão: "arquivo.ext": "conteúdo",
		fileRegex := regexp.MustCompile(`"([^"]+)":\s*"((?:\\"|[^"])*)"(?:,|\s*$)`)
		fileMatches := fileRegex.FindAllStringSubmatch(filesBlock, -1)
		
		for _, match := range fileMatches {
			if len(match) > 2 {
				fileName := match[1]
				fileContent := match[2]
				
				// Desescapa aspas no conteúdo
				fileContent = strings.ReplaceAll(fileContent, "\\\"", "\"")
				fileContent = strings.ReplaceAll(fileContent, "\\n", "\n")
				fileContent = strings.ReplaceAll(fileContent, "\\t", "\t")
				
				// Corrige o problema específico com package.json
				if fileName == "package.json" {
					fileContent = FixPackageJsonContent(fileContent)
				}
				
				scaffoldResp.Structure.Files[fileName] = fileContent
				fmt.Printf("Arquivo encontrado: %s (%d bytes)\n", fileName, len(fileContent))
			}
		}
	}
	
	return scaffoldResp
}

// FixPackageJsonContent corrige o conteúdo do package.json para lidar com o problema do @
func FixPackageJsonContent(content string) string {
	// Primeiro, vamos verificar se o conteúdo já é um JSON válido
	var testObj interface{}
	if err := json.Unmarshal([]byte(content), &testObj); err == nil {
		// Se não houver erro, o JSON já é válido
		return content
	}
	
	// Corrige o problema específico com @types
	// Busca por linhas como: "@types/express": "^4.17.17",
	atTypesRegex := regexp.MustCompile(`"(@[^"]+)":\s*"([^"]+)",`)
	fixed := atTypesRegex.ReplaceAllString(content, "\"$1\": \"$2\",")
	
	// Corrige outros problemas comuns
	fixed = strings.ReplaceAll(fixed, "\n,", ",")
	fixed = strings.ReplaceAll(fixed, ",\n,", ",")
	
	return fixed
}
