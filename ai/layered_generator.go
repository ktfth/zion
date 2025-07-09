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
	provider              providers.Provider
	maxTokens             int
	language              string
	projectName           string
	description           string
	llmsContext           *LLMsContext
	outputLayers          []LayerResult
	instructionController *AdaptiveInstructionController
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

	// Criar controlador adaptativo de instruções
	instructionController := NewAdaptiveInstructionController(detectProjectType(description), language, description)

	return &LayeredGenerator{
		provider:              provider,
		maxTokens:             maxTokens,
		language:              language,
		projectName:           projectName,
		description:           description,
		llmsContext:           llmsContext,
		outputLayers:          make([]LayerResult, 0),
		instructionController: instructionController,
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
	// Construir prompt adaptativo usando o controlador de instruções
	basePrompt := fmt.Sprintf(`SISTEMA CAMALEÃO - PLANEJAMENTO ADAPTATIVO DE CAMADAS

🎯 CONTEXTO DO PROJETO:
- NOME: %s
- LINGUAGEM: %s
- OBJETIVO FINAL: %s

🎭 MISSÃO CAMALEÃO:
Analise o projeto e determine as camadas MÍNIMAS necessárias para atingir o objetivo final.
ADAPTE-SE ao contexto específico. Seja um camaleão que:
- MUDA sua estrutura baseada no propósito
- ELIMINA camadas desnecessárias
- MANTÉM foco no objetivo final
- PRIORIZA funcionalidade sobre estrutura

IMPORTANTE: Siga rigorosamente as instruções de escopo e restrições.

📋 FORMATO DE RESPOSTA:
{
  "layers": [
    {
      "name": "nome-da-camada",
      "description": "descrição específica da camada",
      "priority": 1,
      "focus": ["elemento1", "elemento2"]
    }
  ]
}

🎪 REGRAS DE PLANEJAMENTO CAMALEÃO:
- Máximo 6 camadas (menos é melhor)
- Cada camada deve ter foco específico e justificado
- Prioridade 1 = mais crítica para o objetivo
- Focus deve listar elementos essenciais da camada
- ELIMINE camadas desnecessárias para o propósito
- ADAPTE as camadas ao contexto específico do projeto
- MANTENHA consistência com o objetivo final declarado`, lg.projectName, lg.language, lg.description)

	// Aplicar controle adaptativo de instruções
	adaptivePrompt := lg.instructionController.BuildAdaptivePrompt(basePrompt)

	// Adicionar contexto se disponível, mas limitado e relevante
	if lg.llmsContext != nil && lg.llmsContext.HasLLMsFile {
		contextSummary := lg.llmsContext.GetProjectDescription()
		if len(contextSummary) > 1200 {
			// Extrair apenas partes mais relevantes
			lines := strings.Split(contextSummary, "\n")
			relevantLines := make([]string, 0)

			for _, line := range lines {
				if strings.Contains(strings.ToLower(line), "objetivo") ||
					strings.Contains(strings.ToLower(line), "propósito") ||
					strings.Contains(strings.ToLower(line), "funcionalidade") ||
					strings.Contains(strings.ToLower(line), "requisito") ||
					strings.Contains(strings.ToLower(line), "feature") {
					relevantLines = append(relevantLines, line)
				}
			}

			if len(relevantLines) > 0 {
				contextSummary = strings.Join(relevantLines, "\n")
			} else {
				contextSummary = contextSummary[:1200]
			}
		}
		adaptivePrompt += fmt.Sprintf("\n\n📝 CONTEXTO RELEVANTE:\n%s", contextSummary)
	}

	response, err := lg.provider.GenerateContent(adaptivePrompt)
	if err != nil {
		return nil, err
	}

	// Validar conformidade com as instruções
	compliance, err := lg.instructionController.ValidateInstructionCompliance(response)
	if err != nil {
		fmt.Printf("⚠️  Erro ao validar conformidade do planejamento: %v\n", err)
	} else if !compliance.IsCompliant {
		fmt.Printf("🔴 Planejamento NÃO CONFORME (score: %.1f%%):\n", compliance.ComplianceScore)
		for _, violation := range compliance.ViolatedRules {
			fmt.Printf("   • ⚠️  Violação: %s\n", violation)
		}
		for _, missing := range compliance.MissingRequirements {
			fmt.Printf("   • ❌ Faltando: %s\n", missing)
		}

		// Se muito baixo, usar fallback
		if compliance.ComplianceScore < 30 {
			fmt.Printf("🚨 Score muito baixo, usando planejamento adaptativo fallback...\n")
			return lg.getDefaultLayersWithAdaptation(), nil
		}
	} else {
		fmt.Printf("✅ Planejamento CONFORME - Sistema Camaleão validou (score: %.1f%%)\n", compliance.ComplianceScore)
	}

	// Extrair JSON da resposta
	jsonContent := extractJSONContent(response)

	var planResponse struct {
		Layers []LayerPlan `json:"layers"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &planResponse); err != nil {
		fmt.Printf("⚠️  Erro ao processar JSON do planejamento, usando fallback adaptativo...\n")
		return lg.getDefaultLayersWithAdaptation(), nil
	}

	// Validar se as camadas fazem sentido
	if len(planResponse.Layers) == 0 {
		fmt.Printf("⚠️  Nenhuma camada planejada, usando fallback adaptativo...\n")
		return lg.getDefaultLayersWithAdaptation(), nil
	}

	// Validar e filtrar camadas baseado no contexto
	validatedLayers := lg.validateAndFilterLayers(planResponse.Layers)

	// Ordenar por prioridade
	sort.Slice(validatedLayers, func(i, j int) bool {
		return validatedLayers[i].Priority < validatedLayers[j].Priority
	})

	return validatedLayers, nil
}

// validateAndFilterLayers valida e filtra camadas baseado no contexto
func (lg *LayeredGenerator) validateAndFilterLayers(layers []LayerPlan) []LayerPlan {
	profile := lg.instructionController.GetInstructionProfile()
	validatedLayers := make([]LayerPlan, 0)

	// Limitar número de camadas baseado no escopo
	maxLayers := 6
	if profile.ScopeControl == "minimal" {
		maxLayers = 3
	} else if profile.ScopeControl == "comprehensive" {
		maxLayers = 6
	}

	// Camadas essenciais que sempre devem existir
	coreFound := false

	for _, layer := range layers {
		// Verificar se a camada é válida para o contexto
		if lg.isLayerValidForContext(layer) {
			validatedLayers = append(validatedLayers, layer)

			if layer.Name == "core" {
				coreFound = true
			}
		}

		// Limitar número máximo de camadas
		if len(validatedLayers) >= maxLayers {
			break
		}
	}

	// Garantir que sempre há uma camada core
	if !coreFound {
		coreLayer := LayerPlan{
			Name:        "core",
			Description: "Estrutura básica e configurações essenciais",
			Priority:    1,
			Focus:       []string{"main", "config", "setup"},
		}
		validatedLayers = append([]LayerPlan{coreLayer}, validatedLayers...)
	}

	return validatedLayers
}

// isLayerValidForContext verifica se uma camada é válida para o contexto atual
func (lg *LayeredGenerator) isLayerValidForContext(layer LayerPlan) bool {
	// Verificar exclusões explícitas
	if lg.instructionController.Adaptations["exclude_tests"] == true &&
		strings.Contains(strings.ToLower(layer.Name), "test") {
		return false
	}

	if lg.instructionController.Adaptations["exclude_docker"] == true &&
		(strings.Contains(strings.ToLower(layer.Name), "deploy") ||
			strings.Contains(strings.ToLower(layer.Name), "docker")) {
		return false
	}

	if lg.instructionController.Adaptations["exclude_frontend"] == true &&
		strings.Contains(strings.ToLower(layer.Name), "frontend") {
		return false
	}

	if lg.instructionController.Adaptations["exclude_api"] == true &&
		strings.Contains(strings.ToLower(layer.Name), "api") {
		return false
	}

	if lg.instructionController.Adaptations["exclude_database"] == true &&
		(strings.Contains(strings.ToLower(layer.Name), "database") ||
			strings.Contains(strings.ToLower(layer.Name), "data")) {
		return false
	}

	// Para escopo mínimo, permitir apenas camadas essenciais
	if lg.instructionController.Scope == "minimal" {
		essentialLayers := []string{"core", "business", "main"}
		for _, essential := range essentialLayers {
			if strings.Contains(strings.ToLower(layer.Name), essential) {
				return true
			}
		}
		return false
	}

	return true
}

// getDefaultLayersWithAdaptation retorna camadas padrão adaptadas ao contexto
func (lg *LayeredGenerator) getDefaultLayersWithAdaptation() []LayerPlan {
	profile := lg.instructionController.GetInstructionProfile()

	// Camadas base
	defaultLayers := []LayerPlan{
		{
			Name:        "core",
			Description: "Estrutura básica e arquivos de configuração",
			Priority:    1,
			Focus:       []string{"main", "config", "package"},
		},
	}

	// Adaptar camadas com base no perfil de instruções
	if profile.ScopeControl != "minimal" {
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "business",
			Description: "Lógica de negócio e modelos",
			Priority:    2,
			Focus:       []string{"models", "services", "utils"},
		})
	}

	// Adicionar camadas baseadas em adaptações específicas
	if lg.instructionController.Adaptations["include_api"] == true {
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "api",
			Description: "Endpoints e controladores",
			Priority:    3,
			Focus:       []string{"controllers", "routes", "middleware"},
		})
	}

	if lg.instructionController.Adaptations["include_frontend"] == true {
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "frontend",
			Description: "Componentes de interface e assets",
			Priority:    3,
			Focus:       []string{"components", "pages", "styles"},
		})
	}

	if lg.instructionController.Adaptations["include_tests"] == true {
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "tests",
			Description: "Testes unitários e de integração",
			Priority:    4,
			Focus:       []string{"tests", "mocks", "fixtures"},
		})
	}

	if lg.instructionController.Adaptations["include_docker"] == true {
		defaultLayers = append(defaultLayers, LayerPlan{
			Name:        "deployment",
			Description: "Configuração de deploy e containerização",
			Priority:    5,
			Focus:       []string{"docker", "ci-cd", "deployment"},
		})
	}

	// Adaptar para linguagens específicas apenas se não for escopo mínimo
	if profile.ScopeControl != "minimal" {
		switch strings.ToLower(lg.language) {
		case "javascript", "typescript":
			// Adicionar camada frontend apenas se não foi explicitamente incluída
			if lg.instructionController.Adaptations["include_frontend"] != true {
				defaultLayers = append(defaultLayers, LayerPlan{
					Name:        "frontend",
					Description: "Componentes de interface e assets",
					Priority:    3,
					Focus:       []string{"components", "pages", "styles"},
				})
			}
		case "go":
			// Adaptar foco para Go
			for i := range defaultLayers {
				if defaultLayers[i].Name == "business" {
					defaultLayers[i].Focus = append(defaultLayers[i].Focus, "handlers", "internal")
				}
			}
		case "python":
			// Adaptar foco para Python
			for i := range defaultLayers {
				if defaultLayers[i].Name == "business" {
					defaultLayers[i].Focus = append(defaultLayers[i].Focus, "modules", "packages")
				}
			}
		}
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

	// NOVA FUNCIONALIDADE: Aplicar filtro de Ultimate Goal na camada
	if lg.instructionController != nil && lg.instructionController.UltimateGoalController != nil {
		filteredResponse, filterErr := lg.instructionController.UltimateGoalController.FilterGeneratedContent(response)
		if filterErr != nil {
			fmt.Printf("⚠️  Erro ao filtrar camada %s: %v\n", plan.Name, filterErr)
			// Continuar com resposta original se filtro falhar
		} else {
			response = filteredResponse
			fmt.Printf("🎯 Camada %s filtrada baseada no objetivo final\n", plan.Name)
		}
	}

	// Processar resposta da camada
	return lg.parseLayerResponse(plan.Name, response)
}

// buildLayerPrompt constrói o prompt para uma camada específica
func (lg *LayeredGenerator) buildLayerPrompt(plan LayerPlan, previousLayers []LayerResult) string {
	var prompt strings.Builder

	// Prompt base para a camada com foco no objetivo final
	basePrompt := fmt.Sprintf(`SISTEMA CAMALEÃO - GERAÇÃO ADAPTATIVA DE CAMADA

CONTEXTO DO PROJETO:
- NOME: %s
- LINGUAGEM: %s
- OBJETIVO FINAL: %s

CAMADA ATUAL - FOCO ABSOLUTO:
- Nome: %s
- Descrição: %s
- Elementos específicos: %v

INSTRUÇÃO CRÍTICA - PRINCÍPIO CAMALEÃO:
Esta camada deve conter EXCLUSIVAMENTE arquivos relacionados a: %s

🎯 ULTIMATE GOAL FOCUS ATIVADO:
- Gere APENAS arquivos que contribuem para: %s
- Cada arquivo deve ter JUSTIFICATIVA clara no contexto do objetivo final
- REJEITE tentativas de adicionar recursos não solicitados
- MANTENHA laser focus no propósito específico declarado
- ELIMINE qualquer componente que não seja DIRETAMENTE necessário

ADAPTE-SE precisamente ao contexto e propósito. Seja um camaleão que:
- MUDA sua estrutura baseada no objetivo final
- ELIMINA componentes desnecessários
- MANTÉM consistência e coerência
- FOCA no valor entregue

`, lg.projectName, lg.language, lg.description, plan.Name, plan.Description, plan.Focus, strings.Join(plan.Focus, ", "), lg.description)

	// Aplicar controle adaptativo de instruções
	adaptivePrompt := lg.instructionController.BuildAdaptivePrompt(basePrompt)
	prompt.WriteString(adaptivePrompt)

	// Adicionar contexto das camadas anteriores de forma mais inteligente
	if len(previousLayers) > 0 {
		prompt.WriteString("\n🔗 CONTEXTO DE CAMADAS ANTERIORES:\n")
		allCreatedFiles := make([]string, 0)
		layerSummary := make([]string, 0)

		for _, layer := range previousLayers {
			layerSummary = append(layerSummary, fmt.Sprintf("- %s: %d arquivos", layer.LayerName, len(layer.Files)))
			for fileName := range layer.Files {
				allCreatedFiles = append(allCreatedFiles, fileName)
			}
		}

		prompt.WriteString(strings.Join(layerSummary, "\n"))
		prompt.WriteString("\n")

		if len(allCreatedFiles) > 0 {
			prompt.WriteString("\n🚫 ARQUIVOS JÁ EXISTENTES (NUNCA RECRIAR):\n")
			for _, fileName := range allCreatedFiles {
				prompt.WriteString(fmt.Sprintf("- %s\n", fileName))
			}
		}
		prompt.WriteString("\n")
	}

	// Adicionar instruções específicas por camada adaptadas ao contexto
	profile := lg.instructionController.GetInstructionProfile()
	prompt.WriteString(lg.buildLayerSpecificInstructions(plan, profile))

	// Adicionar contexto limitado do projeto com inteligência
	if lg.llmsContext != nil && lg.llmsContext.HasLLMsFile {
		summary := lg.llmsContext.GetProjectDescription()
		if len(summary) > 800 {
			// Extrair apenas partes relevantes para esta camada
			summary = lg.extractRelevantContext(summary, plan)
		}
		prompt.WriteString(fmt.Sprintf("📋 CONTEXTO RELEVANTE DO PROJETO:\n%s\n\n", summary))
	}

	prompt.WriteString(`📝 FORMATO DE SAÍDA OBRIGATÓRIO:
{
  "layer_name": "` + plan.Name + `",
  "description": "descrição específica da camada implementada",
  "directories": ["lista", "de", "diretórios", "necessários"],
  "files": {
    "arquivo.ext": {
      "content": "conteúdo funcional e realista"
    }
  },
  "dependencies": ["apenas", "deps", "essenciais"],
  "next_steps": ["próximos", "passos", "específicos"]
}

⚡ INSTRUÇÕES FINAIS - MODO CAMALEÃO:
1. 🎯 FOQUE EXCLUSIVAMENTE nos elementos desta camada: ` + strings.Join(plan.Focus, ", ") + `
2. 🚫 NUNCA recrie arquivos listados como "JÁ EXISTENTES"
3. 📂 Se um arquivo já existe, PULE e NÃO o inclua nesta camada
4. ✨ Crie APENAS arquivos novos específicos e necessários
5. 💡 Use conteúdo realista, funcional e diretamente relacionado ao objetivo
6. 🔗 Mantenha coerência total com o projeto e camadas anteriores
7. 📄 Retorne APENAS JSON válido e bem formatado
8. 🎭 ADAPTE-SE como um camaleão: mude baseado no contexto específico
9. 🚀 PRIORIZE funcionalidade sobre estrutura elaborada desnecessária
10. 🎪 OBEDEÇA RIGOROSAMENTE às instruções de escopo e restrições`)

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

// buildLayerSpecificInstructions constrói instruções específicas para cada camada
func (lg *LayeredGenerator) buildLayerSpecificInstructions(plan LayerPlan, profile InstructionProfile) string {
	var instructions strings.Builder

	// Instruções base por camada com foco adaptativo
	instructions.WriteString(fmt.Sprintf("🎯 INSTRUÇÕES ESPECÍFICAS - CAMADA %s:\n", strings.ToUpper(plan.Name)))

	switch plan.Name {
	case "core":
		instructions.WriteString(`
CAMADA CORE - FUNDAÇÃO ADAPTATIVA:
✅ CRIE APENAS:
- Arquivo principal/entry point (main.js, index.js, app.js, main.go)
- Configurações essenciais (package.json, go.mod, requirements.txt)
- Variáveis de ambiente (.env.example)
- Estrutura básica de diretórios
- Configurações core do projeto
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: documentação extensa, exemplos, testes, middleware, rotas específicas\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: logging básico, error handling, configurações avançadas\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: configurações básicas, error handling simples\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Adapte a estrutura core baseada no propósito específico do projeto\n")

	case "business":
		instructions.WriteString(`
CAMADA BUSINESS - LÓGICA ADAPTATIVA:
✅ CRIE APENAS:
- Modelos de dados (models/ ou types/)
- Serviços de negócio (services/ ou business/)
- Utilitários específicos (utils/ ou helpers/)
- Interfaces e contratos
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: middleware complexo, validações extensas, padrões avançados\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: validações robustas, padrões de design, error handling específico\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: validações básicas, error handling essencial\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Modele a lógica baseada no domínio específico do projeto\n")

	case "api":
		instructions.WriteString(`
CAMADA API - INTERFACE ADAPTATIVA:
✅ CRIE APENAS:
- Controladores/handlers (controllers/ ou handlers/)
- Rotas HTTP (routes/ ou endpoints/)
- Middleware de API
- Validadores de entrada
- Serializadores/formatadores
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: documentação API extensa, middleware complexo, rate limiting\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: documentação OpenAPI, middleware avançado, rate limiting\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: middleware básico, validação de entrada\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Defina endpoints baseados nos requisitos específicos do projeto\n")

	case "frontend":
		instructions.WriteString(`
CAMADA FRONTEND - INTERFACE VISUAL ADAPTATIVA:
✅ CRIE APENAS:
- Componentes de interface (components/)
- Páginas/views (pages/ ou views/)
- Estilos (styles/ ou css/)
- Assets básicos (images/, icons/)
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: animações complexas, múltiplos temas, assets pesados\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: sistema de temas, animações, responsividade avançada\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: responsividade básica, estilos organizados\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Crie interface baseada no público-alvo e propósito do projeto\n")

	case "tests":
		instructions.WriteString(`
CAMADA TESTS - VALIDAÇÃO ADAPTATIVA:
✅ CRIE APENAS:
- Testes unitários (.test.js, .spec.js, _test.go)
- Configuração de testes (jest.config.js, vitest.config.js)
- Mocks e fixtures (mocks/, fixtures/)
- Helpers de teste (test-utils/)
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: testes e2e, coverage complexo, múltiplos frameworks\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: testes e2e, coverage reports, performance tests\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: testes unitários essenciais, basic coverage\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Teste apenas o que é crítico para o funcionamento do projeto\n")

	case "deployment":
		instructions.WriteString(`
CAMADA DEPLOYMENT - INFRAESTRUTURA ADAPTATIVA:
✅ CRIE APENAS:
- Dockerfile (se containerização necessária)
- docker-compose.yml (se ambiente local necessário)
- Scripts de deploy (deploy.sh, deploy.yml)
- Configurações de CI/CD (.github/workflows/, .gitlab-ci.yml)
`)
		if profile.ScopeControl == "minimal" {
			instructions.WriteString("❌ MODO MÍNIMO - NÃO CRIE: orquestração complexa, múltiplos ambientes, scripts avançados\n")
		} else if profile.ScopeControl == "comprehensive" {
			instructions.WriteString("✅ MODO COMPLETO - ADICIONE: múltiplos ambientes, monitoring, backup strategies\n")
		} else {
			instructions.WriteString("⚖️ MODO PADRÃO - INCLUA: configuração básica de produção\n")
		}
		instructions.WriteString("🎭 CAMALEÃO: Configure deployment baseado no ambiente de destino do projeto\n")

	default:
		instructions.WriteString(fmt.Sprintf(`
CAMADA %s - IMPLEMENTAÇÃO ADAPTATIVA:
✅ CRIE APENAS:
- Arquivos diretamente relacionados a: %s
- Funcionalidades específicas desta camada
- Configurações necessárias para o funcionamento
`, strings.ToUpper(plan.Name), strings.Join(plan.Focus, ", ")))
		instructions.WriteString("🎭 CAMALEÃO: Adapte completamente baseado no contexto específico\n")
	}

	// Adicionar adaptações específicas baseadas no contexto
	if profile.StrictnessLevel >= 8 {
		instructions.WriteString("\n🔒 MODO ULTRA-RESTRITIVO ATIVADO:\n")
		instructions.WriteString("- Seja extremamente seletivo - apenas o ABSOLUTAMENTE essencial\n")
		instructions.WriteString("- Rejeite qualquer tentativa de adicionar recursos extras\n")
		instructions.WriteString("- Cada arquivo deve ter justificativa CLARA e ESPECÍFICA\n")
		instructions.WriteString("- Priorize funcionalidade direta sobre estrutura elaborada\n")
	}

	// Adicionar restrições específicas por adaptação
	if lg.instructionController.Scope == "minimal" {
		instructions.WriteString("\n⚡ RESTRIÇÕES DE ESCOPO ULTRA-MÍNIMO:\n")
		instructions.WriteString("- Evite QUALQUER funcionalidade não essencial\n")
		instructions.WriteString("- Mantenha arquivos simples, diretos e focados\n")
		instructions.WriteString("- NÃO adicione comentários extensos ou documentação elaborada\n")
		instructions.WriteString("- Foque APENAS no que é absolutamente necessário para o funcionamento\n")
		instructions.WriteString("- Elimine qualquer abstração desnecessária\n")
	}

	// Adaptações específicas baseadas no projeto
	if lg.instructionController.Adaptations["chameleon_focus"] == true {
		instructions.WriteString("\n🎯 FOCO CAMALEÃO LASER ATIVADO:\n")
		instructions.WriteString("- CONCENTRE-SE exclusivamente no objetivo final desta camada\n")
		instructions.WriteString("- ELIMINE qualquer elemento que não contribua diretamente\n")
		instructions.WriteString("- ADAPTE estrutura e conteúdo ao propósito específico\n")
		instructions.WriteString("- MANTENHA coerência absoluta com o objetivo do projeto\n")
	}

	return instructions.String()
}

// parseLayerResponse processa a resposta de uma camada
func (lg *LayeredGenerator) parseLayerResponse(layerName, response string) (*LayerResult, error) {
	jsonContent := extractJSONContent(response)

	// Validar conformidade com as instruções antes de processar
	compliance, err := lg.instructionController.ValidateInstructionCompliance(jsonContent)
	if err != nil {
		fmt.Printf("⚠️  Erro ao validar conformidade da camada %s: %v\n", layerName, err)
	} else if !compliance.IsCompliant {
		fmt.Printf("🔴 Camada %s NÃO CONFORME (score: %.1f%%):\n", layerName, compliance.ComplianceScore)
		for _, violation := range compliance.ViolatedRules {
			fmt.Printf("   • ⚠️  Violação: %s\n", violation)
		}
		for _, missing := range compliance.MissingRequirements {
			fmt.Printf("   • ❌ Faltando: %s\n", missing)
		}
		for _, deviation := range compliance.ScopeDeviations {
			fmt.Printf("   • 🚫 Desvio de escopo: %s\n", deviation)
		}

		// Se o nível de rigidez for alto, falhar em caso de não conformidade
		profile := lg.instructionController.GetInstructionProfile()
		if profile.StrictnessLevel >= 8 && compliance.ComplianceScore < profile.QualityThreshold {
			return nil, fmt.Errorf("🚨 CAMADA %s REJEITADA - Não conformidade crítica (score: %.1f%%, requerido: %.1f%%). Sistema Camaleão detectou incompatibilidade com objetivo final",
				layerName, compliance.ComplianceScore, profile.QualityThreshold)
		}

		// Para níveis menores, mostrar avisos mas continuar
		fmt.Printf("⚠️  Continuando com avisos (modo tolerante)...\n")
	} else {
		fmt.Printf("✅ Camada %s CONFORME - Sistema Camaleão validou (score: %.1f%%)\n", layerName, compliance.ComplianceScore)
	}

	// Validar a estrutura da camada
	validation := ValidateProjectStructure(jsonContent, lg.language)
	if !validation.IsValid {
		fmt.Printf("🔴 Camada %s apresenta problemas estruturais:\n", layerName)
		for _, issue := range validation.Issues {
			fmt.Printf("   • ⚠️  %s\n", issue)
		}
		fmt.Printf("📊 Pontuação estrutural: %.1f/100\n", validation.Score)

		// Se o score for muito baixo, falhar
		if validation.Score < 30 {
			return nil, fmt.Errorf("🚨 CAMADA %s REJEITADA - Falha estrutural crítica (score: %.1f/100). Sistema Camaleão detectou estrutura inadequada para o objetivo", layerName, validation.Score)
		}

		// Para scores baixos mas não críticos, mostrar avisos
		if validation.Score < 60 {
			fmt.Printf("⚠️  Continuando com avisos estruturais...\n")
		}
	} else {
		fmt.Printf("✅ Camada %s estruturalmente válida (score: %.1f/100)\n", layerName, validation.Score)
		if len(validation.Suggestions) > 0 {
			fmt.Printf("💡 Sugestões de melhoria para camada %s:\n", layerName)
			for _, suggestion := range validation.Suggestions {
				fmt.Printf("   • 💡 %s\n", suggestion)
			}
		}
	}

	// Processar JSON
	var layerResult LayerResult
	if err := json.Unmarshal([]byte(jsonContent), &layerResult); err != nil {
		return nil, fmt.Errorf("🚨 ERRO ao processar JSON da camada %s: %v", layerName, err)
	}

	// Garantir que o nome da camada está correto
	layerResult.LayerName = layerName

	// Validação adicional específica do sistema camaleão
	if err := lg.validateChameleonCompliance(&layerResult); err != nil {
		return nil, fmt.Errorf("🚨 CAMADA %s REJEITADA - Falha na validação Camaleão: %v", layerName, err)
	}

	return &layerResult, nil
}

// validateChameleonCompliance valida se a camada está em conformidade com o sistema camaleão
func (lg *LayeredGenerator) validateChameleonCompliance(layer *LayerResult) error {
	profile := lg.instructionController.GetInstructionProfile()

	// Verificar se há arquivos desnecessários para escopo mínimo
	if profile.ScopeControl == "minimal" {
		if len(layer.Files) > 8 {
			return fmt.Errorf("escopo mínimo violado: %d arquivos criados (máximo recomendado: 8)", len(layer.Files))
		}

		// Verificar se há arquivos típicos de escopo expandido
		unnecessaryPatterns := []string{
			"example", "sample", "demo", "template", "boilerplate",
			"readme", "license", "changelog", "contributing", "docs",
		}

		for filename := range layer.Files {
			for _, pattern := range unnecessaryPatterns {
				if strings.Contains(strings.ToLower(filename), pattern) {
					return fmt.Errorf("arquivo desnecessário para escopo mínimo: %s", filename)
				}
			}
		}
	}

	// Verificar exclusões específicas
	if lg.instructionController.Adaptations["exclude_tests"] == true {
		for filename := range layer.Files {
			if strings.Contains(strings.ToLower(filename), "test") ||
				strings.Contains(strings.ToLower(filename), "spec") {
				return fmt.Errorf("arquivo de teste encontrado apesar de exclusão explícita: %s", filename)
			}
		}
	}

	if lg.instructionController.Adaptations["exclude_docker"] == true {
		for filename := range layer.Files {
			if strings.Contains(strings.ToLower(filename), "docker") ||
				strings.Contains(strings.ToLower(filename), "compose") {
				return fmt.Errorf("arquivo Docker encontrado apesar de exclusão explícita: %s", filename)
			}
		}
	}

	// Verificar se há consistência nos nomes dos arquivos
	if layer.LayerName == "core" {
		hasMainFile := false
		for filename := range layer.Files {
			if strings.Contains(strings.ToLower(filename), "main") ||
				strings.Contains(strings.ToLower(filename), "index") ||
				strings.Contains(strings.ToLower(filename), "app") {
				hasMainFile = true
				break
			}
		}
		if !hasMainFile && profile.ScopeControl != "minimal" {
			return fmt.Errorf("camada core deve ter pelo menos um arquivo principal (main, index, app)")
		}
	}

	// Verificar se há dependências coerentes
	if len(layer.Dependencies) > 20 && profile.ScopeControl == "minimal" {
		return fmt.Errorf("muitas dependências para escopo mínimo: %d (máximo recomendado: 20)", len(layer.Dependencies))
	}

	return nil
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

// detectProjectType detecta o tipo de projeto baseado na descrição
func detectProjectType(description string) string {
	desc := strings.ToLower(description)

	if containsAny(desc, []string{"api", "rest", "servidor", "backend", "microserviço", "microservice"}) {
		return "backend_api"
	}

	if containsAny(desc, []string{"frontend", "ui", "interface", "react", "vue", "angular", "web"}) {
		return "frontend_web"
	}

	if containsAny(desc, []string{"cli", "command", "linha de comando", "terminal", "script"}) {
		return "cli_tool"
	}

	if containsAny(desc, []string{"biblioteca", "library", "package", "lib", "sdk"}) {
		return "library"
	}

	if containsAny(desc, []string{"bot", "chatbot", "automação", "automation"}) {
		return "automation_bot"
	}

	if containsAny(desc, []string{"jogo", "game", "gaming"}) {
		return "game"
	}

	if containsAny(desc, []string{"mobile", "app", "android", "ios"}) {
		return "mobile_app"
	}

	if containsAny(desc, []string{"dashboard", "admin", "painel", "gerenciamento"}) {
		return "admin_dashboard"
	}

	// Tipo padrão
	return "general_application"
}

// extractRelevantContext extrai partes do contexto relevantes para uma camada específica
func (lg *LayeredGenerator) extractRelevantContext(fullContext string, plan LayerPlan) string {
	var relevantParts []string
	contextLines := strings.Split(fullContext, "\n")

	// Palavras-chave relevantes para cada camada
	relevantKeywords := make([]string, 0)
	relevantKeywords = append(relevantKeywords, plan.Focus...)
	relevantKeywords = append(relevantKeywords, strings.ToLower(plan.Name))

	// Adicionar palavras-chave específicas por camada
	switch plan.Name {
	case "core":
		relevantKeywords = append(relevantKeywords, "main", "config", "setup", "init", "base")
	case "business":
		relevantKeywords = append(relevantKeywords, "logic", "service", "model", "data", "business")
	case "api":
		relevantKeywords = append(relevantKeywords, "route", "endpoint", "controller", "handler", "http")
	case "frontend":
		relevantKeywords = append(relevantKeywords, "component", "ui", "interface", "page", "view")
	case "tests":
		relevantKeywords = append(relevantKeywords, "test", "spec", "unit", "integration", "mock")
	case "deployment":
		relevantKeywords = append(relevantKeywords, "docker", "deploy", "container", "ci", "cd")
	}

	// Extrair linhas relevantes
	for _, line := range contextLines {
		lineLower := strings.ToLower(line)
		for _, keyword := range relevantKeywords {
			if strings.Contains(lineLower, strings.ToLower(keyword)) {
				relevantParts = append(relevantParts, line)
				break
			}
		}
	}

	// Se não encontrou nada relevante, retornar o início do contexto
	if len(relevantParts) == 0 {
		if len(fullContext) > 400 {
			return fullContext[:400] + "... (contexto adaptado para camada " + plan.Name + ")"
		}
		return fullContext
	}

	// Juntar partes relevantes
	result := strings.Join(relevantParts, "\n")

	// Limitar tamanho
	if len(result) > 600 {
		result = result[:600] + "... (contexto adaptado para camada " + plan.Name + ")"
	}

	return result
}

// GetInstructionController retorna o controlador de instruções (para testes)
func (lg *LayeredGenerator) GetInstructionController() *AdaptiveInstructionController {
	return lg.instructionController
}
