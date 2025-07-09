package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ktfth/zion/ai"
	"github.com/ktfth/zion/evaluator"
	"github.com/ktfth/zion/plugins"

	"github.com/spf13/cobra"
)

var language string
var projectName string
var description string
var aiProvider string
var apiKey string
var model string
var maxRetries int
var skipEvaluation bool
var enableAIEvaluation bool
var contextualMode bool

// scaffoldCmd define o comando "scaffold".
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Gera a estrutura de um projeto com a ajuda de AI",
	Long: `Gera a estrutura completa de um projeto utilizando inteligência artificial.
	
O comando inclui um sistema de retry inteligente que tenta até 3 vezes (configurável) 
caso a geração inicial falhe, garantindo maior robustez na criação dos projetos.

MODO CONTEXTUAL:
Quando executado em um diretório com arquivo 'llms.txt', o comando automaticamente
entra em modo contextual, lendo informações do projeto existente para gerar
recursos adicionais de forma mais precisa e contextualizada.

Exemplos:
  zion scaffold -l go -n meu-projeto -d "API REST com PostgreSQL"
  zion scaffold -l typescript -n webapp -d "App React com TypeScript" -r 5
  zion scaffold -l python -n ml-project -d "Projeto de Machine Learning" --retries 2
  zion scaffold --contextual -d "Adicionar sistema de autenticação" (em projeto existente)`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		// Verificar se estamos em modo contextual
		llmsContext, _ := ai.ReadLLMsContext(".")
		isContextualMode := contextualMode || llmsContext.HasLLMsFile

		if isContextualMode {
			if llmsContext.HasLLMsFile {
				fmt.Printf("\n� MODO CONTEXTUAL DETECTADO\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
				fmt.Printf("📖 Arquivo llms.txt encontrado\n")

				// Auto-detectar parâmetros do contexto
				if language == "" {
					detectedLang := llmsContext.DetectProjectLanguage()
					if detectedLang != "" {
						language = detectedLang
						fmt.Printf("🔍 Linguagem detectada: %s\n", language)
					}
				}

				if projectName == "" {
					currentDir, _ := os.Getwd()
					projectName = filepath.Base(currentDir)
					fmt.Printf("📁 Nome do projeto: %s\n", projectName)
				}

				if description == "" {
					projectDesc := llmsContext.GetProjectDescription()
					if projectDesc != "" {
						description = "Expandir projeto baseado no contexto existente"
						fmt.Printf("📝 Modo: Expansão contextual\n")
					}
				}

				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
			} else if contextualMode {
				// Modo contextual forçado sem llms.txt - usar análise estrutural
				fmt.Printf("\n🔍 MODO CONTEXTUAL - ANÁLISE ESTRUTURAL\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
				fmt.Printf("📄 Arquivo llms.txt não encontrado - usando análise da estrutura\n")

				// Ainda tentar detectar parâmetros básicos
				if projectName == "" {
					currentDir, _ := os.Getwd()
					projectName = filepath.Base(currentDir)
					fmt.Printf("📁 Nome do projeto: %s\n", projectName)
				}

				// Tentar detectar linguagem pelos arquivos existentes
				if language == "" {
					detectedLang := llmsContext.DetectProjectLanguage()
					if detectedLang != "" {
						language = detectedLang
						fmt.Printf("🔍 Linguagem detectada: %s\n", language)
					}
				}

				// Configurar modo de expansão estrutural
				if description == "" {
					description = "Expandir projeto baseado na análise da estrutura existente"
				}

				fmt.Printf("🧠 Modo: Análise estrutural inteligente\n")
				fmt.Printf("💡 Dica: Crie um arquivo llms.txt para melhor contextualização\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
			}
		}

		// Validação final de parâmetros obrigatórios
		if language == "" || projectName == "" {
			fmt.Printf("❌ Parâmetros obrigatórios em falta:\n")
			if language == "" {
				fmt.Printf("   - Linguagem (-l/--language) não pôde ser detectada automaticamente\n")
			}
			if projectName == "" {
				fmt.Printf("   - Nome do projeto (-n/--name) não pôde ser determinado\n")
			}
			fmt.Printf("\n💡 Use 'zion scaffold --help' para mais informações\n")
			if isContextualMode {
				fmt.Printf("💡 Em modo contextual, tente especificar os parâmetros manualmente\n")
			}
			os.Exit(1)
		}

		fmt.Printf("\n�🚀 Iniciando geração do projeto\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📦 Projeto: %s\n", projectName)
		fmt.Printf("🔧 Linguagem: %s\n", language)
		fmt.Printf("📝 Descrição: %s\n", description)
		if isContextualMode {
			if llmsContext.HasLLMsFile {
				fmt.Printf("🎯 Modo: Contextual (baseado em llms.txt)\n")
			} else {
				fmt.Printf("🎯 Modo: Contextual (análise estrutural)\n")
			}
		}
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Lista plugins ativos
		pluginsList := plugins.ListPlugins()
		if len(pluginsList) > 0 {
			fmt.Printf("🔌 Plugins ativos: %v\n\n", pluginsList)
		}

		fmt.Print("🤖 Gerando estrutura com IA...")

		// Preparar configurações do provider
		var response string
		var err error
		var attempt int

		// Lógica de retry para geração do projeto
		for attempt = 1; attempt <= maxRetries; attempt++ {
			if attempt > 1 {
				fmt.Printf("\n🔄 Tentativa %d/%d...", attempt, maxRetries)
			}

			if aiProvider != "" || apiKey != "" || model != "" {
				// Usar configurações customizadas
				response, err = ai.GenerateProjectScaffoldingWithProvider(language, projectName, description, pluginsList, aiProvider, apiKey, model)
			} else {
				// Usar configurações padrão
				response, err = ai.GenerateProjectScaffolding(language, projectName, description, pluginsList)
			}

			if err == nil {
				break
			}

			// Verificar se é erro de contexto e se ainda não tentamos camadas
			if ai.IsContextOverflowError(err) && attempt == 1 {
				fmt.Printf("\n⚠️  Erro de contexto detectado na tentativa %d\n", attempt)
				fmt.Printf("🔄 O sistema automaticamente tentará geração em camadas...\n")
				// O erro será tratado automaticamente dentro dos métodos de geração
			} else if attempt < maxRetries {
				fmt.Printf("\n⚠️  Erro na tentativa %d: %v", attempt, err)
				fmt.Printf("\n🔄 Tentando novamente em 2 segundos...")
				time.Sleep(2 * time.Second)
			}
		}

		if err != nil {
			fmt.Printf("\n❌ Falha após %d tentativas:\n%v\n", maxRetries, err)
			if response != "" {
				fmt.Printf("\nResposta da API:\n%s\n", response)
			}
			os.Exit(1)
		}
		fmt.Println(" ✅")

		// Avaliar qualidade do projeto antes de materializar
		if !skipEvaluation {
			if enableAIEvaluation {
				fmt.Print("🔍 Avaliando qualidade do projeto com IA...")
			} else {
				fmt.Print("🔍 Avaliando qualidade do projeto...")
			}

			evaluator := evaluator.NewProjectEvaluator()
			evaluator.EnableAIEvaluation(enableAIEvaluation)
			evalResult, err := evaluator.EvaluateProject(response, language)

			if err != nil {
				fmt.Printf("\n⚠️  Aviso: Erro na avaliação: %v\n", err)
				fmt.Println("Continuando com a criação do projeto...")
			} else {
				fmt.Println(" ✅")

				// Exibir resumo da avaliação
				fmt.Printf("\n📊 RESULTADO DA AVALIAÇÃO\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
				fmt.Printf("🎯 Score: %.1f/100\n", evalResult.Score)
				fmt.Printf("⭐ Qualidade: %s\n", getQualityDisplay(evalResult.Quality))
				fmt.Printf("⚠️  Issues: %d\n", len(evalResult.Issues))

				if !evalResult.Valid {
					fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
					fmt.Printf("❌ PROJETO COM ISSUES CRÍTICOS\n")

					// Mostrar apenas issues críticos
					for _, issue := range evalResult.Issues {
						if issue.Severity == "critical" {
							fmt.Printf("🚨 %s: %s\n", issue.Category, issue.Description)
							if issue.Suggestion != "" {
								fmt.Printf("   💡 %s\n", issue.Suggestion)
							}
						}
					}

					fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
					fmt.Printf("⚠️  Para ignorar a avaliação, use --skip-evaluation\n")
					fmt.Printf("📋 Para ver relatório completo, use: zion evaluate -f response.txt -l %s\n", language)
					os.Exit(1)
				}

				// Mostrar principais sugestões se score for baixo
				if evalResult.Score < 80 && len(evalResult.Suggestions) > 0 {
					fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
					fmt.Printf("💡 PRINCIPAIS SUGESTÕES:\n")
					for i, suggestion := range evalResult.Suggestions {
						if i >= 3 { // Mostrar só as 3 primeiras
							break
						}
						fmt.Printf("• %s\n", suggestion)
					}
				}

				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
			}
		}

		fmt.Print("📂 Criando estrutura do projeto...")

		// Lógica de retry para criação da estrutura
		var createAttempt int
		for createAttempt = 1; createAttempt <= maxRetries; createAttempt++ {
			if createAttempt > 1 {
				fmt.Printf("\n🔄 Tentativa de criação %d/%d...", createAttempt, maxRetries)
			}

			// Usar sistema unificado de criação de projetos
			err = ai.CreateProjectWithUnifiedStrategy(projectName, response, llmsContext, isContextualMode)

			if err == nil {
				break
			}

			if createAttempt < maxRetries {
				fmt.Printf("\n⚠️  Erro na tentativa %d: %v", createAttempt, err)
				fmt.Printf("\n🔄 Tentando novamente em 1 segundo...")
				time.Sleep(1 * time.Second)
			}
		}

		if err != nil {
			fmt.Printf("\n⚠️  Erro ao criar estrutura após %d tentativas, tentando método alternativo...\n", maxRetries)
			if isContextualMode {
				// Para modo contextual (com ou sem llms.txt), salvar em arquivo de saída
				err = ai.SaveRawResponse("zion_contextual_output", response)
				if err != nil {
					fmt.Printf("\n❌ Erro ao salvar resposta:\n%v\n", err)
					os.Exit(1)
				}
				fmt.Println("💡 Resposta salva em zion_contextual_output.md para revisão manual.")
			} else {
				// Para modo normal, usar método padrão
				err = ai.SaveRawResponse(projectName, response)
				if err != nil {
					fmt.Printf("\n❌ Erro ao salvar resposta:\n%v\n", err)
					os.Exit(1)
				}
				fmt.Println("💡 Resposta salva em README.md no diretório do projeto.")
			}
		}
		fmt.Println(" ✅")

		// Executa plugins
		if len(pluginsList) > 0 {
			fmt.Print("🔌 Executando plugins...")
			plugins.ExecutePlugins()
			fmt.Println(" ✅")
		}
		// INTEGRAÇÃO DOS SISTEMAS INTELIGENTES
		fmt.Print("🧠 Registrando aprendizado...")

		// 1. Registrar sessão no sistema de aprendizado
		learningPath := ai.GetLearningDirectory()
		ai.EnsureDirectoryExists(learningPath)
		learningSystem := ai.NewLearningSystem(learningPath)

		currentDir, _ := os.Getwd()
		session := ai.GenerationSession{
			ID:          fmt.Sprintf("session_%d", time.Now().Unix()),
			Timestamp:   startTime,
			Language:    language,
			ProjectType: ai.DetectProjectTypeFromPath(currentDir),
			Description: description,
			Success:     err == nil,
			Duration:    time.Since(startTime),
			TokensUsed:  0, // TODO: implementar contagem de tokens
			Provider:    aiProvider,
			LayeredMode: false, // TODO: detectar se foi usado modo em camadas
			Quality:     1.0,   // TODO: obter qualidade da avaliação
			UserRating:  0,     // Será preenchido pelo feedback
			Feedback:    "",
			Patterns:    []string{},
			Adaptations: make(map[string]interface{}),
			Context: ai.SessionContext{
				PrimaryLanguage: language,
				HasLLMsFile:     llmsContext.HasLLMsFile,
				UserIntent:      description,
			},
			Outcomes: ai.GenerationOutcomes{
				ExecutionTime: time.Since(startTime).Seconds(),
			},
			Improvements: []string{},
		}
		learningSystem.RecordSession(session)

		// 2. Atualizar cache inteligente com a resposta
		if err == nil {
			cachePath := ai.GetCacheDirectory()
			ai.EnsureDirectoryExists(cachePath)
			cache := ai.NewIntelligentCache(cachePath, 1024*1024*100) // 100MB

			cacheRequest := ai.CacheRequest{
				Language:    language,
				ProjectType: ai.DetectProjectTypeFromPath(currentDir),
				Description: description,
				Context:     make(map[string]interface{}),
			}

			cache.Put(cacheRequest, response, 1.0) // Qualidade alta se passou na avaliação
		}

		fmt.Println(" ✅")

		elapsedTime := time.Since(startTime)

		fmt.Printf("\n✨ Projeto criado com sucesso! ✨\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📁 Local: %s\n", projectName)
		fmt.Printf("⏱️  Tempo total: %.2f segundos\n", elapsedTime.Seconds())
		if attempt > 1 {
			fmt.Printf("🔄 Tentativas de geração: %d/%d\n", attempt, maxRetries)
		}
		if createAttempt > 1 {
			fmt.Printf("🔄 Tentativas de criação: %d/%d\n", createAttempt, maxRetries)
		}
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		fmt.Printf("💡 Para começar a desenvolver:\n")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   Consulte o README.md para instruções detalhadas\n\n")
	},
}

// getQualityDisplay retorna uma representação visual da qualidade
func getQualityDisplay(quality evaluator.QualityLevel) string {
	switch quality {
	case "excellent":
		return "🏆 EXCELENTE"
	case "good":
		return "✅ BOM"
	case "fair":
		return "⚡ REGULAR"
	case "poor":
		return "⚠️ RUIM"
	case "critical":
		return "❌ CRÍTICO"
	default:
		return "❓ INDEFINIDO"
	}
}

func init() {
	// Configura flags para o comando scaffold
	scaffoldCmd.Flags().StringVarP(&language, "language", "l", "", "Linguagem para o scaffold (ex: go, python, etc)")
	scaffoldCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	scaffoldCmd.Flags().StringVarP(&description, "description", "d", "", "Descrição objetiva da estrutura desejada")
	scaffoldCmd.Flags().StringVarP(&aiProvider, "provider", "p", "", "Provider de IA (gemini, gpt, openrouter)")
	scaffoldCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API Key do provider")
	scaffoldCmd.Flags().StringVarP(&model, "model", "m", "", "Modelo específico do provider")
	scaffoldCmd.Flags().IntVarP(&maxRetries, "retries", "r", 3, "Número máximo de tentativas em caso de falha (padrão: 3)")
	scaffoldCmd.Flags().BoolVar(&skipEvaluation, "skip-evaluation", false, "Pular avaliação de qualidade do projeto")
	scaffoldCmd.Flags().BoolVar(&enableAIEvaluation, "ai-evaluation", false, "Habilitar avaliação avançada por IA")
	scaffoldCmd.Flags().BoolVar(&contextualMode, "contextual", false, "Força modo contextual (usar arquivo llms.txt)")

	// Flags obrigatórias são validadas dentro do comando, não aqui
	// para permitir modo contextual sem todas as flags

	// Registra o comando scaffold no comando raiz
	rootCmd.AddCommand(scaffoldCmd)

	// Registra o comando de teste das camadas
	rootCmd.AddCommand(testLayersCmd)
}

// testLayersCmd testa o sistema de geração em camadas
var testLayersCmd = &cobra.Command{
	Use:   "test-layers",
	Short: "Testa o sistema de geração em camadas",
	Long: `Executa testes do sistema de geração em camadas para verificar
se a funcionalidade está operando corretamente.

Este comando é útil para:
- Verificar se o sistema de detecção de overflow funciona
- Testar o planejamento automático de camadas
- Validar a conversão entre formatos
- Diagnosticar problemas no sistema de camadas

Exemplo:
  zion test-layers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n🧪 TESTE DO SISTEMA DE GERAÇÃO EM CAMADAS\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Teste 1: Detecção de overflow
		fmt.Printf("1️⃣ Testando detecção de overflow...\n")

		// Criar um prompt muito grande
		largePrompt := "Este é um prompt de teste "
		for i := 0; i < 100000; i++ {
			largePrompt += "muito grande com muitas palavras "
		}

		isOverflow := ai.DetectContextOverflow(largePrompt, "openrouter")
		if isOverflow {
			fmt.Printf("✅ Detecção de overflow funcionando\n")
		} else {
			fmt.Printf("❌ Falha na detecção de overflow\n")
		}

		// Teste 2: Detecção de erro de contexto
		fmt.Printf("\n2️⃣ Testando detecção de erros de contexto...\n")

		testErrors := []string{
			"This endpoint's maximum context length is 200000 tokens",
			"API retornou status 400: token limit exceeded",
			"context too long",
		}

		for _, errMsg := range testErrors {
			testErr := fmt.Errorf(errMsg)
			if ai.IsContextOverflowError(testErr) {
				fmt.Printf("✅ Detectou erro de contexto: %s\n", errMsg[:50]+"}")
			} else {
				fmt.Printf("❌ Falha na detecção: %s\n", errMsg[:50]+"}")
			}
		}

		// Teste 3: Criação do gerador em camadas
		fmt.Printf("\n3️⃣ Testando criação do gerador em camadas...\n")

		llmsContext := &ai.LLMsContext{}
		layeredGen, err := ai.NewLayeredGenerator("go", "test-project", "Test project description", llmsContext)
		if err != nil {
			fmt.Printf("❌ Erro ao criar gerador: %v\n", err)
		} else {
			fmt.Printf("✅ Gerador em camadas criado com sucesso\n")
			_ = layeredGen
		}

		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("✅ Teste concluído\n\n")
	},
}

func init() {
	// Registra o comando test-layers no comando raiz
	rootCmd.AddCommand(testLayersCmd)
}
