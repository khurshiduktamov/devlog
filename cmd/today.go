package cmd

import (
	"fmt"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's git commits",
	RunE:  runToday,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

func runToday(cmd *cobra.Command, args []string) error {
	commits, err := git.GetCommitsSince("24 hours ago")
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	activities := activity.FromCommits(commits)

	if len(activities) == 0 {
		fmt.Println("No activity found in the last 24 hours.")
		return nil
	}

	fmt.Println("Today:")
	for _, a := range activities {
		fmt.Printf("  • %s\n", a.Message)
	}

	return nil
}