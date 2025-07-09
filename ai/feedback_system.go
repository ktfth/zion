package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FeedbackSystem implementa sistema de feedback e melhoria contínua
type FeedbackSystem struct {
	feedbackPath   string
	feedbacks      []UserFeedback
	improvements   []SystemImprovement
	analytics      FeedbackAnalytics
	autoTuning     *AutoTuningEngine
	qualityMonitor *QualityMonitor
	predictor      *OutcomePredictor
}

// UserFeedback representa feedback do usuário
type UserFeedback struct {
	ID             string                 `json:"id"`
	SessionID      string                 `json:"session_id"`
	Timestamp      time.Time              `json:"timestamp"`
	UserID         string                 `json:"user_id"`
	Rating         int                    `json:"rating"` // 1-5
	Category       FeedbackCategory       `json:"category"`
	Type           FeedbackType           `json:"type"`
	Content        string                 `json:"content"`
	Tags           []string               `json:"tags"`
	Context        FeedbackContext        `json:"context"`
	Sentiment      float64                `json:"sentiment"` // -1 a +1
	Priority       int                    `json:"priority"`  // 1-5
	Status         FeedbackStatus         `json:"status"`
	ActionItems    []ActionItem           `json:"action_items"`
	Resolution     *FeedbackResolution    `json:"resolution"`
	Impact         float64                `json:"impact"`
	Effort         float64                `json:"effort"`
	AutoClassified bool                   `json:"auto_classified"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// FeedbackCategory define categorias de feedback
type FeedbackCategory string

const (
	CategoryGeneral       FeedbackCategory = "general"
	CategoryQuality       FeedbackCategory = "quality"
	CategoryPerformance   FeedbackCategory = "performance"
	CategoryUsability     FeedbackCategory = "usability"
	CategoryBug           FeedbackCategory = "bug"
	CategoryFeature       FeedbackCategory = "feature"
	CategoryDocumentation FeedbackCategory = "documentation"
	CategorySecurity      FeedbackCategory = "security"
	CategoryArchitecture  FeedbackCategory = "architecture"
	CategoryUI            FeedbackCategory = "ui"
)

// FeedbackType define tipos de feedback
type FeedbackType string

const (
	TypePositive       FeedbackType = "positive"
	TypeNegative       FeedbackType = "negative"
	TypeSuggestion     FeedbackType = "suggestion"
	TypeBugReport      FeedbackType = "bug_report"
	TypeFeatureRequest FeedbackType = "feature_request"
	TypeQuestion       FeedbackType = "question"
	TypeCompliment     FeedbackType = "compliment"
	TypeComplaint      FeedbackType = "complaint"
)

// FeedbackStatus define status do feedback
type FeedbackStatus string

const (
	StatusNew        FeedbackStatus = "new"
	StatusReviewed   FeedbackStatus = "reviewed"
	StatusInProgress FeedbackStatus = "in_progress"
	StatusResolved   FeedbackStatus = "resolved"
	StatusClosed     FeedbackStatus = "closed"
	StatusDeferred   FeedbackStatus = "deferred"
)

// FeedbackContext contém contexto do feedback
type FeedbackContext struct {
	Language       string                 `json:"language"`
	ProjectType    string                 `json:"project_type"`
	Provider       string                 `json:"provider"`
	LayeredMode    bool                   `json:"layered_mode"`
	GenerationTime time.Duration          `json:"generation_time"`
	TokensUsed     int                    `json:"tokens_used"`
	Quality        float64                `json:"quality"`
	Success        bool                   `json:"success"`
	Errors         []string               `json:"errors"`
	UserProfile    UserProfile            `json:"user_profile"`
	Environment    map[string]interface{} `json:"environment"`
}

// UserProfile representa perfil do usuário
type UserProfile struct {
	ExperienceLevel    string                 `json:"experience_level"` // beginner, intermediate, advanced
	PreferredLanguages []string               `json:"preferred_languages"`
	UsageFrequency     string                 `json:"usage_frequency"` // daily, weekly, monthly
	PrimaryUseCase     string                 `json:"primary_use_case"`
	PainPoints         []string               `json:"pain_points"`
	Preferences        map[string]interface{} `json:"preferences"`
}

// ActionItem representa ação a ser tomada
type ActionItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Effort      string    `json:"effort"` // low, medium, high
	Impact      string    `json:"impact"` // low, medium, high
	Category    string    `json:"category"`
	Assignee    string    `json:"assignee"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Progress    float64   `json:"progress"`
}

// FeedbackResolution representa resolução do feedback
type FeedbackResolution struct {
	ResolvedBy   string    `json:"resolved_by"`
	ResolvedAt   time.Time `json:"resolved_at"`
	Resolution   string    `json:"resolution"`
	Changes      []string  `json:"changes"`
	TestResults  []string  `json:"test_results"`
	UserNotified bool      `json:"user_notified"`
	Satisfaction int       `json:"satisfaction"` // 1-5
}

// SystemImprovement representa melhorias do sistema
type SystemImprovement struct {
	ID              string              `json:"id"`
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	Category        ImprovementCategory `json:"category"`
	Priority        int                 `json:"priority"`
	Impact          float64             `json:"impact"`
	Effort          float64             `json:"effort"`
	ROI             float64             `json:"roi"`
	Status          ImprovementStatus   `json:"status"`
	Implementation  *Implementation     `json:"implementation"`
	BasedOnFeedback []string            `json:"based_on_feedback"`
	Metrics         ImprovementMetrics  `json:"metrics"`
	Timeline        Timeline            `json:"timeline"`
	Dependencies    []string            `json:"dependencies"`
	RiskLevel       string              `json:"risk_level"`
	Testing         TestingPlan         `json:"testing"`
}

// ImprovementCategory define categorias de melhoria
type ImprovementCategory string

const (
	ImprovementAlgorithm      ImprovementCategory = "algorithm"
	ImprovementPerformance    ImprovementCategory = "performance"
	ImprovementQuality        ImprovementCategory = "quality"
	ImprovementUsability      ImprovementCategory = "usability"
	ImprovementFeature        ImprovementCategory = "feature"
	ImprovementBugFix         ImprovementCategory = "bug_fix"
	ImprovementSecurity       ImprovementCategory = "security"
	ImprovementInfrastructure ImprovementCategory = "infrastructure"
)

// ImprovementStatus define status da melhoria
type ImprovementStatus string

const (
	ImprovementProposed      ImprovementStatus = "proposed"
	ImprovementApproved      ImprovementStatus = "approved"
	ImprovementInDevelopment ImprovementStatus = "in_development"
	ImprovementTesting       ImprovementStatus = "testing"
	ImprovementDeployed      ImprovementStatus = "deployed"
	ImprovementCompleted     ImprovementStatus = "completed"
	ImprovementRejected      ImprovementStatus = "rejected"
)

// Implementation detalhes de implementação
type Implementation struct {
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Progress      float64   `json:"progress"`
	Developer     string    `json:"developer"`
	CodeChanges   []string  `json:"code_changes"`
	TestsAdded    []string  `json:"tests_added"`
	Documentation []string  `json:"documentation"`
	Rollback      *Rollback `json:"rollback"`
}

// Rollback informações de rollback
type Rollback struct {
	Available  bool      `json:"available"`
	Reason     string    `json:"reason"`
	ExecutedAt time.Time `json:"executed_at"`
	RollbackBy string    `json:"rollback_by"`
	Recovery   string    `json:"recovery"`
}

// ImprovementMetrics métricas da melhoria
type ImprovementMetrics struct {
	Before           map[string]float64 `json:"before"`
	After            map[string]float64 `json:"after"`
	Improvement      map[string]float64 `json:"improvement"`
	UserSatisfaction float64            `json:"user_satisfaction"`
	Performance      float64            `json:"performance"`
	Quality          float64            `json:"quality"`
	Adoption         float64            `json:"adoption"`
}

// Timeline cronograma da melhoria
type Timeline struct {
	Planning    TimelinePhase `json:"planning"`
	Development TimelinePhase `json:"development"`
	Testing     TimelinePhase `json:"testing"`
	Deployment  TimelinePhase `json:"deployment"`
	Monitoring  TimelinePhase `json:"monitoring"`
}

// TimelinePhase fase do cronograma
type TimelinePhase struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status"`
	Progress     float64   `json:"progress"`
	Deliverables []string  `json:"deliverables"`
	Blockers     []string  `json:"blockers"`
}

// TestingPlan plano de testes
type TestingPlan struct {
	UnitTests        []string `json:"unit_tests"`
	IntegrationTests []string `json:"integration_tests"`
	E2ETests         []string `json:"e2e_tests"`
	PerformanceTests []string `json:"performance_tests"`
	SecurityTests    []string `json:"security_tests"`
	UserAcceptance   []string `json:"user_acceptance"`
	Coverage         float64  `json:"coverage"`
	PassRate         float64  `json:"pass_rate"`
}

// FeedbackAnalytics analytics do feedback
type FeedbackAnalytics struct {
	TotalFeedbacks        int                      `json:"total_feedbacks"`
	AverageRating         float64                  `json:"average_rating"`
	SentimentDistribution map[string]int           `json:"sentiment_distribution"`
	CategoryDistribution  map[FeedbackCategory]int `json:"category_distribution"`
	TypeDistribution      map[FeedbackType]int     `json:"type_distribution"`
	TrendAnalysis         TrendAnalysis            `json:"trend_analysis"`
	TopIssues             []IssueFrequency         `json:"top_issues"`
	TopSuggestions        []SuggestionFrequency    `json:"top_suggestions"`
	UserSegmentation      UserSegmentation         `json:"user_segmentation"`
	SatisfactionTrends    []SatisfactionPoint      `json:"satisfaction_trends"`
	ResolutionMetrics     ResolutionMetrics        `json:"resolution_metrics"`
	PredictiveInsights    PredictiveInsights       `json:"predictive_insights"`
}

// TrendAnalysis análise de tendências
type TrendAnalysis struct {
	RatingTrend      []RatingPoint        `json:"rating_trend"`
	VolumeTrend      []VolumePoint        `json:"volume_trend"`
	CategoryTrends   map[string][]float64 `json:"category_trends"`
	SentimentTrend   []SentimentPoint     `json:"sentiment_trend"`
	SeasonalPatterns map[string]float64   `json:"seasonal_patterns"`
	GrowthRate       float64              `json:"growth_rate"`
	Forecasting      ForecastData         `json:"forecasting"`
}

// Tipos de pontos para análise
type RatingPoint struct {
	Date   time.Time `json:"date"`
	Rating float64   `json:"rating"`
	Count  int       `json:"count"`
}

type VolumePoint struct {
	Date   time.Time `json:"date"`
	Volume int       `json:"volume"`
}

type SentimentPoint struct {
	Date      time.Time `json:"date"`
	Sentiment float64   `json:"sentiment"`
	Positive  int       `json:"positive"`
	Negative  int       `json:"negative"`
	Neutral   int       `json:"neutral"`
}

type SatisfactionPoint struct {
	Date         time.Time `json:"date"`
	Satisfaction float64   `json:"satisfaction"`
	NPS          float64   `json:"nps"`
	CSAT         float64   `json:"csat"`
}

// ForecastData dados de previsão
type ForecastData struct {
	NextMonth   ForecastPeriod `json:"next_month"`
	NextQuarter ForecastPeriod `json:"next_quarter"`
	NextYear    ForecastPeriod `json:"next_year"`
	Confidence  float64        `json:"confidence"`
	Methodology string         `json:"methodology"`
}

// ForecastPeriod período de previsão
type ForecastPeriod struct {
	ExpectedVolume     int      `json:"expected_volume"`
	ExpectedRating     float64  `json:"expected_rating"`
	ExpectedSentiment  float64  `json:"expected_sentiment"`
	PotentialIssues    []string `json:"potential_issues"`
	RecommendedActions []string `json:"recommended_actions"`
}

// IssueFrequency frequência de problemas
type IssueFrequency struct {
	Issue      string    `json:"issue"`
	Frequency  int       `json:"frequency"`
	Severity   float64   `json:"severity"`
	Trend      string    `json:"trend"` // increasing, decreasing, stable
	LastSeen   time.Time `json:"last_seen"`
	FirstSeen  time.Time `json:"first_seen"`
	Resolution string    `json:"resolution"`
}

// SuggestionFrequency frequência de sugestões
type SuggestionFrequency struct {
	Suggestion    string    `json:"suggestion"`
	Frequency     int       `json:"frequency"`
	Priority      float64   `json:"priority"`
	Feasibility   float64   `json:"feasibility"`
	Impact        float64   `json:"impact"`
	LastSuggested time.Time `json:"last_suggested"`
	Status        string    `json:"status"`
}

// UserSegmentation segmentação de usuários
type UserSegmentation struct {
	ByExperience   map[string]UserSegment `json:"by_experience"`
	ByUsage        map[string]UserSegment `json:"by_usage"`
	ByLanguage     map[string]UserSegment `json:"by_language"`
	BySatisfaction map[string]UserSegment `json:"by_satisfaction"`
}

// UserSegment segmento de usuário
type UserSegment struct {
	Count           int                `json:"count"`
	AverageRating   float64            `json:"average_rating"`
	TopIssues       []string           `json:"top_issues"`
	TopSuggestions  []string           `json:"top_suggestions"`
	Satisfaction    float64            `json:"satisfaction"`
	RetentionRate   float64            `json:"retention_rate"`
	Characteristics map[string]float64 `json:"characteristics"`
}

// ResolutionMetrics métricas de resolução
type ResolutionMetrics struct {
	AverageResolutionTime       time.Duration            `json:"average_resolution_time"`
	ResolutionRate              float64                  `json:"resolution_rate"`
	FirstContactResolution      float64                  `json:"first_contact_resolution"`
	ReopenRate                  float64                  `json:"reopen_rate"`
	SatisfactionAfterResolution float64                  `json:"satisfaction_after_resolution"`
	ResolutionByCategory        map[string]time.Duration `json:"resolution_by_category"`
	EscalationRate              float64                  `json:"escalation_rate"`
}

// PredictiveInsights insights preditivos
type PredictiveInsights struct {
	ChurnRisk           []UserChurnRisk     `json:"churn_risk"`
	QualityPrediction   QualityPrediction   `json:"quality_prediction"`
	CapacityForecasting CapacityForecast    `json:"capacity_forecasting"`
	TrendPredictions    []TrendPrediction   `json:"trend_predictions"`
	AnomalyDetection    []Anomaly           `json:"anomaly_detection"`
	RecommendedActions  []RecommendedAction `json:"recommended_actions"`
}

// UserChurnRisk risco de churn do usuário
type UserChurnRisk struct {
	UserID          string        `json:"user_id"`
	RiskScore       float64       `json:"risk_score"`
	Indicators      []string      `json:"indicators"`
	Recommendations []string      `json:"recommendations"`
	TimeToChurn     time.Duration `json:"time_to_churn"`
}

// QualityPrediction previsão de qualidade
type QualityPrediction struct {
	PredictedQuality        float64   `json:"predicted_quality"`
	ConfidenceInterval      []float64 `json:"confidence_interval"`
	KeyFactors              []string  `json:"key_factors"`
	RiskFactors             []string  `json:"risk_factors"`
	RecommendedImprovements []string  `json:"recommended_improvements"`
}

// CapacityForecast previsão de capacidade
type CapacityForecast struct {
	PredictedLoad          float64            `json:"predicted_load"`
	CurrentCapacity        float64            `json:"current_capacity"`
	CapacityUtilization    float64            `json:"capacity_utilization"`
	BottleneckPrediction   []string           `json:"bottleneck_prediction"`
	ScalingRecommendations []string           `json:"scaling_recommendations"`
	ResourceRequirements   map[string]float64 `json:"resource_requirements"`
}

// TrendPrediction previsão de tendência
type TrendPrediction struct {
	Trend       string        `json:"trend"`
	Probability float64       `json:"probability"`
	Timeline    time.Duration `json:"timeline"`
	Impact      string        `json:"impact"`
	Preparation []string      `json:"preparation"`
}

// Anomaly anomalia detectada
type Anomaly struct {
	Type          string                 `json:"type"`
	Description   string                 `json:"description"`
	Severity      string                 `json:"severity"`
	DetectedAt    time.Time              `json:"detected_at"`
	Data          map[string]interface{} `json:"data"`
	Investigation string                 `json:"investigation"`
	Actions       []string               `json:"actions"`
}

// RecommendedAction ação recomendada
type RecommendedAction struct {
	Action       string        `json:"action"`
	Priority     int           `json:"priority"`
	Impact       float64       `json:"impact"`
	Effort       float64       `json:"effort"`
	Timeline     time.Duration `json:"timeline"`
	Dependencies []string      `json:"dependencies"`
	Success      float64       `json:"success_probability"`
}

// AutoTuningEngine motor de auto-ajuste
type AutoTuningEngine struct {
	parameters  map[string]Parameter
	experiments []Experiment
	bestConfigs map[string]interface{}
	performance PerformanceTracker
}

// Parameter parâmetro do sistema
type Parameter struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	CurrentValue interface{} `json:"current_value"`
	Range        Range       `json:"range"`
	Impact       float64     `json:"impact"`
	Tunable      bool        `json:"tunable"`
	LastChanged  time.Time   `json:"last_changed"`
}

// Range faixa de valores
type Range struct {
	Min  interface{} `json:"min"`
	Max  interface{} `json:"max"`
	Step interface{} `json:"step"`
}

// Experiment experimento de auto-ajuste
type Experiment struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
	Results    ExperimentResults      `json:"results"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Status     string                 `json:"status"`
	Hypothesis string                 `json:"hypothesis"`
	Conclusion string                 `json:"conclusion"`
}

// ExperimentResults resultados do experimento
type ExperimentResults struct {
	Metrics          map[string]float64 `json:"metrics"`
	Quality          float64            `json:"quality"`
	Performance      float64            `json:"performance"`
	UserSatisfaction float64            `json:"user_satisfaction"`
	Errors           []string           `json:"errors"`
	Observations     []string           `json:"observations"`
	Significance     float64            `json:"significance"`
}

// PerformanceTracker rastreador de performance
type PerformanceTracker struct {
	Metrics   map[string][]MetricPoint `json:"metrics"`
	Baselines map[string]float64       `json:"baselines"`
	Targets   map[string]float64       `json:"targets"`
	Alerts    []PerformanceAlert       `json:"alerts"`
}

// MetricPoint ponto métrico
type MetricPoint struct {
	Timestamp time.Time              `json:"timestamp"`
	Value     float64                `json:"value"`
	Context   map[string]interface{} `json:"context"`
}

// PerformanceAlert alerta de performance
type PerformanceAlert struct {
	Metric       string    `json:"metric"`
	Threshold    float64   `json:"threshold"`
	CurrentValue float64   `json:"current_value"`
	Severity     string    `json:"severity"`
	Timestamp    time.Time `json:"timestamp"`
	Action       string    `json:"action"`
}

// QualityMonitor monitor de qualidade
type QualityMonitor struct {
	thresholds   map[string]float64
	trends       []QualityTrend
	alerts       []QualityAlert
	improvements []QualityImprovement
}

// QualityTrend tendência de qualidade
type QualityTrend struct {
	Metric     string        `json:"metric"`
	Direction  string        `json:"direction"` // improving, degrading, stable
	Rate       float64       `json:"rate"`
	Period     time.Duration `json:"period"`
	Confidence float64       `json:"confidence"`
}

// QualityAlert alerta de qualidade
type QualityAlert struct {
	Type      string    `json:"type"`
	Metric    string    `json:"metric"`
	Current   float64   `json:"current"`
	Threshold float64   `json:"threshold"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
	Suggested []string  `json:"suggested_actions"`
}

// QualityImprovement melhoria de qualidade
type QualityImprovement struct {
	Area     string        `json:"area"`
	Current  float64       `json:"current"`
	Target   float64       `json:"target"`
	Actions  []string      `json:"actions"`
	Timeline time.Duration `json:"timeline"`
	Priority int           `json:"priority"`
	ROI      float64       `json:"roi"`
}

// OutcomePredictor preditor de resultados
type OutcomePredictor struct {
	models      map[string]PredictionModel
	features    []Feature
	predictions []Prediction
	accuracy    map[string]float64
}

// PredictionModel modelo de predição
type PredictionModel struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Accuracy    float64                `json:"accuracy"`
	LastTrained time.Time              `json:"last_trained"`
	Features    []string               `json:"features"`
	Performance ModelPerformance       `json:"performance"`
}

// Feature característica para predição
type Feature struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Importance  float64 `json:"importance"`
	Correlation float64 `json:"correlation"`
	Description string  `json:"description"`
}

// Prediction predição
type Prediction struct {
	Target      string                 `json:"target"`
	Value       interface{}            `json:"value"`
	Confidence  float64                `json:"confidence"`
	Probability float64                `json:"probability"`
	Features    map[string]interface{} `json:"features"`
	Timestamp   time.Time              `json:"timestamp"`
	Context     map[string]interface{} `json:"context"`
}

// ModelPerformance performance do modelo
type ModelPerformance struct {
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1Score   float64 `json:"f1_score"`
	AUC       float64 `json:"auc"`
	MSE       float64 `json:"mse"`
	MAE       float64 `json:"mae"`
}

// NewFeedbackSystem cria um novo sistema de feedback
func NewFeedbackSystem(feedbackPath string) *FeedbackSystem {
	fs := &FeedbackSystem{
		feedbackPath: feedbackPath,
		feedbacks:    make([]UserFeedback, 0),
		improvements: make([]SystemImprovement, 0),
		analytics: FeedbackAnalytics{
			SentimentDistribution: make(map[string]int),
			CategoryDistribution:  make(map[FeedbackCategory]int),
			TypeDistribution:      make(map[FeedbackType]int),
		},
		autoTuning: &AutoTuningEngine{
			parameters:  make(map[string]Parameter),
			experiments: make([]Experiment, 0),
			bestConfigs: make(map[string]interface{}),
		},
		qualityMonitor: &QualityMonitor{
			thresholds:   make(map[string]float64),
			trends:       make([]QualityTrend, 0),
			alerts:       make([]QualityAlert, 0),
			improvements: make([]QualityImprovement, 0),
		},
		predictor: &OutcomePredictor{
			models:      make(map[string]PredictionModel),
			features:    make([]Feature, 0),
			predictions: make([]Prediction, 0),
			accuracy:    make(map[string]float64),
		},
	}

	// Carregar dados existentes
	if err := fs.loadData(); err != nil {
		fmt.Printf("⚠️  Aviso: erro ao carregar dados de feedback: %v\n", err)
	}

	// Inicializar auto-tuning
	fs.initializeAutoTuning()

	return fs
}

// SubmitFeedback submete feedback do usuário
func (fs *FeedbackSystem) SubmitFeedback(feedback UserFeedback) error {
	// Gerar ID único
	feedback.ID = fmt.Sprintf("fb_%d", time.Now().Unix())
	feedback.Timestamp = time.Now()

	// Classificar automaticamente
	fs.classifyFeedback(&feedback)

	// Analisar sentimento
	feedback.Sentiment = fs.analyzeSentiment(feedback.Content)

	// Extrair action items
	feedback.ActionItems = fs.extractActionItems(feedback)

	// Determinar prioridade
	feedback.Priority = fs.calculatePriority(feedback)

	// Adicionar à lista
	fs.feedbacks = append(fs.feedbacks, feedback)

	// Atualizar analytics
	fs.updateAnalytics(feedback)

	// Verificar se gera melhorias automáticas
	fs.generateImprovements(feedback)

	// Salvar dados
	if err := fs.saveData(); err != nil {
		return fmt.Errorf("erro ao salvar feedback: %v", err)
	}

	fmt.Printf("📝 Feedback registrado: %s (Rating: %d/5)\n", feedback.ID, feedback.Rating)
	return nil
}

// classifyFeedback classifica feedback automaticamente
func (fs *FeedbackSystem) classifyFeedback(feedback *UserFeedback) {
	content := strings.ToLower(feedback.Content)

	// Classificar tipo
	if strings.Contains(content, "bug") || strings.Contains(content, "erro") || strings.Contains(content, "problem") {
		feedback.Type = TypeBugReport
		feedback.Category = CategoryBug
	} else if strings.Contains(content, "sugest") || strings.Contains(content, "melhori") || strings.Contains(content, "poderia") {
		feedback.Type = TypeSuggestion
		feedback.Category = CategoryFeature
	} else if strings.Contains(content, "lento") || strings.Contains(content, "demora") || strings.Contains(content, "performance") {
		feedback.Type = TypeNegative
		feedback.Category = CategoryPerformance
	} else if strings.Contains(content, "qualidade") || strings.Contains(content, "result") {
		feedback.Category = CategoryQuality
	} else if strings.Contains(content, "segur") || strings.Contains(content, "vulnerab") {
		feedback.Category = CategorySecurity
	} else if strings.Contains(content, "document") || strings.Contains(content, "help") {
		feedback.Category = CategoryDocumentation
	} else if strings.Contains(content, "interface") || strings.Contains(content, "usabilidade") {
		feedback.Category = CategoryUsability
	} else {
		feedback.Category = CategoryGeneral
	}

	// Determinar tipo baseado no rating
	if feedback.Rating >= 4 {
		if feedback.Type == "" {
			feedback.Type = TypePositive
		}
	} else if feedback.Rating <= 2 {
		if feedback.Type == "" {
			feedback.Type = TypeNegative
		}
	}

	feedback.AutoClassified = true
}

// analyzeSentiment analisa sentimento do feedback
func (fs *FeedbackSystem) analyzeSentiment(content string) float64 {
	// Análise simplificada de sentimento
	positiveWords := []string{
		"bom", "ótimo", "excelente", "perfeito", "incrível",
		"amor", "gosto", "funciona", "rápido", "fácil",
		"útil", "prático", "eficiente", "quality", "amazing",
	}

	negativeWords := []string{
		"ruim", "péssimo", "horrível", "lento", "difícil",
		"problema", "erro", "bug", "falha", "frustrante",
		"complicado", "confuso", "terrible", "awful", "slow",
	}

	content = strings.ToLower(content)
	positiveCount := 0
	negativeCount := 0

	for _, word := range positiveWords {
		positiveCount += strings.Count(content, word)
	}

	for _, word := range negativeWords {
		negativeCount += strings.Count(content, word)
	}

	totalWords := positiveCount + negativeCount
	if totalWords == 0 {
		return 0.0 // Neutro
	}

	sentiment := float64(positiveCount-negativeCount) / float64(totalWords)

	// Normalizar entre -1 e 1
	if sentiment > 1 {
		sentiment = 1
	} else if sentiment < -1 {
		sentiment = -1
	}

	return sentiment
}

// extractActionItems extrai itens de ação do feedback
func (fs *FeedbackSystem) extractActionItems(feedback UserFeedback) []ActionItem {
	actionItems := make([]ActionItem, 0)

	content := strings.ToLower(feedback.Content)

	// Identificar ações baseadas no conteúdo
	if strings.Contains(content, "bug") || strings.Contains(content, "erro") {
		actionItems = append(actionItems, ActionItem{
			ID:          fmt.Sprintf("action_%d_1", time.Now().Unix()),
			Description: "Investigar e corrigir bug reportado",
			Priority:    1,
			Effort:      "medium",
			Impact:      "high",
			Category:    "bug_fix",
			Status:      "pending",
		})
	}

	if strings.Contains(content, "performance") || strings.Contains(content, "lento") {
		actionItems = append(actionItems, ActionItem{
			ID:          fmt.Sprintf("action_%d_2", time.Now().Unix()),
			Description: "Otimizar performance do sistema",
			Priority:    2,
			Effort:      "high",
			Impact:      "high",
			Category:    "performance",
			Status:      "pending",
		})
	}

	if strings.Contains(content, "document") {
		actionItems = append(actionItems, ActionItem{
			ID:          fmt.Sprintf("action_%d_3", time.Now().Unix()),
			Description: "Melhorar documentação",
			Priority:    3,
			Effort:      "low",
			Impact:      "medium",
			Category:    "documentation",
			Status:      "pending",
		})
	}

	return actionItems
}

// calculatePriority calcula prioridade do feedback
func (fs *FeedbackSystem) calculatePriority(feedback UserFeedback) int {
	priority := 3 // Padrão médio

	// Ajustar baseado no rating
	if feedback.Rating <= 2 {
		priority = 1 // Alta prioridade para ratings baixos
	} else if feedback.Rating >= 4 {
		priority = 4 // Baixa prioridade para ratings altos (exceto se for sugestão)
	}

	// Ajustar baseado no tipo
	switch feedback.Type {
	case TypeBugReport:
		priority = 1
	case TypeFeatureRequest:
		priority = 3
	case TypeSuggestion:
		priority = 2
	case TypeComplaint:
		priority = 1
	}

	// Ajustar baseado na categoria
	switch feedback.Category {
	case CategorySecurity:
		priority = 1
	case CategoryBug:
		priority = 1
	case CategoryPerformance:
		priority = 2
	case CategoryUsability:
		priority = 2
	}

	// Ajustar baseado no sentimento
	if feedback.Sentiment < -0.5 {
		priority = min(priority, 2)
	}

	return priority
}

// updateAnalytics atualiza analytics do feedback
func (fs *FeedbackSystem) updateAnalytics(feedback UserFeedback) {
	fs.analytics.TotalFeedbacks++

	// Atualizar rating médio
	totalRating := fs.analytics.AverageRating * float64(fs.analytics.TotalFeedbacks-1)
	fs.analytics.AverageRating = (totalRating + float64(feedback.Rating)) / float64(fs.analytics.TotalFeedbacks)

	// Atualizar distribuições
	sentimentCategory := fs.getSentimentCategory(feedback.Sentiment)
	fs.analytics.SentimentDistribution[sentimentCategory]++
	fs.analytics.CategoryDistribution[feedback.Category]++
	fs.analytics.TypeDistribution[feedback.Type]++

	// Atualizar tendências de satisfação
	fs.analytics.SatisfactionTrends = append(fs.analytics.SatisfactionTrends, SatisfactionPoint{
		Date:         feedback.Timestamp,
		Satisfaction: float64(feedback.Rating) / 5.0,
		NPS:          fs.calculateNPS(feedback.Rating),
		CSAT:         float64(feedback.Rating) / 5.0 * 100,
	})
}

// generateImprovements gera melhorias baseadas no feedback
func (fs *FeedbackSystem) generateImprovements(feedback UserFeedback) {
	// Gerar melhorias automáticas baseadas no feedback
	if feedback.Rating <= 2 && feedback.Category == CategoryPerformance {
		improvement := SystemImprovement{
			ID:              fmt.Sprintf("imp_%d", time.Now().Unix()),
			Title:           "Otimização de Performance",
			Description:     "Melhorar performance baseado em feedback de usuário",
			Category:        ImprovementPerformance,
			Priority:        1,
			Impact:          0.8,
			Effort:          0.6,
			ROI:             0.8 / 0.6,
			Status:          ImprovementProposed,
			BasedOnFeedback: []string{feedback.ID},
			Timeline: Timeline{
				Planning: TimelinePhase{
					StartDate: time.Now(),
					EndDate:   time.Now().Add(7 * 24 * time.Hour),
					Status:    "pending",
				},
				Development: TimelinePhase{
					StartDate: time.Now().Add(7 * 24 * time.Hour),
					EndDate:   time.Now().Add(21 * 24 * time.Hour),
					Status:    "pending",
				},
			},
			RiskLevel: "medium",
		}

		fs.improvements = append(fs.improvements, improvement)
	}
}

// GetAnalytics retorna analytics do feedback
func (fs *FeedbackSystem) GetAnalytics() FeedbackAnalytics {
	// Atualizar analytics em tempo real
	fs.calculateTrends()
	fs.identifyTopIssues()
	fs.performPredictiveAnalysis()

	return fs.analytics
}

// calculateTrends calcula tendências
func (fs *FeedbackSystem) calculateTrends() {
	// Implementar cálculo de tendências
	// Simplificado para exemplo
	if len(fs.feedbacks) < 2 {
		return
	}

	// Calcular tendência de rating
	recent := fs.feedbacks[len(fs.feedbacks)-10:]
	if len(recent) < 10 {
		recent = fs.feedbacks
	}

	totalRecent := 0.0
	for _, f := range recent {
		totalRecent += float64(f.Rating)
	}
	recentAvg := totalRecent / float64(len(recent))

	fs.analytics.TrendAnalysis.GrowthRate = (recentAvg - fs.analytics.AverageRating) / fs.analytics.AverageRating
}

// identifyTopIssues identifica principais problemas
func (fs *FeedbackSystem) identifyTopIssues() {
	issueCount := make(map[string]int)

	for _, feedback := range fs.feedbacks {
		if feedback.Rating <= 2 {
			// Extrair problemas do conteúdo
			content := strings.ToLower(feedback.Content)
			if strings.Contains(content, "lento") {
				issueCount["Performance lenta"]++
			}
			if strings.Contains(content, "erro") {
				issueCount["Erros frequentes"]++
			}
			if strings.Contains(content, "difícil") {
				issueCount["Dificuldade de uso"]++
			}
		}
	}

	// Converter para slice ordenado
	topIssues := make([]IssueFrequency, 0)
	for issue, count := range issueCount {
		topIssues = append(topIssues, IssueFrequency{
			Issue:     issue,
			Frequency: count,
			Severity:  0.8,
			Trend:     "stable",
			LastSeen:  time.Now(),
		})
	}

	fs.analytics.TopIssues = topIssues
}

// performPredictiveAnalysis executa análise preditiva
func (fs *FeedbackSystem) performPredictiveAnalysis() {
	// Implementar análise preditiva
	// Simplificado para exemplo

	// Predição de qualidade
	fs.analytics.PredictiveInsights.QualityPrediction = QualityPrediction{
		PredictedQuality:   fs.analytics.AverageRating / 5.0,
		ConfidenceInterval: []float64{0.6, 0.9},
		KeyFactors:         []string{"User feedback", "Performance metrics"},
		RiskFactors:        []string{"Increasing negative feedback"},
	}

	// Detecção de anomalias
	if fs.analytics.AverageRating < 3.0 {
		fs.analytics.PredictiveInsights.AnomalyDetection = append(
			fs.analytics.PredictiveInsights.AnomalyDetection,
			Anomaly{
				Type:        "quality_drop",
				Description: "Rating médio abaixo do threshold",
				Severity:    "high",
				DetectedAt:  time.Now(),
				Actions:     []string{"Investigar causas", "Implementar melhorias"},
			},
		)
	}
}

// getSentimentCategory categoriza sentimento
func (fs *FeedbackSystem) getSentimentCategory(sentiment float64) string {
	if sentiment > 0.3 {
		return "positive"
	} else if sentiment < -0.3 {
		return "negative"
	} else {
		return "neutral"
	}
}

// calculateNPS calcula Net Promoter Score
func (fs *FeedbackSystem) calculateNPS(rating int) float64 {
	if rating >= 4 {
		return 100.0 // Promoter
	} else if rating >= 3 {
		return 0.0 // Passive
	} else {
		return -100.0 // Detractor
	}
}

// initializeAutoTuning inicializa sistema de auto-ajuste
func (fs *FeedbackSystem) initializeAutoTuning() {
	// Configurar parâmetros ajustáveis
	fs.autoTuning.parameters["quality_threshold"] = Parameter{
		Name:         "quality_threshold",
		Type:         "float",
		CurrentValue: 0.7,
		Range:        Range{Min: 0.5, Max: 0.9, Step: 0.1},
		Impact:       0.8,
		Tunable:      true,
		LastChanged:  time.Now(),
	}

	fs.autoTuning.parameters["max_tokens"] = Parameter{
		Name:         "max_tokens",
		Type:         "int",
		CurrentValue: 150000,
		Range:        Range{Min: 100000, Max: 200000, Step: 10000},
		Impact:       0.6,
		Tunable:      true,
		LastChanged:  time.Now(),
	}

	fs.autoTuning.parameters["cache_ttl"] = Parameter{
		Name:         "cache_ttl",
		Type:         "duration",
		CurrentValue: 24 * time.Hour,
		Range:        Range{Min: 6 * time.Hour, Max: 72 * time.Hour, Step: 6 * time.Hour},
		Impact:       0.4,
		Tunable:      true,
		LastChanged:  time.Now(),
	}
}

// loadData carrega dados do disco
func (fs *FeedbackSystem) loadData() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(fs.feedbackPath, 0755); err != nil {
		return err
	}

	// Carregar feedbacks
	feedbackFile := filepath.Join(fs.feedbackPath, "feedbacks.json")
	if data, err := os.ReadFile(feedbackFile); err == nil {
		if err := json.Unmarshal(data, &fs.feedbacks); err != nil {
			return err
		}
	}

	// Carregar melhorias
	improvementsFile := filepath.Join(fs.feedbackPath, "improvements.json")
	if data, err := os.ReadFile(improvementsFile); err == nil {
		if err := json.Unmarshal(data, &fs.improvements); err != nil {
			return err
		}
	}

	// Carregar analytics
	analyticsFile := filepath.Join(fs.feedbackPath, "analytics.json")
	if data, err := os.ReadFile(analyticsFile); err == nil {
		if err := json.Unmarshal(data, &fs.analytics); err != nil {
			return err
		}
	}

	return nil
}

// saveData salva dados no disco
func (fs *FeedbackSystem) saveData() error {
	// Criar diretório se não existir
	if err := os.MkdirAll(fs.feedbackPath, 0755); err != nil {
		return err
	}

	// Salvar feedbacks
	feedbackFile := filepath.Join(fs.feedbackPath, "feedbacks.json")
	if data, err := json.MarshalIndent(fs.feedbacks, "", "  "); err == nil {
		if err := os.WriteFile(feedbackFile, data, 0644); err != nil {
			return err
		}
	}

	// Salvar melhorias
	improvementsFile := filepath.Join(fs.feedbackPath, "improvements.json")
	if data, err := json.MarshalIndent(fs.improvements, "", "  "); err == nil {
		if err := os.WriteFile(improvementsFile, data, 0644); err != nil {
			return err
		}
	}

	// Salvar analytics
	analyticsFile := filepath.Join(fs.feedbackPath, "analytics.json")
	if data, err := json.MarshalIndent(fs.analytics, "", "  "); err == nil {
		if err := os.WriteFile(analyticsFile, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GenerateReport gera relatório do sistema de feedback
func (fs *FeedbackSystem) GenerateReport() string {
	analytics := fs.GetAnalytics()

	report := "📊 RELATÓRIO DO SISTEMA DE FEEDBACK\n\n"

	// Estatísticas gerais
	report += fmt.Sprintf("Total de Feedbacks: %d\n", analytics.TotalFeedbacks)
	report += fmt.Sprintf("Rating Médio: %.2f/5\n", analytics.AverageRating)

	// Distribuição por categoria
	report += "\n📋 Distribuição por Categoria:\n"
	for category, count := range analytics.CategoryDistribution {
		report += fmt.Sprintf("   • %s: %d\n", category, count)
	}

	// Principais problemas
	report += "\n❗ Principais Problemas:\n"
	for _, issue := range analytics.TopIssues {
		report += fmt.Sprintf("   • %s: %d ocorrências\n", issue.Issue, issue.Frequency)
	}

	// Melhorias propostas
	report += fmt.Sprintf("\n🔧 Melhorias Propostas: %d\n", len(fs.improvements))

	// Análise preditiva
	if analytics.PredictiveInsights.QualityPrediction.PredictedQuality > 0 {
		report += fmt.Sprintf("\n🔮 Qualidade Predita: %.2f\n", analytics.PredictiveInsights.QualityPrediction.PredictedQuality)
	}

	return report
}
