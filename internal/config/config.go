package config

import "os"

type Config struct {
	Addr               string
	DataDir            string
	TraefikDynamicPath string
	DockerSocketPath   string
	DockerProxyNetwork string
	SessionSecret      string
	Username           string
	Password           string
}

func FromEnv() Config {
	return Config{
		Addr:               env("NERDGATE_ADDR", ":8080"),
		DataDir:            env("NERDGATE_DATA_DIR", "./data/app"),
		TraefikDynamicPath: env("NERDGATE_TRAEFIK_DYNAMIC_PATH", "./data/traefik/routes.yml"),
		DockerSocketPath:   env("NERDGATE_DOCKER_SOCKET", "/var/run/docker.sock"),
		DockerProxyNetwork: env("NERDGATE_DOCKER_PROXY_NETWORK", "nerdgate-proxy"),
		SessionSecret:      env("NERDGATE_SESSION_SECRET", ""),
		Username:           env("NERDGATE_USERNAME", "admin"),
		Password:           env("NERDGATE_PASSWORD", "change-me"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
