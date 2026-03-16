package notes

// Note represents a single developer note entry.
type Note struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

// SaveNote persists a new note to ~/.devlog/notes.json.
func SaveNote(text string) error {
	return nil // will be implemented in Step 6
}

// LoadNotes reads all notes from ~/.devlog/notes.json.
func LoadNotes() ([]Note, error) {
	return []Note{}, nil // will be implemented in Step 6
}