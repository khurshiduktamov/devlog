package blockers

import (
	"time"

	"github.com/khurshiduktamov/devlog/internal/store"
)

// FileName is the data file blockers are stored in, inside ~/.devlog.
const FileName = "blockers.json"

// Blocker represents a single impediment noted by the developer.
type Blocker struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// loadAll reads all blockers. Returns an empty slice if none are saved yet.
func loadAll() ([]Blocker, error) {
	blockers := []Blocker{}
	if err := store.Load(FileName, &blockers); err != nil {
		return nil, err
	}
	return blockers, nil
}

// AddBlocker appends a new blocker with the current timestamp.
func AddBlocker(message string) error {
	all, err := loadAll()
	if err != nil {
		return err
	}

	all = append(all, Blocker{
		Message: message,
		Time:    time.Now(),
	})

	return store.Save(FileName, all)
}

// GetActiveBlockers returns all blockers saved in the last 24 hours.
func GetActiveBlockers() ([]Blocker, error) {
	cutoff := time.Now().Add(-24 * time.Hour)

	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	filtered := make([]Blocker, 0)
	for _, b := range all {
		if b.Time.After(cutoff) {
			filtered = append(filtered, b)
		}
	}

	return filtered, nil
}
