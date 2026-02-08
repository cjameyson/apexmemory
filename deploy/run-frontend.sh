#!/usr/bin/env bash
set -euo pipefail

# Apex Memory - Frontend runner for launchd
# Sets required env vars then exec's the Node.js production build.

# launchd uses minimal PATH; add common node locations
export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"

set -a
source /opt/apexmemory/.env
set +a

export PORT=3000
export HOST=127.0.0.1
export COOKIES_SECURE=false  # Tailscale provides network-layer encryption
export BODY_SIZE_LIMIT=10485760  # 10MB for image uploads

exec node /opt/apexmemory/frontend/build/index.js
