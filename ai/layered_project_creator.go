package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateLayeredProject cria a estrutura do projeto baseada na resposta em camadas
func CreateLayeredProject(projectName string, layeredResponse *LayeredResponse) error {
	// Validar estrutura antes de criar o projeto
	responseBytes, err := json.Marshal(layeredResponse)
	if err != nil {
		return fmt.Errorf("erro ao serializar resposta para validação: %v", err)
	}

	validation := ValidateProjectStructure(string(responseBytes), layeredResponse.ProjectInfo.Language)
	if !validation.IsValid {
		fmt.Printf("⚠️  Estrutura do projeto apresenta problemas:\n")
		for _, issue := range validation.Issues {
			fmt.Printf("   • %s\n", issue)
		}
		fmt.Printf("📊 Pontuação de qualidade: %.1f/100\n", validation.Score)

		// Se o score for muito baixo, falhar
		if validation.Score < 50 {
			return fmt.Errorf("projeto não passou na validação (score: %.1f/100)", validation.Score)
		}

		// Caso contrário, mostrar avisos mas continuar
		fmt.Printf("⚠️  Continuando com avisos...\n")
	} else {
		fmt.Printf("✅ Estrutura validada com sucesso (score: %.1f/100)\n", validation.Score)
		if len(validation.Suggestions) > 0 {
			fmt.Printf("💡 Sugestões de melhoria:\n")
			for _, suggestion := range validation.Suggestions {
				fmt.Printf("   • %s\n", suggestion)
			}
		}
	}

	fmt.Printf("\n🏗️  Materializando projeto em camadas: %s\n", projectName)

	// Criar o diretório raiz do projeto
	projectDir := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %v", err)
	}

	// Coletar todos os diretórios únicos de todas as camadas
	allDirs := make(map[string]bool)
	for _, layer := range layeredResponse.Layers {
		for _, dir := range layer.Directories {
			allDirs[dir] = true
		}
	}

	// Criar todos os diretórios
	fmt.Printf("📁 Criando %d diretórios...\n", len(allDirs))
	for dir := range allDirs {
		if dir == "" || dir == "." {
			continue
		}
		dirPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
		}
		fmt.Printf("   📂 %s\n", dir)
	}

	// Criar arquivos por camada
	totalFiles := 0
	createdFiles := make(map[string]string) // mapa para rastrear arquivos criados: arquivo -> camada

	for i, layer := range layeredResponse.Layers {
		fmt.Printf("\n⚙️  Materializando camada %d/%d: %s\n", i+1, len(layeredResponse.Layers), layer.LayerName)

		layerFileCount := 0
		for filePath, content := range layer.Files {
			fullPath := filepath.Join(projectDir, filePath)

			// Verificar se o arquivo já foi criado por uma camada anterior
			if existingLayer, exists := createdFiles[filePath]; exists {
				fmt.Printf("   ⚠️  Arquivo %s já foi criado na camada '%s', pulando...\n", filePath, existingLayer)
				continue
			}

			// Garantir que o diretório pai exista
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório pai para %s: %v", filePath, err)
			}

			// Determinar o conteúdo do arquivo
			var contentBytes []byte
			var err error

			// Se o content é um objeto com campo "content"
			if contentMap, ok := content.(map[string]interface{}); ok {
				if contentValue, exists := contentMap["content"]; exists {
					// Se é string, usar diretamente
					if strContent, isString := contentValue.(string); isString {
						contentBytes = []byte(strContent)
					} else {
						// Se é objeto (JSON), serializar
						contentBytes, err = json.MarshalIndent(contentValue, "", "  ")
						if err != nil {
							return fmt.Errorf("erro ao serializar conteúdo JSON para %s: %v", filePath, err)
						}
					}
				} else {
					// Se não tem campo "content", usar o objeto inteiro
					contentBytes, err = json.MarshalIndent(content, "", "  ")
					if err != nil {
						return fmt.Errorf("erro ao serializar conteúdo para %s: %v", filePath, err)
					}
				}
			} else if strContent, ok := content.(string); ok {
				// Se é string direta
				contentBytes = []byte(strContent)
			} else {
				// Se é qualquer outro tipo, serializar como JSON
				contentBytes, err = json.MarshalIndent(content, "", "  ")
				if err != nil {
					return fmt.Errorf("erro ao serializar conteúdo para %s: %v", filePath, err)
				}
			}

			// Escrever o arquivo
			if err := os.WriteFile(fullPath, contentBytes, 0644); err != nil {
				return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
			}

			// Registrar o arquivo como criado
			createdFiles[filePath] = layer.LayerName

			fmt.Printf("   📄 %s (%d bytes)\n", filePath, len(contentBytes))
			layerFileCount++
			totalFiles++
		}

		fmt.Printf("✅ Camada %s: %d arquivos criados\n", layer.LayerName, layerFileCount)

		// Mostrar dependências se existirem
		if len(layer.Dependencies) > 0 {
			fmt.Printf("   📦 Dependências: %v\n", layer.Dependencies)
		}

		// Mostrar próximos passos se existirem
		if len(layer.NextSteps) > 0 {
			fmt.Printf("   📋 Próximos passos: %v\n", layer.NextSteps)
		}
	}

	fmt.Printf("\n🎉 Projeto criado com sucesso!\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📁 Local: %s\n", projectDir)
	fmt.Printf("📂 Diretórios: %d\n", len(allDirs))
	fmt.Printf("📄 Arquivos: %d\n", totalFiles)
	fmt.Printf("🏗️  Camadas: %d\n", len(layeredResponse.Layers))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Sempre gerar ou melhorar o README.md com instruções reais
	fmt.Println("📋 Gerando README.md com instruções reais...")
	readmeContent := generateReadmeContentForLayered(projectName, layeredResponse)

	// Verificar se já existe um README.md
	readmePath := filepath.Join(projectDir, "README.md")
	hasExistingReadme := false

	// Verificar se README.md foi criado em alguma camada
	for _, layer := range layeredResponse.Layers {
		if _, exists := layer.Files["README.md"]; exists {
			hasExistingReadme = true
			break
		}
	}

	if !hasExistingReadme {
		// Criar novo README.md
		if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
			fmt.Printf("⚠️  Erro ao criar README.md: %v\n", err)
		} else {
			fmt.Printf("✅ README.md criado com instruções detalhadas\n")
		}
	} else {
		// Verificar se o README existente é básico demais
		if existingContent, err := os.ReadFile(readmePath); err == nil {
			existingReadme := string(existingContent)

			// Verificar se o README existente é muito básico
			if len(strings.TrimSpace(existingReadme)) < 200 || (!strings.Contains(existingReadme, "##") && !strings.Contains(existingReadme, "Instalação") && !strings.Contains(existingReadme, "Como Executar")) {
				if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
					fmt.Printf("⚠️  Erro ao melhorar README.md: %v\n", err)
				} else {
					fmt.Printf("✅ README.md melhorado com instruções detalhadas\n")
				}
			} else {
				fmt.Printf("ℹ️  README.md já contém instruções adequadas\n")
			}
		}
	}

	// Criar um arquivo de resumo das camadas
	summaryPath := filepath.Join(projectDir, "ZION_LAYERS_SUMMARY.md")
	summaryContent := generateLayersSummary(layeredResponse)
	if err := os.WriteFile(summaryPath, []byte(summaryContent), 0644); err != nil {
		fmt.Printf("⚠️  Aviso: Não foi possível criar resumo das camadas: %v\n", err)
	} else {
		fmt.Printf("📋 Resumo das camadas salvo em: ZION_LAYERS_SUMMARY.md\n")
	}

	return nil
}

// generateReadmeContentForLayered gera conteúdo do README.md para projetos em camadas
func generateReadmeContentForLayered(projectName string, layeredResponse *LayeredResponse) string {
	language := layeredResponse.ProjectInfo.Language
	description := layeredResponse.ProjectInfo.Description

	// Coletar todos os arquivos de todas as camadas
	allFiles := make(map[string]interface{})
	allDirs := make(map[string]bool)

	for _, layer := range layeredResponse.Layers {
		for filePath, content := range layer.Files {
			allFiles[filePath] = content
		}
		for _, dir := range layer.Directories {
			allDirs[dir] = true
		}
	}

	// Converter para lista de diretórios
	var directories []string
	for dir := range allDirs {
		directories = append(directories, dir)
	}

	// Criar estrutura scaffoldResp temporária para reusar as funções existentes
	scaffoldResp := &ScaffoldResponse{
		Structure: struct {
			Directories []string               `json:"directories"`
			Files       map[string]interface{} `json:"files"`
		}{
			Directories: directories,
			Files:       allFiles,
		},
	}

	// Detectar tipo de projeto baseado nos arquivos
	projectType := detectProjectTypeFromFiles(scaffoldResp)

	// Gerar conteúdo do README
	readme := fmt.Sprintf(`# %s

## Descrição

%s

Este projeto foi gerado usando o Zion CLI com sistema de camadas e contém uma estrutura %s para %s.

## Estrutura do Projeto

`, projectName, description, projectType, language)

	// Adicionar informações sobre diretórios
	if len(directories) > 0 {
		readme += "### Diretórios\n\n"
		for _, dir := range directories {
			readme += fmt.Sprintf("- `%s/` - %s\n", dir, getDirectoryDescription(dir, language))
		}
		readme += "\n"
	}

	// Adicionar informações sobre arquivos importantes
	readme += "### Arquivos Principais\n\n"
	for filePath := range allFiles {
		if isImportantFile(filePath) {
			readme += fmt.Sprintf("- `%s` - %s\n", filePath, getFileDescription(filePath, language))
		}
	}

	// Adicionar informações sobre camadas
	readme += "### Camadas do Projeto\n\n"
	readme += "Este projeto foi gerado usando o sistema de camadas do Zion AI:\n\n"

	for i, layer := range layeredResponse.Layers {
		readme += fmt.Sprintf("%d. **%s** - %s\n", i+1, layer.LayerName, layer.Description)
	}
	readme += "\n"

	// Adicionar instruções de instalação/configuração
	readme += generateInstallationInstructions(language, scaffoldResp)

	// Adicionar instruções de execução
	readme += generateRunInstructions(language, scaffoldResp)

	// Adicionar próximos passos
	readme += generateNextSteps(language, scaffoldResp)

	return readme
}

// generateLayersSummary gera um resumo em markdown das camadas criadas
func generateLayersSummary(layeredResponse *LayeredResponse) string {
	summary := fmt.Sprintf(`# Resumo das Camadas - %s

Este projeto foi gerado pelo Zion AI usando o sistema de camadas para gerenciar contextos grandes.

## Informações do Projeto

- **Nome**: %s
- **Linguagem**: %s
- **Descrição**: %s
- **Camadas Criadas**: %d

## Camadas Implementadas

`, layeredResponse.ProjectInfo.Name,
		layeredResponse.ProjectInfo.Name,
		layeredResponse.ProjectInfo.Language,
		layeredResponse.ProjectInfo.Description,
		len(layeredResponse.Layers))

	for i, layer := range layeredResponse.Layers {
		summary += fmt.Sprintf(`### %d. %s

**Descrição**: %s

**Arquivos criados** (%d):
`, i+1, layer.LayerName, layer.Description, len(layer.Files))

		for filePath := range layer.Files {
			summary += fmt.Sprintf("- `%s`\n", filePath)
		}

		if len(layer.Directories) > 0 {
			summary += fmt.Sprintf("\n**Diretórios** (%d):\n", len(layer.Directories))
			for _, dir := range layer.Directories {
				summary += fmt.Sprintf("- `%s/`\n", dir)
			}
		}

		if len(layer.Dependencies) > 0 {
			summary += fmt.Sprintf("\n**Dependências**:\n")
			for _, dep := range layer.Dependencies {
				summary += fmt.Sprintf("- %s\n", dep)
			}
		}

		if len(layer.NextSteps) > 0 {
			summary += fmt.Sprintf("\n**Próximos Passos**:\n")
			for _, step := range layer.NextSteps {
				summary += fmt.Sprintf("- %s\n", step)
			}
		}

		summary += "\n"
	}

	summary += `
## Como foi gerado

Este projeto foi criado usando o sistema de geração em camadas do Zion AI, que divide projetos grandes em múltiplas etapas de geração para evitar limitações de contexto das APIs de IA.

Cada camada foi gerada de forma sequencial, mantendo consistência entre elas e garantindo que o resultado final seja um projeto coeso e funcional.

## Próximos Passos

1. Revise os arquivos gerados em cada camada
2. Instale as dependências listadas nas camadas relevantes
3. Execute os comandos sugeridos nos "Próximos Passos" de cada camada
4. Teste a funcionalidade básica do projeto

---
*Gerado por Zion AI - Sistema de Geração em Camadas*
`

	return summary
}

// ExtractAndCreateLayeredProject é um wrapper que aceita JSON de resposta em camadas
func ExtractAndCreateLayeredProject(projectName, jsonResponse string) error {
	// Tentar decodificar como resposta em camadas primeiro
	var layeredResponse LayeredResponse
	if err := json.Unmarshal([]byte(jsonResponse), &layeredResponse); err == nil {
		// É uma resposta em camadas
		return CreateLayeredProject(projectName, &layeredResponse)
	}

	// Se não for resposta em camadas, tentar o método tradicional
	return ExtractAndCreateProject(projectName, jsonResponse)
}
