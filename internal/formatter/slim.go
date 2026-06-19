package formatter

import (
	"fmt"
	"io"

	"github.com/kumarasakti/dpv/internal/color"
	"github.com/kumarasakti/dpv/internal/docker"
	"github.com/dustin/go-humanize"
)

// SlimFormatter renders a minimal status-dot + name list, with optional detail rows.
type SlimFormatter struct {
	Includes Includes
}

func (f *SlimFormatter) Format(w io.Writer, containers []docker.Container, stats docker.Stats) error {
	printSummaryBar(w, stats)
	fmt.Fprintln(w)

	hasDetails := f.Includes.ContainerID || f.Includes.Image || f.Includes.Command ||
		f.Includes.Created || f.Includes.Status || f.Includes.Ports || f.Includes.Health

	for i, c := range containers {
		nameStyle := color.ForIndex(i)
		fmt.Fprintf(w, "  %s %s\n", color.StatusDot(c.Running), nameStyle.Render(c.Name))

		if !hasDetails {
			continue
		}
		if f.Includes.ContainerID {
			printSlimField(w, "ID", color.Dim.Render(c.ID))
		}
		if f.Includes.Image {
			printSlimField(w, "Image", c.Image)
		}
		if f.Includes.Command {
			printSlimField(w, "Command", color.Dim.Render(c.Command))
		}
		if f.Includes.Created {
			ageStyle := color.AgeStyle(c.Created)
			printSlimField(w, "Created", ageStyle.Render(humanize.Time(c.Created)))
		}
		if f.Includes.Status {
			printSlimField(w, "Status", c.Status)
		}
		if f.Includes.Ports {
			if len(c.Ports) == 0 {
				printSlimField(w, "Ports", color.Dim.Render("none"))
			} else {
				for j, p := range c.Ports {
					if j == 0 {
						printSlimField(w, "Ports", p.String())
					} else {
						fmt.Fprintf(w, "      %-14s %s\n", "", p.String())
					}
				}
			}
		}
		if f.Includes.Health && c.Health != "" {
			printSlimField(w, "Health", c.Health)
		}
		fmt.Fprintln(w)
	}

	if !hasDetails {
		fmt.Fprintln(w)
	}
	printSummaryFooter(w, containers, stats)
	return nil
}

func printSlimField(w io.Writer, label, value string) {
	bold := color.Bold.Render(label + ":")
	fmt.Fprintf(w, "      %-14s %s\n", bold, value)
}
