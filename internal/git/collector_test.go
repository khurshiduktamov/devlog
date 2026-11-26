package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// newRepo creates a throwaway git repository and returns its path.
func newRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	// Keep the user's global git config out of the tests.
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, ".gitconfig-none"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	run(t, dir, "init", "-q", "-b", "main")
	run(t, dir, "config", "user.name", "Me")
	run(t, dir, "config", "user.email", "me@example.com")
	return dir
}

func run(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// commitAs makes an empty commit by the given author.
func commitAs(t *testing.T, dir, name, email, msg string) {
	t.Helper()
	cmd := exec.Command("git", "commit", "-q", "--allow-empty", "-m", msg)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME="+name, "GIT_AUTHOR_EMAIL="+email,
		"GIT_COMMITTER_NAME="+name, "GIT_COMMITTER_EMAIL="+email,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit: %v\n%s", err, out)
	}
}

func messages(cs []Commit) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Message
	}
	return out
}

func TestAuthorFilterExact(t *testing.T) {
	dir := newRepo(t)
	commitAs(t, dir, "Me", "me@example.com", "mine")
	commitAs(t, dir, "Other", "other@example.com", "theirs")
	// Shares a suffix with my email; must not match an exact filter.
	commitAs(t, dir, "Lookalike", "home@example.com", "lookalike")

	if got := CurrentUserEmail(dir); got != "me@example.com" {
		t.Fatalf("CurrentUserEmail = %q", got)
	}

	cs, err := GetCommits(Options{Dir: dir, Since: "1 hour ago", Author: "me@example.com", AuthorExact: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := messages(cs); len(got) != 1 || got[0] != "mine" {
		t.Fatalf("exact author filter: got %v, want [mine]", got)
	}

	all, err := GetCommits(Options{Dir: dir, Since: "1 hour ago"})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("no filter: got %v, want 3 commits", messages(all))
	}
}

func TestParseLog(t *testing.T) {
	cs := parseLog("Fix bug\t2026-11-20T10:00:00+05:00\n\nAdd cache\tgarbage\n")
	if len(cs) != 2 {
		t.Fatalf("got %d commits", len(cs))
	}
	if cs[0].Message != "Fix bug" || cs[0].Time.Day() != 20 {
		t.Fatalf("bad first commit: %+v", cs[0])
	}
	if cs[1].Message != "Add cache" || cs[1].Time.IsZero() {
		t.Fatalf("bad second commit: %+v", cs[1])
	}
}
