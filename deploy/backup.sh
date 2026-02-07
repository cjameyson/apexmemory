#!/usr/bin/env bash
set -euo pipefail

# Apex Memory - Database Backup
# Dumps app + app_code schemas via pg_dump, gzip compresses, prunes files older than 7 days.
# Intended to run via launchd every 6 hours.

BACKUP_DIR="/opt/apexmemory-data/backups"
CONTAINER="apexmemory-pg18"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/backup-${TIMESTAMP}.sql.gz"
RETENTION_DAYS=7
MIN_SIZE_BYTES=500

mkdir -p "${BACKUP_DIR}"

echo "[backup] Starting database backup at $(date)"

# Dump app + app_code schemas, pipe through gzip
docker exec "${CONTAINER}" pg_dump \
  -U apex_migrator \
  -d apexmemory \
  --no-owner \
  --no-privileges \
  -n app \
  -n app_code \
  | gzip > "${BACKUP_FILE}"

# Verify output file exists and is not suspiciously small
FILE_SIZE=$(stat -f%z "${BACKUP_FILE}" 2>/dev/null || stat --printf="%s" "${BACKUP_FILE}" 2>/dev/null)
if [ "${FILE_SIZE}" -lt "${MIN_SIZE_BYTES}" ]; then
  echo "[backup] ERROR: Backup file suspiciously small (${FILE_SIZE} bytes). Keeping file for inspection." >&2
  exit 1
fi

echo "[backup] Backup written: ${BACKUP_FILE} (${FILE_SIZE} bytes)"

# Prune backups older than retention period
PRUNED=$(find "${BACKUP_DIR}" -name "backup-*.sql.gz" -mtime +${RETENTION_DAYS} -print -delete | wc -l | tr -d ' ')
if [ "${PRUNED}" -gt 0 ]; then
  echo "[backup] Pruned ${PRUNED} backup(s) older than ${RETENTION_DAYS} days"
fi

echo "[backup] Done at $(date)"
