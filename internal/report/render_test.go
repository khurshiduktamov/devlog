package report

import (
	"encoding/json"
	"testing"
)

var sample = StandupReport{
	Yesterday: []string{"Fixed bug"},
	Today:     []string{"Ship it"},
	Blockers:  []string{},
}

func TestRenderMarkdown(t *testing.T) {
	got, err := Render(sample, "markdown")
	if err != nil {
		t.Fatal(err)
	}
	want := "**Yesterday**\n- Fixed bug\n\n**Today**\n- Ship it\n\n**Blockers**\n- _none_\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderSlack(t *testing.T) {
	got, err := Render(sample, "slack")
	if err != nil {
		t.Fatal(err)
	}
	want := "*Yesterday*\n• Fixed bug\n\n*Today*\n• Ship it\n\n*Blockers*\n• _none_\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderJSON(t *testing.T) {
	got, err := Render(Build(nil, nil, nil), "json")
	if err != nil {
		t.Fatal(err)
	}
	var back map[string][]string
	if err := json.Unmarshal([]byte(got), &back); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, got)
	}
	// Empty sections must be [] rather than null, for scripts.
	for _, k := range []string{"yesterday", "today", "blockers"} {
		if v, ok := back[k]; !ok || v == nil {
			t.Errorf("%s = %v, want []", k, v)
		}
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	if _, err := Render(sample, "html"); err == nil {
		t.Fatal("expected error")
	}
}
