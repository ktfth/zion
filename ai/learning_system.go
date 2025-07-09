package ai

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LearningSystem implementa um sistema de aprendizado baseado em histórico
type LearningSystem struct {
	historyPath string
	sessions    []GenerationSession
	patterns    map[string]PatternLearning
	preferences UserPreferences
	analytics   GenerationAnalytics
}

// GenerationSession representa uma sessão de geração
type GenerationSession struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	Language     string                 `json:"language"`
	ProjectType  string                 `json:"project_type"`
	Description  string                 `json:"description"`
	Success      bool                   `json:"success"`
	Duration     time.Duration          `json:"duration"`
	TokensUsed   int                    `json:"tokens_used"`
	Provider     string                 `json:"provider"`
	LayeredMode  bool                   `json:"layered_mode"`
	Quality      float64                `json:"quality"`
	UserRating   int                    `json:"user_rating"` // 1-5
	Feedback     string                 `json:"feedback"`
	Patterns     []string               `json:"patterns"`
	Adaptations  map[string]interface{} `json:"adaptations"`
	Context      SessionContext         `json:"context"`
	Outcomes     GenerationOutcomes     `json:"outcomes"`
	Improvements []string               `json:"improvements"`
}

// SessionContext contém contexto da sessão
type SessionContext struct {
	WorkspaceSize    int                    `json:"workspace_size"`
	ExistingFiles    []string               `json:"existing_files"`
	Dependencies     []string               `json:"dependencies"`
	DetectedPatterns []DetectedPattern      `json:"detected_patterns"`
	PrimaryLanguage  string                 `json:"primary_language"`
	Complexity       float64                `json:"complexity"`
	HasLLMsFile      bool                   `json:"has_llms_file"`
	UserIntent       string                 `json:"user_intent"`
	Scope            string                 `json:"scope"`
	Exclusions       []string               `json:"exclusions"`
	Inclusions       []string               `json:"inclusions"`
	CustomRules      map[string]interface{} `json:"custom_rules"`
}

// GenerationOutcomes contém resultados da geração
type GenerationOutcomes struct {
	FilesGenerated     int      `json:"files_generated"`
	DirectoriesCreated int      `json:"directories_created"`
	LinesOfCode        int      `json:"lines_of_code"`
	TestsCoverage      float64  `json:"tests_coverage"`
	Dependencies       []string `json:"dependencies"`
	Warnings           []string `json:"warnings"`
	Errors             []string `json:"errors"`
	ExecutionTime      float64  `json:"execution_time"`
	MemoryUsage        float64  `json:"memory_usage"`
	APICallsCount      int      `json:"api_calls_count"`
	CacheHits          int      `json:"cache_hits"`
	CacheMisses        int      `json:"cache_misses"`
}

// PatternLearning representa aprendizado de padrões
type PatternLearning struct {
	Pattern            string                   `json:"pattern"`
	SuccessRate        float64                  `json:"success_rate"`
	AverageQuality     float64                  `json:"average_quality"`
	TotalUsage         int                      `json:"total_usage"`
	BestConfigurations []map[string]interface{} `json:"best_configurations"`
	CommonAdaptations  map[string]float64       `json:"common_adaptations"`
	OptimalPrompts     []string                 `json:"optimal_prompts"`
	TrendingKeywords   []string                 `json:"trending_keywords"`
	LastImprovement    time.Time                `json:"last_improvement"`
	LearningScore      float64                  `json:"learning_score"`
}

// UserPreferences contém preferências do usuário
type UserPreferences struct {
	PreferredLanguages    []string          `json:"preferred_languages"`
	PreferredProviders    []string          `json:"preferred_providers"`
	PreferredPatterns     []string          `json:"preferred_patterns"`
	QualityThreshold      float64           `json:"quality_threshold"`
	MaxComplexity         int               `json:"max_complexity"`
	PreferMinimal         bool              `json:"prefer_minimal"`
	PreferTests           bool              `json:"prefer_tests"`
	PreferDocker          bool              `json:"prefer_docker"`
	PreferTypeScript      bool              `json:"prefer_typescript"`
	PreferAsync           bool              `json:"prefer_async"`
	CustomPromptRules     []string          `json:"custom_prompt_rules"`
	ExclusionPatterns     []string          `json:"exclusion_patterns"`
	InclusionPatterns     []string          `json:"inclusion_patterns"`
	PreferredArchitecture string            `json:"preferred_architecture"`
	CodingStyle           map[string]string `json:"coding_style"`
	LastUpdated           time.Time         `json:"last_updated"`
}

// GenerationAnalytics contém analytics da geração
type GenerationAnalytics struct {
	TotalSessions          int                 `json:"total_sessions"`
	SuccessRate            float64             `json:"success_rate"`
	AverageQuality         float64             `json:"average_quality"`
	AverageRating          float64             `json:"average_rating"`
	MostUsedLanguages      map[string]int      `json:"most_used_languages"`
	MostUsedPatterns       map[string]int      `json:"most_used_patterns"`
	MostUsedProviders      map[string]int      `json:"most_used_providers"`
	QualityTrends          []QualityPoint      `json:"quality_trends"`
	PerformanceMetrics     PerformanceMetrics  `json:"performance_metrics"`
	ErrorAnalysis          ErrorAnalysis       `json:"error_analysis"`
	UserSatisfaction       SatisfactionMetrics `json:"user_satisfaction"`
	ImprovementSuggestions []string            `json:"improvement_suggestions"`
	LastUpdated            time.Time           `json:"last_updated"`
}

// PerformanceMetrics contém métricas de performance
type PerformanceMetrics struct {
	AverageResponseTime float64 `json:"average_response_time"`
	AverageTokensUsed   int     `json:"average_tokens_used"`
	LayeredModeUsage    float64 `json:"layered_mode_usage"`
	CacheHitRate        float64 `json:"cache_hit_rate"`
	ErrorRate           float64 `json:"error_rate"`
	RetryRate           float64 `json:"retry_rate"`
	ResourceUsage       float64 `json:"resource_usage"`
	ThroughputPerMinute float64 `json:"throughput_per_minute"`
}

// ErrorAnalysis contém análise de erros
type ErrorAnalysis struct {
	MostCommonErrors map[string]int `json:"most_common_errors"`
	ErrorsByProvider map[string]int `json:"errors_by_provider"`
	ErrorsByLanguage map[string]int `json:"errors_by_language"`
	ErrorsByPattern  map[string]int `json:"errors_by_pattern"`
	RecoveryRate     float64        `json:"recovery_rate"`
	ErrorTrends      []ErrorPoint   `json:"error_trends"`
}

// ErrorPoint representa um ponto na tendência de erros
type ErrorPoint struct {
	Timestamp time.Time `json:"timestamp"`
	ErrorType string    `json:"error_type"`
	Count     int       `json:"count"`
	Severity  string    `json:"severity"`
}

// SatisfactionMetrics contém métricas de satisfação
type SatisfactionMetrics struct {
	AverageRating       float64        `json:"average_rating"`
	RatingDistribution  map[string]int `json:"rating_distribution"`
	FeedbackSentiment   float64        `json:"feedback_sentiment"`
	ImprovementRequests map[string]int `json:"improvement_requests"`
	UserRetentionRate   float64        `json:"user_retention_rate"`
	RecommendationScore float64        `json:"recommendation_score"`
}

// LearningStats representa estatísticas do sistema de aprendizado
type LearningStats struct {
	TotalSessions      int               `json:"total_sessions"`
	SuccessfulSessions int               `json:"successful_sessions"`
	SuccessRate        float64           `json:"success_rate"`
	LearnedPatterns    []PatternLearning `json:"learned_patterns"`
	AverageQuality     float64           `json:"average_quality"`
	TopLanguages       []string          `json:"top_languages"`
	TopProviders       []string          `json:"top_providers"`
	RecentTrends       []string          `json:"recent_trends"`
	TotalDuration      time.Duration     `json:"total_duration"`
	AverageDuration    time.Duration     `json:"average_duration"`
	QualityImprovement float64           `json:"quality_improvement"`
	LearningProgress   float64           `json:"learning_progress"`
}

// NewLearningSystem cria um novo sistema de aprendizado
func NewLearningSystem(historyPath string) *LearningSystem {
	ls := &LearningSystem{
		historyPath: historyPath,
		sessions:    make([]GenerationSession, 0),
		patterns:    make(map[string]PatternLearning),
		preferences: UserPreferences{
			QualityThreshold: 0.7,
			MaxComplexity:    50,
			PreferMinimal:    false,
			PreferTests:      true,
			PreferDocker:     false,
			PreferTypeScript: false,
			PreferAsync:      true,
			LastUpdated:      time.Now(),
		},
		analytics: GenerationAnalytics{
			LastUpdated: time.Now(),
		},
	}

	// Carregar histórico existente
	if err := ls.loadHistory(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao carregar histórico: %v\n", err)
	}

	return ls
}

// RecordSession registra uma nova sessão
func (ls *LearningSystem) RecordSession(session GenerationSession) error {
	// Gerar ID único
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%s-%d", session.Language, session.Description, session.Timestamp.Unix())))
	session.ID = fmt.Sprintf("%x", hasher.Sum(nil))[:16]

	// Adicionar à lista
	ls.sessions = append(ls.sessions, session)

	// Atualizar aprendizado de padrões
	ls.updatePatternLearning(session)

	// Atualizar preferências
	ls.updatePreferences(session)

	// Atualizar analytics
	ls.updateAnalytics(session)

	// Salvar histórico
	if err := ls.saveHistory(); err != nil {
		return fmt.Errorf("erro ao salvar histórico: %v", err)
	}

	fmt.Printf("📚 Sessão registrada no histórico: %s\n", session.ID)
	return nil
}

// GetOptimalConfiguration retorna configuração ótima baseada no histórico
func (ls *LearningSystem) GetOptimalConfiguration(language, projectType, description string) (*OptimalConfiguration, error) {
	// Analisar histórico para encontrar configurações similares
	similarSessions := ls.findSimilarSessions(language, projectType, description)

	if len(similarSessions) == 0 {
		return ls.getDefaultConfiguration(language, projectType), nil
	}

	// Agregar configurações bem-sucedidas
	config := &OptimalConfiguration{
		Language:        language,
		ProjectType:     projectType,
		Confidence:      0.0,
		Adaptations:     make(map[string]interface{}),
		Suggestions:     make([]string, 0),
		Warnings:        make([]string, 0),
		BasedOnSessions: make([]string, 0),
	}

	// Calcular configuração ótima
	totalWeight := 0.0
	for _, session := range similarSessions {
		weight := ls.calculateSessionWeight(session)
		totalWeight += weight

		// Agregar adaptações
		for key, value := range session.Adaptations {
			if existing, exists := config.Adaptations[key]; exists {
				// Fazer média ponderada
				config.Adaptations[key] = ls.weightedAverage(existing, value, weight)
			} else {
				config.Adaptations[key] = value
			}
		}

		config.BasedOnSessions = append(config.BasedOnSessions, session.ID)
	}

	// Calcular confiança
	config.Confidence = ls.calculateConfidence(similarSessions)

	// Gerar sugestões
	config.Suggestions = ls.generateSuggestions(similarSessions)

	// Gerar avisos
	config.Warnings = ls.generateWarnings(similarSessions)

	// Aplicar preferências do usuário
	config = ls.applyUserPreferences(config)

	return config, nil
}

// OptimalConfiguration representa uma configuração ótima
type OptimalConfiguration struct {
	Language        string                 `json:"language"`
	ProjectType     string                 `json:"project_type"`
	Confidence      float64                `json:"confidence"`
	Adaptations     map[string]interface{} `json:"adaptations"`
	Suggestions     []string               `json:"suggestions"`
	Warnings        []string               `json:"warnings"`
	BasedOnSessions []string               `json:"based_on_sessions"`
	QualityScore    float64                `json:"quality_score"`
	EstimatedTime   time.Duration          `json:"estimated_time"`
	EstimatedTokens int                    `json:"estimated_tokens"`
	RiskFactors     []string               `json:"risk_factors"`
}

// findSimilarSessions encontra sessões similares
func (ls *LearningSystem) findSimilarSessions(language, projectType, description string) []GenerationSession {
	var similar []GenerationSession

	for _, session := range ls.sessions {
		if !session.Success || session.Quality < 0.6 {
			continue
		}

		similarity := ls.calculateSimilarity(session, language, projectType, description)
		if similarity > 0.5 {
			similar = append(similar, session)
		}
	}

	// Ordenar por qualidade e recência
	sort.Slice(similar, func(i, j int) bool {
		scoreI := similar[i].Quality*0.7 + ls.getRecencyScore(similar[i].Timestamp)*0.3
		scoreJ := similar[j].Quality*0.7 + ls.getRecencyScore(similar[j].Timestamp)*0.3
		return scoreI > scoreJ
	})

	// Limitar aos 10 mais relevantes
	if len(similar) > 10 {
		similar = similar[:10]
	}

	return similar
}

// calculateSimilarity calcula similaridade entre sessões
func (ls *LearningSystem) calculateSimilarity(session GenerationSession, language, projectType, description string) float64 {
	similarity := 0.0

	// Similaridade de linguagem
	if session.Language == language {
		similarity += 0.3
	}

	// Similaridade de tipo de projeto
	if session.ProjectType == projectType {
		similarity += 0.2
	}

	// Similaridade de descrição (baseada em palavras-chave)
	descSimilarity := ls.calculateTextSimilarity(session.Description, description)
	similarity += descSimilarity * 0.3

	// Similaridade de contexto
	contextSimilarity := ls.calculateContextSimilarity(session.Context, language, projectType)
	similarity += contextSimilarity * 0.2

	return similarity
}

// calculateTextSimilarity calcula similaridade entre textos
func (ls *LearningSystem) calculateTextSimilarity(text1, text2 string) float64 {
	words1 := strings.Fields(strings.ToLower(text1))
	words2 := strings.Fields(strings.ToLower(text2))

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Criar sets de palavras
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, word := range words1 {
		set1[word] = true
	}
	for _, word := range words2 {
		set2[word] = true
	}

	// Calcular interseção
	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}

	// Calcular união
	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// calculateContextSimilarity calcula similaridade de contexto
func (ls *LearningSystem) calculateContextSimilarity(context SessionContext, language, projectType string) float64 {
	similarity := 0.0

	// Similaridade de linguagem primária
	if context.PrimaryLanguage == language {
		similarity += 0.3
	}

	// Similaridade de escopo
	if context.Scope != "" {
		similarity += 0.2
	}

	// Similaridade de padrões
	for _, pattern := range context.DetectedPatterns {
		if string(pattern.Type) == projectType {
			similarity += 0.5
		}
	}

	return similarity
}

// calculateSessionWeight calcula peso de uma sessão
func (ls *LearningSystem) calculateSessionWeight(session GenerationSession) float64 {
	weight := 1.0

	// Peso por qualidade
	weight *= session.Quality

	// Peso por rating do usuário
	if session.UserRating > 0 {
		weight *= (float64(session.UserRating) / 5.0)
	}

	// Peso por recência
	weight *= ls.getRecencyScore(session.Timestamp)

	// Peso por sucesso
	if session.Success {
		weight *= 1.2
	}

	return weight
}

// getRecencyScore calcula score de recência
func (ls *LearningSystem) getRecencyScore(timestamp time.Time) float64 {
	daysSince := time.Since(timestamp).Hours() / 24

	// Decaimento exponencial
	return 1.0 / (1.0 + daysSince/30.0)
}

// calculateConfidence calcula confiança da configuração
func (ls *LearningSystem) calculateConfidence(sessions []GenerationSession) float64 {
	if len(sessions) == 0 {
		return 0.0
	}

	// Confiança baseada em quantidade e qualidade
	confidence := 0.0
	for _, session := range sessions {
		confidence += session.Quality * float64(session.UserRating) / 5.0
	}

	confidence /= float64(len(sessions))

	// Boost por quantidade de sessões
	quantityBoost := 1.0 + (float64(len(sessions))-1.0)*0.1
	if quantityBoost > 1.5 {
		quantityBoost = 1.5
	}

	return confidence * quantityBoost
}

// generateSuggestions gera sugestões baseadas no histórico
func (ls *LearningSystem) generateSuggestions(sessions []GenerationSession) []string {
	suggestions := make([]string, 0)

	// Análise de padrões bem-sucedidos
	patternCounts := make(map[string]int)
	for _, session := range sessions {
		for _, pattern := range session.Patterns {
			patternCounts[pattern]++
		}
	}

	// Sugestões baseadas em padrões comuns
	for pattern, count := range patternCounts {
		if count >= len(sessions)/2 {
			suggestions = append(suggestions, fmt.Sprintf("Considere usar padrão %s (usado em %d/%d sessões similares)", pattern, count, len(sessions)))
		}
	}

	// Análise de adaptações bem-sucedidas
	adaptationCounts := make(map[string]int)
	for _, session := range sessions {
		for adaptation := range session.Adaptations {
			adaptationCounts[adaptation]++
		}
	}

	// Sugestões baseadas em adaptações
	for adaptation, count := range adaptationCounts {
		if count >= len(sessions)/2 {
			suggestions = append(suggestions, fmt.Sprintf("Recomenda-se ativar %s (usado em %d/%d sessões similares)", adaptation, count, len(sessions)))
		}
	}

	// Análise de qualidade
	avgQuality := 0.0
	for _, session := range sessions {
		avgQuality += session.Quality
	}
	avgQuality /= float64(len(sessions))

	if avgQuality > 0.8 {
		suggestions = append(suggestions, "Configuração baseada em sessões de alta qualidade")
	}

	return suggestions
}

// generateWarnings gera avisos baseados no histórico
func (ls *LearningSystem) generateWarnings(sessions []GenerationSession) []string {
	warnings := make([]string, 0)

	// Análise de erros comuns
	errorCounts := make(map[string]int)
	for _, session := range sessions {
		for _, error := range session.Outcomes.Errors {
			errorCounts[error]++
		}
	}

	// Avisos baseados em erros
	for error, count := range errorCounts {
		if count >= len(sessions)/3 {
			warnings = append(warnings, fmt.Sprintf("Atenção: %s (ocorreu em %d/%d sessões)", error, count, len(sessions)))
		}
	}

	// Análise de performance
	totalSessions := len(sessions)
	layeredSessions := 0
	for _, session := range sessions {
		if session.LayeredMode {
			layeredSessions++
		}
	}

	if layeredSessions > totalSessions/2 {
		warnings = append(warnings, "Projeto complexo: modo em camadas frequentemente necessário")
	}

	return warnings
}

// applyUserPreferences aplica preferências do usuário
func (ls *LearningSystem) applyUserPreferences(config *OptimalConfiguration) *OptimalConfiguration {
	// Aplicar preferências de qualidade
	if config.QualityScore < ls.preferences.QualityThreshold {
		config.Warnings = append(config.Warnings, "Qualidade abaixo do threshold preferido")
	}

	// Aplicar preferências de minimalismo
	if ls.preferences.PreferMinimal {
		config.Adaptations["scope"] = "minimal"
		config.Adaptations["exclude_extras"] = true
	}

	// Aplicar preferências de testes
	if ls.preferences.PreferTests {
		config.Adaptations["include_tests"] = true
	}

	// Aplicar preferências de Docker
	if ls.preferences.PreferDocker {
		config.Adaptations["include_docker"] = true
	}

	// Aplicar preferências de TypeScript
	if ls.preferences.PreferTypeScript && (config.Language == "javascript" || config.Language == "typescript") {
		config.Adaptations["prefer_typescript"] = true
	}

	// Aplicar estilo de código
	for key, value := range ls.preferences.CodingStyle {
		config.Adaptations[key] = value
	}

	return config
}

// getDefaultConfiguration retorna configuração padrão
func (ls *LearningSystem) getDefaultConfiguration(language, projectType string) *OptimalConfiguration {
	return &OptimalConfiguration{
		Language:     language,
		ProjectType:  projectType,
		Confidence:   0.5,
		Adaptations:  make(map[string]interface{}),
		Suggestions:  []string{"Configuração padrão - sem histórico específico"},
		Warnings:     []string{"Recomenda-se fornecer feedback após a geração"},
		QualityScore: 0.7,
	}
}

// updatePatternLearning atualiza aprendizado de padrões
func (ls *LearningSystem) updatePatternLearning(session GenerationSession) {
	for _, pattern := range session.Patterns {
		learning, exists := ls.patterns[pattern]
		if !exists {
			learning = PatternLearning{
				Pattern:            pattern,
				BestConfigurations: make([]map[string]interface{}, 0),
				CommonAdaptations:  make(map[string]float64),
				OptimalPrompts:     make([]string, 0),
				TrendingKeywords:   make([]string, 0),
			}
		}

		// Atualizar estatísticas
		learning.TotalUsage++

		// Atualizar taxa de sucesso
		if session.Success {
			learning.SuccessRate = (learning.SuccessRate*float64(learning.TotalUsage-1) + 1.0) / float64(learning.TotalUsage)
		} else {
			learning.SuccessRate = (learning.SuccessRate * float64(learning.TotalUsage-1)) / float64(learning.TotalUsage)
		}

		// Atualizar qualidade média
		learning.AverageQuality = (learning.AverageQuality*float64(learning.TotalUsage-1) + session.Quality) / float64(learning.TotalUsage)

		// Atualizar adaptações comuns
		for adaptation := range session.Adaptations {
			learning.CommonAdaptations[adaptation] = learning.CommonAdaptations[adaptation] + 1.0
		}

		// Atualizar score de aprendizado
		learning.LearningScore = learning.SuccessRate*0.4 + learning.AverageQuality*0.6
		learning.LastImprovement = time.Now()

		ls.patterns[pattern] = learning
	}
}

// updatePreferences atualiza preferências do usuário
func (ls *LearningSystem) updatePreferences(session GenerationSession) {
	// Atualizar linguagens preferidas
	ls.updateLanguagePreference(session.Language, session.Success, session.Quality)

	// Atualizar provedores preferidos
	ls.updateProviderPreference(session.Provider, session.Success, session.Quality)

	// Atualizar padrões preferidos
	for _, pattern := range session.Patterns {
		ls.updatePatternPreference(pattern, session.Success, session.Quality)
	}

	// Atualizar threshold de qualidade baseado em ratings
	if session.UserRating > 0 {
		ls.preferences.QualityThreshold = (ls.preferences.QualityThreshold * 0.9) + (session.Quality * 0.1)
	}

	ls.preferences.LastUpdated = time.Now()
}

// updateLanguagePreference atualiza preferência de linguagem
func (ls *LearningSystem) updateLanguagePreference(language string, success bool, quality float64) {
	if success && quality > 0.7 {
		// Adicionar à lista de preferidas se não existir
		found := false
		for _, lang := range ls.preferences.PreferredLanguages {
			if lang == language {
				found = true
				break
			}
		}
		if !found {
			ls.preferences.PreferredLanguages = append(ls.preferences.PreferredLanguages, language)
		}
	}
}

// updateProviderPreference atualiza preferência de provedor
func (ls *LearningSystem) updateProviderPreference(provider string, success bool, quality float64) {
	if success && quality > 0.7 {
		found := false
		for _, prov := range ls.preferences.PreferredProviders {
			if prov == provider {
				found = true
				break
			}
		}
		if !found {
			ls.preferences.PreferredProviders = append(ls.preferences.PreferredProviders, provider)
		}
	}
}

// updatePatternPreference atualiza preferência de padrão
func (ls *LearningSystem) updatePatternPreference(pattern string, success bool, quality float64) {
	if success && quality > 0.7 {
		found := false
		for _, pat := range ls.preferences.PreferredPatterns {
			if pat == pattern {
				found = true
				break
			}
		}
		if !found {
			ls.preferences.PreferredPatterns = append(ls.preferences.PreferredPatterns, pattern)
		}
	}
}

// updateAnalytics atualiza analytics
func (ls *LearningSystem) updateAnalytics(session GenerationSession) {
	ls.analytics.TotalSessions++

	// Atualizar taxa de sucesso
	if session.Success {
		ls.analytics.SuccessRate = (ls.analytics.SuccessRate*float64(ls.analytics.TotalSessions-1) + 1.0) / float64(ls.analytics.TotalSessions)
	} else {
		ls.analytics.SuccessRate = (ls.analytics.SuccessRate * float64(ls.analytics.TotalSessions-1)) / float64(ls.analytics.TotalSessions)
	}

	// Atualizar qualidade média
	ls.analytics.AverageQuality = (ls.analytics.AverageQuality*float64(ls.analytics.TotalSessions-1) + session.Quality) / float64(ls.analytics.TotalSessions)

	// Atualizar rating médio
	if session.UserRating > 0 {
		ls.analytics.AverageRating = (ls.analytics.AverageRating*float64(ls.analytics.TotalSessions-1) + float64(session.UserRating)) / float64(ls.analytics.TotalSessions)
	}

	// Atualizar contadores
	if ls.analytics.MostUsedLanguages == nil {
		ls.analytics.MostUsedLanguages = make(map[string]int)
	}
	ls.analytics.MostUsedLanguages[session.Language]++

	if ls.analytics.MostUsedPatterns == nil {
		ls.analytics.MostUsedPatterns = make(map[string]int)
	}
	for _, pattern := range session.Patterns {
		ls.analytics.MostUsedPatterns[pattern]++
	}

	if ls.analytics.MostUsedProviders == nil {
		ls.analytics.MostUsedProviders = make(map[string]int)
	}
	ls.analytics.MostUsedProviders[session.Provider]++

	// Adicionar ponto de qualidade
	ls.analytics.QualityTrends = append(ls.analytics.QualityTrends, QualityPoint{
		Timestamp: session.Timestamp,
		Quality:   session.Quality,
		Success:   session.Success,
	})

	ls.analytics.LastUpdated = time.Now()
}

// loadHistory carrega histórico do disco
func (ls *LearningSystem) loadHistory() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(ls.historyPath, 0755); err != nil {
		return err
	}

	// Carregar sessões
	sessionsFile := filepath.Join(ls.historyPath, "sessions.json")
	if data, err := os.ReadFile(sessionsFile); err == nil {
		if err := json.Unmarshal(data, &ls.sessions); err != nil {
			return err
		}
	}

	// Carregar padrões
	patternsFile := filepath.Join(ls.historyPath, "patterns.json")
	if data, err := os.ReadFile(patternsFile); err == nil {
		if err := json.Unmarshal(data, &ls.patterns); err != nil {
			return err
		}
	}

	// Carregar preferências
	preferencesFile := filepath.Join(ls.historyPath, "preferences.json")
	if data, err := os.ReadFile(preferencesFile); err == nil {
		if err := json.Unmarshal(data, &ls.preferences); err != nil {
			return err
		}
	}

	// Carregar analytics
	analyticsFile := filepath.Join(ls.historyPath, "analytics.json")
	if data, err := os.ReadFile(analyticsFile); err == nil {
		if err := json.Unmarshal(data, &ls.analytics); err != nil {
			return err
		}
	}

	return nil
}

// saveHistory salva histórico no disco
func (ls *LearningSystem) saveHistory() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(ls.historyPath, 0755); err != nil {
		return err
	}

	// Salvar sessões
	sessionsFile := filepath.Join(ls.historyPath, "sessions.json")
	if data, err := json.MarshalIndent(ls.sessions, "", "  "); err == nil {
		if err := os.WriteFile(sessionsFile, data, 0644); err != nil {
			return err
		}
	}

	// Salvar padrões
	patternsFile := filepath.Join(ls.historyPath, "patterns.json")
	if data, err := json.MarshalIndent(ls.patterns, "", "  "); err == nil {
		if err := os.WriteFile(patternsFile, data, 0644); err != nil {
			return err
		}
	}

	// Salvar preferências
	preferencesFile := filepath.Join(ls.historyPath, "preferences.json")
	if data, err := json.MarshalIndent(ls.preferences, "", "  "); err == nil {
		if err := os.WriteFile(preferencesFile, data, 0644); err != nil {
			return err
		}
	}

	// Salvar analytics
	analyticsFile := filepath.Join(ls.historyPath, "analytics.json")
	if data, err := json.MarshalIndent(ls.analytics, "", "  "); err == nil {
		if err := os.WriteFile(analyticsFile, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetAnalytics retorna analytics do sistema
func (ls *LearningSystem) GetAnalytics() GenerationAnalytics {
	return ls.analytics
}

// GetPreferences retorna preferências do usuário
func (ls *LearningSystem) GetPreferences() UserPreferences {
	return ls.preferences
}

// GetPatternLearning retorna aprendizado de padrões
func (ls *LearningSystem) GetPatternLearning() map[string]PatternLearning {
	return ls.patterns
}

// GetRecentSessions retorna sessões recentes
func (ls *LearningSystem) GetRecentSessions(limit int) []GenerationSession {
	if limit <= 0 || limit > len(ls.sessions) {
		limit = len(ls.sessions)
	}

	// Ordenar por timestamp
	sessions := make([]GenerationSession, len(ls.sessions))
	copy(sessions, ls.sessions)

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Timestamp.After(sessions[j].Timestamp)
	})

	return sessions[:limit]
}

// weightedAverage calcula média ponderada
func (ls *LearningSystem) weightedAverage(existing, new interface{}, weight float64) interface{} {
	// Implementação simplificada - na prática seria mais complexa
	return new
}

// GenerateReport gera relatório do sistema de aprendizado
func (ls *LearningSystem) GenerateReport() string {
	report := "📊 RELATÓRIO DO SISTEMA DE APRENDIZADO\n\n"

	// Estatísticas gerais
	report += fmt.Sprintf("Total de Sessões: %d\n", ls.analytics.TotalSessions)
	report += fmt.Sprintf("Taxa de Sucesso: %.1f%%\n", ls.analytics.SuccessRate*100)
	report += fmt.Sprintf("Qualidade Média: %.2f\n", ls.analytics.AverageQuality)
	report += fmt.Sprintf("Rating Médio: %.1f/5\n", ls.analytics.AverageRating)

	// Linguagens mais usadas
	report += "\n🗣️ Linguagens Mais Usadas:\n"
	for lang, count := range ls.analytics.MostUsedLanguages {
		report += fmt.Sprintf("   • %s: %d sessões\n", lang, count)
	}

	// Padrões mais usados
	report += "\n🏗️ Padrões Mais Usados:\n"
	for pattern, count := range ls.analytics.MostUsedPatterns {
		report += fmt.Sprintf("   • %s: %d sessões\n", pattern, count)
	}

	// Preferências do usuário
	report += "\n⚙️ Preferências Detectadas:\n"
	report += fmt.Sprintf("   • Threshold de Qualidade: %.2f\n", ls.preferences.QualityThreshold)
	report += fmt.Sprintf("   • Prefere Minimal: %v\n", ls.preferences.PreferMinimal)
	report += fmt.Sprintf("   • Prefere Testes: %v\n", ls.preferences.PreferTests)
	report += fmt.Sprintf("   • Prefere Docker: %v\n", ls.preferences.PreferDocker)

	return report
}

// GetLearningStats retorna estatísticas do sistema de aprendizado
func (ls *LearningSystem) GetLearningStats() LearningStats {
	stats := LearningStats{
		TotalSessions:      len(ls.sessions),
		SuccessfulSessions: 0,
		AverageQuality:     0.0,
		TopLanguages:       make([]string, 0),
		TopProviders:       make([]string, 0),
		RecentTrends:       make([]string, 0),
		TotalDuration:      0,
		QualityImprovement: 0.0,
		LearningProgress:   0.0,
	}

	if len(ls.sessions) == 0 {
		return stats
	}

	// Calcular estatísticas básicas
	var totalQuality float64
	languageCount := make(map[string]int)
	providerCount := make(map[string]int)

	for _, session := range ls.sessions {
		if session.Success {
			stats.SuccessfulSessions++
		}
		totalQuality += session.Quality
		stats.TotalDuration += session.Duration

		languageCount[session.Language]++
		providerCount[session.Provider]++
	}

	stats.SuccessRate = float64(stats.SuccessfulSessions) / float64(stats.TotalSessions)
	stats.AverageQuality = totalQuality / float64(stats.TotalSessions)
	stats.AverageDuration = stats.TotalDuration / time.Duration(stats.TotalSessions)

	// Top languages
	type langCount struct {
		lang  string
		count int
	}
	var langs []langCount
	for lang, count := range languageCount {
		langs = append(langs, langCount{lang, count})
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].count > langs[j].count })
	for i, lang := range langs {
		if i >= 5 { // Top 5
			break
		}
		stats.TopLanguages = append(stats.TopLanguages, lang.lang)
	}

	// Top providers
	type provCount struct {
		prov  string
		count int
	}
	var provs []provCount
	for prov, count := range providerCount {
		provs = append(provs, provCount{prov, count})
	}
	sort.Slice(provs, func(i, j int) bool { return provs[i].count > provs[j].count })
	for i, prov := range provs {
		if i >= 3 { // Top 3
			break
		}
		stats.TopProviders = append(stats.TopProviders, prov.prov)
	}

	// Padrões aprendidos
	for _, pattern := range ls.patterns {
		stats.LearnedPatterns = append(stats.LearnedPatterns, pattern)
	}

	// Tendências recentes (baseado nas últimas 10 sessões)
	recentSessions := ls.sessions
	if len(recentSessions) > 10 {
		recentSessions = recentSessions[len(recentSessions)-10:]
	}

	recentLangs := make(map[string]int)
	for _, session := range recentSessions {
		recentLangs[session.Language]++
	}

	for lang, count := range recentLangs {
		if count >= 2 { // Aparece pelo menos 2x nas últimas sessões
			stats.RecentTrends = append(stats.RecentTrends, fmt.Sprintf("Aumento no uso de %s", lang))
		}
	}

	// Melhoria de qualidade (comparar primeira metade com segunda metade)
	if len(ls.sessions) >= 4 {
		midPoint := len(ls.sessions) / 2
		firstHalf := ls.sessions[:midPoint]
		secondHalf := ls.sessions[midPoint:]

		var firstQuality, secondQuality float64
		for _, session := range firstHalf {
			firstQuality += session.Quality
		}
		for _, session := range secondHalf {
			secondQuality += session.Quality
		}

		firstAvg := firstQuality / float64(len(firstHalf))
		secondAvg := secondQuality / float64(len(secondHalf))

		stats.QualityImprovement = ((secondAvg - firstAvg) / firstAvg) * 100
	}

	// Progresso de aprendizado (baseado na diversidade de padrões)
	stats.LearningProgress = float64(len(ls.patterns)) / 10.0 * 100 // Max 10 padrões = 100%
	if stats.LearningProgress > 100 {
		stats.LearningProgress = 100
	}

	return stats
}
