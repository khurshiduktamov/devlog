// Package plans stores what the developer intends to work on, which fills
// the "Today" section of the standup.
package plans

import (
	"time"

	"github.com/khurshiduktamov/devlog/internal/store"
)

// FileName is the data file plans are stored in, inside ~/.devlog.
const FileName = "plans.json"

// dateLayout is the format of Plan.For.
const dateLayout = "2006-01-02"

// Plan is one item the developer intends to work on.
type Plan struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"` // when it was written
	For     string    `json:"for"`  // the day it is planned for, YYYY-MM-DD
}

// DateString formats day as a Plan.For value.
func DateString(day time.Time) string {
	return day.Format(dateLayout)
}

func loadAll() ([]Plan, error) {
	plans := []Plan{}
	if err := store.Load(FileName, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

// AddPlan saves a plan for the given day.
func AddPlan(message string, day time.Time) error {
	all, err := loadAll()
	if err != nil {
		return err
	}

	all = append(all, Plan{
		Message: message,
		Time:    time.Now(),
		For:     DateString(day),
	})

	return store.Save(FileName, all)
}

// GetPlansFor returns the plans for the given day, in the order written.
func GetPlansFor(day time.Time) ([]Plan, error) {
	all, err := loadAll()
	if err != nil {
		return nil, err
	}

	want := DateString(day)
	out := []Plan{}
	for _, p := range all {
		if p.For == want {
			out = append(out, p)
		}
	}
	return out, nil
}
