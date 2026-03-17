package notes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Note represents a single manually entered developer note.
type Note struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// devlogDir returns the path to the ~/.devlog directory.
func devlogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not resolve home directory: %w", err)
	}
	return filepath.Join(home, ".devlog"), nil
}

// notesFilePath returns the full path to ~/.devlog/notes.json.
func notesFilePath() (string, error) {
	dir, err := devlogDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "notes.json"), nil
}

// ensureDevlogDir creates ~/.devlog if it does not exist.
func ensureDevlogDir() error {
	dir, err := devlogDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// loadAll reads all notes from ~/.devlog/notes.json.
// Returns an empty slice if the file does not exist yet.
func loadAll() ([]Note, error) {
	path, err := notesFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Note{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, fmt.Errorf("failed to parse notes file: %w", err)
	}

	return notes, nil
}

// saveAll writes the full notes slice to ~/.devlog/notes.json.
func saveAll(notes []Note) error {
	if err := ensureDevlogDir(); err != nil {
		return err
	}

	path, err := notesFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize notes: %w", err)
	}

	return os.WriteFile(path, data, 0644)
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

	return saveAll(notes)
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