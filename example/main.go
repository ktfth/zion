package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ktfth/zion/internal/core"
	"github.com/ktfth/zion/internal/providers"
)

func main() {
	// Example usage of the refactored Zion CLI components

	// 1. Create project configuration
	config := &core.ProjectConfig{
		Name:        "example-project",
		Language:    "go",
		Description: "A sample Go project with REST API",
		OutputDir:   ".",
	}

	// 2. Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// 3. Create provider factory
	factory := providers.NewProviderFactory()

	// 4. Get default provider (requires API key in environment)
	var aiProvider core.AIProvider
	var err error

	// Try to create a provider with a test key (this would normally come from environment)
	if geminiKey := os.Getenv("GEMINI_API_KEY"); geminiKey != "" {
		providerConfig := &core.ProviderConfig{
			APIKey: geminiKey,
			Model:  "gemini-2.0-flash",
		}
		aiProvider, err = factory.CreateProvider("gemini", providerConfig)
	} else if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		providerConfig := &core.ProviderConfig{
			APIKey: openaiKey,
			Model:  "gpt-3.5-turbo",
		}
		aiProvider, err = factory.CreateProvider("openai", providerConfig)
	} else {
		fmt.Println("No API key found. Please set GEMINI_API_KEY or OPENAI_API_KEY environment variable.")
		fmt.Println("For demonstration purposes, using mock provider...")
		aiProvider = &MockProvider{}
	}

	if err != nil {
		log.Fatalf("Failed to create AI provider: %v", err)
	}

	// 5. Create project generator
	generator := core.NewProjectGenerator(config, aiProvider)

	// 6. Generate project
	fmt.Printf("🚀 Generating project with %s\n", aiProvider.Name())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📦 Project: %s\n", config.Name)
	fmt.Printf("🔧 Language: %s\n", config.Language)
	fmt.Printf("📝 Description: %s\n", config.Description)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	result, err := generator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate project: %v", err)
	}

	// 7. Display result
	fmt.Printf("✨ Project generated successfully! ✨\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📁 Location: %s\n", result.ProjectName)
	fmt.Printf("📊 Files: %d, Directories: %d\n", result.FilesCreated, result.DirsCreated)
	fmt.Printf("⏱️  Duration: %.2f seconds\n", result.Duration.Seconds())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("   1. cd %s\n", result.ProjectName)
	fmt.Printf("   2. Explore the generated files\n")
	fmt.Printf("   3. Start developing!\n")
}

// MockProvider for demonstration when no API key is available
type MockProvider struct{}

func (m *MockProvider) Name() string {
	return "Mock Provider"
}

func (m *MockProvider) GenerateContent(prompt string) (string, error) {
	return `{
  "project_name": "example-project",
  "language": "go",
  "description": "A sample Go project with REST API",
  "structure": {
    "directories": ["cmd", "internal", "pkg", "api", "docs"],
    "files": {
      "main.go": {
        "content": "package main\n\nimport (\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n)\n\nfunc main() {\n\thttp.HandleFunc(\"/\", func(w http.ResponseWriter, r *http.Request) {\n\t\tfmt.Fprintln(w, \"Hello, World!\")\n\t})\n\n\tlog.Println(\"Server starting on :8080...\")\n\tlog.Fatal(http.ListenAndServe(\":8080\", nil))\n}"
      },
      "go.mod": {
        "content": "module example-project\n\ngo 1.21\n"
      },
      "README.md": {
        "content": "# Example Project\n\nA sample Go project with REST API\n\n## Usage\n\ngo run main.go\n\nThen visit http://localhost:8080\n"
      }
    }
  },
  "dependencies": ["go 1.21"],
  "next_steps": ["go mod tidy", "go run main.go"]
}`, nil
}
