package cmd

import (
	"fmt"

	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/spf13/cobra"
)

var blockerCmd = &cobra.Command{
	Use:   "blocker [message]",
	Short: "Save a blocker that will appear in your standup report",
	Args:  cobra.ExactArgs(1),
	RunE:  runBlocker,
}

func init() {
	rootCmd.AddCommand(blockerCmd)
}

func runBlocker(cmd *cobra.Command, args []string) error {
	message := args[0]

	if err := blockers.AddBlocker(message); err != nil {
		return fmt.Errorf("failed to save blocker: %w", err)
	}

	fmt.Println("Blocker saved.")
	return nil
}