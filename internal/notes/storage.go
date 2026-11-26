package notes

import (
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

// GetNotesSince returns all notes created at or after since.
func GetNotesSince(since time.Time) ([]Note, error) {
	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	filtered := make([]Note, 0)
	for _, n := range all {
		if !n.Time.Before(since) {
			filtered = append(filtered, n)
		}
	}

	return filtered, nil
}
