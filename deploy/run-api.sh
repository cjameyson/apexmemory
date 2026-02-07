#!/usr/bin/env bash
set -euo pipefail

# Apex Memory - API runner for launchd
# Sources production env vars then exec's the Go binary.

set -a
source /opt/apexmemory/.env
set +a

exec /opt/apexmemory/backend/bin/api \
  -port 4000 \
  -env production
