package core

import (
	"errors"
	"testing"
)

// MockAIProvider is a mock implementation of AIProvider for testing
type MockAIProvider struct {
	name     string
	response string
	err      error
}

func NewMockAIProvider(name, response string, err error) *MockAIProvider {
	return &MockAIProvider{
		name:     name,
		response: response,
		err:      err,
	}
}

func (m *MockAIProvider) Name() string {
	return m.name
}

func (m *MockAIProvider) GenerateContent(prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func TestProjectGenerator_Generate(t *testing.T) {
	tests := []struct {
		name        string
		config      *ProjectConfig
		aiProvider  AIProvider
		wantErr     bool
		wantSuccess bool
	}{
		{
			name: "successful generation",
			config: &ProjectConfig{
				Name:        "test-project",
				Language:    "go",
				Description: "A test project",
				OutputDir:   ".",
			},
			aiProvider:  NewMockAIProvider("TestAI", "test response", nil),
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name: "invalid config",
			config: &ProjectConfig{
				Language:    "go",
				Description: "A test project",
				OutputDir:   ".",
			},
			aiProvider:  NewMockAIProvider("TestAI", "test response", nil),
			wantErr:     true,
			wantSuccess: false,
		},
		{
			name: "ai provider error",
			config: &ProjectConfig{
				Name:        "test-project",
				Language:    "go",
				Description: "A test project",
				OutputDir:   ".",
			},
			aiProvider:  NewMockAIProvider("TestAI", "", errors.New("AI error")),
			wantErr:     true,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewProjectGenerator(tt.config, tt.aiProvider)
			result, err := generator.Generate()

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectGenerator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result == nil {
				t.Error("ProjectGenerator.Generate() returned nil result")
				return
			}

			if result.Success != tt.wantSuccess {
				t.Errorf("ProjectGenerator.Generate() success = %v, wantSuccess %v", result.Success, tt.wantSuccess)
			}

			if tt.wantSuccess {
				if result.ProjectName != tt.config.Name {
					t.Errorf("ProjectGenerator.Generate() project name = %v, want %v", result.ProjectName, tt.config.Name)
				}
				if result.Language != tt.config.Language {
					t.Errorf("ProjectGenerator.Generate() language = %v, want %v", result.Language, tt.config.Language)
				}
				if result.Duration == 0 {
					t.Error("ProjectGenerator.Generate() duration should be > 0")
				}
			}
		})
	}
}
