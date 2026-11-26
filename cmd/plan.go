package cmd

import (
	"fmt"
	"time"

	"github.com/khurshiduktamov/devlog/internal/plans"
	"github.com/khurshiduktamov/devlog/internal/workday"
	"github.com/spf13/cobra"
)

var planNextDay bool

var planCmd = &cobra.Command{
	Use:   "plan [message...]",
	Short: "Save what you plan to work on (the standup's Today section)",
	Long: `Save something you plan to work on. Plans fill the "Today" section of
devlog standup.

With no message, lists the plans for today (or the next workday with
--tomorrow).`,
	Example: `  devlog plan Finish Redis retry logic
  devlog plan --tomorrow Review Sam's billing PR
  devlog plan`,
	Args: cobra.ArbitraryArgs,
	RunE: runPlan,
}

func init() {
	planCmd.Flags().BoolVarP(&planNextDay, "tomorrow", "t", false,
		"plan for the next workday instead of today (Friday → Monday)")
	rootCmd.AddCommand(planCmd)
}

func runPlan(cmd *cobra.Command, args []string) error {
	day := time.Now()
	label := "today"
	if planNextDay {
		day = workday.Next(day)
		label = day.Format("Monday")
	}

	if len(args) == 0 {
		return listPlans(day, label)
	}

	message, err := joinMessage(args)
	if err != nil {
		return err
	}
	if err := plans.AddPlan(message, day); err != nil {
		return fmt.Errorf("failed to save plan: %w", err)
	}

	fmt.Printf("Plan saved for %s.\n", label)
	return nil
}

func listPlans(day time.Time, label string) error {
	list, err := plans.GetPlansFor(day)
	if err != nil {
		return fmt.Errorf("failed to load plans: %w", err)
	}
	if len(list) == 0 {
		fmt.Printf("No plans for %s. Add one with: devlog plan <message>\n", label)
		return nil
	}
	fmt.Printf("Plans for %s:\n", label)
	for _, p := range list {
		fmt.Printf("  • %s\n", p.Message)
	}
	return nil
}
