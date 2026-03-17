package blockers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Blocker represents a single impediment noted by the developer.
type Blocker struct {
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

// blockersFilePath returns the full path to ~/.devlog/blockers.json.
func blockersFilePath() (string, error) {
	dir, err := devlogDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "blockers.json"), nil
}

// ensureDevlogDir creates ~/.devlog if it does not exist.
func ensureDevlogDir() error {
	dir, err := devlogDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// loadAll reads all blockers from ~/.devlog/blockers.json.
// Returns an empty slice if the file does not exist yet.
func loadAll() ([]Blocker, error) {
	path, err := blockersFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Blocker{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read blockers file: %w", err)
	}

	var blockers []Blocker
	if err := json.Unmarshal(data, &blockers); err != nil {
		return nil, fmt.Errorf("failed to parse blockers file: %w", err)
	}

	return blockers, nil
}

// saveAll writes the full blockers slice to ~/.devlog/blockers.json.
func saveAll(blockers []Blocker) error {
	if err := ensureDevlogDir(); err != nil {
		return err
	}

	path, err := blockersFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(blockers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize blockers: %w", err)
	}

	return os.WriteFile(path, data, 0644)
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

	return saveAll(all)
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