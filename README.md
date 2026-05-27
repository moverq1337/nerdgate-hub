<div align="center">

<img src="site/logo.png" alt="NerdGate Hub" width="96" />

# NerdGate Hub

A tiny **open-source admin panel over [Traefik](https://traefik.io/)** for routing domains to Docker containers, host ports, and remote HTTP services. Built in Go, shipped as a single binary.

[**Docs**](https://moverq1337.github.io/nerdgate-hub/) ·
[**Get started**](https://moverq1337.github.io/nerdgate-hub/en/get-started) ·
[**Russian**](docs/README.ru.md) ·
[**GHCR image**](https://github.com/moverq1337/nerdgate-hub/pkgs/container/nerdgate-hub)

<p>
  <img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-111111?style=flat-square&labelColor=000000" />
  <img alt="Go 1.25" src="https://img.shields.io/badge/go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white&labelColor=000000" />
  <img alt="Traefik" src="https://img.shields.io/badge/proxy-Traefik-24A1C1?style=flat-square&logo=traefikproxy&logoColor=white&labelColor=000000" />
  <img alt="SQLite" src="https://img.shields.io/badge/storage-SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white&labelColor=000000" />
  <img alt="Docker" src="https://img.shields.io/badge/container-Docker-2496ED?style=flat-square&logo=docker&logoColor=white&labelColor=000000" />
</p>

</div>

---

## Why

Routing is a solved problem. The dashboards on top of it are usually too heavy.

```txt
domain  →  Docker container  /  host port  /  remote HTTP service
```

Traefik still owns the hard parts — ports `80`/`443`, HTTPS, Let's Encrypt, reverse proxying. NerdGate Hub owns **routes, users, settings, and dynamic config** — and nothing else.

| | NerdGate Hub | Nginx Proxy Manager |
|---|---|---|
| Reverse proxy | Traefik (TLS via ACME) | nginx (TLS via certbot) |
| Storage | SQLite (modernc, pure Go) | SQLite + filesystem |
| Binary | Single binary, no cgo | Container with node + nginx |
| Scope | Domain → target. That's it. | Hosts, streams, redirections, 404s, access lists, … |
| Mental model | Small | Medium |

---

## Features

- **One-line installer** with DNS preflight.
- **Browser first-run setup** with a one-time setup token.
- **Traefik dynamic config** generation from SQLite.
- **Multiple domains in one form submission** — comma, space, or newline separated.
- **Inline route editing** directly in the routes table.
- **Docker container picker** through the Docker Engine API, including a one-click attach to `nerdgate-proxy`.
- **Published host-port targets** (`http://host.docker.internal:3000`) and **remote targets** (`http://203.0.113.10:8080`).
- **Target health checks** and in-panel diagnostics.
- **Backup download + staged restore** for SQLite data and Let's Encrypt `acme.json`.
- **Account password change** in the panel.
- **Hardening**: CSRF tokens on every POST, rate-limited login/setup, strict session cookies, security headers (CSP, X-Frame-Options, Permissions-Policy), and an audit log.
- **Helper scripts**: install, backup, restore, reset-password, diagnose, uninstall.
- **Tagged releases** on GitHub (`v0.1.0`) with matching GHCR image tags.

---

## Quick install

Create an A record for the panel domain first:

```txt
nerdgate.example.com  →  SERVER_PUBLIC_IP
```

Then:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

or:

```sh
wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

The installer asks for the panel domain and Let's Encrypt email, checks DNS, writes `.env`, bootstraps the panel route, pulls the image, and starts Docker Compose. It prints a **one-time setup token**. Open the panel URL, paste the token, and create the first admin account in the browser.

---

## Operations

| Action | In panel | Helper script |
|---|---|---|
| Reset admin password | `Maintenance → Account password` | `reset-password.sh` |
| Download backup | `Diagnostics → Download backup` | `backup.sh` |
| Restore backup | `Maintenance → Restore backup` | `restore.sh` |
| Diagnose HTTPS / Traefik / DNS / certs | `Diagnostics` panel | `diagnose.sh` |
| Uninstall | — | `uninstall.sh` |

Helper scripts:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/restore.sh | sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

Update:

```sh
cd /opt/nerdgate-hub
docker compose pull
docker compose up -d
```

Tag a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The release workflow creates a GitHub Release, and the Docker workflow publishes the matching GHCR tag.

---

## Local development

Production installs pull `ghcr.io/moverq1337/nerdgate-hub:latest`. For source builds:

```sh
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Replace the example setup token before exposing the panel to the internet.

Docs-only changes do **not** rebuild the service image. The Docker image workflow is limited to `Dockerfile`, `go.mod`, `go.sum`, `cmd/**`, and `internal/**`. Site changes only trigger the GitHub Pages workflow.

The marketing site lives in [`site-next/`](site-next/) — Vite + React + Tailwind + shadcn + [cult-ui](https://www.cult-ui.com/). To run it locally:

```sh
cd site-next
pnpm install
pnpm dev
```

---

## Data layout

```txt
data/app/nerdgate.db       SQLite: users, settings, routes
data/traefik/routes.yml    Traefik dynamic config
data/acme/acme.json        Let's Encrypt certificates
```

Backup archives include `nerdgate.db`, optional legacy `routes.json`, optional `acme.json`, and backup metadata.

---

## Roadmap

- Cloudflare DNS automation
- Multi-server agent mode
- Optional tunnel mode for servers without public 80/443

---

## GitHub Pages

If the Pages workflow fails with `Get Pages site failed`, enable Pages once:

```txt
Settings → Pages → Build and deployment → Source → GitHub Actions
```

Alternatively, add a repository secret named `PAGES_TOKEN` with Pages/admin rights. The workflow can then enable Pages automatically.

---

<div align="center">

**[MIT licensed](LICENSE)** · Built with care, shipped without bloat.

</div>
