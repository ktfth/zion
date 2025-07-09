package cmd

import (
	"fmt"
	"os"

	"github.com/ktfth/zion/internal/core"
	"github.com/ktfth/zion/internal/providers"
	"github.com/spf13/cobra"
)

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage AI providers",
	Long:  `Manage and configure AI providers like Gemini and OpenAI.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available AI providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🤖 Available AI Providers:\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		providers := []string{"gemini", "openai"}
		for _, provider := range providers {
			status := "❌ Not configured"

			switch provider {
			case "gemini":
				if os.Getenv("GEMINI_API_KEY") != "" {
					status = "✅ Configured"
				}
			case "openai":
				if os.Getenv("OPENAI_API_KEY") != "" {
					status = "✅ Configured"
				}
			}

			fmt.Printf("• %s: %s\n", provider, status)
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("💡 To configure a provider, set the appropriate environment variable:\n")
		fmt.Printf("   export GEMINI_API_KEY=\"your-key-here\"\n")
		fmt.Printf("   export OPENAI_API_KEY=\"your-key-here\"\n")
	},
}

var testCmd = &cobra.Command{
	Use:   "test [provider]",
	Short: "Test connection to an AI provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName := args[0]

		fmt.Printf("🧪 Testing %s provider...\n", providerName)

		factory := providers.NewProviderFactory()

		var apiKey string
		switch providerName {
		case "gemini":
			apiKey = os.Getenv("GEMINI_API_KEY")
		case "openai":
			apiKey = os.Getenv("OPENAI_API_KEY")
		default:
			fmt.Printf("❌ Unsupported provider: %s\n", providerName)
			os.Exit(1)
		}

		if apiKey == "" {
			fmt.Printf("❌ API key not configured for %s\n", providerName)
			fmt.Printf("💡 Set the appropriate environment variable\n")
			os.Exit(1)
		}

		config := &core.ProviderConfig{
			APIKey: apiKey,
		}

		provider, err := factory.CreateProvider(providerName, config)
		if err != nil {
			fmt.Printf("❌ Error creating provider: %v\n", err)
			os.Exit(1)
		}

		response, err := provider.GenerateContent("Say 'OK' to confirm you're working.")
		if err != nil {
			fmt.Printf("❌ Test failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Test successful!\n")
		fmt.Printf("📝 Response: %s\n", response)
	},
}

func init() {
	providerCmd.AddCommand(listCmd)
	providerCmd.AddCommand(testCmd)
	rootCmd.AddCommand(providerCmd)
}
