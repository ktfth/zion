package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ktfth/zion/ai"
	"github.com/spf13/cobra"
)

// economyCmd representa o comando de economia de recursos
var economyCmd = &cobra.Command{
	Use:   "economy",
	Short: "Gerencia economia de recursos do sistema",
	Long: `O comando economy fornece insights sobre o uso de recursos do sistema,
incluindo monitoramento de memória, goroutines, tokens de API e outras métricas
relacionadas à economia de recursos.`,
	RunE: runEconomyCommand,
}

// economyStatsCmd mostra estatísticas de economia
var economyStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Mostra estatísticas de economia de recursos",
	Long: `Exibe estatísticas detalhadas sobre o uso de recursos do sistema,
incluindo métricas de performance, uso de memória e análise de custos.`,
	RunE: runEconomyStats,
}

// economyMonitorCmd monitora recursos em tempo real
var economyMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitora recursos em tempo real",
	Long: `Inicia o monitoramento em tempo real do uso de recursos,
exibindo métricas atualizadas sobre economia e performance.`,
	RunE: runEconomyMonitor,
}

// economyOptimizeCmd otimiza o uso de recursos
var economyOptimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Otimiza o uso de recursos",
	Long: `Executa otimizações automáticas para melhorar a economia de recursos,
incluindo limpeza de cache, garbage collection e ajustes de configuração.`,
	RunE: runEconomyOptimize,
}

// economyReportCmd gera relatório de economia
var economyReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Gera relatório de economia de recursos",
	Long: `Gera um relatório detalhado sobre o uso de recursos,
incluindo análise de custos, tendências e recomendações de otimização.`,
	RunE: runEconomyReport,
}

func init() {
	rootCmd.AddCommand(economyCmd)

	// Adicionar subcomandos
	economyCmd.AddCommand(economyStatsCmd)
	economyCmd.AddCommand(economyMonitorCmd)
	economyCmd.AddCommand(economyOptimizeCmd)
	economyCmd.AddCommand(economyReportCmd)

	// Flags para monitor
	economyMonitorCmd.Flags().DurationP("interval", "i", 5*time.Second, "Intervalo de atualização")
	economyMonitorCmd.Flags().IntP("count", "c", 0, "Número de atualizações (0 = infinito)")

	// Flags para relatório
	economyReportCmd.Flags().StringP("format", "f", "text", "Formato do relatório (text, json, html)")
	economyReportCmd.Flags().StringP("output", "o", "", "Arquivo de saída")
}

// runEconomyCommand executa o comando principal de economia
func runEconomyCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("🏦 Sistema de Economia de Recursos - Zion CLI")
	fmt.Println()

	// Inicializar resource manager
	resourceManager := ai.NewResourceManager()
	if err := resourceManager.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar resource manager: %w", err)
	}
	defer resourceManager.Stop()

	// Mostrar status geral
	status := resourceManager.GetResourceStatus()

	fmt.Println("📊 Status Atual dos Recursos:")
	fmt.Println()

	// Memória
	if memStats, ok := status["memory"].(map[string]interface{}); ok {
		fmt.Printf("💾 Memória: %.1f%% (%.2f MB / %.2f MB)\n",
			memStats["usage"].(float64),
			float64(memStats["current"].(int64))/1024/1024,
			float64(memStats["limit"].(int64))/1024/1024)
	}

	// Goroutines
	if goroutineStats, ok := status["goroutines"].(map[string]interface{}); ok {
		fmt.Printf("🔄 Goroutines: %.1f%% (%d / %d)\n",
			goroutineStats["usage"].(float64),
			goroutineStats["current"].(int),
			goroutineStats["limit"].(int))
	}

	// Tokens
	if tokenStats, ok := status["tokens"].(map[string]interface{}); ok {
		fmt.Printf("🎫 Tokens: %.1f%% (%d / %d)\n",
			tokenStats["usage"].(float64),
			tokenStats["limit"].(int)-tokenStats["current"].(int),
			tokenStats["limit"].(int))
	}

	fmt.Println()
	fmt.Println("💡 Comandos disponíveis:")
	fmt.Println("  zion economy stats     - Estatísticas detalhadas")
	fmt.Println("  zion economy monitor   - Monitoramento em tempo real")
	fmt.Println("  zion economy optimize  - Otimizar recursos")
	fmt.Println("  zion economy report    - Gerar relatório")

	return nil
}

// runEconomyStats mostra estatísticas detalhadas
func runEconomyStats(cmd *cobra.Command, args []string) error {
	fmt.Println("📈 Estatísticas de Economia de Recursos")
	fmt.Println()

	resourceManager := ai.NewResourceManager()
	if err := resourceManager.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar resource manager: %w", err)
	}
	defer resourceManager.Stop()

	metrics := resourceManager.GetMetrics()

	fmt.Printf("🕐 Última atualização: %s\n", metrics.LastUpdated.Format("15:04:05"))
	fmt.Printf("💾 Uso de memória: %.2f MB\n", float64(metrics.MemoryUsage)/1024/1024)
	fmt.Printf("🔄 Goroutines ativas: %d\n", metrics.GoroutineCount)
	fmt.Printf("🎫 Tokens utilizados: %d\n", metrics.TokensUsed)
	fmt.Printf("📊 Total de requisições: %d\n", metrics.RequestCount)
	fmt.Printf("❌ Taxa de erro: %.2f%%\n", metrics.ErrorRate*100)

	return nil
}

// runEconomyMonitor monitora recursos em tempo real
func runEconomyMonitor(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	count, _ := cmd.Flags().GetInt("count")

	fmt.Println("🔍 Monitoramento em Tempo Real")
	fmt.Printf("📊 Intervalo: %v\n", interval)
	if count > 0 {
		fmt.Printf("🔢 Atualizações: %d\n", count)
	}
	fmt.Println()

	resourceManager := ai.NewResourceManager()
	if err := resourceManager.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar resource manager: %w", err)
	}
	defer resourceManager.Stop()

	iteration := 0
	for {
		if count > 0 && iteration >= count {
			break
		}

		// Limpar tela (simples)
		if iteration > 0 {
			fmt.Print("\033[2J\033[H")
		}

		fmt.Printf("🔍 Monitoramento - Atualização %d (%s)\n", iteration+1, time.Now().Format("15:04:05"))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		status := resourceManager.GetResourceStatus()

		// Exibir métricas
		if memStats, ok := status["memory"].(map[string]interface{}); ok {
			usage := memStats["usage"].(float64)
			bar := createProgressBar(usage, 40)
			fmt.Printf("💾 Memória   [%s] %.1f%%\n", bar, usage)
		}

		if goroutineStats, ok := status["goroutines"].(map[string]interface{}); ok {
			usage := goroutineStats["usage"].(float64)
			bar := createProgressBar(usage, 40)
			fmt.Printf("🔄 Goroutines [%s] %.1f%%\n", bar, usage)
		}

		if tokenStats, ok := status["tokens"].(map[string]interface{}); ok {
			usage := tokenStats["usage"].(float64)
			bar := createProgressBar(usage, 40)
			fmt.Printf("🎫 Tokens    [%s] %.1f%%\n", bar, usage)
		}

		iteration++
		time.Sleep(interval)
	}

	return nil
}

// runEconomyOptimize otimiza o uso de recursos
func runEconomyOptimize(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Otimizando Recursos...")
	fmt.Println()

	resourceManager := ai.NewResourceManager()
	if err := resourceManager.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar resource manager: %w", err)
	}
	defer resourceManager.Stop()

	// Verificar uso atual
	fmt.Println("📊 Verificando uso atual de recursos...")
	beforeMem, _ := resourceManager.CheckMemoryUsage()
	beforeGoroutines, _ := resourceManager.CheckGoroutines()

	// Executar otimizações
	fmt.Println("🧹 Executando limpeza...")
	resourceManager.Cleanup()

	// Verificar após otimização
	fmt.Println("📈 Verificando resultados...")
	afterMem, _ := resourceManager.CheckMemoryUsage()
	afterGoroutines, _ := resourceManager.CheckGoroutines()

	// Mostrar resultados
	fmt.Println()
	fmt.Println("✅ Otimização concluída!")
	fmt.Printf("💾 Memória: %.2f MB → %.2f MB (%.2f MB liberados)\n",
		float64(beforeMem)/1024/1024,
		float64(afterMem)/1024/1024,
		float64(beforeMem-afterMem)/1024/1024)
	fmt.Printf("🔄 Goroutines: %d → %d (%d reduzidos)\n",
		beforeGoroutines,
		afterGoroutines,
		beforeGoroutines-afterGoroutines)

	return nil
}

// runEconomyReport gera relatório de economia
func runEconomyReport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	fmt.Println("📋 Gerando Relatório de Economia de Recursos...")
	fmt.Println()

	resourceManager := ai.NewResourceManager()
	if err := resourceManager.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar resource manager: %w", err)
	}
	defer resourceManager.Stop()

	// Gerar relatório
	report := generateEconomyReport(resourceManager, format)

	// Salvar ou exibir
	if output != "" {
		if err := os.WriteFile(output, []byte(report), 0644); err != nil {
			return fmt.Errorf("erro ao salvar relatório: %w", err)
		}
		fmt.Printf("✅ Relatório salvo em: %s\n", output)
	} else {
		fmt.Println(report)
	}

	return nil
}

// createProgressBar cria barra de progresso
func createProgressBar(percentage float64, width int) string {
	filled := int(percentage * float64(width) / 100)
	bar := ""

	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return bar
}

// generateEconomyReport gera relatório de economia
func generateEconomyReport(resourceManager *ai.ResourceManager, format string) string {
	metrics := resourceManager.GetMetrics()
	status := resourceManager.GetResourceStatus()

	switch format {
	case "json":
		return generateJSONReport(metrics, status)
	case "html":
		return generateHTMLReport(metrics, status)
	default:
		return generateTextReport(metrics, status)
	}
}

// generateTextReport gera relatório em texto
func generateTextReport(metrics *ai.ResourceMetrics, status map[string]interface{}) string {
	report := "🏦 RELATÓRIO DE ECONOMIA DE RECURSOS\n"
	report += "═════════════════════════════════════\n\n"

	report += fmt.Sprintf("📅 Data: %s\n", time.Now().Format("02/01/2006 15:04:05"))
	report += fmt.Sprintf("🕐 Última atualização: %s\n\n", metrics.LastUpdated.Format("15:04:05"))

	report += "📊 RESUMO EXECUTIVO\n"
	report += "───────────────────\n"

	if memStats, ok := status["memory"].(map[string]interface{}); ok {
		report += fmt.Sprintf("💾 Uso de Memória: %.1f%%\n", memStats["usage"].(float64))
	}

	if goroutineStats, ok := status["goroutines"].(map[string]interface{}); ok {
		report += fmt.Sprintf("🔄 Goroutines: %.1f%%\n", goroutineStats["usage"].(float64))
	}

	if tokenStats, ok := status["tokens"].(map[string]interface{}); ok {
		report += fmt.Sprintf("🎫 Tokens: %.1f%%\n", tokenStats["usage"].(float64))
	}

	report += "\n📈 MÉTRICAS DETALHADAS\n"
	report += "─────────────────────\n"
	report += fmt.Sprintf("📊 Requisições: %d\n", metrics.RequestCount)
	report += fmt.Sprintf("❌ Taxa de erro: %.2f%%\n", metrics.ErrorRate*100)
	report += fmt.Sprintf("🎫 Tokens utilizados: %d\n", metrics.TokensUsed)

	report += "\n💡 RECOMENDAÇÕES\n"
	report += "───────────────\n"
	report += "• Execute 'zion economy optimize' regularmente\n"
	report += "• Monitore o uso de recursos com 'zion economy monitor'\n"
	report += "• Configure alertas para uso excessivo de recursos\n"

	return report
}

// generateJSONReport gera relatório em JSON
func generateJSONReport(metrics *ai.ResourceMetrics, status map[string]interface{}) string {
	return fmt.Sprintf(`{
	"timestamp": "%s",
	"metrics": {
		"cpu_usage": %.2f,
		"memory_usage": %d,
		"goroutine_count": %d,
		"tokens_used": %d,
		"request_count": %d,
		"error_rate": %.4f,
		"last_updated": "%s"
	},
	"status": %v
}`, time.Now().Format(time.RFC3339), metrics.CPUUsage, metrics.MemoryUsage, metrics.GoroutineCount, metrics.TokensUsed, metrics.RequestCount, metrics.ErrorRate, metrics.LastUpdated.Format(time.RFC3339), status)
}

// generateHTMLReport gera relatório em HTML
func generateHTMLReport(metrics *ai.ResourceMetrics, status map[string]interface{}) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Relatório de Economia de Recursos</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 5px; }
        .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 20px 0; }
        .metric { background: #ecf0f1; padding: 15px; border-radius: 5px; }
        .metric h3 { margin: 0 0 10px 0; color: #2c3e50; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🏦 Relatório de Economia de Recursos</h1>
        <p>Gerado em: ` + time.Now().Format("02/01/2006 15:04:05") + `</p>
    </div>
    
    <div class="metrics">
        <div class="metric">
            <h3>💾 Memória</h3>
            <p>Uso atual: ` + fmt.Sprintf("%.2f MB", float64(metrics.MemoryUsage)/1024/1024) + `</p>
        </div>
        <div class="metric">
            <h3>🔄 Goroutines</h3>
            <p>Quantidade: ` + fmt.Sprintf("%d", metrics.GoroutineCount) + `</p>
        </div>
        <div class="metric">
            <h3>🎫 Tokens</h3>
            <p>Utilizados: ` + fmt.Sprintf("%d", metrics.TokensUsed) + `</p>
        </div>
        <div class="metric">
            <h3>📊 Requisições</h3>
            <p>Total: ` + fmt.Sprintf("%d", metrics.RequestCount) + `</p>
        </div>
    </div>
</body>
</html>`
}
