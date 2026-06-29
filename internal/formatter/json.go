package formatter

import (
	"encoding/json"
	"io"

	"github.com/kumarasakti/dpv/internal/docker"
)

// JSONFormatter outputs container data as JSON.
type JSONFormatter struct{}

type jsonOutput struct {
	TotalContainers int                `json:"total_containers"`
	TotalRunning    int                `json:"total_running"`
	Containers      []docker.Container `json:"containers"`
}

func (f *JSONFormatter) Format(w io.Writer, containers []docker.Container, stats docker.Stats) error {
	if containers == nil {
		containers = []docker.Container{}
	}
	out := jsonOutput{
		TotalContainers: stats.Total,
		TotalRunning:    stats.Running,
		Containers:      containers,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
