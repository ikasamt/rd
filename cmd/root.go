package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rd",
	Short: "Redmine CLI tool",
	Long: `rd is a command-line interface tool for Redmine.
It allows you to manage tickets, projects, and users from your terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("url", "", "Redmine URL (overrides REDMINE_URL)")
	rootCmd.PersistentFlags().String("key", "", "Redmine API key (overrides REDMINE_API_KEY)")
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("quiet", false, "Minimal output")
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")
}