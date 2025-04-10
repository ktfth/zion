package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zion/config"
	"zion/plugins"
)

// rootCmd é o comando principal da CLI.
var rootCmd = &cobra.Command{
	Use:   "zion",
	Short: "Zion CLI - Scaffolding de projetos com integração a AI",
	Long: `Zion é uma ferramenta de scaffolding que gera a estrutura
de projetos para qualquer linguagem, integrando-se com serviços de AI (GPT/Gemini)
e reforçando boas práticas de código. Além disso, possui um sistema de plugins para extensão.`,
}

// Execute inicia a CLI.
func Execute() {
	// Carregar a configuração
	cfg := config.LoadConfig()
	
	// Carregar plugins
	if err := plugins.LoadPlugins(cfg); err != nil {
		fmt.Printf("Erro ao carregar plugins: %v\n", err)
	}
	
	// Executar o comando raiz
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

