// Package workday computes the time window a standup report covers.
package workday

import "time"

// PreviousStart returns midnight (in now's location) at the start of the
// most recent weekday before now's date. On Monday, Saturday or Sunday
// that is Friday; on any other day it is the day before.
func PreviousStart(now time.Time) time.Time {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		day = day.AddDate(0, 0, -1)
		if wd := day.Weekday(); wd != time.Saturday && wd != time.Sunday {
			return day
		}
	}
}
