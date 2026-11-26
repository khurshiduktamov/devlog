package report

import (
	"sort"
	"strings"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/khurshiduktamov/devlog/internal/plans"
)

// StandupReport holds the three sections of a daily standup.
type StandupReport struct {
	Yesterday []string `json:"yesterday"` // work done since the previous workday
	Today     []string `json:"today"`     // what you plan to do
	Blockers  []string `json:"blockers"`
}

// Build assembles a standup report:
//   - Yesterday: every commit and note in done, oldest first (on a Monday
//     this covers Friday too; it also includes anything done this morning
//     before standup)
//   - Today: the plans for today
//   - Blockers: the active blockers
func Build(done []activity.Activity, todayPlans []plans.Plan, activeBlockers []blockers.Blocker) StandupReport {
	sorted := make([]activity.Activity, len(done))
	copy(sorted, done)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Time.Before(sorted[j].Time)
	})

	report := StandupReport{
		Yesterday: []string{},
		Today:     []string{},
		Blockers:  []string{},
	}
	for _, a := range sorted {
		report.Yesterday = append(report.Yesterday, a.Message)
	}
	for _, p := range todayPlans {
		report.Today = append(report.Today, p.Message)
	}
	for _, b := range activeBlockers {
		report.Blockers = append(report.Blockers, b.Message)
	}
	return report
}

// GenerateStandup builds the report and renders it as plain text.
func GenerateStandup(done []activity.Activity, todayPlans []plans.Plan, activeBlockers []blockers.Blocker) string {
	return format(Build(done, todayPlans, activeBlockers))
}

// format renders a StandupReport into a human-readable string.
func format(report StandupReport) string {
	var sb strings.Builder

	section := func(title string, items []string, empty string) {
		sb.WriteString(title + ":\n")
		if len(items) == 0 {
			sb.WriteString("  " + empty + "\n")
			return
		}
		for _, item := range items {
			sb.WriteString("  • " + item + "\n")
		}
	}

	section("Yesterday", report.Yesterday, "(empty)")
	sb.WriteString("\n")
	section("Today", report.Today, "(empty)")
	sb.WriteString("\n")
	section("Blockers", report.Blockers, "None")

	return sb.String()
}
