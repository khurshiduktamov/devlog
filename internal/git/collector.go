package git

import (
	"errors"
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

// Options controls which commits GetCommits returns.
type Options struct {
	// Dir is the repository to read; empty means the current directory.
	Dir string

	// Since is a git date expression, e.g. "24 hours ago".
	Since string

	// Author limits commits to a single author. When AuthorExact is true it
	// must equal the author's email; otherwise it is passed to
	// `git log --author` as a pattern matched against "Name <email>".
	// Empty means all authors.
	Author      string
	AuthorExact bool
}

// CurrentUserEmail returns the user.email configured for the repo in dir
// (falling back to the global config), or "" if none is set.
// An empty dir means the current directory.
func CurrentUserEmail(dir string) string {
	cmd := exec.Command("git", "config", "user.email")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GetCommits returns commits from the current git repo matching opts.
func GetCommits(opts Options) ([]Commit, error) {
	args := []string{
		"log",
		"--since=" + opts.Since,
		"--pretty=format:%s%x09%aI",
	}
	if opts.Author != "" {
		if opts.AuthorExact {
			// Wrapping in <> matches the whole email, so "al@x.com"
			// doesn't also match "sal@x.com".
			args = append(args, "--fixed-strings", "--author=<"+opts.Author+">")
		} else {
			args = append(args, "--author="+opts.Author)
		}
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = opts.Dir
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("git error: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}

	return parseLog(string(output)), nil
}

// parseLog parses `git log --pretty=format:%s%x09%aI` output.
func parseLog(output string) []Commit {
	raw := strings.TrimSpace(output)
	if raw == "" {
		return []Commit{}
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

		commitTime := time.Now()
		if len(parts) == 2 {
			if t, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1])); err == nil {
				commitTime = t
			}
		}

		commits = append(commits, Commit{
			Message: message,
			Time:    commitTime,
		})
	}

	return commits
}
