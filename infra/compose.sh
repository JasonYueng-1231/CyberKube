#!/usr/bin/env bash
set -euo pipefail

# Determine repo root and compose file
REPO_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$REPO_DIR/infra/docker-compose.yml"

echo "[compose] using file: $COMPOSE_FILE"

# Prefer containerized compose to avoid local Python/plugin issues

# If user explicitly asks to use local tool
if [[ "${FORCE_LOCAL_COMPOSE:-}" == "1" ]]; then
  echo "[compose] FORCE_LOCAL_COMPOSE=1, trying local tools"
  if docker compose version 2>&1 | grep -qi "Docker Compose version"; then
    exec docker compose -f "$COMPOSE_FILE" "$@"
  elif command -v docker-compose >/dev/null 2>&1; then
    exec docker-compose -f "$COMPOSE_FILE" "$@"
  else
    echo "[compose] no local compose available" >&2
  fi
fi

# Try plugin only if it looks like real compose (not plain Docker version)
if docker compose version 2>&1 | grep -qi "Docker Compose version"; then
  echo "[compose] detected docker compose plugin"
  exec docker compose -f "$COMPOSE_FILE" "$@"
fi

# Containerized compose (robust on old Docker)
echo "[compose] using containerized docker/compose:1.29.2"
docker image inspect docker/compose:1.29.2 >/dev/null 2>&1 || docker pull -q docker/compose:1.29.2
exec docker run --rm -i \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$REPO_DIR":"$REPO_DIR" -w "$REPO_DIR" \
  docker/compose:1.29.2 -f infra/docker-compose.yml "$@"

