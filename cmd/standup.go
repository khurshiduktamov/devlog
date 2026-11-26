package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/khurshiduktamov/devlog/internal/notes"
	"github.com/khurshiduktamov/devlog/internal/plans"
	"github.com/khurshiduktamov/devlog/internal/report"
	"github.com/khurshiduktamov/devlog/internal/workday"
	"github.com/spf13/cobra"
)

var standupCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate a daily standup report from commits and notes",
	RunE:  runStandup,
}

var (
	standupAuthor authorFlags
	standupFormat string
)

func init() {
	standupAuthor.register(standupCmd)
	standupCmd.Flags().StringVarP(&standupFormat, "format", "f", "text",
		"output format: "+strings.Join(report.Formats, ", "))
	rootCmd.AddCommand(standupCmd)
}

func runStandup(cmd *cobra.Command, args []string) error {
	// Cover everything since the start of the previous workday, so a
	// Monday standup includes Friday.
	since := workday.PreviousStart(time.Now())

	opts := git.Options{Since: since}
	standupAuthor.apply(&opts)

	commits, err := git.GetCommits(opts)
	switch {
	case errors.Is(err, git.ErrNotRepository), errors.Is(err, git.ErrGitNotInstalled):
		// Notes and blockers are still useful without commits.
		fmt.Fprintf(os.Stderr, "Note: %v — showing notes and blockers only.\n\n", err)
		commits = nil
	case err != nil:
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	savedNotes, err := notes.GetNotesSince(since)
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	todayPlans, err := plans.GetPlansFor(time.Now())
	if err != nil {
		return fmt.Errorf("failed to load plans: %w", err)
	}

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

	out, err := report.Render(report.Build(activities, todayPlans, activeBlockers), standupFormat)
	if err != nil {
		return err
	}
	fmt.Print(out)

	return nil
}
