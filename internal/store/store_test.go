package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME")) // Windows

	in := []string{"a", "b"}
	if err := Save("x.json", in); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var out []string
	if err := Load("x.json", &out); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(out) != 2 || out[0] != "a" || out[1] != "b" {
		t.Fatalf("got %v, want %v", out, in)
	}

	// No temp files left behind.
	dir, _ := Dir()
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected only x.json in %s, got %d entries", dir, len(entries))
	}
}

func TestLoadMissingFileLeavesValue(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))

	out := []string{"keep"}
	if err := Load("missing.json", &out); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(out) != 1 || out[0] != "keep" {
		t.Fatalf("value changed: %v", out)
	}
}

func TestLoadCorruptFileErrors(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	os.MkdirAll(filepath.Join(home, ".devlog"), 0755)
	os.WriteFile(filepath.Join(home, ".devlog", "bad.json"), []byte("{not json"), 0644)

	var out []string
	if err := Load("bad.json", &out); err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
