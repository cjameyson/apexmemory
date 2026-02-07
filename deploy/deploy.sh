#!/usr/bin/env bash
set -euo pipefail

# Apex Memory - Deploy Script
# Usage: cd /opt/apexmemory && bash deploy/deploy.sh

PROJECT_DIR="/opt/apexmemory"
ENV_FILE="${PROJECT_DIR}/.env"
LOG_DIR="/opt/apexmemory-data/logs"
GUI_DOMAIN="gui/$(id -u)"

cd "${PROJECT_DIR}"

echo "=== Apex Memory Deploy ==="
echo "Started at $(date)"
echo ""

# 0. Pre-flight checks
if [ ! -f "${ENV_FILE}" ]; then
  echo "ERROR: ${ENV_FILE} not found. Create it from .env.example." >&2
  exit 1
fi

if ! docker ps --format '{{.Names}}' | grep -q '^apexmemory-pg18$'; then
  echo "ERROR: PostgreSQL container 'apexmemory-pg18' is not running." >&2
  echo "  Start it with: cd ${PROJECT_DIR} && docker compose up -d" >&2
  exit 1
fi

# 1. Pull latest code
echo "[1/6] Pulling latest code..."
git pull --ff-only
echo ""

# 2. Run database migrations
echo "[2/6] Running database migrations..."
set -a
source "${ENV_FILE}"
set +a
TERN_CONFIG="${PROJECT_DIR}/backend/tern.conf" \
TERN_MIGRATIONS="${PROJECT_DIR}/backend/db/migrations" \
  tern migrate
echo ""

# 3. Build Go API
echo "[3/6] Building Go API..."
cd "${PROJECT_DIR}/backend"
go build -o bin/api ./cmd/api
cd "${PROJECT_DIR}"
echo ""

# 4. Build SvelteKit frontend
echo "[4/6] Building SvelteKit frontend..."
cd "${PROJECT_DIR}/frontend"
npm ci --production=false
npm run build
cd "${PROJECT_DIR}"
echo ""

# 5. Restart services
echo "[5/6] Restarting services..."
launchctl kickstart -k "${GUI_DOMAIN}/com.apexmemory.api" 2>/dev/null || \
  { launchctl bootout "${GUI_DOMAIN}/com.apexmemory.api" 2>/dev/null; \
    launchctl bootstrap "${GUI_DOMAIN}" ~/Library/LaunchAgents/com.apexmemory.api.plist; }
launchctl kickstart -k "${GUI_DOMAIN}/com.apexmemory.frontend" 2>/dev/null || \
  { launchctl bootout "${GUI_DOMAIN}/com.apexmemory.frontend" 2>/dev/null; \
    launchctl bootstrap "${GUI_DOMAIN}" ~/Library/LaunchAgents/com.apexmemory.frontend.plist; }
echo ""

# 6. Health check
echo "[6/6] Health checks..."
sleep 3

API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:4000/v1/healthcheck || echo "000")
if [ "${API_STATUS}" = "200" ]; then
  echo "  API:      OK (HTTP ${API_STATUS})"
else
  echo "  API:      FAILED (HTTP ${API_STATUS})" >&2
  echo "  Check logs: ${LOG_DIR}/api-stdout.log" >&2
fi

FRONTEND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 || echo "000")
if [ "${FRONTEND_STATUS}" = "200" ]; then
  echo "  Frontend: OK (HTTP ${FRONTEND_STATUS})"
else
  echo "  Frontend: FAILED (HTTP ${FRONTEND_STATUS})" >&2
  echo "  Check logs: ${LOG_DIR}/frontend-stdout.log" >&2
fi

echo ""
echo "Deploy completed at $(date)"
