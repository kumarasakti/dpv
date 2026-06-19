package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kumarasakti/dpv/internal/color"
	"github.com/kumarasakti/dpv/internal/docker"
)

// Formatter writes container data to output.
type Formatter interface {
	Format(w io.Writer, containers []docker.Container, stats docker.Stats) error
}

// Includes defines which optional columns to display.
type Includes struct {
	ContainerID bool // n
	Image       bool // i
	Command     bool // m
	Created     bool // c
	Status      bool // s
	Ports       bool // p
	Health      bool // h
}

// ParseIncludes converts a shorthand string like "cp" into an Includes struct.
// An empty string means "show all default columns".
func ParseIncludes(s string, slim bool) Includes {
	if s == "" && !slim {
		return Includes{
			ContainerID: true,
			Image:       true,
			Command:     true,
			Created:     true,
			Status:      true,
			Ports:       true,
		}
	}
	var inc Includes
	for _, ch := range s {
		switch ch {
		case 'n':
			inc.ContainerID = true
		case 'i':
			inc.Image = true
		case 'm':
			inc.Command = true
		case 'c':
			inc.Created = true
		case 's':
			inc.Status = true
		case 'p':
			inc.Ports = true
		case 'h':
			inc.Health = true
		}
	}
	return inc
}

// WriteTo is a convenience that formats to stdout.
func WriteTo(f Formatter, containers []docker.Container, stats docker.Stats) error {
	return f.Format(os.Stdout, containers, stats)
}

// printSummaryBar renders a compact one-line summary at the top.
func printSummaryBar(w io.Writer, stats docker.Stats) {
	stopped := stats.Total - stats.Running
	parts := []string{
		color.Green.Render(fmt.Sprintf("%d running", stats.Running)),
		color.Dim.Render(fmt.Sprintf("%d stopped", stopped)),
		fmt.Sprintf("%d total", stats.Total),
	}
	bar := "▸ " + strings.Join(parts, color.Dim.Render(" · "))
	if stats.HasSearch {
		bar += color.Dim.Render("  filter: ") + strings.Join(stats.SearchTerms, ", ")
	}
	fmt.Fprintln(w, bar)
}

// printSummaryFooter renders a compact search-match line if applicable.
func printSummaryFooter(w io.Writer, containers []docker.Container, stats docker.Stats) {
	if stats.HasSearch {
		fmt.Fprintf(w, "%s %d matched\n",
			color.Dim.Render("▸"),
			len(containers))
	}
}
