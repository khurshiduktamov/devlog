package blockers

import "testing"

func useTempHome(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
}

func activeMessages(t *testing.T) []string {
	t.Helper()
	active, err := GetActiveBlockers()
	if err != nil {
		t.Fatal(err)
	}
	out := []string{}
	for _, b := range active {
		out = append(out, b.Message)
	}
	return out
}

func TestResolve(t *testing.T) {
	useTempHome(t)
	for _, m := range []string{"a", "b", "c"} {
		if err := AddBlocker(m); err != nil {
			t.Fatal(err)
		}
	}

	resolved, err := Resolve([]int{3, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 2 || resolved[0].Message != "a" || resolved[1].Message != "c" {
		t.Fatalf("resolved = %+v", resolved)
	}
	if got := activeMessages(t); len(got) != 1 || got[0] != "b" {
		t.Fatalf("active = %v, want [b]", got)
	}
}

func TestResolveInvalidNumberChangesNothing(t *testing.T) {
	useTempHome(t)
	AddBlocker("a")

	if _, err := Resolve([]int{1, 2}); err == nil {
		t.Fatal("expected error for #2")
	}
	if got := activeMessages(t); len(got) != 1 {
		t.Fatalf("active = %v, want [a] untouched", got)
	}
}

func TestResolveWithNoBlockers(t *testing.T) {
	useTempHome(t)
	if _, err := Resolve([]int{1}); err == nil {
		t.Fatal("expected error")
	}
}
