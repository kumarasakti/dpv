package filter

import (
	"testing"
	"time"

	"github.com/kumarasakti/dpv/internal/docker"
)

func makeContainers() []docker.Container {
	now := time.Now()
	return []docker.Container{
		{Name: "web-app", Running: true, Created: now.Add(-1 * time.Hour)},
		{Name: "db-postgres", Running: true, Created: now.Add(-2 * time.Hour)},
		{Name: "redis-cache", Running: false, Created: now.Add(-3 * time.Hour)},
		{Name: "web-worker", Running: true, Created: now.Add(-30 * time.Minute)},
	}
}

func TestApply_NoFilters_RunningOnly(t *testing.T) {
	result := Apply(makeContainers(), Options{})
	if len(result) != 3 {
		t.Fatalf("expected 3 running containers, got %d", len(result))
	}
	for _, c := range result {
		if !c.Running {
			t.Errorf("container %q should be running", c.Name)
		}
	}
}

func TestApply_ShowAll(t *testing.T) {
	result := Apply(makeContainers(), Options{ShowAll: true})
	if len(result) != 4 {
		t.Fatalf("expected 4 containers with --all, got %d", len(result))
	}
}

func TestApply_SearchFilter(t *testing.T) {
	result := Apply(makeContainers(), Options{
		SearchTerms: []string{"web"},
		ShowAll:     true,
	})
	if len(result) != 2 {
		t.Fatalf("expected 2 'web' containers, got %d", len(result))
	}
	for _, c := range result {
		if c.Name != "web-app" && c.Name != "web-worker" {
			t.Errorf("unexpected container %q in search results", c.Name)
		}
	}
}

func TestApply_SearchAndRunning(t *testing.T) {
	result := Apply(makeContainers(), Options{
		SearchTerms: []string{"redis"},
	})
	if len(result) != 0 {
		t.Fatalf("expected 0 (redis is stopped), got %d", len(result))
	}
}

func TestApply_OrderByName(t *testing.T) {
	result := Apply(makeContainers(), Options{ShowAll: true, OrderBy: "name"})
	if result[0].Name != "db-postgres" {
		t.Errorf("expected db-postgres first, got %s", result[0].Name)
	}
}

func TestApply_Reverse(t *testing.T) {
	result := Apply(makeContainers(), Options{ShowAll: true, OrderBy: "name", Reverse: true})
	if result[0].Name != "web-worker" {
		t.Errorf("expected web-worker first when reversed, got %s", result[0].Name)
	}
}
