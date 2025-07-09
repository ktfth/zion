package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zion",
	Short: "Zion CLI - AI-powered project scaffolding",
	Long: `Zion is a command-line tool that generates project structures
using artificial intelligence. It supports multiple programming languages
and AI providers to create clean, well-structured projects.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
