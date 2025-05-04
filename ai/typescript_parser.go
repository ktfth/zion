package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ParseAndCreateTypeScriptProject é uma função especializada para criar projetos TypeScript
// que contêm pacotes npm com o caractere @
func ParseAndCreateTypeScriptProject(projectName, response string) error {
	fmt.Println("Usando parser especializado para projetos TypeScript...")

	// Extrair o conteúdo JSON da resposta
	jsonContent := extractJSONFromMarkdown(response)
	fmt.Println("Conteúdo JSON extraído do markdown")

	// Extrair diretórios
	directories := extractDirectories(jsonContent)
	if len(directories) == 0 {
		fmt.Println("Nenhum diretório encontrado, tentando extrair de forma alternativa...")
		directories = extractDirectoriesAlternative(jsonContent)
	}

	// Extrair arquivos
	files := extractFilesWithAtSymbol(jsonContent)
	if len(files) == 0 {
		fmt.Println("Nenhum arquivo encontrado, tentando extrair de forma alternativa...")
		files = extractFilesAlternative(jsonContent)
	}

	// Verificar se conseguimos extrair algo
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

	fmt.Printf("\nEstrutura do projeto TypeScript criada com sucesso em: %s\n", projectDir)
	return nil
}

// extractJSONFromMarkdown extrai o conteúdo JSON de um bloco de código markdown
func extractJSONFromMarkdown(response string) string {
	// Verificar se a resposta está em um bloco de código markdown
	if strings.Contains(response, "```json") {
		fmt.Println("Encontrado bloco de código JSON...")
		parts := strings.Split(response, "```json")
		if len(parts) > 1 {
			endParts := strings.Split(parts[1], "```")
			if len(endParts) > 0 {
				return strings.TrimSpace(endParts[0])
			}
		}
	} else if strings.Contains(response, "```") {
		fmt.Println("Encontrado bloco de código genérico...")
		parts := strings.Split(response, "```")
		if len(parts) > 1 {
			endParts := strings.Split(parts[1], "```")
			if len(endParts) > 0 {
				return strings.TrimSpace(endParts[0])
			}
		}
	}

	// Se não encontrar blocos de código, retorna a resposta original
	return response
}

// extractDirectories extrai os diretórios do conteúdo JSON
func extractDirectories(jsonContent string) []string {
	directories := []string{}

	// Regex para extrair o array de diretórios
	dirRegex := regexp.MustCompile(`"directories":\s*\[\s*([\s\S]*?)\s*\]`)
	dirMatches := dirRegex.FindStringSubmatch(jsonContent)

	if len(dirMatches) > 1 {
		dirList := dirMatches[1]
		// Regex para extrair cada diretório
		dirItemRegex := regexp.MustCompile(`"([^"]+)"`)
		dirItemMatches := dirItemRegex.FindAllStringSubmatch(dirList, -1)

		for _, match := range dirItemMatches {
			if len(match) > 1 {
				directories = append(directories, match[1])
				fmt.Printf("Diretório encontrado: %s\n", match[1])
			}
		}
	}

	return directories
}

// extractDirectoriesAlternative é uma abordagem alternativa para extrair diretórios
func extractDirectoriesAlternative(jsonContent string) []string {
	// Para projetos TypeScript, vamos adicionar alguns diretórios padrão
	// se não conseguirmos extrair da resposta
	return []string{
		"src",
		"dist",
		"node_modules",
	}
}

// extractFilesWithAtSymbol extrai os arquivos do conteúdo JSON, lidando com o caractere @
func extractFilesWithAtSymbol(jsonContent string) map[string]string {
	files := make(map[string]string)

	// Extrair o bloco de arquivos
	filesRegex := regexp.MustCompile(`"files":\s*\{([\s\S]*?)\}\s*\}`)
	filesMatches := filesRegex.FindStringSubmatch(jsonContent)

	if len(filesMatches) > 1 {
		filesBlock := filesMatches[1]

		// Para cada arquivo, verificar se o conteúdo está em formato wrapper
		fileRegex := regexp.MustCompile(`"([^"]+)":\s*(\{[^}]+\}|"[^"]+")\s*,?`)
		fileMatches := fileRegex.FindAllStringSubmatch(filesBlock, -1)

		for _, match := range fileMatches {
			if len(match) > 2 {
				fileName := match[1]
				content := match[2]

				// Se o conteúdo está em formato de objeto
				if strings.HasPrefix(content, "{") {
					// Extrair o conteúdo do wrapper se existir
					contentRegex := regexp.MustCompile(`{\s*"content":\s*"(.*?)"\s*}`)
					contentMatches := contentRegex.FindStringSubmatch(content)
					if len(contentMatches) > 1 {
						content = contentMatches[1]
					}
				} else {
					// Remover aspas do início e fim
					content = strings.Trim(content, "\"")
					content = ProcessEscapedChars(content)
				}

				// Caso especial para package.json
				if fileName == "package.json" {
					var jsonContent map[string]interface{}
					if err := json.Unmarshal([]byte(content), &jsonContent); err == nil {
						ProcessPackageJsonContent(jsonContent)
						jsonBytes, err := json.MarshalIndent(jsonContent, "", "  ")
						if err == nil {
							content = string(jsonBytes)
						}
					}
				}

				files[fileName] = content
				fmt.Printf("Arquivo processado: %s\n", fileName)
			}
		}
	}

	return files
}

// extractFilesAlternative é uma abordagem alternativa para extrair arquivos
func extractFilesAlternative(jsonContent string) map[string]string {
	files := make(map[string]string)

	// Extrair arquivos individualmente
	filePatterns := []string{
		`"(package\.json)":\s*"([\s\S]*?)"`,
		`"(tsconfig\.json)":\s*"([\s\S]*?)"`,
		`"(README\.md)":\s*"([\s\S]*?)"`,
		`"(src/index\.ts)":\s*"([\s\S]*?)"`,
	}

	for _, pattern := range filePatterns {
		fileRegex := regexp.MustCompile(pattern)
		fileMatches := fileRegex.FindStringSubmatch(jsonContent)

		if len(fileMatches) > 2 {
			fileName := fileMatches[1]
			fileContent := fileMatches[2]

			// Corrigir o problema com @types em package.json
			if fileName == "package.json" {
				// Agora usa a versão centralizada da função
				var jsonContent map[string]interface{}
				if err := json.Unmarshal([]byte(fileContent), &jsonContent); err == nil {
					ProcessPackageJsonContent(jsonContent)
					if jsonBytes, err := json.MarshalIndent(jsonContent, "", "  "); err == nil {
						fileContent = string(jsonBytes)
					}
				}
			}

			files[fileName] = fileContent
			fmt.Printf("Arquivo encontrado (alternativo): %s\n", fileName)
		}
	}

	// Se não encontrou package.json, criar um básico
	if _, ok := files["package.json"]; !ok {
		files["package.json"] = `{
  "name": "hello",
  "version": "1.0.0",
  "description": "Simple Hello World API",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  },
  "devDependencies": {
    "@types/express": "^4.17.17",
    "@types/node": "^20.4.5",
    "typescript": "^5.1.6"
  }
}`
		fmt.Println("Criado package.json padrão")
	}

	// Se não encontrou tsconfig.json, criar um básico
	if _, ok := files["tsconfig.json"]; !ok {
		files["tsconfig.json"] = `{
  "compilerOptions": {
    "target": "es6",
    "module": "commonjs",
    "rootDir": "./src",
    "outDir": "./dist",
    "esModuleInterop": true,
    "strict": true
  }
}`
		fmt.Println("Criado tsconfig.json padrão")
	}

	return files
}
