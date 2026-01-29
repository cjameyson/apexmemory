# Paths - hardcoded relative paths for safety (never use env vars for rm -rf targets)
BACKEND := backend
FRONTEND := frontend
BACKEND_BIN := $(BACKEND)/bin
FRONTEND_BIN := $(FRONTEND)/bin
SCHEMA_SNAPSHOT := $(BACKEND)/db/schema-dump.sql
BACKUPS_DIR := $(BACKEND)/db/backups
DATE := $(shell date +%Y%m%d-%H%M%S)
DB_DUMP_FILE := $(BACKUPS_DIR)/backup-$(DATE).sql

# Safety check: abort if we're not in the project root (must have CLAUDE.md)
PROJECT_ROOT_CHECK := $(shell test -f CLAUDE.md && echo "ok")
ifneq ($(PROJECT_ROOT_CHECK),ok)
$(error "Safety check failed: CLAUDE.md not found. Are you in the project root?")
endif

# Safe rm helper: validates path is non-empty and relative before deletion
# Usage: $(call safe_rm,path)
define safe_rm
	@if [ -z "$(1)" ]; then \
		echo "❌ SAFETY: Refusing to rm empty path"; \
		exit 1; \
	elif echo "$(1)" | grep -qE '^(/|\.\./)'; then \
		echo "❌ SAFETY: Refusing to rm absolute or parent path: $(1)"; \
		exit 1; \
	elif [ -e "$(1)" ]; then \
		echo "   Removing $(1)"; \
		rm -rf "$(1)"; \
	fi
endef

# --- Tern ---
TERN_CONF := $(BACKEND)/tern.conf
TERN_TEST_CONF := $(BACKEND)/tern.test.conf
TERN_CODE_DIR := $(BACKEND)/db/code
# Tern doesn't support a migrations dir in the config file, so we use a separate variable
MIGR_DIR := $(BACKEND)/db/migrations

# Docker container name (should match your docker-compose.yml)
PG_CONTAINER := apexmemory-pg18

# Helper to load .env for each recipe
ENV_SH := set -a; . ./.env; set +a

.PHONY: help
help:
	@echo "Common targets:"
	@echo "  docker.up        - start docker services (PostgreSQL)"
	@echo "  docker.down      - stop docker services (keep volume)"
	@echo "  docker.nuke      - stop & remove volumes (DESTROYS DATA)"
	@echo "  docker.logs      - follow postgres logs"
	@echo ""
	@echo "  build.api        - build Go API binary"
	@echo "  build.frontend   - build SvelteKit for production"
	@echo "  build.clean      - clean build artifacts"
	@echo ""
	@echo "  dev.backend      - restart Go API dev server (background)"
	@echo "  dev.frontend     - restart SvelteKit dev server (background)"
	@echo "  dev.up           - start both API and Frontend"
	@echo "  dev.mobile       - start servers for mobile testing (exposes to network)"
	@echo "  dev.stop         - stop API and Frontend"
	@echo "  dev.check        - quick checks (Go vet, Svelte check)"
	@echo ""
	@echo "  test.unit          - run backend unit tests (no Docker)"
	@echo "  test.backend       - run all backend tests (needs Docker daemon)"
	@echo "  test.backend.v     - run backend tests (verbose)"
	@echo "  test.backend.cover - run backend tests with coverage"
	@echo "  test.frontend      - run frontend unit tests (Vitest)"
	@echo "  test.frontend.e2e  - run frontend e2e tests (Playwright)"
	@echo "  test.all           - run all tests (backend + frontend)"
	@echo ""
	@echo "  lint.frontend      - run ESLint on frontend"
	@echo "  format.frontend    - run Prettier on frontend"
	@echo "  validate.frontend  - run all frontend checks (check + lint + test)"
	@echo ""
	@echo "  tern.new name=...  - create a new tern migration"
	@echo "  tern.migrate       - apply pending migrations"
	@echo "  tern.status        - show migration status"
	@echo "  tern.rollback      - roll back one migration"
	@echo "  tern.reset         - migrate down to 0 then up (dev only)"
	@echo "  db.psql.app        - psql as app role"
	@echo "  db.psql.migrator   - psql as migrator role"
	@echo "  db.psql.super      - psql as superuser"
	@echo "  db.schema          - dump schema snapshot for sqlc"
	@echo "  db.dump            - dump full app schema+data to timestamped SQL"
	@echo "  db.sqlc            - run sqlc generate (uses schema snapshot)"

# ---------- Build ----------
.PHONY: build.api build.frontend build.clean
build.api:
	@echo "🔨 Building Go API binary..."
	@mkdir -p $(BACKEND_BIN)
	@cd $(BACKEND); go build -o bin/api ./cmd/api
	@echo "✅ API binary built: $(BACKEND_BIN)/api"

build.frontend:
	@echo "🔨 Building SvelteKit for production..."
	@cd $(FRONTEND); npm run build
	@echo "✅ Frontend built"

build.clean:
	@echo "🧹 Cleaning build artifacts..."
	$(call safe_rm,$(BACKEND_BIN))
	$(call safe_rm,$(FRONTEND_BIN))
	$(call safe_rm,$(FRONTEND)/.svelte-kit)
	$(call safe_rm,$(FRONTEND)/build)
	@echo "✅ Build artifacts cleaned"

# ---------- Docker ----------
.PHONY: docker.up docker.down docker.nuke docker.logs
docker.up:
	docker compose up -d
	@echo "⏳ Waiting for database to be ready..."
	@sleep 3
	@$(ENV_SH); docker exec $(PG_CONTAINER) pg_isready -U "$$PG_SUPER_USER" -d "$$PG_DATABASE" > /dev/null 2>&1 && echo "✅ Database is ready!" || echo "⚠️  Database might still be starting..."

docker.down:
	docker compose down

docker.nuke:
	@echo "⚠️  This will DESTROY ALL POSTGRES DATA. Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	docker compose down -v

docker.logs:
	docker logs -f $(PG_CONTAINER)

# ---------- Development helpers ----------
.PHONY: dev.backend dev.frontend dev.up dev.stop dev.check

# Ports (overridable): make API_PORT=8080 FRONTEND_PORT=5174 dev.up
API_PORT ?= 4000
FRONTEND_PORT ?= 5173

dev.backend:
	@echo "🚀 (API) Restarting Go dev server on :$(API_PORT)"
	@mkdir -p $(BACKEND_BIN)
	@PID=$$(lsof -ti :$(API_PORT) || true); \
	if [ -n "$$PID" ]; then \
	  echo "🔪 Killing process on port $(API_PORT) (pid $$PID)"; \
	  kill -9 $$PID || true; \
	fi
	@$(ENV_SH); cd $(BACKEND); \
	LOG_FILE=/tmp/apexmemory-api.log; \
	nohup go run ./cmd/api -port $(API_PORT) > $$LOG_FILE 2>&1 & echo $$! > bin/api.pid; \
	sleep 1; \
	if ! kill -0 $$(cat bin/api.pid) > /dev/null 2>&1; then \
		echo "❌ API failed to start. Logs:"; \
		cat $$LOG_FILE; \
		rm -f bin/api.pid; \
		exit 1; \
	fi; \
	echo "✅ API started (pid $$(cat bin/api.pid)) — logs: $$LOG_FILE"

dev.frontend:
	@echo "🚀 (Web) Restarting SvelteKit dev server on :$(FRONTEND_PORT)"
	@mkdir -p $(FRONTEND_BIN)
	@PID=$$(lsof -ti :$(FRONTEND_PORT) || true); \
	if [ -n "$$PID" ]; then \
	  echo "🔪 Killing process on port $(FRONTEND_PORT) (pid $$PID)"; \
	  kill -9 $$PID || true; \
	fi
	@cd $(FRONTEND); \
	if [ ! -d node_modules ]; then \
	  echo "📦 Installing frontend deps..."; \
	  npm ci || npm install; \
	fi; \
	nohup npm run dev -- --port $(FRONTEND_PORT) --strictPort > /tmp/apexmemory-frontend.log 2>&1 & echo $$! > bin/frontend.pid; \
	echo "✅ Frontend started (pid $$(cat bin/frontend.pid)) — logs: /tmp/apexmemory-frontend.log"

dev.up: dev.backend dev.frontend
	@echo "✅ Dev environment up (API :$(API_PORT), Web :$(FRONTEND_PORT))"

# Mobile testing - exposes servers to local network (auto-detects Tailscale)
dev.mobile:
	@echo "📱 Starting dev servers for mobile testing..."
	@# Try Tailscale first, fall back to local WiFi IP
	@MOBILE_IP=$$(tailscale ip -4 2>/dev/null); \
	if [ -n "$$MOBILE_IP" ]; then \
		echo "🔗 Using Tailscale IP: $$MOBILE_IP"; \
	else \
		MOBILE_IP=$$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null); \
		if [ -z "$$MOBILE_IP" ]; then \
			echo "❌ Could not detect IP. Are you connected to WiFi or Tailscale?"; \
			exit 1; \
		fi; \
		echo "🌐 Using local WiFi IP: $$MOBILE_IP"; \
	fi; \
	echo ""
	@# Start backend with network binding
	@echo "🚀 (API) Starting Go dev server on 0.0.0.0:$(API_PORT)"
	@mkdir -p $(BACKEND_BIN)
	@PID=$$(lsof -ti :$(API_PORT) || true); \
	if [ -n "$$PID" ]; then \
	  echo "🔪 Killing process on port $(API_PORT) (pid $$PID)"; \
	  kill -9 $$PID || true; \
	fi
	@$(ENV_SH); cd $(BACKEND); \
	LOG_FILE=/tmp/apexmemory-api.log; \
	nohup go run ./cmd/api -port $(API_PORT) > $$LOG_FILE 2>&1 & echo $$! > bin/api.pid; \
	sleep 1; \
	if ! kill -0 $$(cat bin/api.pid) > /dev/null 2>&1; then \
		echo "❌ API failed to start. Logs:"; \
		cat $$LOG_FILE; \
		rm -f bin/api.pid; \
		exit 1; \
	fi; \
	echo "✅ API started (pid $$(cat bin/api.pid))"
	@# Start frontend with --host flag
	@echo "🚀 (Web) Starting SvelteKit dev server with --host"
	@mkdir -p $(FRONTEND_BIN)
	@PID=$$(lsof -ti :$(FRONTEND_PORT) || true); \
	if [ -n "$$PID" ]; then \
	  echo "🔪 Killing process on port $(FRONTEND_PORT) (pid $$PID)"; \
	  kill -9 $$PID || true; \
	fi
	@cd $(FRONTEND); \
	if [ ! -d node_modules ]; then \
	  echo "📦 Installing frontend deps..."; \
	  npm ci || npm install; \
	fi; \
	nohup npm run dev -- --host --port $(FRONTEND_PORT) --strictPort > /tmp/apexmemory-frontend.log 2>&1 & echo $$! > bin/frontend.pid; \
	echo "✅ Frontend started (pid $$(cat bin/frontend.pid))"
	@echo ""
	@MOBILE_IP=$$(tailscale ip -4 2>/dev/null || ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null); \
	echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
	echo "📱 Mobile URL: http://$$MOBILE_IP:$(FRONTEND_PORT)"; \
	echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
	echo ""; \
	echo "Logs: /tmp/apexmemory-api.log, /tmp/apexmemory-frontend.log"; \
	echo "Stop with: make dev.stop"

dev.stop:
	@echo "🛑 Stopping API and Frontend (if running)"
	@PID=$$(lsof -ti :$(API_PORT) || true); if [ -n "$$PID" ]; then echo "🔪 Killing API pid $$PID"; kill -9 $$PID || true; fi
	@PID=$$(lsof -ti :$(FRONTEND_PORT) || true); if [ -n "$$PID" ]; then echo "🔪 Killing Web pid $$PID"; kill -9 $$PID || true; fi
	@[ -f $(BACKEND_BIN)/api.pid ] && rm -f $(BACKEND_BIN)/api.pid || true
	@[ -f $(FRONTEND_BIN)/frontend.pid ] && rm -f $(FRONTEND_BIN)/frontend.pid || true
	@echo "✅ Stopped"

dev.check:
	@echo "🔎 Running quick checks..."
	@cd $(BACKEND); go vet ./... || true
	@cd $(FRONTEND); npm run --silent check || true

# ---------- Testing ----------
.PHONY: test.unit test.backend test.backend.v test.backend.cover test.frontend test.frontend.e2e test.all

# Unit tests only (no Docker required)
test.unit:
	@echo "Running unit tests..."
	@cd $(BACKEND) && go test ./...

# All tests including integration (requires Docker daemon running)
test.backend:
	@echo "Running backend tests (unit + integration)..."
	@cd $(BACKEND) && go test -tags integration ./...

test.backend.v:
	@echo "Running backend tests (verbose)..."
	@cd $(BACKEND) && go test -tags integration -v ./...

test.backend.cover:
	@echo "Running backend tests with coverage..."
	@cd $(BACKEND) && go test -tags integration -coverprofile=coverage.out ./...
	@cd $(BACKEND) && go tool cover -func=coverage.out

test.frontend:
	@echo "Running frontend unit tests..."
	@cd $(FRONTEND) && npm run test

test.frontend.e2e:
	@echo "Running frontend e2e tests..."
	@cd $(FRONTEND) && npm run test:e2e

test.all: test.backend test.frontend
	@echo "All tests completed"

# ---------- Frontend Linting & Formatting ----------
.PHONY: lint.frontend format.frontend validate.frontend

lint.frontend:
	@echo "🔍 Running ESLint on frontend..."
	@cd $(FRONTEND); npm run lint

format.frontend:
	@echo "✨ Running Prettier on frontend..."
	@cd $(FRONTEND); npm run format

validate.frontend:
	@echo "🔎 Running all frontend validations..."
	@cd $(FRONTEND); npm run validate

# ---------- Tern ----------
.PHONY: tern.new tern.migrate tern.status tern.rollback tern.reset
tern.new:
	@if [ -z "$(name)" ]; then echo "❌ Usage: make tern.new name=my_change"; exit 2; fi
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern new "$(name)"
	@echo "✅ Created migration in $(MIGR_DIR)"
	@echo "💡 Edit the file and then run: make tern.migrate"

tern.migrate:
	@echo "⬆️  Running database migrations..."
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern migrate
	$(ENV_SH); make db.schema

tern.migrate.to:
	@echo "🎯 Migrating to destination: $(DEST)"
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern migrate --destination $(DEST)

tern.code.install:
	@echo "⬆️  Running database migrations..."
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern code install $(TERN_CODE_DIR)

tern.code.snapshot:
	@echo "⬆️  Running database migrations..."
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern code snapshot $(TERN_CODE_DIR)

tern.status:
	@echo "📊 Migration status:"
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern status

# Step back one migration
tern.rollback:
	@echo "⬇️  Rolling back one migration..."
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern migrate --destination -1

# Dev-only: go to version 0 then back up
tern.reset:
	@echo "🔄 Resetting database (migrate to 0, then back up)..."
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern migrate --destination 0
	$(ENV_SH); TERN_CONFIG=$(TERN_CONF) TERN_MIGRATIONS=$(MIGR_DIR) tern migrate

# ---------- psql helpers (using docker exec) ----------
.PHONY: db.psql.app db.psql.migrator db.psql.super db.psql.claude
db.psql.app:
	@echo "🐘 Connecting to PostgreSQL as app user..."
	@$(ENV_SH); docker exec -it $(PG_CONTAINER) psql -U "$$PG_APP_USER" -d "$$PG_DATABASE"

db.psql.migrator:
	@echo "🐘 Connecting to PostgreSQL as migrator user..."
	@$(ENV_SH); docker exec -it $(PG_CONTAINER) psql -U "$$PG_MIGRATOR_USER" -d "$$PG_DATABASE"

db.psql.super:
	@echo "🐘 Connecting to PostgreSQL as superuser..."
	@$(ENV_SH); docker exec -it $(PG_CONTAINER) psql -U "$$PG_SUPER_USER" -d "$$PG_DATABASE"

# Non-interactive psql for scripts/agents: make db.psql.claude SQL="SELECT 1"
db.psql.claude:
ifndef SQL
	$(error SQL is required. Usage: make db.psql.claude SQL="SELECT 1")
endif
	@$(ENV_SH); docker exec $(PG_CONTAINER) psql -U "$$PG_APP_USER" -d "$$PG_DATABASE" -c "$(SQL)"

# Quick database inspection commands
.PHONY: db.tables db.users db.schemas
db.tables:
	@echo "📋 Tables in app schema:"
	@$(ENV_SH); docker exec -t $(PG_CONTAINER) psql -U "$$PG_APP_USER" -d "$$PG_DATABASE" -c "\dt $$PG_APP_SCHEMA.*"

db.users:
	@echo "👥 Database users:"
	@$(ENV_SH); docker exec -t $(PG_CONTAINER) psql -U "$$PG_SUPER_USER" -d "$$PG_DATABASE" -c "\du"

db.schemas:
	@echo "📁 Database schemas:"
	@$(ENV_SH); docker exec -t $(PG_CONTAINER) psql -U "$$PG_SUPER_USER" -d "$$PG_DATABASE" -c "\dn+"

# ---------- sqlc ----------
.PHONY: db.schema db.sqlc db.dump
db.schema:
	@echo "📸 Dumping database schema to file..."
	@$(ENV_SH); \
	docker exec $(PG_CONTAINER) pg_dump -U "$$PG_MIGRATOR_USER" -d "$$PG_DATABASE" \
	  --schema-only --no-owner --no-privileges -n "$$PG_APP_SCHEMA" > "$(SCHEMA_SNAPSHOT)"
	@echo "✅ Schema dumped to $(SCHEMA_SNAPSHOT)"

.PHONY: db.dump
db.dump:
	@echo "💾 Dumping database (schema+data) to $(DB_DUMP_FILE)..."
	@mkdir -p "$(BACKUPS_DIR)"
	@$(ENV_SH); \
	docker exec $(PG_CONTAINER) pg_dump -U "$$PG_MIGRATOR_USER" -d "$$PG_DATABASE" \
	  --no-owner --no-privileges -n "$$PG_APP_SCHEMA" > "$(DB_DUMP_FILE)"
	@echo "✅ Dump written to $(DB_DUMP_FILE)"

db.sqlc:
	@echo "🔄 Generating Go code from SQL..."
	@$(ENV_SH); cd $(BACKEND); sqlc generate
	@echo "✅ sqlc generation complete"


# ---------- Seed Data ----------
.PHONY: seed.notebooks seed.reviews

# Seed notebooks for a user: make seed.notebooks EMAIL=user@example.com [CLEAR=1]
# CLEAR=1: Wipe existing notebooks before seeding
seed.notebooks:
ifndef EMAIL
	$(error EMAIL is required. Usage: make seed.notebooks EMAIL=user@example.com [CLEAR=1])
endif
	@echo "🌱 Seeding notebooks for $(EMAIL)..."
	@$(ENV_SH); cd $(BACKEND); go run ./scripts/seed-notebooks -email $(EMAIL) $(if $(CLEAR),-clear,)

# Generate review history: make seed.reviews EMAIL=user@example.com [DAYS=30] [CONSISTENCY=high]
# DAYS: number of days of history (default: 30)
# CONSISTENCY: high (no misses), medium (1-2/week misses), low (4-5/week misses)
seed.reviews:
ifndef EMAIL
	$(error EMAIL is required. Usage: make seed.reviews EMAIL=user@example.com DAYS=30 CONSISTENCY=high)
endif
	@echo "📊 Generating reviews for $(EMAIL)..."
	@$(ENV_SH); cd $(BACKEND); go run ./scripts/generate-reviews -email $(EMAIL) -days $(or $(DAYS),30) -consistency $(or $(CONSISTENCY),high)
