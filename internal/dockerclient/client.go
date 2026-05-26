package dockerclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	socketPath   string
	proxyNetwork string
	httpClient   *http.Client
}

type Container struct {
	ID     string
	Name   string
	Image  string
	State  string
	Status string
	Ports  []Port
}

type Port struct {
	PrivatePort      int
	PublicPort       int
	Type             string
	NetworkReachable bool
	NetworkHost      string
}

type TargetOption struct {
	Label     string
	Value     string
	Available bool
	Kind      string
}

type apiContainer struct {
	ID     string    `json:"Id"`
	Names  []string  `json:"Names"`
	Image  string    `json:"Image"`
	State  string    `json:"State"`
	Status string    `json:"Status"`
	Ports  []apiPort `json:"Ports"`
}

type apiPort struct {
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

type inspectContainer struct {
	Config struct {
		ExposedPorts map[string]struct{} `json:"ExposedPorts"`
	} `json:"Config"`
	NetworkSettings struct {
		Networks map[string]inspectNetwork `json:"Networks"`
		Ports    map[string][]inspectPort  `json:"Ports"`
	} `json:"NetworkSettings"`
}

type inspectNetwork struct {
	Aliases []string `json:"Aliases"`
}

type inspectPort struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

func New(socketPath, proxyNetwork string) *Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "unix", socketPath)
		},
	}

	return &Client{
		socketPath:   socketPath,
		proxyNetwork: strings.TrimSpace(proxyNetwork),
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   3 * time.Second,
		},
	}
}

func (c *Client) ListContainers(ctx context.Context) ([]Container, error) {
	if strings.TrimSpace(c.socketPath) == "" {
		return nil, fmt.Errorf("docker socket path is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/containers/json?all=0", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("docker API returned %s", resp.Status)
	}

	var apiContainers []apiContainer
	if err := json.NewDecoder(resp.Body).Decode(&apiContainers); err != nil {
		return nil, err
	}

	containers := make([]Container, 0, len(apiContainers))
	for _, api := range apiContainers {
		ports := make([]Port, 0, len(api.Ports))
		networkHost := ""

		inspect, err := c.inspectContainer(ctx, api.ID)
		if err == nil {
			networkHost = c.networkHost(cleanName(api.Names), inspect)
			ports = append(ports, inspectedPorts(inspect, networkHost)...)
		}

		for _, port := range api.Ports {
			ports = append(ports, Port{
				PrivatePort:      port.PrivatePort,
				PublicPort:       port.PublicPort,
				Type:             port.Type,
				NetworkReachable: networkHost != "",
				NetworkHost:      networkHost,
			})
		}

		containers = append(containers, Container{
			ID:     api.ID,
			Name:   cleanName(api.Names),
			Image:  api.Image,
			State:  api.State,
			Status: api.Status,
			Ports:  ports,
		})
	}

	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Name < containers[j].Name
	})

	return containers, nil
}

func TargetOptions(containers []Container) []TargetOption {
	options := make([]TargetOption, 0)
	seen := make(map[string]bool)

	for _, container := range containers {
		if container.State != "running" {
			continue
		}

		for _, port := range container.Ports {
			if port.Type != "tcp" {
				continue
			}

			if port.PublicPort > 0 {
				value := fmt.Sprintf("http://host.docker.internal:%d", port.PublicPort)
				key := fmt.Sprintf("published:%s:%d:%d:%s", container.ID, port.PrivatePort, port.PublicPort, port.Type)
				if !seen[key] {
					seen[key] = true
					options = append(options, TargetOption{
						Label:     fmt.Sprintf("%s - published :%d -> :%d/%s", container.Name, port.PublicPort, port.PrivatePort, port.Type),
						Value:     value,
						Available: true,
						Kind:      "published",
					})
				}
			}

			if port.NetworkReachable && port.NetworkHost != "" {
				value := fmt.Sprintf("http://%s:%d", port.NetworkHost, port.PrivatePort)
				key := fmt.Sprintf("internal:%s:%d:%s", container.ID, port.PrivatePort, port.Type)
				if !seen[key] {
					seen[key] = true
					options = append(options, TargetOption{
						Label:     fmt.Sprintf("%s - internal :%d/%s", container.Name, port.PrivatePort, port.Type),
						Value:     value,
						Available: true,
						Kind:      "internal",
					})
				}
				continue
			}

			if port.PublicPort == 0 && port.PrivatePort > 0 {
				key := fmt.Sprintf("unavailable:%s:%d:%s", container.ID, port.PrivatePort, port.Type)
				if !seen[key] {
					seen[key] = true
					options = append(options, TargetOption{
						Label:     fmt.Sprintf("%s - :%d/%s - attach to nerdgate-proxy", container.Name, port.PrivatePort, port.Type),
						Available: false,
						Kind:      "unavailable",
					})
				}
			}
		}
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Label < options[j].Label
	})

	return options
}

func (c *Client) inspectContainer(ctx context.Context, id string) (inspectContainer, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/containers/"+id+"/json", nil)
	if err != nil {
		return inspectContainer{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return inspectContainer{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return inspectContainer{}, fmt.Errorf("docker inspect returned %s", resp.Status)
	}

	var inspect inspectContainer
	if err := json.NewDecoder(resp.Body).Decode(&inspect); err != nil {
		return inspectContainer{}, err
	}

	return inspect, nil
}

func (c *Client) networkHost(containerName string, inspect inspectContainer) string {
	if c.proxyNetwork == "" {
		return ""
	}

	network, ok := inspect.NetworkSettings.Networks[c.proxyNetwork]
	if !ok {
		return ""
	}

	for _, alias := range network.Aliases {
		alias = strings.TrimSpace(alias)
		if alias != "" {
			return alias
		}
	}

	return containerName
}

func inspectedPorts(inspect inspectContainer, networkHost string) []Port {
	ports := make([]Port, 0)

	for key := range inspect.Config.ExposedPorts {
		privatePort, portType, ok := parsePortKey(key)
		if !ok {
			continue
		}
		ports = append(ports, Port{
			PrivatePort:      privatePort,
			Type:             portType,
			NetworkReachable: networkHost != "",
			NetworkHost:      networkHost,
		})
	}

	for key, bindings := range inspect.NetworkSettings.Ports {
		privatePort, portType, ok := parsePortKey(key)
		if !ok {
			continue
		}
		if len(bindings) == 0 {
			ports = append(ports, Port{
				PrivatePort:      privatePort,
				Type:             portType,
				NetworkReachable: networkHost != "",
				NetworkHost:      networkHost,
			})
			continue
		}
		for _, binding := range bindings {
			publicPort, err := strconv.Atoi(binding.HostPort)
			if err != nil {
				continue
			}
			ports = append(ports, Port{
				PrivatePort:      privatePort,
				PublicPort:       publicPort,
				Type:             portType,
				NetworkReachable: networkHost != "",
				NetworkHost:      networkHost,
			})
		}
	}

	return ports
}

func parsePortKey(key string) (int, string, bool) {
	port, portType, ok := strings.Cut(key, "/")
	if !ok {
		return 0, "", false
	}
	privatePort, err := strconv.Atoi(port)
	if err != nil {
		return 0, "", false
	}
	return privatePort, portType, true
}

func cleanName(names []string) string {
	if len(names) == 0 {
		return "unknown"
	}
	name := strings.TrimSpace(names[0])
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return "unknown"
	}
	return name
}
