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
- Browser first-run setup with a one-time setup token.
- Traefik dynamic config generation.
- SQLite storage for routes, users, and settings.
- Admin login with session cookies.
- Multiple domains can be created in one route form submission.
- Inline route editing.
- Docker container picker through Docker Engine API.
- Published host-port targets: `http://host.docker.internal:3000`.
- Internal Docker-network targets through `nerdgate-proxy`.
- Attach containers to `nerdgate-proxy` from the panel.
- Remote targets: `http://203.0.113.10:8080`.
- Target health checks and in-panel diagnostics.
- Password reset, shell diagnostics, and uninstall scripts.
- Tagged GitHub releases like `v0.1.0`.
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

The installer asks for the panel domain and Let's Encrypt email, checks DNS, creates `.env`, bootstraps the panel route, pulls the published image, and starts Docker Compose. It prints a one-time setup token. Open the panel URL, paste that token, and create the first admin account in the browser.

### Operations

Reset admin password:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh
```

Uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

Diagnose HTTPS, Traefik, DNS, and certificate problems:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh
```

Update:

```sh
cd /opt/nerdgate-hub
docker compose pull
docker compose up -d
```

Create a tagged release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The release workflow creates a GitHub Release, and the Docker workflow publishes the matching GHCR tag.

### Local Development

Production installs use the published image:

```txt
ghcr.io/moverq1337/nerdgate-hub:latest
```

For local source builds:

```sh
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

For anything reachable from the internet, replace the example setup token before starting.

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
- Первый запуск через одноразовый setup token в браузере.
- Генерация Traefik dynamic config.
- SQLite для routes, users и settings.
- Login-панель с cookie-сессиями.
- Несколько доменов можно создать одной отправкой формы.
- Inline-редактирование routes.
- Выбор Docker-контейнера через Docker Engine API.
- Маршруты на опубликованные host ports: `http://host.docker.internal:3000`.
- Маршруты на контейнеры внутри Docker-сети `nerdgate-proxy`.
- Подключение контейнера к `nerdgate-proxy` из панели.
- Маршруты на другой сервер: `http://203.0.113.10:8080`.
- Health checks для targets и diagnostics прямо в панели.
- Скрипты сброса пароля, shell-диагностики и удаления.
- Tagged GitHub releases вроде `v0.1.0`.
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

Инсталлер спросит домен панели и email для Let's Encrypt, проверит DNS, создаст `.env`, добавит стартовый route панели, скачает опубликованный Docker image и запустит Docker Compose. В конце он покажет одноразовый setup token. Открой URL панели, вставь token и создай первого admin-пользователя в браузере.

### Операции

Сбросить admin password:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/reset-password.sh | sh
```

Удалить NerdGate Hub:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

Диагностика HTTPS, Traefik, DNS и сертификатов:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/diagnose.sh | sh
```

Обновить:

```sh
cd /opt/nerdgate-hub
docker compose pull
docker compose up -d
```

Создать tagged release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

Release workflow создаст GitHub Release, а Docker workflow опубликует такой же GHCR tag.

### Разработка

Production-установка использует готовый image:

```txt
ghcr.io/moverq1337/nerdgate-hub:latest
```

Локальный build из исходников:

```sh
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Если окружение доступно из интернета, замени пример setup token перед запуском.

Изменения в документации не пересобирают Docker image. Workflow сборки image ограничен файлами приложения: `Dockerfile`, `go.mod`, `go.sum`, `cmd/**`, `internal/**`. Изменения в `site/**` запускают только GitHub Pages.

### Где лежат данные

```txt
data/app/nerdgate.db       SQLite: users, settings, routes
data/traefik/routes.yml    Traefik dynamic config
data/acme/acme.json        Let's Encrypt certificates
```

### Roadmap

- Cloudflare DNS automation;
- multi-server agent mode;
- optional tunnel mode for servers without public 80/443.
