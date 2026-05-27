<p align="center">
  <img src="../site/logo.png" alt="NerdGate Hub" width="128">
</p>

<h1 align="center">NerdGate Hub</h1>

<p align="center">
  Маленькая open-source панель поверх Traefik для маршрутизации доменов на Docker-контейнеры, порты сервера и внешние HTTP-сервисы.
</p>

<p align="center">
  <a href="../README.md">English README</a>
  ·
  <a href="https://moverq1337.github.io/nerdgate-hub/ru/get-started.html">Документация</a>
  ·
  <a href="https://github.com/moverq1337/nerdgate-hub/pkgs/container/nerdgate-hub">GHCR Image</a>
</p>

---

## Идея

NerdGate Hub не пытается быть большой платформой. Его задача специально узкая:

```txt
домен -> Docker-контейнер / порт сервера / внешний HTTP-сервис
```

Traefik занимается трафиком: `80/443`, HTTPS, Let's Encrypt и reverse proxy. NerdGate Hub управляет маршрутами, пользователями, настройками и dynamic config для Traefik.

## Возможности

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
- CSRF protection, rate limit для login/setup, строгие cookies, security headers и audit log.
- Backup и restore SQLite-данных вместе с Let's Encrypt `acme.json`.
- Скрипты сброса пароля, shell-диагностики и удаления.
- Tagged GitHub releases вроде `v0.1.0`.
- Двуязычная документация на GitHub Pages.

## Быстрый старт

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

## Операции

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

Создать backup:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/backup.sh | sh
```

Восстановиться из backup:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/restore.sh | sh
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

## Разработка

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

Изменения в документации не пересобирают Docker image. Workflow сборки image ограничен файлами приложения: `Dockerfile`, `go.mod`, `go.sum`, `cmd/**`, и `internal/**`. Изменения в `site/**` запускают только GitHub Pages.

## Где лежат данные

```txt
data/app/nerdgate.db       SQLite: users, settings, routes
data/traefik/routes.yml    Traefik dynamic config
data/acme/acme.json        Let's Encrypt certificates
```

Backup-архив содержит `nerdgate.db`, optional legacy `routes.json`, optional `acme.json` и metadata.

## Roadmap

- Cloudflare DNS automation;
- multi-server agent mode;
- optional tunnel mode for servers without public 80/443.
