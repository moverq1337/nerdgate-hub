<p align="center">
  <img src="site/logo.png" alt="NerdGate Hub" width="128">
</p>

<h1 align="center">NerdGate Hub</h1>

<p align="center">
  A tiny open-source Traefik dashboard for routing domains to Docker containers, host ports, and remote HTTP services.
</p>

<p align="center">
  <a href="https://moverq1337.github.io/nerdgate-hub/">Docs</a>
  ·
  <a href="https://moverq1337.github.io/nerdgate-hub/en/get-started.html">Get Started</a>
  ·
  <a href="https://moverq1337.github.io/nerdgate-hub/ru/get-started.html">Русская документация</a>
  ·
  <a href="https://github.com/moverq1337/nerdgate-hub/pkgs/container/nerdgate-hub">GHCR Image</a>
</p>

<p align="center">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-black">
  <img alt="Go" src="https://img.shields.io/badge/go-1.25-black">
  <img alt="Traefik" src="https://img.shields.io/badge/proxy-Traefik-black">
  <img alt="SQLite" src="https://img.shields.io/badge/storage-SQLite-black">
</p>

---

## English

NerdGate Hub is a small admin panel over Traefik. It does not try to become a platform. Its job is intentionally narrow:

```txt
domain -> Docker container / host port / remote HTTP service
```

Traefik still owns the hard traffic work: ports `80/443`, HTTPS, Let's Encrypt, and reverse proxying. NerdGate Hub manages routes, users, settings, and Traefik dynamic config.

### Features

- One-line installer with DNS preflight.
- Traefik dynamic config generation.
- SQLite storage for routes, users, and settings.
- Admin login with session cookies.
- Docker container picker through Docker Engine API.
- Published host-port targets: `http://host.docker.internal:3000`.
- Internal Docker-network targets through `nerdgate-proxy`.
- Remote targets: `http://203.0.113.10:8080`.
- Password reset and uninstall scripts.
- Static bilingual docs site for GitHub Pages.

### Quick Install

Before installing, create an A record for the panel domain:

```txt
nerdgate.example.com -> SERVER_PUBLIC_IP
```

Then run:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

or:

```sh
wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

The installer checks DNS, creates `.env`, bootstraps the panel route, pulls the published image, and starts Docker Compose.

### Operations

Reset admin password:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh
```

Uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

Update:

```sh
cd /opt/nerdgate-hub
docker compose pull
docker compose up -d
```

### Local Development

Production installs use the published image:

```txt
ghcr.io/moverq1337/nerdgate-hub:latest
```

For local source builds:

```sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Docs-only changes do not rebuild the service image. The Docker image workflow is limited to app files such as `Dockerfile`, `go.mod`, `go.sum`, `cmd/**`, and `internal/**`. Site changes only trigger the GitHub Pages workflow.

### Data Layout

```txt
data/app/nerdgate.db       SQLite: users, settings, routes
data/traefik/routes.yml    Traefik dynamic config
data/acme/acme.json        Let's Encrypt certificates
```

### GitHub Pages

If the Pages workflow fails with `Get Pages site failed`, enable Pages once:

```txt
Settings -> Pages -> Build and deployment -> Source -> GitHub Actions
```

Alternatively, add a repository secret named `PAGES_TOKEN` with Pages/admin rights. The workflow can then enable Pages automatically.

---

## Русский

NerdGate Hub - маленькая open-source панель поверх Traefik. Она не пытается быть большой платформой. Ее задача специально узкая:

```txt
домен -> Docker-контейнер / порт сервера / внешний HTTP-сервис
```

Traefik занимается трафиком: `80/443`, HTTPS, Let's Encrypt и reverse proxy. NerdGate Hub управляет маршрутами, пользователями, настройками и dynamic config для Traefik.

### Возможности

- Установка одной командой с DNS-проверкой.
- Генерация Traefik dynamic config.
- SQLite для routes, users и settings.
- Login-панель с cookie-сессиями.
- Выбор Docker-контейнера через Docker Engine API.
- Маршруты на опубликованные host ports: `http://host.docker.internal:3000`.
- Маршруты на контейнеры внутри Docker-сети `nerdgate-proxy`.
- Маршруты на другой сервер: `http://203.0.113.10:8080`.
- Скрипты сброса пароля и удаления.
- Двуязычная документация на GitHub Pages.

### Быстрый старт

Перед установкой создай A-запись для домена панели:

```txt
nerdgate.example.com -> SERVER_PUBLIC_IP
```

Запусти:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

или:

```sh
wget -qO- https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

Инсталлер проверит DNS, создаст `.env`, добавит стартовый route панели, скачает опубликованный Docker image и запустит Docker Compose.

### Операции

Сбросить admin password:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh
```

Удалить NerdGate Hub:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

Обновить:

```sh
cd /opt/nerdgate-hub
docker compose pull
docker compose up -d
```

### Разработка

Production-установка использует готовый image:

```txt
ghcr.io/moverq1337/nerdgate-hub:latest
```

Локальный build из исходников:

```sh
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Изменения в документации не пересобирают Docker image. Workflow сборки image ограничен файлами приложения: `Dockerfile`, `go.mod`, `go.sum`, `cmd/**`, `internal/**`. Изменения в `site/**` запускают только GitHub Pages.

### Где лежат данные

```txt
data/app/nerdgate.db       SQLite: users, settings, routes
data/traefik/routes.yml    Traefik dynamic config
data/acme/acme.json        Let's Encrypt certificates
```

### Roadmap

- first-run setup tokens;
- route editing;
- target health checks;
- Cloudflare DNS automation;
- attach-container-to-network action;
- tagged releases like `v0.1.0`;
- multi-server agent mode.
