#!/bin/sh
set -eu

REPO_TARBALL_URL="${NERDGATE_REPO_TARBALL_URL:-https://github.com/moverq1337/nerdgate-hub/archive/refs/heads/main.tar.gz}"

say() {
  printf '%s\n' "$*"
}

die() {
  say "Error: $*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "$1 is required"
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

read_secret_tty() {
  prompt="$1"
  printf '%s: ' "$prompt" > /dev/tty
  stty -echo < /dev/tty
  IFS= read -r value < /dev/tty
  stty echo < /dev/tty
  printf '\n' > /dev/tty
  printf '%s' "$value"
}

download_file() {
  url="$1"
  output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
    return
  fi
  die "curl or wget is required"
}

fetch_text() {
  url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url"
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO- "$url"
    return
  fi
  die "curl or wget is required"
}

random_hex() {
  bytes="$1"
  od -An -N"$bytes" -tx1 /dev/urandom | tr -d ' \n'
}

detect_public_ip() {
  fetch_text "https://api.ipify.org" 2>/dev/null || fetch_text "https://icanhazip.com" 2>/dev/null
}

resolve_domain_ips() {
  domain="$1"
  if command -v getent >/dev/null 2>&1; then
    getent ahosts "$domain" | awk '{ print $1 }' | sort -u
    return
  fi
  if command -v dig >/dev/null 2>&1; then
    {
      dig +short A "$domain"
      dig +short AAAA "$domain"
    } | awk 'NF { print $1 }' | sort -u
    return
  fi
  die "getent or dig is required for DNS checks"
}

validate_domain() {
  domain="$1"
  case "$domain" in
    ""|http://*|https://*|*/*|*:*|.*|*.)
      die "enter a bare domain like nerdgate.example.com"
      ;;
  esac
  case "$domain" in
    *.*) ;;
    *) die "domain must include a dot, for example nerdgate.example.com" ;;
  esac
}

validate_safe_env_value() {
  name="$1"
  value="$2"
  case "$value" in
    ""|*[!A-Za-z0-9_.:@%+=,-]*)
      die "$name contains unsupported characters; use letters, numbers, and ._:@%+=,-"
      ;;
  esac
}

dns_matches() {
  expected="$1"
  ips="$2"
  for ip in $ips; do
    if [ "$ip" = "$expected" ]; then
      return 0
    fi
  done
  return 1
}

write_env() {
  install_dir="$1"
  email="$2"
  username="$3"
  password="$4"
  session_secret="$5"

  cat > "$install_dir/.env" <<EOF_ENV
ACME_EMAIL=$email
NERDGATE_USERNAME=$username
NERDGATE_PASSWORD=$password
NERDGATE_SESSION_SECRET=$session_secret
EOF_ENV
}

write_bootstrap_route() {
  install_dir="$1"
  panel_domain="$2"
  mkdir -p "$install_dir/data/app" "$install_dir/data/traefik" "$install_dir/data/acme"

  if [ -f "$install_dir/data/app/routes.json" ]; then
    return
  fi

  route_id="route-$(date +%s)"
  now="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  cat > "$install_dir/data/app/routes.json" <<EOF_JSON
[
  {
    "id": "$route_id",
    "domain": "$panel_domain",
    "target_url": "http://nerdgate-hub:8080",
    "tls": true,
    "created_at": "$now",
    "updated_at": "$now"
  }
]
EOF_JSON

  touch "$install_dir/data/acme/acme.json"
  chmod 600 "$install_dir/data/acme/acme.json"
}

download_project() {
  install_dir="$1"

  if [ -f "$install_dir/.env" ]; then
    die "$install_dir already contains .env; refusing to overwrite an existing install"
  fi

  tmp_dir="$(mktemp -d)"
  archive="$tmp_dir/nerdgate-hub.tar.gz"
  download_file "$REPO_TARBALL_URL" "$archive"
  tar -xzf "$archive" -C "$tmp_dir"
  src_dir="$(find "$tmp_dir" -mindepth 1 -maxdepth 1 -type d | head -n 1)"

  mkdir -p "$install_dir"
  cp -R "$src_dir/." "$install_dir/"
  rm -rf "$tmp_dir"
}

main() {
  [ -r /dev/tty ] || die "interactive terminal is required"

  need_cmd docker
  need_cmd tar
  need_cmd awk
  need_cmd od
  docker compose version >/dev/null 2>&1 || die "Docker Compose plugin is required"

  say "NerdGate Hub installer"
  say ""

  if [ "$(id -u)" -eq 0 ]; then
    default_install_dir="/opt/nerdgate-hub"
  else
    default_install_dir="$HOME/nerdgate-hub"
  fi

  install_dir="$(read_tty "Install directory" "$default_install_dir")"
  panel_domain="$(read_tty "Panel domain" "${NERDGATE_PANEL_DOMAIN:-}")"
  validate_domain "$panel_domain"

  email="$(read_tty "Let's Encrypt email" "${ACME_EMAIL:-}")"
  validate_safe_env_value "email" "$email"

  username="$(read_tty "Admin username" "${NERDGATE_USERNAME:-admin}")"
  validate_safe_env_value "username" "$username"

  password="$(read_secret_tty "Admin password (empty = generate)")"
  if [ -z "$password" ]; then
    password="$(random_hex 16)"
    say "Generated admin password: $password"
  fi
  validate_safe_env_value "password" "$password"

  say ""
  say "Checking public IP..."
  public_ip="$(detect_public_ip | tr -d ' \n\r')"
  [ -n "$public_ip" ] || die "could not detect public IP"
  say "Server public IP: $public_ip"

  say "Checking DNS for $panel_domain..."
  resolved_ips="$(resolve_domain_ips "$panel_domain" | tr '\n' ' ')"
  [ -n "$resolved_ips" ] || die "$panel_domain has no A/AAAA records"
  say "Resolved IPs: $resolved_ips"

  if ! dns_matches "$public_ip" "$resolved_ips"; then
    die "$panel_domain must point to $public_ip before installation"
  fi

  say ""
  say "DNS check passed."
  say "Downloading NerdGate Hub..."
  download_project "$install_dir"

  session_secret="$(random_hex 32)"
  write_env "$install_dir" "$email" "$username" "$password" "$session_secret"
  write_bootstrap_route "$install_dir" "$panel_domain"

  say "Starting containers..."
  (
    cd "$install_dir"
    docker compose up -d --build
  )

  say ""
  say "NerdGate Hub is starting."
  say "URL: https://$panel_domain"
  say "Username: $username"
  say "Password: $password"
}

main "$@"

