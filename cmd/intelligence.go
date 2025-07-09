package cmd

import (
	"fmt"
	"time"

	"github.com/ktfth/zion/ai"
	"github.com/spf13/cobra"
)

// intelligenceCmd define o comando principal para funcionalidades inteligentes
var intelligenceCmd = &cobra.Command{
	Use:   "intelligence",
	Short: "Sistema inteligente do Zion - analytics, contexto e aprendizado",
	Long: `Acesso às funcionalidades inteligentes do Zion CLI, incluindo:
- Analytics e métricas de uso
- Análise de contexto de projetos
- Sistema de aprendizado
- Gerenciamento de cache inteligente
- Sistema de feedback`,
}

// analyticsCmd define o comando "analytics"
var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Exibe analytics e insights do sistema inteligente",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n📊 ZION ANALYTICS DASHBOARD\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Inicializar sistema de aprendizado
		learningPath := ai.GetLearningDirectory()
		if err := ai.EnsureDirectoryExists(learningPath); err != nil {
			fmt.Printf("❌ Erro ao criar diretório de aprendizado: %v\n", err)
			return
		}

		learningSystem := ai.NewLearningSystem(learningPath)
		stats := learningSystem.GetLearningStats()

		fmt.Printf("\n🧠 SISTEMA DE APRENDIZADO\n")
		fmt.Printf("─────────────────────────────────\n")
		fmt.Printf("📚 Total de sessões: %d\n", stats.TotalSessions)
		fmt.Printf("✅ Sessões bem-sucedidas: %d\n", stats.SuccessfulSessions)
		fmt.Printf("📊 Taxa de sucesso: %.1f%%\n", stats.SuccessRate*100)
		fmt.Printf("🎯 Padrões aprendidos: %d\n", len(stats.LearnedPatterns))

		if len(stats.TopLanguages) > 0 {
			fmt.Printf("\n💻 Linguagens mais usadas:\n")
			for i, lang := range stats.TopLanguages {
				if i >= 3 {
					break
				}
				fmt.Printf("  • %s\n", lang)
			}
		}

		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("✅ Relatório gerado com sucesso!\n")
	},
}

// contextAnalysisCmd define o comando para análise de contexto
var contextAnalysisCmd = &cobra.Command{
	Use:   "context [path]",
	Short: "Analisa o contexto de um projeto",
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("\n🔍 ANÁLISE DE CONTEXTO\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📁 Analisando: %s\n", targetPath)

		contextAnalyzer := ai.NewContextAnalyzer(targetPath)
		analysis, err := contextAnalyzer.AnalyzeContext()
		if err != nil {
			fmt.Printf("❌ Erro na análise: %v\n", err)
			return
		}

		fmt.Printf("\n📊 RESULTADOS\n")
		fmt.Printf("─────────────────────────────────\n")
		fmt.Printf("🎯 Tipo do projeto: %s\n", analysis.ProjectType)
		fmt.Printf("💻 Linguagem principal: %s\n", analysis.PrimaryLanguage)
		fmt.Printf("📁 Total de arquivos: %d\n", analysis.TotalFiles)
		fmt.Printf("📈 Complexidade: %.1f\n", analysis.ComplexityScore)

		if len(analysis.Technologies) > 0 {
			fmt.Printf("\n🔧 Tecnologias identificadas:\n")
			for i, tech := range analysis.Technologies {
				if i >= 5 {
					break
				}
				fmt.Printf("  • %s\n", tech)
			}
		}

		fmt.Printf("\n✅ Análise concluída!\n")
	},
}

// feedbackCmd define o comando para feedback
var feedbackSubmitCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Envia feedback sobre o Zion",
	Run: func(cmd *cobra.Command, args []string) {
		rating, _ := cmd.Flags().GetInt("rating")
		message, _ := cmd.Flags().GetString("message")

		fmt.Printf("\n💬 FEEDBACK DO USUÁRIO\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		if rating == 0 {
			fmt.Printf("⭐ Avalie sua experiência (1-5): ")
			fmt.Scanln(&rating)
		}

		if message == "" {
			fmt.Printf("💭 Comentário (opcional): ")
			fmt.Scanln(&message)
		}

		// Inicializar sistema de feedback
		feedbackPath := ai.GetFeedbackDirectory()
		if err := ai.EnsureDirectoryExists(feedbackPath); err != nil {
			fmt.Printf("❌ Erro ao criar diretório de feedback: %v\n", err)
			return
		}

		feedbackSystem := ai.NewFeedbackSystem(feedbackPath)

		// Criar feedback
		feedback := ai.UserFeedback{
			ID:        fmt.Sprintf("feedback_%d", time.Now().Unix()),
			Timestamp: time.Now(),
			Rating:    rating,
			Content:   message,
			Context: ai.FeedbackContext{
				Language:    ai.DetectProjectTypeFromPath("."),
				ProjectType: ai.DetectProjectTypeFromPath("."),
			},
		}

		err := feedbackSystem.SubmitFeedback(feedback)
		if err != nil {
			fmt.Printf("❌ Erro ao enviar feedback: %v\n", err)
			return
		}

		fmt.Printf("✅ Feedback enviado com sucesso!\n")
		fmt.Printf("🙏 Obrigado por contribuir para melhorar o Zion!\n")
	},
}

func init() {
	// Comando principal de inteligência
	rootCmd.AddCommand(intelligenceCmd)

	// Subcomandos
	intelligenceCmd.AddCommand(analyticsCmd)
	intelligenceCmd.AddCommand(contextAnalysisCmd)
	intelligenceCmd.AddCommand(feedbackSubmitCmd)

	// Flags para feedback
	feedbackSubmitCmd.Flags().IntP("rating", "r", 0, "Rating de 1 a 5")
	feedbackSubmitCmd.Flags().StringP("message", "m", "", "Mensagem de feedback")
}
