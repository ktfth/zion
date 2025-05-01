package cmd

import (
	"fmt"
	"os"
	"time"
	"github.com/ktfth/zion/ai"
	"github.com/ktfth/zion/plugins"

	"github.com/spf13/cobra"
)

var language string
var projectName string
var description string

// scaffoldCmd define o comando "scaffold".
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Gera a estrutura de um projeto com a ajuda de AI",
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		fmt.Printf("\n🚀 Iniciando geração do projeto\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📦 Projeto: %s\n", projectName)
		fmt.Printf("🔧 Linguagem: %s\n", language)
		fmt.Printf("📝 Descrição: %s\n", description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		// Lista plugins ativos
		pluginsList := plugins.ListPlugins()
		if len(pluginsList) > 0 {
			fmt.Printf("🔌 Plugins ativos: %v\n\n", pluginsList)
		}

		fmt.Print("🤖 Gerando estrutura com IA...")
		response, err := ai.GenerateProjectScaffolding(language, projectName, description, pluginsList)
		if err != nil {
			fmt.Printf("\n❌ Erro na geração da estrutura:\n%v\n", err)
			if response != "" {
				fmt.Printf("\nResposta da API:\n%s\n", response)
			}
			os.Exit(1)
		}
		fmt.Println(" ✅")

		fmt.Print("📂 Criando estrutura do projeto...")
		err = ai.ExtractAndCreateProject(projectName, response)
		if err != nil {
			fmt.Printf("\n⚠️  Erro ao criar estrutura padrão, tentando método alternativo...\n")
			err = ai.SaveRawResponse(projectName, response)
			if err != nil {
				fmt.Printf("\n❌ Erro ao salvar resposta:\n%v\n", err)
				os.Exit(1)
			}
			fmt.Println("💡 Resposta salva em README.md no diretório do projeto.")
		}
		fmt.Println(" ✅")

		// Executa plugins
		if len(pluginsList) > 0 {
			fmt.Print("🔌 Executando plugins...")
			plugins.ExecutePlugins()
			fmt.Println(" ✅")
		}

		elapsedTime := time.Since(startTime)

		fmt.Printf("\n✨ Projeto criado com sucesso! ✨\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📁 Local: %s\n", projectName)
		fmt.Printf("⏱️  Tempo total: %.2f segundos\n", elapsedTime.Seconds())
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		fmt.Printf("💡 Para começar a desenvolver:\n")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   Consulte o README.md para instruções detalhadas\n\n")
	},
}

func init() {
	// Configura flags para o comando scaffold
	scaffoldCmd.Flags().StringVarP(&language, "language", "l", "", "Linguagem para o scaffold (ex: go, python, etc)")
	scaffoldCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	scaffoldCmd.Flags().StringVarP(&description, "description", "d", "", "Descrição objetiva da estrutura desejada")
	scaffoldCmd.MarkFlagRequired("language")
	scaffoldCmd.MarkFlagRequired("name")

	// Registra o comando scaffold no comando raiz
	rootCmd.AddCommand(scaffoldCmd)
}
