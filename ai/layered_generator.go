package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ktfth/zion/ai/providers"
	"github.com/ktfth/zion/config"
)

// LayeredGenerator gerencia a geração de projetos por camadas para contextos grandes
type LayeredGenerator struct {
	provider     providers.Provider
	maxTokens    int
	language     string
	projectName  string
	description  string
	llmsContext  *LLMsContext
	outputLayers []LayerResult
}

// LayerResult representa o resultado de uma camada de geração
type LayerResult struct {
	LayerName    string                 `json:"layer_name"`
	Description  string                 `json:"description"`
	Directories  []string               `json:"directories"`
	Files        map[string]interface{} `json:"files"`
	Dependencies []string               `json:"dependencies"`
	NextSteps    []string               `json:"next_steps"`
}

// LayeredResponse representa a resposta completa em camadas
type LayeredResponse struct {
	ProjectInfo struct {
		Name        string `json:"name"`
		Language    string `json:"language"`
		Description string `json:"description"`
	} `json:"project_info"`
	Layers []LayerResult `json:"layers"`
}

// NewLayeredGenerator cria uma nova instância do gerador em camadas
func NewLayeredGenerator(language, projectName, description string, llmsContext *LLMsContext) (*LayeredGenerator, error) {
	cfg := config.LoadConfig()
	provider, err := providers.DefaultManager.GetProvider(cfg.AIProvider, cfg.GetAIConfig())
	if err != nil {
		return nil, fmt.Errorf("erro ao obter provedor de IA: %v", err)
	}

	// Determinar limite de tokens baseado no provider
	maxTokens := determineMaxTokens(provider.Name())

	return &LayeredGenerator{
		provider:     provider,
		maxTokens:    maxTokens,
		language:     language,
		projectName:  projectName,
		description:  description,
		llmsContext:  llmsContext,
		outputLayers: make([]LayerResult, 0),
	}, nil
}

// determineMaxTokens retorna o limite seguro de tokens para cada provider
func determineMaxTokens(providerName string) int {
	switch strings.ToLower(providerName) {
	case "openrouter":
		return 120000 // Mais conservador para OpenRouter
	case "gpt":
		return 80000 // Mais conservador para GPT-4
	case "claude":
		return 120000 // Mais conservador para Claude
	case "gemini":
		return 150000 // Mais conservador para Gemini
	default:
		return 40000 // Muito conservador para providers desconhecidos
	}
}

// EstimateTokens estima aproximadamente o número de tokens em um texto (função exportada)
func EstimateTokens(text string) int {
	return estimateTokens(text)
}

// estimateTokens estima aproximadamente o número de tokens em um texto
func estimateTokens(text string) int {
	// Estimativa mais precisa baseada em:
	// - Contagem de caracteres
	// - Formatação JSON
	// - Overhead de instruções

	charCount := len(text)

	// Aproximação: 1 token ≈ 4 caracteres para texto normal
	// Mais overhead para JSON e instruções
	baseTokens := charCount / 4

	// Adicionar overhead para formatação JSON e estruturas
	jsonOverhead := strings.Count(text, "{") + strings.Count(text, "}") + strings.Count(text, "[") + strings.Count(text, "]")

	// Adicionar overhead para instruções e prompts
	instructionOverhead := 3000

	return baseTokens + jsonOverhead*10 + instructionOverhead
}

// GenerateLayeredProject gera o projeto em camadas para evitar overflow de contexto
func (lg *LayeredGenerator) GenerateLayeredProject() (*LayeredResponse, error) {
	fmt.Printf("🏗️ Iniciando geração em camadas (limite: %d tokens)\n", lg.maxTokens)

	// 1. Primeiro, determinar a estratégia de camadas
	layers, err := lg.planLayers()
	if err != nil {
		return nil, fmt.Errorf("erro ao planejar camadas: %v", err)
	}

	fmt.Printf("📋 Planejadas %d camadas de geração\n", len(layers))

	// 2. Gerar cada camada sequencialmente
	response := &LayeredResponse{}
	response.ProjectInfo.Name = lg.projectName
	response.ProjectInfo.Language = lg.language
	response.ProjectInfo.Description = lg.description

	for i, layerPlan := range layers {
		fmt.Printf("🔧 Gerando camada %d/%d: %s...\n", i+1, len(layers), layerPlan.Name)

		layerResult, err := lg.generateLayer(layerPlan, response.Layers)
		if err != nil {
			return nil, fmt.Errorf("erro na camada %d (%s): %v", i+1, layerPlan.Name, err)
		}

		response.Layers = append(response.Layers, *layerResult)
		lg.outputLayers = append(lg.outputLayers, *layerResult)

		fmt.Printf("✅ Camada %s concluída (%d arquivos, %d diretórios)\n",
			layerPlan.Name, len(layerResult.Files), len(layerResult.Directories))
	}

	fmt.Printf("🎉 Geração em camadas concluída: %d camadas, %d arquivos totais\n",
		len(response.Layers), lg.countTotalFiles(response.Layers))

	return response, nil
}

// LayerPlan define o plano para uma camada específica
type LayerPlan struct {
	Name        string
	Description string
	Priority    int
	Focus       []string
}

// planLayers determina quais camadas criar baseado no projeto
func (lg *LayeredGenerator) planLayers() ([]LayerPlan, error) {
	basePrompt := fmt.Sprintf(`Analise este projeto e determine as camadas de desenvolvimento necessárias:

PROJETO: %s
LINGUAGEM: %s
DESCRIÇÃO: %s

Retorne um JSON com as camadas em ordem de prioridade:
{
  "layers": [
    {
      "name": "nome-da-camada",
      "description": "descrição da camada",
      "priority": 1,
      "focus": ["elemento1", "elemento2"]
    }
  ]
}

REGRAS:
- Máximo 6 camadas
- Cada camada deve ter foco específico (core, config, api, frontend, tests, docs)
- Prioridade 1 = mais importante
- Focus deve listar os elementos principais da camada`, lg.projectName, lg.language, lg.description)

	// Adicionar contexto se disponível, mas limitado
	if lg.llmsContext != nil && lg.llmsContext.HasLLMsFile {
		contextSummary := lg.llmsContext.GetProjectDescription()
		if len(contextSummary) > 1000 {
			contextSummary = contextSummary[:1000] + "... (resumido)"
		}
		basePrompt += fmt.Sprintf("\n\nCONTEXTO ADICIONAL:\n%s", contextSummary)
	}

	response, err := lg.provider.GenerateContent(basePrompt)
	if err != nil {
		return nil, err
	}

	// Extrair JSON da resposta
	jsonContent := extractJSONContent(response)

	var planResponse struct {
		Layers []LayerPlan `json:"layers"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &planResponse); err != nil {
		// Fallback para camadas padrão se o parsing falhar
		return lg.getDefaultLayers(), nil
	}

	// Ordenar por prioridade
	sort.Slice(planResponse.Layers, func(i, j int) bool {
		return planResponse.Layers[i].Priority < planResponse.Layers[j].Priority
	})

	if len(planResponse.Layers) == 0 {
		return lg.getDefaultLayers(), nil
	}

	return planResponse.Layers, nil
}

// getDefaultLayers retorna camadas padrão caso o planejamento automático falhe
func (lg *LayeredGenerator) getDefaultLayers() []LayerPlan {
	defaultLayers := []LayerPlan{
		{
			Name:        "core",
			Description: "Estrutura básica e arquivos de configuração",
			Priority:    1,
			Focus:       []string{"main", "config", "package"},
		},
		{
			Name:        "business",
			Description: "Lógica de negócio e modelos",
			Priority:    2,
			Focus:       []string{"models", "services", "utils"},
		},
		{
			Name:        "api",
			Description: "Endpoints e controladores",
			Priority:    3,
			Focus:       []string{"controllers", "routes", "middleware"},
		},
		{
			Name:        "tests",
			Description: "Testes unitários e de integração",
			Priority:    4,
			Focus:       []string{"tests", "mocks", "fixtures"},
		},
	}

	// Adaptar para linguagens específicas
	switch strings.ToLower(lg.language) {
	case "javascript", "typescript":
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "frontend",
			Description: "Componentes de interface e assets",
			Priority:    3,
			Focus:       []string{"components", "pages", "styles"},
		})
	case "go":
		defaultLayers[1].Focus = append(defaultLayers[1].Focus, "handlers", "internal")
	case "python":
		defaultLayers[1].Focus = append(defaultLayers[1].Focus, "modules", "packages")
	}

	return defaultLayers
}

// generateLayer gera uma camada específica
func (lg *LayeredGenerator) generateLayer(plan LayerPlan, previousLayers []LayerResult) (*LayerResult, error) {
	prompt := lg.buildLayerPrompt(plan, previousLayers)

	// Verificar se o prompt está dentro do limite
	if estimateTokens(prompt) > lg.maxTokens {
		// Se ainda assim for grande demais, simplificar o prompt
		prompt = lg.buildSimplifiedLayerPrompt(plan, previousLayers)
	}

	response, err := lg.provider.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	// Processar resposta da camada
	return lg.parseLayerResponse(plan.Name, response)
}

// buildLayerPrompt constrói o prompt para uma camada específica
func (lg *LayeredGenerator) buildLayerPrompt(plan LayerPlan, previousLayers []LayerResult) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf(`Gere EXCLUSIVAMENTE a camada "%s" para o projeto:

PROJETO: %s
LINGUAGEM: %s
DESCRIÇÃO GERAL: %s

CAMADA ATUAL - FOQUE APENAS NISTO:
Nome: %s
Descrição: %s
Elementos desta camada: %v

IMPORTANTE: Esta camada deve conter apenas arquivos relacionados a: %s

`, lg.projectName, lg.language, lg.description, plan.Name, plan.Description, plan.Focus, strings.Join(plan.Focus, ", ")))

	// Adicionar contexto das camadas anteriores
	if len(previousLayers) > 0 {
		prompt.WriteString("CAMADAS JÁ CRIADAS:\n")
		allCreatedFiles := make([]string, 0)
		for _, layer := range previousLayers {
			prompt.WriteString(fmt.Sprintf("- %s: %d arquivos\n", layer.LayerName, len(layer.Files)))
			for fileName := range layer.Files {
				allCreatedFiles = append(allCreatedFiles, fileName)
			}
		}

		if len(allCreatedFiles) > 0 {
			prompt.WriteString("\nARQUIVOS JÁ CRIADOS (NÃO RECRIAR):\n")
			for _, fileName := range allCreatedFiles {
				prompt.WriteString(fmt.Sprintf("- %s\n", fileName))
			}
		}
		prompt.WriteString("\n")
	}

	// Adicionar instruções específicas por camada
	switch plan.Name {
	case "core":
		prompt.WriteString(`CAMADA CORE - CRIE APENAS:
- Arquivo principal (index.js, main.js, app.js)
- package.json com dependências
- Configurações básicas (.env.example, config files)
- Estrutura de diretórios básica
NÃO crie testes, rotas específicas ou lógica de negócio

`)
	case "business":
		prompt.WriteString(`CAMADA BUSINESS - CRIE APENAS:
- Modelos de dados (models/)
- Serviços de negócio (services/)
- Utilitários (utils/)
- Middleware de negócio
NÃO crie testes, rotas HTTP ou configurações

`)
	case "api":
		prompt.WriteString(`CAMADA API - CRIE APENAS:
- Controladores (controllers/)
- Rotas HTTP (routes/)
- Middleware de API
- Validadores de entrada
NÃO crie testes, modelos ou configurações

`)
	case "tests":
		prompt.WriteString(`CAMADA TESTS - CRIE APENAS:
- Arquivos de teste (.test.js, .spec.js)
- Configuração de testes
- Mocks e fixtures
- Scripts de teste
NÃO crie código de produção

`)
	}

	// Adicionar contexto limitado do projeto
	if lg.llmsContext != nil && lg.llmsContext.HasLLMsFile {
		summary := lg.llmsContext.GetProjectDescription()
		if len(summary) > 500 {
			summary = summary[:500] + "..."
		}
		prompt.WriteString(fmt.Sprintf("CONTEXTO DO PROJETO:\n%s\n\n", summary))
	}

	prompt.WriteString(`FORMATO DE SAÍDA OBRIGATÓRIO:
{
  "layer_name": "` + plan.Name + `",
  "description": "descrição da camada implementada",
  "directories": ["lista", "de", "diretórios"],
  "files": {
    "arquivo.ext": {
      "content": "conteúdo do arquivo"
    }
  },
  "dependencies": ["deps", "necessárias"],
  "next_steps": ["próximos", "passos"]
}

INSTRUÇÕES:
1. Foque APENAS nos elementos desta camada especificados acima
2. NÃO recrie NENHUM arquivo listado como "JÁ CRIADOS"
3. Se um arquivo já existe, NÃO o inclua nesta camada
4. Crie apenas arquivos novos específicos desta camada
5. Use conteúdo realista e funcional
6. Mantenha coerência com o projeto
7. Retorne apenas JSON válido`)

	return prompt.String()
}

// buildSimplifiedLayerPrompt constrói um prompt simplificado se o normal for muito grande
func (lg *LayeredGenerator) buildSimplifiedLayerPrompt(plan LayerPlan, previousLayers []LayerResult) string {
	return fmt.Sprintf(`Crie a camada "%s" para projeto %s (%s):

Foco: %v
Camadas anteriores: %d

Retorne JSON:
{
  "layer_name": "%s",
  "description": "...",
  "directories": [...],
  "files": {...},
  "dependencies": [...],
  "next_steps": [...]
}

Apenas arquivos essenciais da camada.`, plan.Name, lg.projectName, lg.language, plan.Focus, len(previousLayers), plan.Name)
}

// parseLayerResponse processa a resposta de uma camada
func (lg *LayeredGenerator) parseLayerResponse(layerName, response string) (*LayerResult, error) {
	jsonContent := extractJSONContent(response)

	// Validar a estrutura da camada
	validation := ValidateProjectStructure(jsonContent, lg.language)
	if !validation.IsValid {
		fmt.Printf("⚠️  Camada %s apresenta problemas:\n", layerName)
		for _, issue := range validation.Issues {
			fmt.Printf("   • %s\n", issue)
		}
		fmt.Printf("📊 Pontuação de qualidade: %.1f/100\n", validation.Score)

		// Se o score for muito baixo, falhar
		if validation.Score < 40 {
			return nil, fmt.Errorf("camada %s não passou na validação (score: %.1f/100)", layerName, validation.Score)
		}

		// Caso contrário, mostrar avisos mas continuar
		fmt.Printf("⚠️  Continuando com avisos...\n")
	} else {
		fmt.Printf("✅ Camada %s validada com sucesso (score: %.1f/100)\n", layerName, validation.Score)
		if len(validation.Suggestions) > 0 {
			fmt.Printf("💡 Sugestões de melhoria para camada %s:\n", layerName)
			for _, suggestion := range validation.Suggestions {
				fmt.Printf("   • %s\n", suggestion)
			}
		}
	}

	var layerResult LayerResult
	if err := json.Unmarshal([]byte(jsonContent), &layerResult); err != nil {
		return nil, fmt.Errorf("erro ao processar resposta da camada %s: %v", layerName, err)
	}

	// Garantir que o nome da camada está correto
	layerResult.LayerName = layerName

	return &layerResult, nil
}

// countTotalFiles conta o total de arquivos em todas as camadas
func (lg *LayeredGenerator) countTotalFiles(layers []LayerResult) int {
	total := 0
	for _, layer := range layers {
		total += len(layer.Files)
	}
	return total
}

// ConvertToScaffoldResponse converte a resposta em camadas para o formato padrão
func (lg *LayeredGenerator) ConvertToScaffoldResponse(layeredResponse *LayeredResponse) *ScaffoldResponse {
	response := &ScaffoldResponse{}
	response.Structure.Directories = make([]string, 0)
	response.Structure.Files = make(map[string]interface{})

	// Coletar todos os diretórios (removendo duplicatas)
	dirSet := make(map[string]bool)
	for _, layer := range layeredResponse.Layers {
		for _, dir := range layer.Directories {
			if !dirSet[dir] {
				dirSet[dir] = true
				response.Structure.Directories = append(response.Structure.Directories, dir)
			}
		}
	}

	// Coletar todos os arquivos
	for _, layer := range layeredResponse.Layers {
		for fileName, content := range layer.Files {
			response.Structure.Files[fileName] = content
		}
	}

	return response
}

// DetectContextOverflow verifica se um prompt causaria overflow
func DetectContextOverflow(prompt string, providerName string) bool {
	maxTokens := determineMaxTokens(providerName)
	estimatedTokens := estimateTokens(prompt)

	fmt.Printf("🔍 Estimativa de tokens: %d/%d\n", estimatedTokens, maxTokens)

	return estimatedTokens > maxTokens
}

// IsContextOverflowError verifica se um erro é de overflow de contexto
func IsContextOverflowError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())

	// Padrões comuns de erro de contexto
	contextErrorPatterns := []string{
		"context length",
		"maximum context",
		"token limit",
		"too many tokens",
		"context too long",
		"input too long",
		"exceeds.*tokens",
		"reduce.*length",
	}

	for _, pattern := range contextErrorPatterns {
		if matched, _ := regexp.MatchString(pattern, errorStr); matched {
			return true
		}
	}

	return false
}
