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
}

// Execute is the entry point called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}