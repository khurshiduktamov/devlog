package git

// Commit represents a single git commit entry.
type Commit struct {
	Message string
	Date    string
}

// GetCommitsSince returns commits from the current git repo within the given time range.
// duration examples: "24 hours ago", "48 hours ago"
func GetCommitsSince(duration string) ([]Commit, error) {
	return []Commit{}, nil // will be implemented in Step 5
}