package report

import (
	"strings"
	"time"

	"github.com/khurshiduktamov/devlog/internal/activity"
)

// StandupReport holds the three sections of a daily standup.
type StandupReport struct {
	Yesterday []string
	Today     []string
	Blockers  string
}

// GenerateStandup splits activities into Yesterday and Today buckets
// based on whether their timestamp falls on the current calendar day.
func GenerateStandup(activities []activity.Activity) string {
	report := StandupReport{
		Blockers: "None",
	}

	today := time.Now().Format("2006-01-02")

	for _, a := range activities {
		activityDate := a.Time.Format("2006-01-02")
		if activityDate == today {
			report.Today = append(report.Today, a.Message)
		} else {
			report.Yesterday = append(report.Yesterday, a.Message)
		}
	}

	return format(report)
}

// format renders a StandupReport into a human-readable string.
func format(report StandupReport) string {
	var sb strings.Builder

	sb.WriteString("Yesterday:\n")
	if len(report.Yesterday) == 0 {
		sb.WriteString("  (empty)\n")
	} else {
		for _, item := range report.Yesterday {
			sb.WriteString("  • " + item + "\n")
		}
	}

	sb.WriteString("\nToday:\n")
	if len(report.Today) == 0 {
		sb.WriteString("  (empty)\n")
	} else {
		for _, item := range report.Today {
			sb.WriteString("  • " + item + "\n")
		}
	}

	sb.WriteString("\nBlockers:\n")
	sb.WriteString("  " + report.Blockers + "\n")

	return sb.String()
}