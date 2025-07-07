package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ktfth/zion/evaluator"
	"github.com/spf13/cobra"
)

var evaluateFile string
var evaluateLanguage string
var outputFormat string
var showDetails bool
var enableAIEvaluationCmd bool

// evaluateCmd define o comando "evaluate"
var evaluateCmd = &cobra.Command{
	Use:   "evaluate",
	Short: "Avalia a qualidade de um projeto antes de materializá-lo",
	Long: `Avalia a qualidade e aderência a boas práticas de um projeto gerado pela IA.

O comando analisa a estrutura do projeto, dependências, configurações e melhores práticas,
fornecendo um relatório detalhado com score de qualidade e sugestões de melhoria.

Exemplos:
  zion evaluate -f project.json -l go
  zion evaluate -f response.txt -l typescript --details
  zion evaluate -f project.json -l python --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		fmt.Printf("\n🔍 Iniciando avaliação do projeto\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📄 Arquivo: %s\n", evaluateFile)
		fmt.Printf("🔧 Linguagem: %s\n", evaluateLanguage)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Verificar se arquivo existe
		if _, err := os.Stat(evaluateFile); os.IsNotExist(err) {
			fmt.Printf("❌ Arquivo não encontrado: %s\n", evaluateFile)
			os.Exit(1)
		}

		// Ler conteúdo do arquivo
		fmt.Print("📖 Lendo arquivo...")
		content, err := os.ReadFile(evaluateFile)
		if err != nil {
			fmt.Printf("\n❌ Erro ao ler arquivo: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(" ✅")

		// Criar avaliador
		if enableAIEvaluationCmd {
			fmt.Print("🔧 Inicializando avaliador com IA...")
		} else {
			fmt.Print("🔧 Inicializando avaliador...")
		}
		eval := evaluator.NewProjectEvaluator()
		eval.EnableAIEvaluation(enableAIEvaluationCmd)
		fmt.Println(" ✅")

		// Executar avaliação
		if enableAIEvaluationCmd {
			fmt.Print("🔍 Analisando projeto com IA...")
		} else {
			fmt.Print("🔍 Analisando projeto...")
		}
		result, err := eval.EvaluateProject(string(content), evaluateLanguage)
		if err != nil {
			fmt.Printf("\n❌ Erro na avaliação: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(" ✅")

		// Exibir resultados
		elapsedTime := time.Since(startTime)

		fmt.Printf("\n")
		switch outputFormat {
		case "json":
			displayJSONResults(result)
		default:
			displayTextResults(result, eval, showDetails)
		}

		fmt.Printf("\n⏱️  Tempo de avaliação: %.2f segundos\n", elapsedTime.Seconds())

		// Indicar se o projeto pode ser materializado
		if result.Valid {
			fmt.Printf("\n✅ Projeto aprovado para materialização!\n")
			fmt.Printf("💡 Para materializar: use o comando scaffold com os mesmos parâmetros\n\n")
		} else {
			fmt.Printf("\n⚠️  Projeto precisa de correções antes da materialização\n")
			fmt.Printf("💡 Corrija os issues críticos e execute a avaliação novamente\n\n")
			os.Exit(1)
		}
	},
}

func displayTextResults(result *evaluator.EvaluationResult, eval *evaluator.ProjectEvaluator, showDetails bool) {
	// Gerar e exibir relatório
	report := eval.GenerateReport(result)
	fmt.Print(report)

	if !showDetails {
		fmt.Printf("💡 Use --details para ver análise completa por categoria\n")
	}
}

func displayJSONResults(result *evaluator.EvaluationResult) {
	// Exibir resultado em formato JSON (implementação básica)
	fmt.Printf("{\n")
	fmt.Printf("  \"valid\": %t,\n", result.Valid)
	fmt.Printf("  \"score\": %.1f,\n", result.Score)
	fmt.Printf("  \"quality\": \"%s\",\n", result.Quality)
	fmt.Printf("  \"issues\": %d,\n", len(result.Issues))
	fmt.Printf("  \"suggestions\": %d\n", len(result.Suggestions))
	fmt.Printf("}\n")
}

func init() {
	// Configura flags para o comando evaluate
	evaluateCmd.Flags().StringVarP(&evaluateFile, "file", "f", "", "Arquivo com a estrutura do projeto (JSON ou texto)")
	evaluateCmd.Flags().StringVarP(&evaluateLanguage, "language", "l", "", "Linguagem do projeto (go, python, typescript, etc)")
	evaluateCmd.Flags().StringVarP(&outputFormat, "format", "o", "text", "Formato de saída (text, json)")
	evaluateCmd.Flags().BoolVarP(&showDetails, "details", "d", false, "Mostrar análise detalhada por categoria")
	evaluateCmd.Flags().BoolVar(&enableAIEvaluationCmd, "ai-evaluation", false, "Habilitar avaliação avançada por IA")

	// Marcar flags obrigatórias
	evaluateCmd.MarkFlagRequired("file")
	evaluateCmd.MarkFlagRequired("language")

	// Registra o comando evaluate no comando raiz
	rootCmd.AddCommand(evaluateCmd)
}
