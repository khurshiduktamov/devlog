package report

import (
	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/khurshiduktamov/devlog/internal/notes"
)

// StandupReport holds all sections of a standup report.
type StandupReport struct {
	Yesterday []string
	Today     []string
	Blockers  string
}

// Generate builds a standup report from commits and notes.
func Generate(commits []git.Commit, notesList []notes.Note) StandupReport {
	return StandupReport{
		Blockers: "None",
	} // will be implemented in Step 9
}