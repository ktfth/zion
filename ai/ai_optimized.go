package ai

import (
	"fmt"
	"strings"
	"time"
)

// OptimizedResponse represents an optimized AI response
type OptimizedResponse struct {
	Content        string        `json:"content"`
	ProcessingTime time.Duration `json:"processing_time"`
	TokenCount     int           `json:"token_count"`
	OptimizedAt    time.Time     `json:"optimized_at"`
}

// OptimizeResponse optimizes AI responses for better performance
func OptimizeResponse(response string) *OptimizedResponse {
	start := time.Now()

	// Basic response optimization
	optimized := strings.TrimSpace(response)

	// Remove excessive whitespace
	optimized = strings.ReplaceAll(optimized, "\n\n\n", "\n\n")
	optimized = strings.ReplaceAll(optimized, "  ", " ")

	// Count tokens (rough estimate)
	tokenCount := estimateTokens(optimized)

	return &OptimizedResponse{
		Content:        optimized,
		ProcessingTime: time.Since(start),
		TokenCount:     tokenCount,
		OptimizedAt:    time.Now(),
	}
}

// CacheKey generates a cache key for AI responses
func CacheKey(language, projectName, description string) string {
	return fmt.Sprintf("%s:%s:%s", language, projectName, description)
}

// IsResponseCacheable determines if a response should be cached
func IsResponseCacheable(response string) bool {
	// Don't cache empty responses
	if strings.TrimSpace(response) == "" {
		return false
	}

	// Don't cache error responses
	if strings.Contains(strings.ToLower(response), "error") {
		return false
	}

	// Don't cache responses that are too small (likely incomplete)
	if len(response) < 100 {
		return false
	}

	return true
}

// OptimizePrompt optimizes prompts for better AI performance
func OptimizePrompt(prompt string) string {
	// Remove excessive whitespace
	optimized := strings.TrimSpace(prompt)
	optimized = strings.ReplaceAll(optimized, "\n\n\n", "\n\n")

	// Ensure consistent formatting
	lines := strings.Split(optimized, "\n")
	var cleanLines []string

	for _, line := range lines {
		cleaned := strings.TrimSpace(line)
		if cleaned != "" {
			cleanLines = append(cleanLines, cleaned)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// EstimateResponseTime estimates response time based on content size
func EstimateResponseTime(contentSize int) time.Duration {
	// Rough estimation: 1 second per 1000 characters
	baseTime := time.Duration(contentSize/1000) * time.Second

	// Add minimum processing time
	if baseTime < 2*time.Second {
		baseTime = 2 * time.Second
	}

	// Add maximum cap
	if baseTime > 30*time.Second {
		baseTime = 30 * time.Second
	}

	return baseTime
}
