package plans

import (
	"testing"
	"time"
)

func TestPlansForDay(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	today := time.Date(2026, 11, 26, 9, 0, 0, 0, time.Local)
	tomorrow := today.AddDate(0, 0, 1)

	AddPlan("today 1", today)
	AddPlan("tomorrow", tomorrow)
	AddPlan("today 2", today)

	got, err := GetPlansFor(today)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Message != "today 1" || got[1].Message != "today 2" {
		t.Fatalf("today = %+v", got)
	}

	got, _ = GetPlansFor(tomorrow)
	if len(got) != 1 || got[0].Message != "tomorrow" {
		t.Fatalf("tomorrow = %+v", got)
	}
}
