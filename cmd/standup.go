package cmd

import (
	"fmt"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/khurshiduktamov/devlog/internal/notes"
	"github.com/khurshiduktamov/devlog/internal/report"
	"github.com/spf13/cobra"
)

var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate a daily standup report from commits and notes",
	RunE:  runStandup,
}

func init() {
	rootCmd.AddCommand(standupCmd)
}

func runStandup(cmd *cobra.Command, args []string) error {
	// Fetch git commits from the last 48 hours.
	commits, err := git.GetCommitsSince("48 hours ago")
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	// Fetch saved notes from the last 48 hours.
	savedNotes, err := notes.GetNotesSince("48 hours ago")
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	// Fetch active blockers from the last 24 hours.
	activeBlockers, err := blockers.GetActiveBlockers()
	if err != nil {
		return fmt.Errorf("failed to load blockers: %w", err)
	}

	// Convert commits to activities.
	activities := activity.FromCommits(commits)

	// Convert notes to activities and merge.
	for _, n := range savedNotes {
		activities = append(activities, activity.Activity{
			Type:    activity.TypeNote,
			Message: n.Message,
			Time:    n.Time,
		})
	}

	// Generate and print the standup report.
	fmt.Print(report.GenerateStandup(activities, activeBlockers))

	return nil
}