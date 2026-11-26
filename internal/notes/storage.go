package notes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/khurshiduktamov/devlog/internal/store"
)

// FileName is the data file notes are stored in, inside ~/.devlog.
const FileName = "notes.json"

// Note represents a single manually entered developer note.
type Note struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// loadAll reads all notes. Returns an empty slice if none are saved yet.
func loadAll() ([]Note, error) {
	notes := []Note{}
	if err := store.Load(FileName, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// AddNote appends a new note with the current timestamp to notes.json.
func AddNote(message string) error {
	notes, err := loadAll()
	if err != nil {
		return err
	}

	notes = append(notes, Note{
		Message: message,
		Time:    time.Now(),
	})

	return store.Save(FileName, notes)
}

// parseDuration converts a human-readable duration string like
// "24 hours ago" or "48 hours ago" into a time.Duration.
func parseDuration(duration string) (time.Duration, error) {
	duration = strings.TrimSuffix(strings.TrimSpace(duration), " ago")
	parts := strings.Fields(duration)

	if len(parts) != 2 {
		return 0, fmt.Errorf("unsupported duration format: %q", duration)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %q", parts[0])
	}

	unit := strings.ToLower(parts[1])
	switch unit {
	case "hour", "hours":
		return time.Duration(value) * time.Hour, nil
	case "day", "days":
		return time.Duration(value) * 24 * time.Hour, nil
	case "minute", "minutes":
		return time.Duration(value) * time.Minute, nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %q", unit)
	}
}

// GetNotesSince returns all notes created within the given duration.
// Example duration values: "24 hours ago", "48 hours ago".
func GetNotesSince(duration string) ([]Note, error) {
	d, err := parseDuration(duration)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-d)

	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	filtered := make([]Note, 0)
	for _, n := range all {
		if n.Time.After(cutoff) {
			filtered = append(filtered, n)
		}
	}

	return filtered, nil
}
