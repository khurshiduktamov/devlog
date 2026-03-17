package cmd

import (
	"fmt"

	"github.com/khurshiduktamov/devlog/internal/notes"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [message]",
	Short: "Save a manual work note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNote,
}

func init() {
	rootCmd.AddCommand(noteCmd)
}

func runNote(cmd *cobra.Command, args []string) error {
	message := args[0]

	if err := notes.AddNote(message); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	fmt.Println("Note saved.")
	return nil
}