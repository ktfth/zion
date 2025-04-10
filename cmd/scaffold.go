package cmd

import (
	"fmt"
	"os"
	"zion/ai"
	"zion/plugins"

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
		fmt.Printf("\nIniciando scaffold para projeto '%s' em '%s'\n", projectName, language)
		fmt.Printf("Descrição: %s\n", description)
		
		// Para todas as linguagens, usa o fluxo normal com a API Gemini
		// Gera o prompt reforçando boas práticas e solicita a estrutura ao serviço de AI
		fmt.Println("Gerando estrutura do projeto com IA...")
		response, err := ai.GenerateProjectScaffolding(language, projectName, description, plugins.ListPlugins())
		if err != nil {
			fmt.Println("Erro na geração da estrutura:", err)
			fmt.Println("Resposta:", response) // Exibe a resposta para debug
			os.Exit(1)
		}

		// Criar a estrutura do projeto com base na resposta JSON
		fmt.Println("Criando estrutura do projeto...")
		
		// Usar o extrator que gera a saída em formato TOML
		fmt.Println("Extraindo e criando arquivos com saída em TOML...")
		err = ai.ExtractAndCreateTomlProject(projectName, response)
		if err != nil {
			fmt.Printf("Erro ao extrair e criar arquivos com TOML: %v\nTentando método alternativo...\n", err)
			
			// Se falhar, tentar o método anterior
			fmt.Println("Tentando extrator direto...")
			err = ai.ExtractAndCreateProject(projectName, response)
			if err != nil {
				fmt.Printf("Erro ao extrair e criar arquivos: %v\nSalvando resposta bruta...\n", err)
				
				// Se falhar novamente, salvar a resposta bruta em um arquivo
				fmt.Println("Salvando resposta bruta...")
				err = ai.SaveRawResponse(projectName, response)
				if err != nil {
					fmt.Println("Erro ao salvar resposta bruta:", err)
					os.Exit(1)
				}
				fmt.Println("Resposta bruta salva. Consulte o README.md no diretório do projeto para instruções.")
			}
		}

		// Executa plugins registrados (para estender as funcionalidades)
		fmt.Println("Executando plugins...")
		plugins.ExecutePlugins()
		
		fmt.Printf("\nProjeto '%s' criado com sucesso!\n", projectName)
	},
}

func init() {
	// Configura flags para o comando scaffold
	scaffoldCmd.Flags().StringVarP(&language, "language", "l", "", "Linguagem para o scaffold (ex: go, python, etc)")
	scaffoldCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	scaffoldCmd.Flags().StringVarP(&description, "description", "d", "", "Descrição objetiva da estrutura desejada (ex: 'aplicação web com autenticação', 'API REST com banco de dados')")
	scaffoldCmd.MarkFlagRequired("language")
	scaffoldCmd.MarkFlagRequired("name")

	// Registra o comando scaffold no comando raiz
	rootCmd.AddCommand(scaffoldCmd)
}

