package cmd

import (
	"fmt"
	"strconv"

	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/spf13/cobra"
)

var blockerCmd = &cobra.Command{
	Use:   "blocker [message...]",
	Short: "Save a blocker that will appear in your standup report",
	Long: `Save a blocker that will appear in your standup report.

With no message, lists your open blockers.`,
	Example: `  devlog blocker Waiting for AWS credentials from DevOps
  devlog blocker list
  devlog blocker resolve 1`,
	Args: cobra.ArbitraryArgs,
	RunE: runBlocker,
}

var blockerListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List open blockers with their numbers",
	Args:    cobra.NoArgs,
	RunE:    runBlockerList,
}

var blockerResolveCmd = &cobra.Command{
	Use:     "resolve <number>...",
	Aliases: []string{"done", "rm"},
	Short:   "Mark one or more blockers as resolved",
	Example: `  devlog blocker resolve 1
  devlog blocker resolve 1 3`,
	Args: cobra.MinimumNArgs(1),
	RunE: runBlockerResolve,
}

func init() {
	blockerCmd.AddCommand(blockerListCmd, blockerResolveCmd)
	rootCmd.AddCommand(blockerCmd)
}

func runBlocker(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runBlockerList(cmd, args)
	}

	message, err := joinMessage(args)
	if err != nil {
		return err
	}

	if err := blockers.AddBlocker(message); err != nil {
		return fmt.Errorf("failed to save blocker: %w", err)
	}

	fmt.Println("Blocker saved.")
	return nil
}

func runBlockerList(cmd *cobra.Command, args []string) error {
	active, err := blockers.GetActiveBlockers()
	if err != nil {
		return fmt.Errorf("failed to load blockers: %w", err)
	}

	if len(active) == 0 {
		fmt.Println("No open blockers.")
		return nil
	}

	fmt.Println("Open blockers:")
	for i, b := range active {
		fmt.Printf("  %d. %s  (since %s)\n", i+1, b.Message, b.Time.Format("Mon Jan 2"))
	}
	fmt.Println("\nResolve one with: devlog blocker resolve <number>")
	return nil
}

func runBlockerResolve(cmd *cobra.Command, args []string) error {
	numbers := make([]int, 0, len(args))
	for _, a := range args {
		n, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("%q is not a blocker number — see devlog blocker list", a)
		}
		numbers = append(numbers, n)
	}

	resolved, err := blockers.Resolve(numbers)
	if err != nil {
		return err
	}

	for _, b := range resolved {
		fmt.Printf("✓ Resolved: %s\n", b.Message)
	}
	return nil
}
