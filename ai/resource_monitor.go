package ai

import (
	"runtime"
	"sync"
	"time"
)

// ResourceMonitor monitora recursos do sistema em tempo real
type ResourceMonitor struct {
	started         bool
	mu              sync.RWMutex
	metrics         *ResourceMetrics
	alertThresholds map[string]float64
	alerts          chan ResourceAlert
	stopChan        chan struct{}
	ticker          *time.Ticker
}

// ResourceAlert representa um alerta de recurso
type ResourceAlert struct {
	Type         string    `json:"type"`
	Message      string    `json:"message"`
	Severity     string    `json:"severity"`
	Timestamp    time.Time `json:"timestamp"`
	CurrentValue float64   `json:"current_value"`
	Threshold    float64   `json:"threshold"`
}

// NewResourceMonitor cria novo monitor de recursos
func NewResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{
		metrics:         &ResourceMetrics{},
		alertThresholds: make(map[string]float64),
		alerts:          make(chan ResourceAlert, 100),
		stopChan:        make(chan struct{}),
		ticker:          time.NewTicker(10 * time.Second),
	}
}

// Start inicia o monitor
func (rm *ResourceMonitor) Start() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.started {
		return nil
	}

	// Configurar thresholds padrão
	rm.setDefaultThresholds()

	// Iniciar goroutine de monitoramento
	go rm.monitorLoop()

	rm.started = true
	return nil
}

// Stop para o monitor
func (rm *ResourceMonitor) Stop() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.started {
		return nil
	}

	close(rm.stopChan)
	rm.ticker.Stop()
	rm.started = false
	return nil
}

// GetMetrics retorna métricas atuais
func (rm *ResourceMonitor) GetMetrics() *ResourceMetrics {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.metrics
}

// GetAlerts retorna canal de alertas
func (rm *ResourceMonitor) GetAlerts() <-chan ResourceAlert {
	return rm.alerts
}

// SetThreshold define threshold para um tipo de recurso
func (rm *ResourceMonitor) SetThreshold(resourceType string, threshold float64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.alertThresholds[resourceType] = threshold
}

// setDefaultThresholds define thresholds padrão
func (rm *ResourceMonitor) setDefaultThresholds() {
	rm.alertThresholds["memory"] = 80.0    // 80% de uso de memória
	rm.alertThresholds["goroutines"] = 500 // 500 goroutines
	rm.alertThresholds["cpu"] = 90.0       // 90% de uso de CPU
}

// monitorLoop loop principal de monitoramento
func (rm *ResourceMonitor) monitorLoop() {
	for {
		select {
		case <-rm.stopChan:
			return
		case <-rm.ticker.C:
			rm.collectMetrics()
			rm.checkAlerts()
		}
	}
}

// collectMetrics coleta métricas do sistema
func (rm *ResourceMonitor) collectMetrics() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Atualizar métricas
	rm.metrics.MemoryUsage = int64(m.Sys)
	rm.metrics.GoroutineCount = runtime.NumGoroutine()
	rm.metrics.LastUpdated = time.Now()

	// Calcular uso de CPU (simplificado)
	rm.metrics.CPUUsage = rm.calculateCPUUsage()
}

// calculateCPUUsage calcula uso de CPU (implementação simplificada)
func (rm *ResourceMonitor) calculateCPUUsage() float64 {
	// Implementação simplificada - em produção seria mais complexa
	return float64(runtime.NumGoroutine()) / 100.0
}

// checkAlerts verifica se algum threshold foi ultrapassado
func (rm *ResourceMonitor) checkAlerts() {
	// Verificar memória
	if threshold, exists := rm.alertThresholds["memory"]; exists {
		memoryUsagePercent := float64(rm.metrics.MemoryUsage) / (1024 * 1024 * 1024) * 100 // Assumindo 1GB como max
		if memoryUsagePercent > threshold {
			alert := ResourceAlert{
				Type:         "memory",
				Message:      "Uso de memória acima do threshold",
				Severity:     "high",
				Timestamp:    time.Now(),
				CurrentValue: memoryUsagePercent,
				Threshold:    threshold,
			}

			select {
			case rm.alerts <- alert:
			default:
				// Canal cheio, descartar alerta
			}
		}
	}

	// Verificar goroutines
	if threshold, exists := rm.alertThresholds["goroutines"]; exists {
		goroutineCount := float64(rm.metrics.GoroutineCount)
		if goroutineCount > threshold {
			alert := ResourceAlert{
				Type:         "goroutines",
				Message:      "Número de goroutines acima do threshold",
				Severity:     "medium",
				Timestamp:    time.Now(),
				CurrentValue: goroutineCount,
				Threshold:    threshold,
			}

			select {
			case rm.alerts <- alert:
			default:
				// Canal cheio, descartar alerta
			}
		}
	}

	// Verificar CPU
	if threshold, exists := rm.alertThresholds["cpu"]; exists {
		if rm.metrics.CPUUsage > threshold {
			alert := ResourceAlert{
				Type:         "cpu",
				Message:      "Uso de CPU acima do threshold",
				Severity:     "high",
				Timestamp:    time.Now(),
				CurrentValue: rm.metrics.CPUUsage,
				Threshold:    threshold,
			}

			select {
			case rm.alerts <- alert:
			default:
				// Canal cheio, descartar alerta
			}
		}
	}
}

// GetResourceStatus retorna status dos recursos
func (rm *ResourceMonitor) GetResourceStatus() map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return map[string]interface{}{
		"memory_usage_mb":   float64(rm.metrics.MemoryUsage) / 1024 / 1024,
		"goroutine_count":   rm.metrics.GoroutineCount,
		"cpu_usage_percent": rm.metrics.CPUUsage,
		"last_updated":      rm.metrics.LastUpdated,
		"alerts_pending":    len(rm.alerts),
	}
}

// GetDetailedReport retorna relatório detalhado
func (rm *ResourceMonitor) GetDetailedReport() string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	report := "📊 RELATÓRIO DETALHADO DE RECURSOS\n"
	report += "═════════════════════════════════════\n\n"

	report += "💾 MEMÓRIA:\n"
	report += "─────────────\n"
	report += "• Uso atual: " + formatBytes(rm.metrics.MemoryUsage) + "\n"
	report += "• Threshold: " + formatFloat(rm.alertThresholds["memory"]) + "%\n\n"

	report += "🔄 GOROUTINES:\n"
	report += "─────────────\n"
	report += "• Quantidade: " + formatInt(rm.metrics.GoroutineCount) + "\n"
	report += "• Threshold: " + formatFloat(rm.alertThresholds["goroutines"]) + "\n\n"

	report += "🖥️ CPU:\n"
	report += "─────────────\n"
	report += "• Uso atual: " + formatFloat(rm.metrics.CPUUsage) + "%\n"
	report += "• Threshold: " + formatFloat(rm.alertThresholds["cpu"]) + "%\n\n"

	report += "⏰ ÚLTIMA ATUALIZAÇÃO:\n"
	report += "─────────────────────\n"
	report += "• " + rm.metrics.LastUpdated.Format("02/01/2006 15:04:05") + "\n"

	return report
}

// formatBytes formata bytes para exibição
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return formatInt64(bytes) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return formatFloat(float64(bytes)/float64(div)) + " " + "KMGTPE"[exp:exp+1] + "B"
}

// formatFloat formata float para exibição
func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return formatInt64(int64(f))
	}
	return sprintf("%.2f", f)
}

// formatInt formata int para exibição
func formatInt(i int) string {
	return formatInt64(int64(i))
}

// formatInt64 formata int64 para exibição
func formatInt64(i int64) string {
	return sprintf("%d", i)
}

// sprintf simples implementação de sprintf
func sprintf(format string, args ...interface{}) string {
	// Implementação básica - em produção usaria fmt.Sprintf
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case int64:
			// Substituir %d por valor
			result = replaceFirst(result, "%d", itoa(int(v)))
		case float64:
			// Substituir %.2f por valor
			result = replaceFirst(result, "%.2f", ftoa(v))
		}
	}
	return result
}

// replaceFirst substitui primeira ocorrência
func replaceFirst(s, old, new string) string {
	// Implementação básica
	idx := indexOf(s, old)
	if idx == -1 {
		return s
	}
	return s[:idx] + new + s[idx+len(old):]
}

// indexOf encontra índice de substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// itoa converte int para string
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	negative := i < 0
	if negative {
		i = -i
	}

	var result []byte
	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}

	if negative {
		result = append([]byte{'-'}, result...)
	}

	return string(result)
}

// ftoa converte float para string (simplificado)
func ftoa(f float64) string {
	// Implementação muito simplificada
	integer := int64(f)
	decimal := int64((f - float64(integer)) * 100)

	return itoa(int(integer)) + "." + padLeft(itoa(int(decimal)), 2, '0')
}

// padLeft adiciona padding à esquerda
func padLeft(s string, length int, pad byte) string {
	for len(s) < length {
		s = string(pad) + s
	}
	return s
}
