package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Commit represents a single git commit entry.
type Commit struct {
	Message string
}

// GetCommitsSince returns commit messages from the current git repo
// since the given duration. Example duration values: "24 hours ago", "48 hours ago".
func GetCommitsSince(duration string) ([]Commit, error) {
	cmd := exec.Command(
		"git", "log",
		"--since="+duration,
		"--pretty=format:%s",
	)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git error: %s", string(exitErr.Stderr))
		}
		return nil, err
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return []Commit{}, nil
	}

	lines := strings.Split(raw, "\n")
	commits := make([]Commit, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			commits = append(commits, Commit{Message: line})
		}
	}

	return commits, nil
}