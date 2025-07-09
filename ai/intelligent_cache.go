package ai

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// IntelligentCache implementa um sistema de cache inteligente
type IntelligentCache struct {
	cachePath           string
	entries             map[string]*CacheEntry
	metadata            CacheMetadata
	mu                  sync.RWMutex
	maxSize             int64
	currentSize         int64
	maxEntries          int
	evictionPolicy      EvictionPolicy
	compressionEnabled  bool
	encryptionEnabled   bool
	analytics           CacheAnalytics
	ttl                 time.Duration
	similarityThreshold float64
}

// CacheEntry representa uma entrada no cache
type CacheEntry struct {
	Key              string                 `json:"key"`
	Value            interface{}            `json:"value"`
	Timestamp        time.Time              `json:"timestamp"`
	LastAccessed     time.Time              `json:"last_accessed"`
	AccessCount      int                    `json:"access_count"`
	Size             int64                  `json:"size"`
	TTL              time.Duration          `json:"ttl"`
	Tags             []string               `json:"tags"`
	Metadata         map[string]interface{} `json:"metadata"`
	CompressionRatio float64                `json:"compression_ratio"`
	IsEncrypted      bool                   `json:"is_encrypted"`
	Hash             string                 `json:"hash"`
	Dependencies     []string               `json:"dependencies"`
	QualityScore     float64                `json:"quality_score"`
	UsagePattern     UsagePattern           `json:"usage_pattern"`
	Request          CacheRequest           `json:"request"`
	CreatedAt        time.Time              `json:"created_at"`
}

// CacheMetadata contém metadados do cache
type CacheMetadata struct {
	Version       string             `json:"version"`
	CreatedAt     time.Time          `json:"created_at"`
	LastCleanup   time.Time          `json:"last_cleanup"`
	TotalEntries  int                `json:"total_entries"`
	TotalSize     int64              `json:"total_size"`
	HitRate       float64            `json:"hit_rate"`
	MissRate      float64            `json:"miss_rate"`
	EvictionCount int                `json:"eviction_count"`
	Configuration CacheConfiguration `json:"configuration"`
}

// CacheConfiguration representa configuração do cache
type CacheConfiguration struct {
	MaxSize             int64          `json:"max_size"`
	DefaultTTL          time.Duration  `json:"default_ttl"`
	EvictionPolicy      EvictionPolicy `json:"eviction_policy"`
	CompressionEnabled  bool           `json:"compression_enabled"`
	EncryptionEnabled   bool           `json:"encryption_enabled"`
	CleanupInterval     time.Duration  `json:"cleanup_interval"`
	PrefetchEnabled     bool           `json:"prefetch_enabled"`
	IntelligentEviction bool           `json:"intelligent_eviction"`
}

// EvictionPolicy define políticas de evicção
type EvictionPolicy string

const (
	EvictionLRU         EvictionPolicy = "lru"         // Least Recently Used
	EvictionLFU         EvictionPolicy = "lfu"         // Least Frequently Used
	EvictionTTL         EvictionPolicy = "ttl"         // Time To Live
	EvictionSize        EvictionPolicy = "size"        // Tamanho da entrada
	EvictionIntelligent EvictionPolicy = "intelligent" // Baseado em ML
)

// UsagePattern representa padrões de uso
type UsagePattern struct {
	AccessTimes     []time.Time `json:"access_times"`
	AccessFrequency float64     `json:"access_frequency"`
	PeakHours       []int       `json:"peak_hours"`
	Seasonality     float64     `json:"seasonality"`
	Predictability  float64     `json:"predictability"`
}

// CacheAnalytics contém analytics do cache
type CacheAnalytics struct {
	TotalRequests           int64                  `json:"total_requests"`
	CacheHits               int64                  `json:"cache_hits"`
	CacheMisses             int64                  `json:"cache_misses"`
	HitRate                 float64                `json:"hit_rate"`
	MissRate                float64                `json:"miss_rate"`
	AverageResponseTime     time.Duration          `json:"average_response_time"`
	PopularKeys             map[string]int         `json:"popular_keys"`
	EvictionStats           map[EvictionPolicy]int `json:"eviction_stats"`
	SizeDistribution        map[string]int         `json:"size_distribution"`
	TimeDistribution        map[string]int         `json:"time_distribution"`
	QualityMetrics          CacheQualityMetrics    `json:"quality_metrics"`
	PredictionAccuracy      float64                `json:"prediction_accuracy"`
	OptimizationSuggestions []string               `json:"optimization_suggestions"`
	ResponseTimes           []time.Duration        `json:"response_times"`
	QualityDistribution     map[string]int         `json:"quality_distribution"`
	EvictionCount           int64                  `json:"eviction_count"`
	CompressionRatio        float64                `json:"compression_ratio"`
	LastEvictionTime        time.Time              `json:"last_eviction_time"`
}

// CacheQualityMetrics contém métricas de qualidade do cache
type CacheQualityMetrics struct {
	AverageQuality      float64        `json:"average_quality"`
	QualityDistribution map[string]int `json:"quality_distribution"`
	LowQualityCount     int            `json:"low_quality_count"`
	HighQualityCount    int            `json:"high_quality_count"`
	QualityTrend        []QualityPoint `json:"quality_trend"`
}

// QualityPoint representa um ponto na tendência de qualidade
type QualityPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Quality   float64   `json:"quality"`
	Success   bool      `json:"success"`
}

// CacheRequest representa uma requisição ao cache
type CacheRequest struct {
	Key         string                 `json:"key"`
	Language    string                 `json:"language"`
	ProjectType string                 `json:"project_type"`
	Description string                 `json:"description"`
	Context     interface{}            `json:"context"`
	Timestamp   time.Time              `json:"timestamp"`
}

// CacheResponse representa uma resposta do cache
type CacheResponse struct {
	Hit         bool                   `json:"hit"`
	Entry       *CacheEntry            `json:"entry"`
	Suggestions []CacheSuggestion      `json:"suggestions"`
	Metadata    map[string]interface{} `json:"metadata"`
	Performance CachePerformance       `json:"performance"`
}

// CacheSuggestion representa sugestões do cache
type CacheSuggestion struct {
	Type        SuggestionType `json:"type"`
	Key         string         `json:"key"`
	Similarity  float64        `json:"similarity"`
	Confidence  float64        `json:"confidence"`
	Reason      string         `json:"reason"`
	Adaptations []string       `json:"adaptations"`
}

// SuggestionType define tipos de sugestões
type SuggestionType string

const (
	SuggestionSimilar     SuggestionType = "similar"
	SuggestionPrefetch    SuggestionType = "prefetch"
	SuggestionOptimize    SuggestionType = "optimize"
	SuggestionAlternative SuggestionType = "alternative"
)

// CachePerformance contém métricas de performance
type CachePerformance struct {
	AccessTime        time.Duration `json:"access_time"`
	LookupTime        time.Duration `json:"lookup_time"`
	SerializationTime time.Duration `json:"serialization_time"`
	CompressionTime   time.Duration `json:"compression_time"`
	DecompressionTime time.Duration `json:"decompression_time"`
	NetworkTime       time.Duration `json:"network_time"`
	TotalTime         time.Duration `json:"total_time"`
}

// NewIntelligentCache cria um novo cache inteligente
func NewIntelligentCache(cachePath string, maxSize int64) *IntelligentCache {
	cache := &IntelligentCache{
		cachePath:           cachePath,
		entries:             make(map[string]*CacheEntry),
		maxSize:             maxSize,
		maxEntries:          1000,
		ttl:                 24 * time.Hour,
		similarityThreshold: 0.8,
		evictionPolicy:      EvictionIntelligent,
		compressionEnabled:  true,
		encryptionEnabled:   false,
		analytics: CacheAnalytics{
			PopularKeys:         make(map[string]int),
			EvictionStats:       make(map[EvictionPolicy]int),
			SizeDistribution:    make(map[string]int),
			TimeDistribution:    make(map[string]int),
			ResponseTimes:       make([]time.Duration, 0),
			QualityDistribution: make(map[string]int),
			QualityMetrics: CacheQualityMetrics{
				QualityDistribution: make(map[string]int),
				QualityTrend:        make([]QualityPoint, 0),
			},
		},
		metadata: CacheMetadata{
			Version:   "1.0",
			CreatedAt: time.Now(),
			Configuration: CacheConfiguration{
				MaxSize:             maxSize,
				DefaultTTL:          24 * time.Hour,
				EvictionPolicy:      EvictionIntelligent,
				CompressionEnabled:  true,
				EncryptionEnabled:   false,
				CleanupInterval:     time.Hour,
				PrefetchEnabled:     true,
				IntelligentEviction: true,
			},
		},
	}

	// Carregar cache existente
	if err := cache.loadFromDisk(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao carregar cache: %v\n", err)
	}

	// Iniciar limpeza automática
	go cache.startCleanupRoutine()

	return cache
}

// Get busca um item no cache
func (ic *IntelligentCache) Get(request CacheRequest) *CacheResponse {
	ic.mu.RLock()
	defer ic.mu.RUnlock()

	startTime := time.Now()

	// Gerar chave de cache
	key := ic.generateKey(request)

	// Buscar entrada
	entry, hit := ic.entries[key]

	// Atualizar analytics
	ic.updateAnalytics(key, hit)

	response := &CacheResponse{
		Hit:         hit,
		Suggestions: make([]CacheSuggestion, 0),
		Metadata:    make(map[string]interface{}),
		Performance: CachePerformance{
			TotalTime: time.Since(startTime),
		},
	}

	if hit {
		// Verificar se não expirou
		if !ic.isExpired(entry) {
			// Atualizar estatísticas de acesso
			entry.LastAccessed = time.Now()
			entry.AccessCount++

			response.Entry = entry
			response.Metadata["cache_age"] = time.Since(entry.Timestamp)
			response.Metadata["access_count"] = entry.AccessCount

			return response
		} else {
			// Entrada expirada
			delete(ic.entries, key)
			hit = false
		}
	}

	return response
}

// Put armazena um item no cache
func (ic *IntelligentCache) Put(request CacheRequest, value interface{}, quality float64) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	key := ic.generateKey(request)

	// Serializar valor
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("erro ao serializar valor: %v", err)
	}

	size := int64(len(data))

	// Comprimir se habilitado
	var compressionRatio float64 = 1.0
	if ic.compressionEnabled {
		// Implementar compressão
		compressionRatio = 0.7 // Estimativa
	}

	entry := &CacheEntry{
		Key:              key,
		Value:            value,
		Timestamp:        time.Now(),
		LastAccessed:     time.Now(),
		AccessCount:      1,
		Size:             size,
		TTL:              ic.ttl,
		Tags:             ic.generateTags(request),
		Metadata:         make(map[string]interface{}),
		CompressionRatio: compressionRatio,
		IsEncrypted:      ic.encryptionEnabled,
		Hash:             ic.generateHash(data),
		Dependencies:     make([]string, 0),
		QualityScore:     quality,
		Request:          request,
		CreatedAt:        time.Now(),
		UsagePattern: UsagePattern{
			AccessTimes:     []time.Time{time.Now()},
			AccessFrequency: 1.0,
			PeakHours:       []int{time.Now().Hour()},
		},
	}

	// Verificar se precisa de evicção
	if ic.currentSize+size > ic.maxSize {
		if err := ic.evictEntries(size); err != nil {
			return fmt.Errorf("erro ao fazer evicção: %v", err)
		}
	}

	// Armazenar entrada
	ic.entries[key] = entry
	ic.currentSize += size
	ic.metadata.TotalEntries++

	return nil
}

// generateKey gera uma chave única para a requisição
func (ic *IntelligentCache) generateKey(request CacheRequest) string {
	// Criar string única baseada nos parâmetros
	keyData := fmt.Sprintf("%s|%s|%s|%v",
		request.Language,
		request.ProjectType,
		request.Description,
		request.Context)

	// Gerar hash
	hash := sha256.Sum256([]byte(keyData))
	return fmt.Sprintf("%x", hash)
}

// generateHash gera hash dos dados
func (ic *IntelligentCache) generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// generateTags gera tags para a entrada
func (ic *IntelligentCache) generateTags(request CacheRequest) []string {
	tags := []string{
		request.Language,
		request.ProjectType,
	}

	// Extrair palavras-chave da descrição
	keywords := ic.extractKeywords(request.Description)
	tags = append(tags, keywords...)

	// Adicionar contexto como tags
	if request.Context != nil {
		if contextMap, ok := request.Context.(map[string]interface{}); ok {
			for key := range contextMap {
				tags = append(tags, key)
			}
		}
	}

	return tags
}

// extractKeywords extrai palavras-chave de uma string
func (ic *IntelligentCache) extractKeywords(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	keywords := make([]string, 0)

	for _, word := range words {
		// Filtrar palavras muito curtas ou stop words
		if len(word) > 3 && !ic.isStopWord(word) {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// isStopWord verifica se uma palavra é uma stop word
func (ic *IntelligentCache) isStopWord(word string) bool {
	stopWords := []string{
		"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by", "this", "that", "these", "those",
		"a", "an", "is", "are", "was", "were", "be", "been", "have", "has", "had", "do", "does", "did", "will", "would",
		"could", "should", "may", "might", "can", "must", "shall", "from", "up", "out", "down", "over", "under", "above",
		"below", "into", "onto", "upon", "off", "through", "across", "along", "around", "between", "among", "within",
		"without", "during", "before", "after", "while", "when", "where", "why", "how", "what", "which", "who", "whom",
		"para", "com", "sem", "por", "de", "em", "na", "no", "da", "do", "dos", "das", "que", "o", "a", "os", "as",
		"um", "uma", "uns", "umas", "este", "esta", "estes", "estas", "esse", "essa", "esses", "essas", "aquele",
		"aquela", "aqueles", "aquelas", "eu", "tu", "ele", "ela", "nos", "vos", "eles", "elas", "meu", "minha",
		"seus", "suas", "nosso", "nossa", "vosso", "vossa", "seu", "sua", "dele", "dela", "deles", "delas",
	}

	for _, stopWord := range stopWords {
		if word == stopWord {
			return true
		}
	}
	return false
}

// isExpired verifica se uma entrada expirou
func (ic *IntelligentCache) isExpired(entry *CacheEntry) bool {
	return time.Since(entry.CreatedAt) > ic.ttl
}

// evictEntries remove entradas para liberar espaço
func (ic *IntelligentCache) evictEntries(requiredSize int64) error {
	entriesToEvict := ic.selectLRUEntries(requiredSize)
	
	var freedSize int64
	for _, entry := range entriesToEvict {
		for key, cachedEntry := range ic.entries {
			if cachedEntry == entry {
				freedSize += entry.Size
				delete(ic.entries, key)
				break
			}
		}
	}
	
	ic.currentSize -= freedSize
	
	if freedSize < requiredSize {
		return fmt.Errorf("não foi possível liberar espaço suficiente")
	}
	
	return nil
}

// selectLRUEntries seleciona entradas para remoção baseado em LRU
func (ic *IntelligentCache) selectLRUEntries(requiredSize int64) []*CacheEntry {
	entries := make([]*CacheEntry, 0)
	for _, entry := range ic.entries {
		entries = append(entries, entry)
	}

	// Ordenar por último acesso (mais antigos primeiro)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].LastAccessed.After(entries[j].LastAccessed) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	var selectedSize int64
	selected := make([]*CacheEntry, 0)

	for _, entry := range entries {
		selected = append(selected, entry)
		selectedSize += entry.Size
		if selectedSize >= requiredSize {
			break
		}
	}

	return selected
}

// updateAnalytics atualiza analytics do cache
func (ic *IntelligentCache) updateAnalytics(key string, hit bool) {
	ic.analytics.TotalRequests++

	if hit {
		ic.analytics.CacheHits++
	} else {
		ic.analytics.CacheMisses++
	}

	if ic.analytics.TotalRequests > 0 {
		ic.analytics.HitRate = float64(ic.analytics.CacheHits) / float64(ic.analytics.TotalRequests)
		ic.analytics.MissRate = float64(ic.analytics.CacheMisses) / float64(ic.analytics.TotalRequests)
	}

	// Atualizar chaves populares
	if hit {
		ic.analytics.PopularKeys[key]++
	}
}

// startCleanupRoutine inicia rotina de limpeza
func (ic *IntelligentCache) startCleanupRoutine() {
	ticker := time.NewTicker(ic.metadata.Configuration.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ic.cleanup()
		}
	}
}

// cleanup remove entradas expiradas
func (ic *IntelligentCache) cleanup() {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	expiredKeys := make([]string, 0)

	for key, entry := range ic.entries {
		if ic.isExpired(entry) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remover entradas expiradas
	for _, key := range expiredKeys {
		entry := ic.entries[key]
		delete(ic.entries, key)
		ic.currentSize -= entry.Size
		ic.metadata.TotalEntries--
	}

	ic.metadata.LastCleanup = time.Now()
}

// loadFromDisk carrega cache do disco
func (ic *IntelligentCache) loadFromDisk() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(ic.cachePath, 0755); err != nil {
		return err
	}

	// Carregar entradas
	entriesFile := filepath.Join(ic.cachePath, "entries.json")
	if data, err := os.ReadFile(entriesFile); err == nil {
		if err := json.Unmarshal(data, &ic.entries); err != nil {
			return err
		}
	}

	// Recalcular tamanho atual
	ic.currentSize = 0
	for _, entry := range ic.entries {
		ic.currentSize += entry.Size
	}

	return nil
}

// saveToDisk salva cache no disco
func (ic *IntelligentCache) saveToDisk() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(ic.cachePath, 0755); err != nil {
		return err
	}

	// Salvar entradas
	entriesFile := filepath.Join(ic.cachePath, "entries.json")
	if data, err := json.MarshalIndent(ic.entries, "", "  "); err == nil {
		if err := os.WriteFile(entriesFile, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetStats retorna estatísticas do cache
func (ic *IntelligentCache) GetStats() map[string]interface{} {
	ic.mu.RLock()
	defer ic.mu.RUnlock()

	return map[string]interface{}{
		"total_entries":   len(ic.entries),
		"current_size":    ic.currentSize,
		"max_size":        ic.maxSize,
		"hit_rate":        ic.analytics.HitRate,
		"miss_rate":       ic.analytics.MissRate,
		"total_requests":  ic.analytics.TotalRequests,
		"cache_hits":      ic.analytics.CacheHits,
		"cache_misses":    ic.analytics.CacheMisses,
	}
}

// CacheStats representa estatísticas do cache
type CacheStats struct {
	HitRate             float64       `json:"hit_rate"`
	TotalRequests       int64         `json:"total_requests"`
	Hits                int64         `json:"hits"`
	Misses              int64         `json:"misses"`
	ItemCount           int           `json:"item_count"`
	TotalSize           int64         `json:"total_size"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	LastEvictionTime    time.Time     `json:"last_eviction_time"`
	EvictionCount       int64         `json:"eviction_count"`
	CompressionRatio    float64       `json:"compression_ratio"`
}
