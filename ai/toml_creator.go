package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ExtractAndCreateTomlProject extrai diretamente os diretórios e arquivos do JSON
// e cria a estrutura do projeto com saída em formato TOML
func ExtractAndCreateTomlProject(projectName string, jsonStr string) error {
	// Remover blocos de código markdown se presentes
	if strings.HasPrefix(jsonStr, "```json\n") && strings.HasSuffix(jsonStr, "\n```") {
		jsonStr = strings.TrimPrefix(jsonStr, "```json\n")
		jsonStr = strings.TrimSuffix(jsonStr, "\n```")
	}

	// Criar o diretório raiz do projeto
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	// Coletar diretórios e arquivos
	directories := []string{}
	fileContents := make(map[string]string)

	// Extrair diretórios usando regex
	dirRegex := regexp.MustCompile(`"directories"\s*:\s*\[\s*((?:"[^"]+"\s*,?\s*)*)\]`)
	dirMatch := dirRegex.FindStringSubmatch(jsonStr)
	
	if len(dirMatch) > 1 {
		// Extrair cada diretório
		dirItemRegex := regexp.MustCompile(`"([^"]+)"`)
		dirItems := dirItemRegex.FindAllStringSubmatch(dirMatch[1], -1)
		
		for _, item := range dirItems {
			if len(item) > 1 {
				dirPath := filepath.Join(projectName, item[1])
				fmt.Println("Criando diretório:", item[1])
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("erro ao criar diretório %s: %v", item[1], err)
				}
				directories = append(directories, item[1])
			}
		}
	}
	
	// Extrair arquivos usando regex
	filesStartRegex := regexp.MustCompile(`"files"\s*:\s*\{`)
	filesStartMatch := filesStartRegex.FindStringIndex(jsonStr)
	
	if len(filesStartMatch) < 2 {
		return fmt.Errorf("não foi possível encontrar a seção 'files' no JSON")
	}
	
	// Posição do início da seção "files"
	filesStart := filesStartMatch[1]
	
	// Encontrar todos os nomes de arquivos
	fileNameRegex := regexp.MustCompile(`"([^"]+)"\s*:`)
	fileNameMatches := fileNameRegex.FindAllStringSubmatchIndex(jsonStr[filesStart:], -1)
	
	for i, match := range fileNameMatches {
		// Extrair o nome do arquivo
		fileName := jsonStr[filesStart+match[2]:filesStart+match[3]]
		
		// Determinar onde começa o conteúdo (após o ":")
		contentStartPos := filesStart + match[1]
		
		// Encontrar onde começa o conteúdo real (após o ":")
		quotePosRegex := regexp.MustCompile(`:\s*"`)
		quoteMatch := quotePosRegex.FindStringIndex(jsonStr[contentStartPos:])
		
		if len(quoteMatch) < 2 {
			// Verificar se é um objeto JSON em vez de uma string
			objStartRegex := regexp.MustCompile(`:\s*\{`)
			objMatch := objStartRegex.FindStringIndex(jsonStr[contentStartPos:])
			
			if len(objMatch) < 2 {
				fmt.Printf("Aviso: Não foi possível encontrar o início do conteúdo para o arquivo %s\n", fileName)
				continue
			}
			
			// É um objeto JSON, encontrar o fechamento correspondente
			objStart := contentStartPos + objMatch[1]
			objEnd := FindMatchingBrace(jsonStr, objStart)
			
			if objEnd == -1 {
				fmt.Printf("Aviso: Não foi possível encontrar o final do objeto para o arquivo %s\n", fileName)
				continue
			}
			
			// Extrair o conteúdo do objeto
			content := jsonStr[objStart-1:objEnd+1] // Incluir as chaves
			
			// Criar o arquivo
			CreateFile(projectName, fileName, content)
			fileContents[fileName] = content
			continue
		}
		
		// Posição após a aspas de abertura
		contentStart := contentStartPos + quoteMatch[1]
		
		// Encontrar o final do conteúdo (a próxima aspas não escapada)
		contentEnd := FindUnescapedQuote(jsonStr, contentStart)
		
		if contentEnd == -1 {
			fmt.Printf("Aviso: Não foi possível encontrar o final do conteúdo para o arquivo %s\n", fileName)
			continue
		}
		
		// Extrair o conteúdo
		content := jsonStr[contentStart:contentEnd]
		
		// Processar caracteres escapados apenas para arquivos que não são JSON
		if !isJSONFile(fileName) {
			content = ProcessEscapedChars(content)
		} else {
			// Para arquivos JSON, apenas remover escapes desnecessários
			content = PreserveJSONFormat(content)
			fmt.Printf("Preservado formato original do arquivo JSON: %s\n", fileName)
		}
		
		// Determinar onde começa o próximo arquivo ou o final da seção "files"
		var nextPos int
		if i < len(fileNameMatches)-1 {
			nextPos = filesStart + fileNameMatches[i+1][0]
		} else {
			// Último arquivo, encontrar o final da seção "files"
			endBraceRegex := regexp.MustCompile(`\}\s*\}`)
			endMatch := endBraceRegex.FindStringIndex(jsonStr[contentEnd:])
			if len(endMatch) < 2 {
				nextPos = len(jsonStr)
			} else {
				nextPos = contentEnd + endMatch[0]
			}
		}
		
		// Verificar se há uma vírgula após o conteúdo
		commaRegex := regexp.MustCompile(`"\s*,`)
		commaMatch := commaRegex.FindStringIndex(jsonStr[contentEnd:nextPos])
		
		if len(commaMatch) < 2 {
			// Se não houver vírgula e não for o último arquivo, pode haver um problema
			if i < len(fileNameMatches)-1 {
				fmt.Printf("Aviso: Formato inesperado após o conteúdo do arquivo %s\n", fileName)
			}
		}
		
		// Criar o arquivo
		CreateFile(projectName, fileName, content)
		fileContents[fileName] = content
	}
	
	// Gerar o arquivo TOML
	generateTomlFile(projectName, directories, fileContents)
	
	fmt.Println("\nEstrutura do projeto criada com sucesso em:", projectName)
	return nil
}

// generateTomlFile gera um arquivo TOML com a estrutura do projeto
func generateTomlFile(projectName string, directories []string, files map[string]string) error {
	// Criar o arquivo TOML
	tomlPath := filepath.Join(projectName, "project_structure.toml")
	
	// Construir o conteúdo TOML
	var tomlContent strings.Builder
	
	// Adicionar cabeçalho
	tomlContent.WriteString("# Estrutura do projeto gerada pelo Zion\n")
	tomlContent.WriteString(fmt.Sprintf("# Gerado em: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	// Adicionar informações do projeto
	tomlContent.WriteString("[project]\n")
	tomlContent.WriteString(fmt.Sprintf("name = \"%s\"\n", projectName))
	tomlContent.WriteString(fmt.Sprintf("created_at = \"%s\"\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	// Adicionar diretórios
	tomlContent.WriteString("# Diretórios\n")
	tomlContent.WriteString("directories = [\n")
	for _, dir := range directories {
		tomlContent.WriteString(fmt.Sprintf("  \"%s\",\n", dir))
	}
	tomlContent.WriteString("]\n\n")
	
	// Adicionar arquivos
	tomlContent.WriteString("# Arquivos\n")
	tomlContent.WriteString("[files]\n")
	for filePath := range files {
		// Usar o caminho do arquivo como chave e um comentário indicando que o conteúdo está no arquivo
		tomlContent.WriteString(fmt.Sprintf("\"%s\" = \"<conteúdo no arquivo>\"\n", filePath))
	}
	
	// Escrever o conteúdo no arquivo
	if err := os.WriteFile(tomlPath, []byte(tomlContent.String()), 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo TOML: %v", err)
	}
	
	fmt.Println("Arquivo TOML gerado em:", tomlPath)
	return nil
}

// As funções auxiliares foram movidas para o arquivo file_utils.go
