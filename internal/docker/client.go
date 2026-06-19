package docker

import (
	"context"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// ContainerLister abstracts Docker container listing for testability.
type ContainerLister interface {
	ListContainers(ctx context.Context, all bool) ([]Container, error)
}

// DockerClient implements ContainerLister using the real Docker Engine API.
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient creates a client that talks to the local Docker daemon.
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

// Close releases the underlying Docker client resources.
func (d *DockerClient) Close() error {
	return d.cli.Close()
}

// ListContainers returns all (or only running) containers from the Docker daemon.
func (d *DockerClient) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	result, err := d.cli.ContainerList(ctx, client.ContainerListOptions{All: all})
	if err != nil {
		return nil, err
	}
	return convertSummaries(result.Items), nil
}

func convertSummaries(raw []container.Summary) []Container {
	out := make([]Container, 0, len(raw))
	for _, r := range raw {
		c := Container{
			ID:      r.ID[:12],
			Name:    cleanName(r.Names),
			Image:   r.Image,
			Command: r.Command,
			Created: time.Unix(r.Created, 0),
			Status:  r.Status,
			State:   string(r.State),
			Health:  extractHealthFromSummary(r.Health),
			Ports:   convertPorts(r.Ports),
			Labels:  r.Labels,
			Running: r.State == container.StateRunning,
		}
		out = append(out, c)
	}
	return out
}

func cleanName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

func extractHealthFromSummary(h *container.HealthSummary) string {
	if h == nil {
		return ""
	}
	return string(h.Status)
}

func convertPorts(raw []container.PortSummary) []Port {
	ports := make([]Port, 0, len(raw))
	for _, p := range raw {
		ip := ""
		if p.IP.IsValid() {
			ip = p.IP.String()
		}
		ports = append(ports, Port{
			IP:          ip,
			PrivatePort: p.PrivatePort,
			PublicPort:  p.PublicPort,
			Type:        p.Type,
		})
	}
	return ports
}
