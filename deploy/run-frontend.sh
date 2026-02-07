#!/usr/bin/env bash
set -euo pipefail

# Apex Memory - Frontend runner for launchd
# Sets required env vars then exec's the Node.js production build.

set -a
source /opt/apexmemory/.env
set +a

export PORT=3000
export HOST=127.0.0.1

exec node /opt/apexmemory/frontend/build/index.js
