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

confirm() {
  prompt="$1"
  answer="$(read_tty "$prompt" "no")"
  case "$answer" in
    y|Y|yes|YES) return 0 ;;
    *) return 1 ;;
  esac
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

  say "NerdGate Hub restore"
  say ""

  install_dir="$(read_tty "Install directory" "$default_install_dir")"
  [ -f "$install_dir/docker-compose.yml" ] || die "$install_dir/docker-compose.yml not found"

  backup_file="$(read_tty "Backup zip" "")"
  [ -f "$backup_file" ] || die "$backup_file not found"

  say ""
  say "This will stop the Compose stack and replace SQLite data plus acme.json from the backup."
  if ! confirm "Continue"; then
    say "Aborted."
    exit 0
  fi

  (
    cd "$install_dir"
    docker compose down --remove-orphans
    docker compose run --rm --no-deps -v "$backup_file:/backup.zip:ro" nerdgate-hub restore /data /acme/acme.json /backup.zip
    docker compose up -d
  )

  say ""
  say "Restore complete."
}

main "$@"
