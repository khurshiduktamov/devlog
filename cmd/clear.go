package cmd

import (
	"fmt"

	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/khurshiduktamov/devlog/internal/notes"
	"github.com/khurshiduktamov/devlog/internal/plans"
	"github.com/khurshiduktamov/devlog/internal/store"
	"github.com/spf13/cobra"
)

var (
	clearNotes    bool
	clearBlockers bool
	clearPlans    bool
	clearAll      bool
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear stored notes, plans and/or blockers",
	Long: `Clear locally stored notes, plans and blockers.

Examples:
  devlog clear --notes
  devlog clear --blockers
  devlog clear --plans
  devlog clear --all`,
	RunE: runClear,
}

func init() {
	clearCmd.Flags().BoolVar(&clearNotes, "notes", false, "Clear all saved notes")
	clearCmd.Flags().BoolVar(&clearBlockers, "blockers", false, "Clear all saved blockers")
	clearCmd.Flags().BoolVar(&clearPlans, "plans", false, "Clear all saved plans")
	clearCmd.Flags().BoolVar(&clearAll, "all", false, "Clear all saved notes, plans and blockers")
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	// If no flags provided, show usage hint.
	if !clearNotes && !clearBlockers && !clearPlans && !clearAll {
		return fmt.Errorf("specify at least one flag: --notes, --plans, --blockers, or --all")
	}

	if clearAll || clearNotes {
		if err := clearFile(notes.FileName, "notes"); err != nil {
			return err
		}
	}

	if clearAll || clearPlans {
		if err := clearFile(plans.FileName, "plans"); err != nil {
			return err
		}
	}

	if clearAll || clearBlockers {
		if err := clearFile(blockers.FileName, "blockers"); err != nil {
			return err
		}
	}

	return nil
}

// clearFile replaces the given data file with an empty JSON array.
func clearFile(name, label string) error {
	exists, err := store.Exists(name)
	if err != nil {
		return err
	}
	// If file does not exist yet, nothing to clear.
	if !exists {
		fmt.Printf("No %s to clear.\n", label)
		return nil
	}

	if err := store.Save(name, []struct{}{}); err != nil {
		return fmt.Errorf("failed to clear %s: %w", label, err)
	}

	fmt.Printf("✓ %s cleared.\n", label)
	return nil
}
