package activity

import (
	"github.com/khurshiduktamov/devlog/internal/git"
)

// FromCommits converts a slice of git.Commit into a slice of Activity.
// Each commit becomes an Activity with TypeCommit, its message, and its timestamp.
func FromCommits(commits []git.Commit) []Activity {
	activities := make([]Activity, 0, len(commits))

	for _, c := range commits {
		activities = append(activities, Activity{
			Type:    TypeCommit,
			Message: c.Message,
			Time:    c.Time,
		})
	}

	return activities
}