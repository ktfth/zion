package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ProjectGenerator is responsible for generating projects
type ProjectGenerator struct {
	config *ProjectConfig
	ai     AIProvider
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator(config *ProjectConfig, ai AIProvider) *ProjectGenerator {
	return &ProjectGenerator{
		config: config,
		ai:     ai,
	}
}

// Generate creates a new project based on the configuration
func (g *ProjectGenerator) Generate() (*ProjectResult, error) {
	start := time.Now()
	result := &ProjectResult{
		ProjectName: g.config.Name,
		Language:    g.config.Language,
		Description: g.config.Description,
	}

	// Validate configuration
	if err := g.config.Validate(); err != nil {
		result.Error = err
		return result, err
	}

	// Generate project structure using AI
	response, err := g.ai.GenerateContent(g.buildPrompt())
	if err != nil {
		result.Error = fmt.Errorf("failed to generate project structure: %w", err)
		return result, result.Error
	}

	// Create project structure
	filesCreated, dirsCreated, err := g.createProjectStructure(response)
	if err != nil {
		result.Error = fmt.Errorf("failed to create project structure: %w", err)
		return result, result.Error
	}

	result.FilesCreated = filesCreated
	result.DirsCreated = dirsCreated
	result.Duration = time.Since(start)
	result.Success = true

	return result, nil
}

// buildPrompt creates the AI prompt for project generation
func (g *ProjectGenerator) buildPrompt() string {
	return fmt.Sprintf(`Generate a complete project structure for:
- Language: %s
- Project: %s
- Description: %s

Return a JSON response with the following structure:
{
  "project_name": "%s",
  "language": "%s",
  "description": "%s",
  "structure": {
    "directories": ["dir1", "dir2"],
    "files": {
      "filename": {
        "content": "file content here"
      }
    }
  },
  "dependencies": ["dep1", "dep2"],
  "next_steps": ["step1", "step2"]
}`, g.config.Language, g.config.Name, g.config.Description, g.config.Name, g.config.Language, g.config.Description)
}

// createProjectStructure creates the actual project files and directories
func (g *ProjectGenerator) createProjectStructure(response string) (int, int, error) {
	// This is a simplified version - in the real implementation,
	// we would parse the JSON response and create the actual files

	// Create project directory
	projectDir := filepath.Join(".", g.config.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create project directory: %w", err)
	}

	// For now, just create a basic structure
	dirs := []string{"src", "tests", "docs"}
	files := map[string]string{
		"README.md": fmt.Sprintf("# %s\n\n%s", g.config.Name, g.config.Description),
		"go.mod":    fmt.Sprintf("module %s\n\ngo 1.21", g.config.Name),
		"main.go":   "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}",
	}

	// Create directories
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return 0, 0, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create files
	for filename, content := range files {
		filePath := filepath.Join(projectDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return 0, 0, fmt.Errorf("failed to create file %s: %w", filename, err)
		}
	}

	return len(files), len(dirs), nil
}
