package ai

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ResourceManager gerencia recursos do sistema
type ResourceManager struct {
	mu               sync.RWMutex
	maxMemory        int64
	maxGoroutines    int
	maxTokensPerHour int
	tokensBucket     *TokenBucket
	memoryUsage      int64
	goroutineCount   int
	resourceMonitor  *ResourceMonitor
	resourcePools    map[string]*ResourcePool
	resourceQuotas   map[string]*ResourceQuota
	resourceMetrics  *ResourceMetrics
	started          bool
}

// ResourcePool pool de recursos
type ResourcePool struct {
	Name        string
	MaxSize     int
	CurrentSize int
	Resources   chan interface{}
	mu          sync.RWMutex
}

// ResourceQuota quota de recursos
type ResourceQuota struct {
	Name      string
	Limit     int64
	Used      int64
	Period    time.Duration
	LastReset time.Time
	mu        sync.RWMutex
}

// ResourceMetrics métricas de recursos
type ResourceMetrics struct {
	CPUUsage       float64
	MemoryUsage    int64
	GoroutineCount int
	TokensUsed     int64
	RequestCount   int64
	ErrorRate      float64
	LastUpdated    time.Time
}

// TokenBucket implementa rate limiting
type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

// NewResourceManager cria novo gerenciador de recursos
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		maxMemory:        1024 * 1024 * 1024, // 1GB
		maxGoroutines:    1000,
		maxTokensPerHour: 10000,
		tokensBucket:     NewTokenBucket(10000, time.Hour),
		resourcePools:    make(map[string]*ResourcePool),
		resourceQuotas:   make(map[string]*ResourceQuota),
		resourceMetrics:  &ResourceMetrics{},
		resourceMonitor:  NewResourceMonitor(),
	}
}

// NewTokenBucket cria novo token bucket
func NewTokenBucket(tokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     tokens,
		maxTokens:  tokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Start inicia o gerenciador de recursos
func (rm *ResourceManager) Start() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.started {
		return fmt.Errorf("resource manager already started")
	}

	// Iniciar monitor de recursos
	if err := rm.resourceMonitor.Start(); err != nil {
		return fmt.Errorf("failed to start resource monitor: %w", err)
	}

	// Iniciar pools padrão
	rm.initializePools()

	// Iniciar quotas padrão
	rm.initializeQuotas()

	// Iniciar coleta de métricas
	go rm.collectMetrics()

	rm.started = true
	return nil
}

// Stop para o gerenciador de recursos
func (rm *ResourceManager) Stop() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.started {
		return nil
	}

	// Parar monitor de recursos
	if err := rm.resourceMonitor.Stop(); err != nil {
		return fmt.Errorf("failed to stop resource monitor: %w", err)
	}

	rm.started = false
	return nil
}

// AcquireToken adquire um token do bucket
func (rm *ResourceManager) AcquireToken() bool {
	return rm.tokensBucket.Take()
}

// Take retira um token do bucket
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens se necessário
	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// refill reabastece tokens
func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefill)

	if duration >= tb.refillRate {
		tb.tokens = tb.maxTokens
		tb.lastRefill = now
	}
}

// CheckMemoryUsage verifica uso de memória
func (rm *ResourceManager) CheckMemoryUsage() (int64, error) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)

	currentMemory := int64(m.Sys)

	rm.mu.Lock()
	rm.memoryUsage = currentMemory
	rm.mu.Unlock()

	if currentMemory > rm.maxMemory {
		return currentMemory, fmt.Errorf("memory usage (%d) exceeded limit (%d)", currentMemory, rm.maxMemory)
	}

	return currentMemory, nil
}

// CheckGoroutines verifica número de goroutines
func (rm *ResourceManager) CheckGoroutines() (int, error) {
	count := runtime.NumGoroutine()

	rm.mu.Lock()
	rm.goroutineCount = count
	rm.mu.Unlock()

	if count > rm.maxGoroutines {
		return count, fmt.Errorf("goroutine count (%d) exceeded limit (%d)", count, rm.maxGoroutines)
	}

	return count, nil
}

// CreatePool cria um pool de recursos
func (rm *ResourceManager) CreatePool(name string, size int) *ResourcePool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	pool := &ResourcePool{
		Name:      name,
		MaxSize:   size,
		Resources: make(chan interface{}, size),
	}

	rm.resourcePools[name] = pool
	return pool
}

// GetPool obtém um pool de recursos
func (rm *ResourceManager) GetPool(name string) *ResourcePool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.resourcePools[name]
}

// CreateQuota cria uma quota de recursos
func (rm *ResourceManager) CreateQuota(name string, limit int64, period time.Duration) *ResourceQuota {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	quota := &ResourceQuota{
		Name:      name,
		Limit:     limit,
		Period:    period,
		LastReset: time.Now(),
	}

	rm.resourceQuotas[name] = quota
	return quota
}

// UseQuota consome quota de recursos
func (rm *ResourceManager) UseQuota(name string, amount int64) error {
	rm.mu.RLock()
	quota, exists := rm.resourceQuotas[name]
	rm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("quota %s not found", name)
	}

	quota.mu.Lock()
	defer quota.mu.Unlock()

	// Reset quota se período expirou
	if time.Since(quota.LastReset) > quota.Period {
		quota.Used = 0
		quota.LastReset = time.Now()
	}

	if quota.Used+amount > quota.Limit {
		return fmt.Errorf("quota %s exceeded: used %d + %d > limit %d", name, quota.Used, amount, quota.Limit)
	}

	quota.Used += amount
	return nil
}

// GetMetrics retorna métricas atuais
func (rm *ResourceManager) GetMetrics() *ResourceMetrics {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.resourceMetrics
}

// initializePools inicializa pools padrão
func (rm *ResourceManager) initializePools() {
	rm.CreatePool("ai_requests", 10)
	rm.CreatePool("cache_entries", 100)
	rm.CreatePool("evaluations", 5)
}

// initializeQuotas inicializa quotas padrão
func (rm *ResourceManager) initializeQuotas() {
	rm.CreateQuota("ai_tokens", 10000, time.Hour)
	rm.CreateQuota("api_requests", 1000, time.Hour)
	rm.CreateQuota("cache_size", 1024*1024*100, time.Hour) // 100MB
}

// collectMetrics coleta métricas periodicamente
func (rm *ResourceManager) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rm.updateMetrics()
		}
	}
}

// updateMetrics atualiza métricas
func (rm *ResourceManager) updateMetrics() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Atualizar métricas de CPU
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	rm.resourceMetrics.MemoryUsage = int64(m.Sys)
	rm.resourceMetrics.GoroutineCount = runtime.NumGoroutine()
	rm.resourceMetrics.LastUpdated = time.Now()
}

// GetResourceStatus retorna status dos recursos
func (rm *ResourceManager) GetResourceStatus() map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	status := make(map[string]interface{})

	// Status de memória
	status["memory"] = map[string]interface{}{
		"current": rm.memoryUsage,
		"limit":   rm.maxMemory,
		"usage":   float64(rm.memoryUsage) / float64(rm.maxMemory) * 100,
	}

	// Status de goroutines
	status["goroutines"] = map[string]interface{}{
		"current": rm.goroutineCount,
		"limit":   rm.maxGoroutines,
		"usage":   float64(rm.goroutineCount) / float64(rm.maxGoroutines) * 100,
	}

	// Status de tokens
	status["tokens"] = map[string]interface{}{
		"current": rm.tokensBucket.tokens,
		"limit":   rm.tokensBucket.maxTokens,
		"usage":   float64(rm.tokensBucket.maxTokens-rm.tokensBucket.tokens) / float64(rm.tokensBucket.maxTokens) * 100,
	}

	return status
}

// Cleanup limpa recursos não utilizados
func (rm *ResourceManager) Cleanup() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Forçar garbage collection
	runtime.GC()

	// Limpar pools vazios
	for name, pool := range rm.resourcePools {
		if pool.CurrentSize == 0 {
			delete(rm.resourcePools, name)
		}
	}
}
