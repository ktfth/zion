package cmd

import (
	"fmt"

	"github.com/ktfth/zion/ai"
	"github.com/spf13/cobra"
)

var testLayeredCmd = &cobra.Command{
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
  zion test-layered`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n🧪 TESTE DO SISTEMA DE GERAÇÃO EM CAMADAS\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		ai.TestLayeredGeneration()

		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("✅ Teste concluído\n\n")
	},
}

func init() {
	// Registra o comando de teste no comando raiz
	rootCmd.AddCommand(testLayeredCmd)
}
