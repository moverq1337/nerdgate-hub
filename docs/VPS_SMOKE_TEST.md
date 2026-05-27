# VPS Smoke Test

Use this checklist before tagging a public release.

## Prerequisites

- Fresh VPS with Docker and Docker Compose plugin.
- A real domain with DNS access.
- Public TCP ports `80` and `443` open to the VPS.
- A test backend container or service.

## Flow

1. Create DNS record:

```txt
nerdgate.example.com -> VPS_PUBLIC_IP
```

2. Run installer:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/install.sh | sh
```

3. Open the panel:

```txt
https://nerdgate.example.com
```

4. Paste the setup token and create the first admin account.

5. Create a route to a public or local test backend:

```txt
app.example.com -> http://host.docker.internal:3000
```

6. Verify:

- `https://app.example.com` opens.
- The route health chip is `Up`.
- Diagnostics has no certificate errors.
- Traefik issued a certificate.

7. Change admin password from `Maintenance -> Account password`.

8. Download backup from `Diagnostics -> Download backup`.

9. Restore backup from `Maintenance -> Restore backup`.

10. Wait for NerdGate Hub to restart, sign in again, and verify routes still exist.

11. Run diagnostics:

```sh
cd /opt/nerdgate-hub
docker compose run --rm --no-deps nerdgate-hub check-domain nerdgate.example.com
docker compose logs --tail=160 traefik
```

12. Test uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/moverq1337/nerdgate-hub/main/scripts/uninstall.sh | sh
```

## Pass Criteria

- Installer stops on wrong DNS and succeeds on correct DNS.
- Setup token works once.
- HTTPS works for the panel and a test route.
- Password change works in the panel.
- Backup downloads from the panel.
- Restore upload stages successfully and survives restart.
- Diagnostics explain target, Docker, Traefik, and certificate problems clearly enough to act.
