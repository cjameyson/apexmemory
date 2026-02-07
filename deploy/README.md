# Apex Memory -- Production Deployment

Mac Mini running Go API + SvelteKit + Caddy + PostgreSQL (Docker), accessed via Tailscale.

```
Browser --> Tailscale --> Caddy (:80) --> SvelteKit (:3000) --> Go API (:4000) --> PostgreSQL (:5432)
```

Caddy handles gzip/zstd compression and reverse-proxies to the SvelteKit frontend. No HTTPS -- Tailscale provides WireGuard encryption at the network layer. The SvelteKit BFF proxies API calls to the Go backend.

## Directory Layout

```
/opt/apexmemory/              # Git repo (code + deploy scripts)
  .env                        # Environment variables (gitignored, create from .env.example)
/opt/apexmemory-data/
  assets/                     # User-uploaded files
  backups/                    # Database backups (gzipped SQL)
  logs/                       # Service logs (stdout + stderr per service)
```

## File Descriptions

| File | Purpose |
|------|---------|
| `Caddyfile` | Caddy reverse proxy config -- listens on `:80`, proxies to `localhost:3000`. |
| `deploy.sh` | Full deploy: git pull, migrate, build Go + SvelteKit, restart services, health check. |
| `run-api.sh` | Launchd wrapper -- sources env vars, exec's the Go binary on port 4000. |
| `run-frontend.sh` | Launchd wrapper -- sources env vars, exec's the Node.js build on port 3000. |
| `backup.sh` | Dumps `app` + `app_code` schemas via `pg_dump`, gzips output, prunes backups older than 7 days. |
| `launchd/com.apexmemory.api.plist` | Launchd agent for the Go API. KeepAlive + RunAtLoad. |
| `launchd/com.apexmemory.frontend.plist` | Launchd agent for SvelteKit. KeepAlive + RunAtLoad. |
| `launchd/com.apexmemory.caddy.plist` | Launchd agent for Caddy. KeepAlive + RunAtLoad. |
| `launchd/com.apexmemory.backup.plist` | Launchd agent for backups. Runs every 6 hours (21600s). |

## First-Time Setup

### Prerequisites

Install on the Mac Mini:

- Go 1.25+
- Node.js (LTS)
- Docker (for PostgreSQL)
- Caddy (`brew install caddy`)
- tern (`go install github.com/jackc/tern/v2@latest`)
- Tailscale

### Steps

```bash
# 1. Clone the repo
sudo mkdir -p /opt/apexmemory /opt/apexmemory-data/{assets,backups,logs}
sudo chown -R $(whoami) /opt/apexmemory /opt/apexmemory-data
git clone <repo-url> /opt/apexmemory

# 2. Configure environment
cp /opt/apexmemory/.env.example /opt/apexmemory/.env
# Edit .env -- replace all CHANGE_ME values with real passwords
# Generate passwords with: openssl rand -base64 24

# 3. Start PostgreSQL
cd /opt/apexmemory && docker compose up -d

# 4. Run migrations
set -a && source /opt/apexmemory/.env && set +a
TERN_CONFIG=/opt/apexmemory/backend/tern.conf \
TERN_MIGRATIONS=/opt/apexmemory/backend/db/migrations \
  tern migrate

# 5. Initial build
cd /opt/apexmemory/backend && go build -o bin/api ./cmd/api
cd /opt/apexmemory/frontend && npm ci --production=false && npm run build

# 6. Install launchd plists
cp /opt/apexmemory/deploy/launchd/*.plist ~/Library/LaunchAgents/
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.apexmemory.api.plist
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.apexmemory.frontend.plist
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.apexmemory.caddy.plist
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.apexmemory.backup.plist
```

## Deploying

From your laptop:

```bash
ssh mac-mini "cd /opt/apexmemory && bash deploy/deploy.sh"
```

Or SSH in and run directly:

```bash
cd /opt/apexmemory && bash deploy/deploy.sh
```

The script runs pre-flight checks (`.env` exists, PostgreSQL is running), then: git pull, migrate, build API, build frontend, restart services, health check.

## Service Management

All services run as launchd user agents. Replace `SERVICE` with one of: `api`, `frontend`, `caddy`, `backup`.

```bash
# Restart a service
launchctl kickstart -k "gui/$(id -u)/com.apexmemory.SERVICE"

# Stop a service
launchctl bootout "gui/$(id -u)/com.apexmemory.SERVICE"

# Start a service
launchctl bootstrap "gui/$(id -u)" ~/Library/LaunchAgents/com.apexmemory.SERVICE.plist

# Check if running
launchctl print "gui/$(id -u)/com.apexmemory.SERVICE"

# List all apex services
launchctl list | grep apexmemory
```

## Logs

All logs are in `/opt/apexmemory-data/logs/`:

| Log file | Source |
|----------|--------|
| `api-stdout.log` / `api-stderr.log` | Go API |
| `frontend-stdout.log` / `frontend-stderr.log` | SvelteKit |
| `caddy-stdout.log` / `caddy-stderr.log` | Caddy |
| `backup-stdout.log` / `backup-stderr.log` | Backup script |

```bash
# Tail API logs
tail -f /opt/apexmemory-data/logs/api-stdout.log

# Check recent backup output
tail -20 /opt/apexmemory-data/logs/backup-stdout.log
```

## Backups

- **Schedule:** Every 6 hours via launchd
- **Location:** `/opt/apexmemory-data/backups/backup-YYYYMMDD-HHMMSS.sql.gz`
- **Scope:** `app` + `app_code` schemas
- **Retention:** 7 days (older files auto-pruned)
- **Validation:** Fails if output is < 500 bytes

```bash
# Manual backup
bash /opt/apexmemory/deploy/backup.sh

# List backups
ls -lh /opt/apexmemory-data/backups/

# Restore (example)
gunzip -c /opt/apexmemory-data/backups/backup-20260207-120000.sql.gz \
  | docker exec -i apexmemory-pg18 psql -U apex_migrator -d apexmemory
```

## Troubleshooting

```bash
# Health check endpoints
curl -s http://localhost:4000/v1/healthcheck    # Go API
curl -s http://localhost:3000                    # SvelteKit

# Is PostgreSQL running?
docker ps | grep apexmemory-pg18

# Restart everything
launchctl kickstart -k "gui/$(id -u)/com.apexmemory.api"
launchctl kickstart -k "gui/$(id -u)/com.apexmemory.frontend"
launchctl kickstart -k "gui/$(id -u)/com.apexmemory.caddy"

# Check for port conflicts
lsof -i :80 -i :3000 -i :4000 -i :5432
```
