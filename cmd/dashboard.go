package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dashboardCmd define o comando "dashboard".
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Exibe um dashboard de informações do projeto",
	Long: `O comando dashboard exibe informações úteis sobre o projeto atual,
incluindo estatísticas, métricas e status do sistema.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 ZION DASHBOARD")
		fmt.Println("================")
		fmt.Println("🚧 Dashboard em desenvolvimento...")
		fmt.Println("Em breve: métricas do projeto, status do sistema e estatísticas detalhadas.")
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
