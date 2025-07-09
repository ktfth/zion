package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ktfth/zion/internal/core"
	"github.com/ktfth/zion/internal/providers"
	"github.com/spf13/cobra"
)

var (
	language    string
	projectName string
	description string
	provider    string
	apiKey      string
	model       string
	maxRetries  int
)

// scaffoldCmd represents the scaffold command
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate a project structure using AI",
	Long: `Generate a complete project structure using artificial intelligence.
	
This command creates a new project with the specified language, name, and description,
using AI to generate appropriate files and directory structure.`,
	Run: runScaffold,
}

func runScaffold(cmd *cobra.Command, args []string) {
	// Validate required flags
	if projectName == "" || language == "" {
		fmt.Println("❌ Error: --name and --language are required")
		cmd.Help()
		os.Exit(1)
	}

	// Create project configuration
	config := &core.ProjectConfig{
		Name:        projectName,
		Language:    language,
		Description: description,
		OutputDir:   ".",
	}

	// Create provider factory and get AI provider
	factory := providers.NewProviderFactory()
	var aiProvider core.AIProvider
	var err error

	if provider != "" || apiKey != "" {
		// Use custom provider configuration
		providerConfig := &core.ProviderConfig{
			APIKey: apiKey,
			Model:  model,
		}
		aiProvider, err = factory.CreateProvider(provider, providerConfig)
	} else {
		// Use default provider
		aiProvider, err = factory.GetDefaultProvider()
	}

	if err != nil {
		fmt.Printf("❌ Error creating AI provider: %v\n", err)
		os.Exit(1)
	}

	// Create project generator
	generator := core.NewProjectGenerator(config, aiProvider)

	// Show generation info
	fmt.Printf("🚀 Generating project with %s\n", aiProvider.Name())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📦 Project: %s\n", projectName)
	fmt.Printf("🔧 Language: %s\n", language)
	fmt.Printf("📝 Description: %s\n", description)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Generate project with retry logic
	var result *core.ProjectResult
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			fmt.Printf("🔄 Attempt %d/%d...\n", attempt, maxRetries)
		}

		result, err = generator.Generate()
		if err == nil {
			break
		}

		if attempt < maxRetries {
			fmt.Printf("⚠️  Attempt %d failed: %v\n", attempt, err)
			fmt.Printf("🔄 Retrying in 2 seconds...\n")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		fmt.Printf("❌ Failed to generate project after %d attempts: %v\n", maxRetries, err)
		os.Exit(1)
	}

	// Show success message
	fmt.Printf("✨ Project created successfully! ✨\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📁 Location: %s\n", result.ProjectName)
	fmt.Printf("📊 Files: %d, Directories: %d\n", result.FilesCreated, result.DirsCreated)
	fmt.Printf("⏱️  Duration: %.2f seconds\n", result.Duration.Seconds())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func init() {
	// Add flags
	scaffoldCmd.Flags().StringVarP(&language, "language", "l", "", "Project language (required)")
	scaffoldCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	scaffoldCmd.Flags().StringVarP(&description, "description", "d", "", "Project description")
	scaffoldCmd.Flags().StringVarP(&provider, "provider", "p", "", "AI provider (gemini, openai)")
	scaffoldCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for the provider")
	scaffoldCmd.Flags().StringVarP(&model, "model", "m", "", "Specific model to use")
	scaffoldCmd.Flags().IntVarP(&maxRetries, "retries", "r", 3, "Maximum number of retries")

	// Mark required flags
	scaffoldCmd.MarkFlagRequired("language")
	scaffoldCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(scaffoldCmd)
}
