package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Formats lists the supported output formats, for help text and errors.
var Formats = []string{"text", "markdown", "slack", "json"}

// Render renders the report in the named format.
func Render(r StandupReport, formatName string) (string, error) {
	switch strings.ToLower(formatName) {
	case "", "text":
		return format(r), nil
	case "markdown", "md":
		return renderMarkup(r, "**%s**", "- ", "_none_"), nil
	case "slack":
		// Slack's message syntax: *bold*, and "•" renders as a bullet when
		// pasted into the message box.
		return renderMarkup(r, "*%s*", "• ", "_none_"), nil
	case "json":
		data, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data) + "\n", nil
	default:
		return "", fmt.Errorf("unknown format %q (use one of: %s)",
			formatName, strings.Join(Formats, ", "))
	}
}

// renderMarkup renders headings with headingFmt and list items with bullet.
func renderMarkup(r StandupReport, headingFmt, bullet, empty string) string {
	var sb strings.Builder
	section := func(title string, items []string) {
		sb.WriteString(fmt.Sprintf(headingFmt, title) + "\n")
		if len(items) == 0 {
			sb.WriteString(bullet + empty + "\n")
			return
		}
		for _, item := range items {
			sb.WriteString(bullet + item + "\n")
		}
	}

	section("Yesterday", r.Yesterday)
	sb.WriteString("\n")
	section("Today", r.Today)
	sb.WriteString("\n")
	section("Blockers", r.Blockers)
	return sb.String()
}
