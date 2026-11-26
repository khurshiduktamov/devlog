package workday

import "time"

// Next returns midnight at the start of the first weekday after now's date.
// On Friday, Saturday or Sunday that is Monday.
func Next(now time.Time) time.Time {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for {
		day = day.AddDate(0, 0, 1)
		if wd := day.Weekday(); wd != time.Saturday && wd != time.Sunday {
			return day
		}
	}
}
