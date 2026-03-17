package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Commit represents a single git commit entry.
type Commit struct {
	Message string
	Time    time.Time
}

// GetCommitsSince returns commits from the current git repo
// since the given duration. Example: "24 hours ago", "48 hours ago".
// Format: <subject><TAB><ISO8601 date>
func GetCommitsSince(duration string) ([]Commit, error) {
	cmd := exec.Command(
		"git", "log",
		"--since="+duration,
		"--pretty=format:%s%x09%aI",
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
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)
		message := parts[0]

		var commitTime time.Time
		if len(parts) == 2 {
			commitTime, err = time.Parse(time.RFC3339, strings.TrimSpace(parts[1]))
			if err != nil {
				commitTime = time.Now()
			}
		} else {
			commitTime = time.Now()
		}

		commits = append(commits, Commit{
			Message: message,
			Time:    commitTime,
		})
	}

	return commits, nil
}