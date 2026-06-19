package docker

import (
	"context"
	"testing"
	"time"
)

// MockLister is a test double for ContainerLister.
type MockLister struct {
	Containers []Container
	Err        error
}

func (m *MockLister) ListContainers(_ context.Context, _ bool) ([]Container, error) {
	return m.Containers, m.Err
}

func TestMockLister_ReturnsContainers(t *testing.T) {
	mock := &MockLister{
		Containers: []Container{
			{ID: "abc123", Name: "web", Running: true, Created: time.Now()},
			{ID: "def456", Name: "db", Running: false, Created: time.Now()},
		},
	}
	got, err := mock.ListContainers(context.Background(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(got))
	}
}

func TestMockLister_ReturnsError(t *testing.T) {
	mock := &MockLister{Err: context.DeadlineExceeded}
	_, err := mock.ListContainers(context.Background(), false)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestCleanName(t *testing.T) {
	tests := []struct {
		input []string
		want  string
	}{
		{nil, ""},
		{[]string{"/my-container"}, "my-container"},
		{[]string{"no-slash"}, "no-slash"},
	}
	for _, tt := range tests {
		got := cleanName(tt.input)
		if got != tt.want {
			t.Errorf("cleanName(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPortString(t *testing.T) {
	tests := []struct {
		port Port
		want string
	}{
		{Port{PrivatePort: 80, Type: "tcp"}, "80/tcp"},
		{Port{IP: "0.0.0.0", PublicPort: 8080, PrivatePort: 80, Type: "tcp"}, "0.0.0.0:8080->80/tcp"},
		{Port{PublicPort: 443, PrivatePort: 443, Type: "tcp"}, "443->443/tcp"},
	}
	for _, tt := range tests {
		got := tt.port.String()
		if got != tt.want {
			t.Errorf("Port.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestComposeProject(t *testing.T) {
	c := Container{Labels: map[string]string{"com.docker.compose.project": "myapp"}}
	if c.ComposeProject() != "myapp" {
		t.Errorf("expected 'myapp', got %q", c.ComposeProject())
	}
	empty := Container{}
	if empty.ComposeProject() != "" {
		t.Errorf("expected empty string for no labels, got %q", empty.ComposeProject())
	}
}
