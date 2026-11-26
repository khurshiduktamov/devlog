package cmd

import (
	"fmt"
	"time"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's git commits",
	RunE:  runToday,
}

var todayAuthor authorFlags

func init() {
	todayAuthor.register(todayCmd)
	rootCmd.AddCommand(todayCmd)
}

func runToday(cmd *cobra.Command, args []string) error {
	opts := git.Options{Since: time.Now().Add(-24 * time.Hour)}
	todayAuthor.apply(&opts)

	commits, err := git.GetCommits(opts)
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
