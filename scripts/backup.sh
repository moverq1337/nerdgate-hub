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

main() {
  [ -r /dev/tty ] || die "interactive terminal is required"
  command -v docker >/dev/null 2>&1 || die "docker is required"
  docker compose version >/dev/null 2>&1 || die "Docker Compose plugin is required"

  if [ "$(id -u)" -eq 0 ]; then
    default_install_dir="/opt/nerdgate-hub"
  else
    default_install_dir="$HOME/nerdgate-hub"
  fi

  say "NerdGate Hub backup"
  say ""

  install_dir="$(read_tty "Install directory" "$default_install_dir")"
  [ -f "$install_dir/docker-compose.yml" ] || die "$install_dir/docker-compose.yml not found"

  default_output="$install_dir/backups/nerdgate-backup-$(date -u +"%Y%m%d-%H%M%S").zip"
  output_path="$(read_tty "Backup file" "$default_output")"
  output_dir="$(dirname "$output_path")"
  output_name="$(basename "$output_path")"
  mkdir -p "$output_dir"

  (
    cd "$install_dir"
    docker compose run --rm --no-deps -v "$output_dir:/backup" nerdgate-hub backup /data /acme/acme.json "/backup/$output_name"
  )

  say ""
  say "Backup written to $output_path"
}

main "$@"
