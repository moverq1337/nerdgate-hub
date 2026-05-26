package dockerclient

import "testing"

func TestTargetOptionsUsesPublishedTCPPorts(t *testing.T) {
	containers := []Container{
		{
			ID:    "container-one",
			Name:  "app",
			State: "running",
			Ports: []Port{
				{PrivatePort: 3000, PublicPort: 3000, Type: "tcp"},
				{PrivatePort: 3000, PublicPort: 3000, Type: "tcp"},
				{PrivatePort: 5353, PublicPort: 5353, Type: "udp"},
				{PrivatePort: 8080, Type: "tcp"},
			},
		},
		{
			ID:    "container-two",
			Name:  "stopped",
			State: "exited",
			Ports: []Port{
				{PrivatePort: 9000, PublicPort: 9000, Type: "tcp"},
			},
		},
	}

	options := TargetOptions(containers)
	if len(options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(options))
	}

	published := findOptionByKind(options, "published")
	if published == nil {
		t.Fatal("expected published option")
	}

	if published.Value != "http://host.docker.internal:3000" {
		t.Fatalf("unexpected target value: %s", published.Value)
	}
}

func TestTargetOptionsUsesInternalPortsOnProxyNetwork(t *testing.T) {
	containers := []Container{
		{
			ID:    "container-one",
			Name:  "app",
			State: "running",
			Ports: []Port{
				{
					PrivatePort:      3000,
					Type:             "tcp",
					NetworkReachable: true,
					NetworkHost:      "app",
				},
			},
		},
	}

	options := TargetOptions(containers)
	if len(options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(options))
	}
	if options[0].Value != "http://app:3000" {
		t.Fatalf("unexpected target value: %s", options[0].Value)
	}
	if options[0].Kind != "internal" {
		t.Fatalf("unexpected target kind: %s", options[0].Kind)
	}
}

func TestTargetOptionsMarksUnreachableInternalPortsUnavailable(t *testing.T) {
	containers := []Container{
		{
			ID:    "container-one",
			Name:  "app",
			State: "running",
			Ports: []Port{
				{PrivatePort: 3000, Type: "tcp"},
			},
		},
	}

	options := TargetOptions(containers)
	if len(options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(options))
	}
	if options[0].Available {
		t.Fatal("expected unavailable option")
	}
}

func findOptionByKind(options []TargetOption, kind string) *TargetOption {
	for i := range options {
		if options[i].Kind == kind {
			return &options[i]
		}
	}
	return nil
}
