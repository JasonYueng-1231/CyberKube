#!/usr/bin/env bash
set -euo pipefail

# Determine repo root and compose file
REPO_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$REPO_DIR/infra/docker-compose.yml"

echo "[compose] using file: $COMPOSE_FILE"

# Try docker compose plugin first (Docker 20.10+)
if docker compose version >/dev/null 2>&1; then
  # Some older Dockers print Docker version instead of Compose version; try a harmless command
  if docker compose ls >/dev/null 2>&1; then
    echo "[compose] detected docker compose plugin"
    exec docker compose -f "$COMPOSE_FILE" "$@"
  fi
fi

# Fallback to docker-compose binary if usable
if command -v docker-compose >/dev/null 2>&1; then
  if docker-compose version >/dev/null 2>&1; then
    echo "[compose] using docker-compose binary"
    exec docker-compose -f "$COMPOSE_FILE" "$@"
  fi
fi

# Final fallback: containerized docker/compose (works on old Docker)
echo "[compose] using containerized docker/compose:1.29.2"
exec docker run --rm -i \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$REPO_DIR":"$REPO_DIR" -w "$REPO_DIR" \
  docker/compose:1.29.2 -f infra/docker-compose.yml "$@"

