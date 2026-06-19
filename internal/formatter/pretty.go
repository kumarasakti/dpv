package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/kumarasakti/dpv/internal/color"
	"github.com/kumarasakti/dpv/internal/docker"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var cardBorder = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1).
	MarginBottom(1)

// PrettyFormatter renders containers as bordered cards with status indicators.
type PrettyFormatter struct {
	Includes Includes
}

func (f *PrettyFormatter) Format(w io.Writer, containers []docker.Container, stats docker.Stats) error {
	printSummaryBar(w, stats)
	fmt.Fprintln(w)

	for i, c := range containers {
		var lines []string

		nameStyle := color.ForIndex(i)
		title := fmt.Sprintf("%s %s", color.StatusDot(c.Running), nameStyle.Bold(true).Render(c.Name))
		lines = append(lines, title)

		if f.Includes.Status {
			lines = append(lines, formatField("Status", c.Status))
		}
		if f.Includes.Health && c.Health != "" {
			val := c.Health
			switch c.Health {
			case "healthy":
				val = color.Green.Render(c.Health)
			case "unhealthy":
				val = color.Red.Render(c.Health)
			}
			lines = append(lines, formatField("Health", val))
		}
		if f.Includes.Created {
			ageStyle := color.AgeStyle(c.Created)
			lines = append(lines, formatField("Created", ageStyle.Render(humanize.Time(c.Created))))
		}
		if f.Includes.Ports {
			lines = append(lines, formatPorts(c.Ports)...)
		}
		if f.Includes.ContainerID {
			lines = append(lines, formatField("ID", color.Dim.Render(c.ID)))
		}
		if f.Includes.Image {
			lines = append(lines, formatField("Image", c.Image))
		}
		if f.Includes.Command {
			lines = append(lines, formatField("Command", color.Dim.Render(c.Command)))
		}
		if stats.ShowAll {
			stateLabel := color.Green.Render("running")
			if !c.Running {
				stateLabel = color.Red.Render(string(c.State))
			}
			lines = append(lines, formatField("State", stateLabel))
		}

		card := cardBorder.Render(strings.Join(lines, "\n"))
		fmt.Fprintln(w, card)
	}

	printSummaryFooter(w, containers, stats)
	return nil
}

const fieldWidth = 20

func formatField(label, value string) string {
	bold := color.Bold.Render(label + ":")
	return fmt.Sprintf("  %-*s %s", fieldWidth, bold, value)
}

func formatPorts(ports []docker.Port) []string {
	if len(ports) == 0 {
		return []string{formatField("Ports", color.Dim.Render("none"))}
	}
	var lines []string
	lines = append(lines, formatField("Ports", ports[0].String()))
	for _, p := range ports[1:] {
		lines = append(lines, fmt.Sprintf("  %-*s %s", fieldWidth, "", p.String()))
	}
	return lines
}
