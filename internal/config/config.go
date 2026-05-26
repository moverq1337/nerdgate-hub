package config

import "os"

type Config struct {
	Addr               string
	DataDir            string
	TraefikDynamicPath string
	Username           string
	Password           string
}

func FromEnv() Config {
	return Config{
		Addr:               env("NERDGATE_ADDR", ":8080"),
		DataDir:            env("NERDGATE_DATA_DIR", "./data/app"),
		TraefikDynamicPath: env("NERDGATE_TRAEFIK_DYNAMIC_PATH", "./data/traefik/routes.yml"),
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
