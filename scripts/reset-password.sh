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

read_secret_tty() {
  prompt="$1"
  printf '%s: ' "$prompt" > /dev/tty
  stty -echo < /dev/tty
  IFS= read -r value < /dev/tty
  stty echo < /dev/tty
  printf '\n' > /dev/tty
  printf '%s' "$value"
}

random_hex() {
  bytes="$1"
  od -An -N"$bytes" -tx1 /dev/urandom | tr -d ' \n'
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

env_value() {
  file="$1"
  key="$2"
  awk -F= -v key="$key" '$1 == key { print substr($0, length(key) + 2); exit }' "$file"
}

set_env_value() {
  file="$1"
  key="$2"
  value="$3"
  tmp="$file.tmp"
  awk -v key="$key" -v value="$value" '
    BEGIN { found = 0 }
    $0 ~ "^" key "=" {
      print key "=" value
      found = 1
      next
    }
    { print }
    END {
      if (!found) {
        print key "=" value
      }
    }
  ' "$file" > "$tmp"
  mv "$tmp" "$file"
}

delete_env_key() {
  file="$1"
  key="$2"
  tmp="$file.tmp"
  awk -v key="$key" '$0 !~ "^" key "=" { print }' "$file" > "$tmp"
  mv "$tmp" "$file"
}

main() {
  [ -r /dev/tty ] || die "interactive terminal is required"
  command -v docker >/dev/null 2>&1 || die "docker is required"
  docker compose version >/dev/null 2>&1 || die "Docker Compose plugin is required"

  if [ "$(id -u)" -eq 0 ]; then
    default_install_dir="/opt/nerdgate-hub"
  else
    default_install_dir="$HOME/nerdgate-hub"
  fi

  say "NerdGate Hub password reset"
  say ""

  install_dir="$(read_tty "Install directory" "$default_install_dir")"
  env_file="$install_dir/.env"
  [ -f "$env_file" ] || die "$env_file not found"
  [ -f "$install_dir/docker-compose.yml" ] || die "$install_dir/docker-compose.yml not found"

  current_user="$(env_value "$env_file" "NERDGATE_USERNAME")"
  username="$(read_tty "Admin username" "${current_user:-admin}")"
  validate_safe_env_value "username" "$username"

  password="$(read_secret_tty "New admin password (empty = generate)")"
  if [ -z "$password" ]; then
    password="$(random_hex 16)"
    say "Generated admin password: $password"
  fi
  validate_safe_env_value "password" "$password"

  say "Updating SQLite credentials..."
  (
    cd "$install_dir"
    docker compose run --rm --no-deps nerdgate-hub reset-password /data "$username" "$password"
  )

  session_secret="$(random_hex 32)"
  set_env_value "$env_file" "NERDGATE_USERNAME" "$username"
  set_env_value "$env_file" "NERDGATE_SESSION_SECRET" "$session_secret"
  delete_env_key "$env_file" "NERDGATE_PASSWORD"

  say "Recreating NerdGate Hub container..."
  (
    cd "$install_dir"
    docker compose up -d --force-recreate nerdgate-hub
  )

  say ""
  say "Password updated."
  say "Username: $username"
  say "Password: $password"
}

main "$@"
