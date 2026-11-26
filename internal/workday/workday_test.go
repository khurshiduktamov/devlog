package workday

import (
	"testing"
	"time"
)

func TestPreviousStart(t *testing.T) {
	loc := time.FixedZone("UTC+5", 5*3600)
	at := func(d int, h int) time.Time { return time.Date(2026, 11, d, h, 30, 0, 0, loc) }

	// November 2026: 20 = Fri, 21 = Sat, 22 = Sun, 23 = Mon, 24 = Tue.
	cases := []struct {
		name string
		now  time.Time
		want int // day of month at midnight
	}{
		{"Tuesday morning -> Monday", at(24, 9), 23},
		{"Monday morning -> Friday", at(23, 9), 20},
		{"Monday just after midnight -> Friday", at(23, 0), 20},
		{"Sunday -> Friday", at(22, 15), 20},
		{"Saturday -> Friday", at(21, 15), 20},
		{"Friday -> Thursday", at(20, 9), 19},
	}
	for _, c := range cases {
		got := PreviousStart(c.now)
		want := time.Date(2026, 11, c.want, 0, 0, 0, 0, loc)
		if !got.Equal(want) {
			t.Errorf("%s: got %v, want %v", c.name, got, want)
		}
	}
}

func TestPreviousStartAcrossMonth(t *testing.T) {
	// Monday 2 Nov 2026 -> Friday 30 Oct 2026.
	got := PreviousStart(time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC))
	if want := time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNext(t *testing.T) {
	// November 2026: 20 = Fri, 23 = Mon, 24 = Tue.
	cases := map[int]int{24: 25, 20: 23, 21: 23, 22: 23, 23: 24}
	for from, want := range cases {
		got := Next(time.Date(2026, 11, from, 17, 0, 0, 0, time.UTC))
		if got.Day() != want || got.Hour() != 0 {
			t.Errorf("Next(Nov %d) = %v, want Nov %d 00:00", from, got, want)
		}
	}
}
