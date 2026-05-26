package traefik

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nerdgatehub/nerdgate-hub/internal/store"
)

var yamlIDPattern = regexp.MustCompile(`[^a-zA-Z0-9-]`)

type Renderer struct {
	path string
}

func NewRenderer(path string) *Renderer {
	return &Renderer{path: path}
}

func (r *Renderer) Render(routes []store.Route) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("http:\n")
	b.WriteString("  routers:\n")
	if len(routes) == 0 {
		b.WriteString("    {}\n")
	} else {
		for _, route := range routes {
			id := safeID(route.ID)
			entryPoint := "web"
			if route.TLS {
				entryPoint = "websecure"
			}

			fmt.Fprintf(&b, "    %s:\n", id)
			fmt.Fprintf(&b, "      rule: \"Host(`%s`)\"\n", escapeDoubleQuoted(route.Domain))
			fmt.Fprintf(&b, "      service: %s\n", id)
			fmt.Fprintf(&b, "      entryPoints:\n")
			fmt.Fprintf(&b, "        - %s\n", entryPoint)
			if route.TLS {
				fmt.Fprintf(&b, "      tls:\n")
				fmt.Fprintf(&b, "        certResolver: letsencrypt\n")
			}
		}
	}

	b.WriteString("  services:\n")
	if len(routes) == 0 {
		b.WriteString("    {}\n")
	} else {
		for _, route := range routes {
			id := safeID(route.ID)
			fmt.Fprintf(&b, "    %s:\n", id)
			fmt.Fprintf(&b, "      loadBalancer:\n")
			fmt.Fprintf(&b, "        servers:\n")
			fmt.Fprintf(&b, "          - url: \"%s\"\n", escapeDoubleQuoted(route.TargetURL))
		}
	}

	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}

func safeID(value string) string {
	id := strings.Trim(yamlIDPattern.ReplaceAllString(value, "-"), "-")
	if id == "" {
		return "route"
	}
	return id
}

func escapeDoubleQuoted(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return value
}
