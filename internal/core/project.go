package core

import (
	"fmt"
	"time"
)

// ProjectConfig represents the configuration for a project
type ProjectConfig struct {
	Name        string
	Language    string
	Description string
	OutputDir   string
}

// Validate checks if the project configuration is valid
func (p *ProjectConfig) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if p.Language == "" {
		return fmt.Errorf("project language cannot be empty")
	}
	if p.Description == "" {
		return fmt.Errorf("project description cannot be empty")
	}
	return nil
}

// ProjectResult represents the result of a project generation
type ProjectResult struct {
	ProjectName  string
	Language     string
	Description  string
	FilesCreated int
	DirsCreated  int
	Duration     time.Duration
	Success      bool
	Error        error
}

// String returns a formatted string representation of the result
func (r *ProjectResult) String() string {
	if !r.Success {
		return fmt.Sprintf("Project '%s' failed: %v", r.ProjectName, r.Error)
	}
	return fmt.Sprintf("Project '%s' created successfully: %d files, %d directories in %.2fs",
		r.ProjectName, r.FilesCreated, r.DirsCreated, r.Duration.Seconds())
}
