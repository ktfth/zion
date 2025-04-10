package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateTomlOutput gera um arquivo TOML com a estrutura do projeto
func GenerateTomlOutput(projectName string, jsonStr string, directories []string, files map[string]string) error {
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
