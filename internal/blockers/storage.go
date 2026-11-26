package blockers

import (
	"fmt"
	"sort"
	"time"

	"github.com/khurshiduktamov/devlog/internal/store"
)

// FileName is the data file blockers are stored in, inside ~/.devlog.
const FileName = "blockers.json"

// Blocker statuses.
const (
	StatusOpen     = "open"
	StatusResolved = "resolved"
)

// Blocker represents a single impediment noted by the developer.
type Blocker struct {
	Message    string     `json:"message"`
	Time       time.Time  `json:"time"`
	Status     string     `json:"status,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// isActive reports whether b should appear in standups and `blocker list`.
func isActive(b Blocker, now time.Time) bool {
	if b.Status == StatusResolved {
		return false
	}
	return b.Time.After(now.Add(-24 * time.Hour))
}

// loadAll reads all blockers. Returns an empty slice if none are saved yet.
func loadAll() ([]Blocker, error) {
	blockers := []Blocker{}
	if err := store.Load(FileName, &blockers); err != nil {
		return nil, err
	}
	return blockers, nil
}

// AddBlocker appends a new open blocker with the current timestamp.
func AddBlocker(message string) error {
	all, err := loadAll()
	if err != nil {
		return err
	}

	all = append(all, Blocker{
		Message: message,
		Time:    time.Now(),
		Status:  StatusOpen,
	})

	return store.Save(FileName, all)
}

// activeIndexes returns the positions in all of the active blockers,
// oldest first. Its order defines the numbers shown by `blocker list`.
func activeIndexes(all []Blocker, now time.Time) []int {
	idx := []int{}
	for i, b := range all {
		if isActive(b, now) {
			idx = append(idx, i)
		}
	}
	return idx
}

// GetActiveBlockers returns the blockers that are still active, oldest
// first. Its order matches the numbers used by Resolve.
func GetActiveBlockers() ([]Blocker, error) {
	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	active := []Blocker{}
	for _, i := range activeIndexes(all, time.Now()) {
		active = append(active, all[i])
	}
	return active, nil
}

// Resolve marks the active blockers with the given 1-based numbers (as
// shown by GetActiveBlockers) as resolved, and returns them. Either all
// numbers are valid and resolved, or nothing changes.
func Resolve(numbers []int) ([]Blocker, error) {
	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	active := activeIndexes(all, now)

	seen := map[int]bool{}
	targets := []int{}
	for _, n := range numbers {
		if n < 1 || n > len(active) {
			if len(active) == 0 {
				return nil, fmt.Errorf("there are no open blockers")
			}
			return nil, fmt.Errorf("no open blocker #%d (valid: 1-%d)", n, len(active))
		}
		if !seen[n] {
			seen[n] = true
			targets = append(targets, active[n-1])
		}
	}
	sort.Ints(targets)

	resolved := make([]Blocker, 0, len(targets))
	for _, i := range targets {
		all[i].Status = StatusResolved
		all[i].ResolvedAt = &now
		resolved = append(resolved, all[i])
	}

	if err := store.Save(FileName, all); err != nil {
		return nil, err
	}
	return resolved, nil
}
