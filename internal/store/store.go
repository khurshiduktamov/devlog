// Package store reads and writes devlog's JSON data files in ~/.devlog.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Dir returns the path to the ~/.devlog directory.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not resolve home directory: %w", err)
	}
	return filepath.Join(home, ".devlog"), nil
}

// Path returns the full path to a data file inside ~/.devlog.
func Path(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

// Load decodes the JSON file ~/.devlog/<name> into v.
// A missing file is not an error: v is left untouched.
func Load(name string, v any) error {
	path, err := Path(name)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}
	return nil
}

// Save encodes v as indented JSON and writes it to ~/.devlog/<name>.
// The write is atomic: data goes to a temp file that is then renamed over
// the target, so an interrupted write never leaves a half-written file.
func Save(name string, v any) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize %s: %w", name, err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(dir, name+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // no-op once renamed

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write %s: %w", name, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to write %s: %w", name, err)
	}
	if err := os.Chmod(tmpPath, 0644); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", name, err)
	}

	return os.Rename(tmpPath, filepath.Join(dir, name))
}

// Exists reports whether ~/.devlog/<name> exists.
func Exists(name string) (bool, error) {
	path, err := Path(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}
