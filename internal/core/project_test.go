package core

import (
	"fmt"
	"testing"
	"time"
)

func TestProjectConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ProjectConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &ProjectConfig{
				Name:        "test-project",
				Language:    "go",
				Description: "A test project",
				OutputDir:   ".",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &ProjectConfig{
				Language:    "go",
				Description: "A test project",
				OutputDir:   ".",
			},
			wantErr: true,
		},
		{
			name: "missing language",
			config: &ProjectConfig{
				Name:        "test-project",
				Description: "A test project",
				OutputDir:   ".",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			config: &ProjectConfig{
				Name:      "test-project",
				Language:  "go",
				OutputDir: ".",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectResult_String(t *testing.T) {
	tests := []struct {
		name   string
		result *ProjectResult
		want   string
	}{
		{
			name: "successful result",
			result: &ProjectResult{
				ProjectName:  "test-project",
				FilesCreated: 5,
				DirsCreated:  3,
				Duration:     2 * time.Second,
				Success:      true,
			},
			want: "Project 'test-project' created successfully: 5 files, 3 directories in 2.00s",
		},
		{
			name: "failed result",
			result: &ProjectResult{
				ProjectName: "test-project",
				Success:     false,
				Error:       fmt.Errorf("some error"),
			},
			want: "Project 'test-project' failed: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.String()
			if got != tt.want {
				t.Errorf("ProjectResult.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
