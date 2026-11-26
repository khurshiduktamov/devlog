package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devlog",
	Short: "devlog — daily developer productivity log",
	Long: `devlog helps you track your work by collecting git commits and notes,
then generating a clean standup report.`,
	// Errors are printed once by Execute; usage is only shown for --help
	// or flag/argument mistakes, not for runtime failures.
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute is the entry point called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
