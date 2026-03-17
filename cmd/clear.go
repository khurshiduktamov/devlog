package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	clearNotes    bool
	clearBlockers bool
	clearAll      bool
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear stored notes and/or blockers",
	Long: `Clear locally stored notes and blockers.

Examples:
  devlog clear --notes
  devlog clear --blockers
  devlog clear --all`,
	RunE: runClear,
}

func init() {
	clearCmd.Flags().BoolVar(&clearNotes, "notes", false, "Clear all saved notes")
	clearCmd.Flags().BoolVar(&clearBlockers, "blockers", false, "Clear all saved blockers")
	clearCmd.Flags().BoolVar(&clearAll, "all", false, "Clear all saved notes and blockers")
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	// If no flags provided, show usage hint.
	if !clearNotes && !clearBlockers && !clearAll {
		return fmt.Errorf("specify at least one flag: --notes, --blockers, or --all")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not resolve home directory: %w", err)
	}

	devlogDir := filepath.Join(home, ".devlog")

	if clearAll || clearNotes {
		if err := clearFile(filepath.Join(devlogDir, "notes.json"), "notes"); err != nil {
			return err
		}
	}

	if clearAll || clearBlockers {
		if err := clearFile(filepath.Join(devlogDir, "blockers.json"), "blockers"); err != nil {
			return err
		}
	}

	return nil
}

// clearFile writes an empty JSON array to the given file path.
func clearFile(path, label string) error {
	// If file does not exist yet, nothing to clear.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("No %s to clear.\n", label)
		return nil
	}

	if err := os.WriteFile(path, []byte("[]\n"), 0644); err != nil {
		return fmt.Errorf("failed to clear %s: %w", label, err)
	}

	fmt.Printf("✓ %s cleared.\n", label)
	return nil
}