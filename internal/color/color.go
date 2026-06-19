package color

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ponytail: simple cycle through a fixed palette; extend the slice if more colors needed.
var palette = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("12")),  // blue
	lipgloss.NewStyle().Foreground(lipgloss.Color("10")),  // green
	lipgloss.NewStyle().Foreground(lipgloss.Color("9")),   // red
	lipgloss.NewStyle().Foreground(lipgloss.Color("14")),  // cyan
	lipgloss.NewStyle().Foreground(lipgloss.Color("11")),  // yellow
	lipgloss.NewStyle().Foreground(lipgloss.Color("13")),  // magenta
}

var (
	Bold  = lipgloss.NewStyle().Bold(true)
	Green = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	Red   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	Dim   = lipgloss.NewStyle().Faint(true)
)

// ForIndex returns a style from the rotating palette for the given index.
func ForIndex(i int) lipgloss.Style {
	return palette[i%len(palette)]
}

// StatusDot returns a colored bullet: green ● for running, red ○ for stopped.
func StatusDot(running bool) string {
	if running {
		return Green.Render("●")
	}
	return Red.Render("○")
}

// AgeStyle returns a style based on how old a timestamp is:
// green for <1h, yellow for <24h, dim for older.
func AgeStyle(created time.Time) lipgloss.Style {
	age := time.Since(created)
	switch {
	case age < time.Hour:
		return Green
	case age < 24*time.Hour:
		return Yellow
	default:
		return Dim
	}
}
