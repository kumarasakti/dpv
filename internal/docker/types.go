package docker

import (
	"fmt"
	"time"
)

// Port represents a container port mapping.
type Port struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port,omitempty"`
	Type        string `json:"type"`
}

func (p Port) String() string {
	if p.PublicPort == 0 {
		return fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
	}
	if p.IP != "" {
		return fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type)
	}
	return fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type)
}

// Container holds normalized data for a single Docker container.
type Container struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Command string            `json:"command"`
	Created time.Time         `json:"created"`
	Status  string            `json:"status"`
	State   string            `json:"state"`
	Health  string            `json:"health,omitempty"`
	Ports   []Port            `json:"ports"`
	Labels  map[string]string `json:"labels,omitempty"`
	Running bool              `json:"running"`
}

// ComposeProject returns the Docker Compose project name, if any.
func (c Container) ComposeProject() string {
	return c.Labels["com.docker.compose.project"]
}

// Stats holds aggregate container counts for display.
type Stats struct {
	Total       int
	Running     int
	Matched     int
	HasSearch   bool
	SearchTerms []string
	ShowAll     bool
}
