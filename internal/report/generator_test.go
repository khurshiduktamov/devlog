package report

import (
	"reflect"
	"testing"
	"time"

	"github.com/khurshiduktamov/devlog/internal/activity"
	"github.com/khurshiduktamov/devlog/internal/blockers"
	"github.com/khurshiduktamov/devlog/internal/plans"
)

func TestBuildSections(t *testing.T) {
	now := time.Now()
	done := []activity.Activity{
		// git log returns newest first; notes are appended after.
		{Type: activity.TypeCommit, Message: "commit this morning", Time: now.Add(-time.Hour)},
		{Type: activity.TypeCommit, Message: "commit yesterday", Time: now.Add(-20 * time.Hour)},
		{Type: activity.TypeNote, Message: "note yesterday", Time: now.Add(-18 * time.Hour)},
	}
	todayPlans := []plans.Plan{{Message: "ship it"}}
	active := []blockers.Blocker{{Message: "need creds"}}

	r := Build(done, todayPlans, active)

	wantYesterday := []string{"commit yesterday", "note yesterday", "commit this morning"}
	if !reflect.DeepEqual(r.Yesterday, wantYesterday) {
		t.Errorf("Yesterday = %v, want %v", r.Yesterday, wantYesterday)
	}
	if !reflect.DeepEqual(r.Today, []string{"ship it"}) {
		t.Errorf("Today = %v, want only plans", r.Today)
	}
	if !reflect.DeepEqual(r.Blockers, []string{"need creds"}) {
		t.Errorf("Blockers = %v", r.Blockers)
	}
}

func TestGenerateStandupEmpty(t *testing.T) {
	got := GenerateStandup(nil, nil, nil)
	want := "Yesterday:\n  (empty)\n\nToday:\n  (empty)\n\nBlockers:\n  None\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}
