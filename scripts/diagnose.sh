#!/bin/sh
set -eu

say() {
  printf '%s\n' "$*"
}

die() {
  say "Error: $*" >&2
  exit 1
}

read_tty() {
  prompt="$1"
  default="$2"
  printf '%s' "$prompt" > /dev/tty
  if [ -n "$default" ]; then
    printf ' [%s]' "$default" > /dev/tty
  fi
  printf ': ' > /dev/tty
  IFS= read -r value < /dev/tty
  if [ -z "$value" ]; then
    value="$default"
  fi
  printf '%s' "$value"
}

env_value() {
  file="$1"
  key="$2"
  awk -F= -v key="$key" '$1 == key { print substr($0, length(key) + 2); exit }' "$file"
}

bootstrap_domain() {
  file="$1"
  [ -f "$file" ] || return 0
  sed -n 's/.*"domain"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$file" | head -n 1
}

hint() {
  found_hints=1
  say "  - $*"
}

main() {
  command -v docker >/dev/null 2>&1 || die "docker is required"
  docker compose version >/dev/null 2>&1 || die "Docker Compose plugin is required"

  if [ "$(id -u)" -eq 0 ]; then
    default_install_dir="/opt/nerdgate-hub"
  else
    default_install_dir="$HOME/nerdgate-hub"
  fi

  install_dir="${1:-}"
  if [ -z "$install_dir" ]; then
    if [ -r /dev/tty ]; then
      install_dir="$(read_tty "Install directory" "$default_install_dir")"
    else
      install_dir="$default_install_dir"
    fi
  fi

  [ -f "$install_dir/docker-compose.yml" ] || die "$install_dir/docker-compose.yml not found"
  [ -f "$install_dir/.env" ] || die "$install_dir/.env not found"

  say "NerdGate Hub diagnostics"
  say ""
  say "Install directory: $install_dir"
  say "ACME email: $(env_value "$install_dir/.env" "ACME_EMAIL")"
  say ""

  panel_domain="$(bootstrap_domain "$install_dir/data/app/routes.json" || true)"
  if [ -z "$panel_domain" ] && [ -r /dev/tty ]; then
    panel_domain="$(read_tty "Panel domain for DNS check (empty = skip)" "")"
  fi

  say "Compose services:"
  (
    cd "$install_dir"
    docker compose ps
  ) || true
  say ""

  if [ -n "$panel_domain" ]; then
    say "DNS check for $panel_domain:"
    (
      cd "$install_dir"
      docker compose run --rm --no-deps nerdgate-hub check-domain "$panel_domain"
    ) || true
    say ""
  fi

  logs="$(
    cd "$install_dir"
    docker compose logs --no-color --tail=250 traefik 2>&1 || true
  )"

  say "Recent Traefik logs:"
  say "$logs"
  say ""

  found_hints=0
  say "Detected hints:"

  case "$logs" in
    *"unable to obtain ACME certificate"*|*"error getting certificate"*|*"acme: error"*)
      hint "Let's Encrypt could not issue a certificate. Check that DNS points to this server and ports 80/443 are reachable from the public internet."
      ;;
  esac

  case "$logs" in
    *"connection refused"*|*"timeout"*|*"Timeout during connect"*)
      hint "HTTP-01 validation likely cannot reach Traefik. Open TCP 80 and 443 in firewall/security group and make sure no CDN proxy blocks validation."
      ;;
  esac

  case "$logs" in
    *"rateLimited"*|*"too many certificates"*|*"too many failed authorizations"*)
      hint "Let's Encrypt rate limit was hit. Fix DNS/ports first, then wait for the limit window before retrying."
      ;;
  esac

  case "$logs" in
    *"permission denied"*|*"acme.json"*)
      hint "ACME storage may have wrong permissions. Run: chmod 600 $install_dir/data/acme/acme.json"
      ;;
  esac

  case "$logs" in
    *"bind: address already in use"*|*"address already in use"*)
      hint "Port 80 or 443 is already occupied. Stop nginx/apache/another proxy or move it away from public 80/443."
      ;;
  esac

  if [ "$found_hints" -eq 0 ]; then
    say "  - No known Traefik certificate pattern detected in the last 250 lines."
  fi

  say ""
  say "Manual checks:"
  say "  - DNS: A/AAAA for the domain must point to this server public IP."
  say "  - Ports: TCP 80 and 443 must be open from the internet."
  say "  - Cloudflare: use DNS-only mode while issuing HTTP-01 certificates."
  say "  - Email: ACME_EMAIL in .env must be a real email address."
  say "  - Live logs: cd $install_dir && docker compose logs -f traefik"
}

main "$@"
