package cmd

import (
	"fmt"
	"strings"

	"github.com/khurshiduktamov/devlog/internal/notes"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note <message...>",
	Short: "Save a manual work note",
	Example: `  devlog note Investigated Redis timeout issue
  devlog note "Paired with Sam on the billing migration"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runNote,
}

func init() {
	rootCmd.AddCommand(noteCmd)
}

func runNote(cmd *cobra.Command, args []string) error {
	message, err := joinMessage(args)
	if err != nil {
		return err
	}

	if err := notes.AddNote(message); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	fmt.Println("Note saved.")
	return nil
}

// joinMessage turns command arguments into a single message, so both
// `devlog note fixed the bug` and `devlog note "fixed the bug"` work.
func joinMessage(args []string) (string, error) {
	message := strings.TrimSpace(strings.Join(args, " "))
	if message == "" {
		return "", fmt.Errorf("message cannot be empty")
	}
	return message, nil
}
