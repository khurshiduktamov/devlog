package activity

import "time"

// ActivityType identifies the source of an activity entry.
type ActivityType string

const (
	// TypeCommit represents a git commit activity.
	TypeCommit ActivityType = "commit"

	// TypeNote represents a manually entered developer note.
	TypeNote ActivityType = "note"
)

// Activity is the unified representation of any developer activity.
// Both the git collector and notes storage will produce []Activity.
// The standup generator consumes []Activity regardless of source.
type Activity struct {
	Type    ActivityType
	Message string
	Time    time.Time
}