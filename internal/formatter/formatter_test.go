package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kumarasakti/dpv/internal/docker"
)

func sampleContainers() []docker.Container {
	return []docker.Container{
		{
			ID:      "abc123def456",
			Name:    "my-web",
			Image:   "nginx:latest",
			Command: "nginx -g daemon off;",
			Created: time.Now().Add(-2 * time.Hour),
			Status:  "Up 2 hours",
			State:   "running",
			Running: true,
			Ports: []docker.Port{
				{IP: "0.0.0.0", PublicPort: 8080, PrivatePort: 80, Type: "tcp"},
				{PrivatePort: 443, Type: "tcp"},
			},
		},
		{
			ID:      "789xyz000111",
			Name:    "my-db",
			Image:   "postgres:16",
			Command: "docker-entrypoint.sh postgres",
			Created: time.Now().Add(-24 * time.Hour),
			Status:  "Up 24 hours",
			State:   "running",
			Health:  "healthy",
			Running: true,
			Ports: []docker.Port{
				{IP: "0.0.0.0", PublicPort: 5432, PrivatePort: 5432, Type: "tcp"},
			},
		},
	}
}

func sampleStats() docker.Stats {
	return docker.Stats{Total: 5, Running: 2, Matched: 2}
}

func TestPrettyFormatter_ContainsNames(t *testing.T) {
	var buf bytes.Buffer
	f := &PrettyFormatter{Includes: ParseIncludes("", false)}
	_ = f.Format(&buf, sampleContainers(), sampleStats())
	out := buf.String()
	if !strings.Contains(out, "my-web") {
		t.Error("pretty output missing container name 'my-web'")
	}
	if !strings.Contains(out, "my-db") {
		t.Error("pretty output missing container name 'my-db'")
	}
}

func TestPrettyFormatter_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	f := &PrettyFormatter{Includes: ParseIncludes("", false)}
	_ = f.Format(&buf, sampleContainers(), sampleStats())
	out := buf.String()
	for _, want := range []string{"ID", "Image", "Command", "Status", "Ports", "Created"} {
		if !strings.Contains(out, want) {
			t.Errorf("pretty output missing field %q", want)
		}
	}
}

func TestPrettyFormatter_HasSummaryBar(t *testing.T) {
	var buf bytes.Buffer
	f := &PrettyFormatter{Includes: ParseIncludes("", false)}
	_ = f.Format(&buf, sampleContainers(), sampleStats())
	out := buf.String()
	if !strings.Contains(out, "running") || !strings.Contains(out, "total") {
		t.Error("pretty output missing summary bar with running/total counts")
	}
}

func TestPrettyFormatter_HasStatusDots(t *testing.T) {
	var buf bytes.Buffer
	f := &PrettyFormatter{Includes: ParseIncludes("", false)}
	containers := sampleContainers()
	containers = append(containers, docker.Container{
		ID: "stopped00000", Name: "stopped-one", Running: false, State: "exited",
		Created: time.Now().Add(-48 * time.Hour),
	})
	stats := docker.Stats{Total: 6, Running: 2, ShowAll: true}
	_ = f.Format(&buf, containers, stats)
	out := buf.String()
	if !strings.Contains(out, "●") {
		t.Error("expected green dot ● for running containers")
	}
	if !strings.Contains(out, "○") {
		t.Error("expected open dot ○ for stopped containers")
	}
}

func TestPrettyFormatter_PortsMultiline(t *testing.T) {
	var buf bytes.Buffer
	f := &PrettyFormatter{Includes: ParseIncludes("p", true)}
	_ = f.Format(&buf, sampleContainers()[:1], sampleStats())
	out := buf.String()
	if !strings.Contains(out, "0.0.0.0:8080->80/tcp") {
		t.Error("missing first port")
	}
	if !strings.Contains(out, "443/tcp") {
		t.Error("missing second port")
	}
}

func TestSlimFormatter_NamesWithDots(t *testing.T) {
	var buf bytes.Buffer
	f := &SlimFormatter{Includes: ParseIncludes("", true)}
	_ = f.Format(&buf, sampleContainers(), sampleStats())
	out := buf.String()
	if !strings.Contains(out, "my-web") || !strings.Contains(out, "my-db") {
		t.Error("slim output missing container names")
	}
	if !strings.Contains(out, "●") {
		t.Error("slim output should show status dots")
	}
	if strings.Contains(out, "ID:") {
		t.Error("slim output should not show ID without includes")
	}
}

func TestSlimFormatter_WithIncludes(t *testing.T) {
	var buf bytes.Buffer
	f := &SlimFormatter{Includes: ParseIncludes("cp", true)}
	_ = f.Format(&buf, sampleContainers(), sampleStats())
	out := buf.String()
	if !strings.Contains(out, "Created") {
		t.Error("slim -i=cp should include Created")
	}
	if !strings.Contains(out, "Ports") {
		t.Error("slim -i=cp should include Ports")
	}
}

func TestJSONFormatter_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	if err := f.Format(&buf, sampleContainers(), sampleStats()); err != nil {
		t.Fatalf("json format error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["total_containers"] != float64(5) {
		t.Errorf("expected total_containers=5, got %v", result["total_containers"])
	}
}

func TestParseIncludes_Empty(t *testing.T) {
	inc := ParseIncludes("", false)
	if !inc.ContainerID || !inc.Image || !inc.Command || !inc.Created || !inc.Status || !inc.Ports {
		t.Error("empty include with slim=false should enable all defaults")
	}
	if inc.Health {
		t.Error("health should not be in defaults")
	}
}

func TestParseIncludes_Specific(t *testing.T) {
	inc := ParseIncludes("cph", false)
	if !inc.Created || !inc.Ports || !inc.Health {
		t.Error("include 'cph' should enable created, ports, health")
	}
	if inc.Image || inc.ContainerID {
		t.Error("include 'cph' should not enable image or container_id")
	}
}

func TestSummaryBar_Content(t *testing.T) {
	var buf bytes.Buffer
	stats := docker.Stats{Total: 10, Running: 7, HasSearch: true, SearchTerms: []string{"web", "api"}}
	printSummaryBar(&buf, stats)
	out := buf.String()
	if !strings.Contains(out, "7 running") {
		t.Error("summary bar missing running count")
	}
	if !strings.Contains(out, "3 stopped") {
		t.Error("summary bar missing stopped count")
	}
	if !strings.Contains(out, "10 total") {
		t.Error("summary bar missing total count")
	}
	if !strings.Contains(out, "web") || !strings.Contains(out, "api") {
		t.Error("summary bar missing search terms")
	}
}
