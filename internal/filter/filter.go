package filter

import (
	"sort"
	"strings"

	"github.com/kumarasakti/dpv/internal/docker"
)

// Options controls how containers are filtered and sorted.
type Options struct {
	SearchTerms []string
	ShowAll     bool
	OrderBy     string // "status" (default), "name", "created"
	Reverse     bool
}

// Apply filters containers by search terms and running state, then sorts.
func Apply(containers []docker.Container, opts Options) []docker.Container {
	result := filterBySearch(containers, opts.SearchTerms)
	if !opts.ShowAll {
		result = filterRunning(result)
	}
	sortContainers(result, opts.OrderBy)
	if opts.Reverse {
		reverseSlice(result)
	}
	return result
}

func filterBySearch(containers []docker.Container, terms []string) []docker.Container {
	if len(terms) == 0 {
		return containers
	}
	var out []docker.Container
	for _, c := range containers {
		for _, t := range terms {
			if strings.Contains(c.Name, t) {
				out = append(out, c)
				break
			}
		}
	}
	return out
}

func filterRunning(containers []docker.Container) []docker.Container {
	var out []docker.Container
	for _, c := range containers {
		if c.Running {
			out = append(out, c)
		}
	}
	return out
}

func sortContainers(containers []docker.Container, field string) {
	sort.SliceStable(containers, func(i, j int) bool {
		switch field {
		case "name":
			return containers[i].Name < containers[j].Name
		case "created":
			return containers[i].Created.Before(containers[j].Created)
		default:
			return containers[i].Created.After(containers[j].Created)
		}
	})
}

func reverseSlice(containers []docker.Container) {
	for i, j := 0, len(containers)-1; i < j; i, j = i+1, j-1 {
		containers[i], containers[j] = containers[j], containers[i]
	}
}
