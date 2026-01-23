#!/usr/bin/env bash
set -euo pipefail

# This script runs inside the PostgreSQL container during first-time init.
# PostgreSQL automatically provides: $POSTGRES_USER, $POSTGRES_PASSWORD, $POSTGRES_DB
# We get our custom variables from docker-compose environment section

echo "Setting up ApexMemory database users and schema..."

# Use the built-in PostgreSQL superuser
export PGPASSWORD="${POSTGRES_PASSWORD}"

psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" <<SQL
-- Ensure roles exist (no superuser privileges)
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '${PG_MIGRATOR_USER}') THEN
    CREATE ROLE ${PG_MIGRATOR_USER} LOGIN PASSWORD '${PG_MIGRATOR_PASSWORD}';
    RAISE NOTICE 'Created migrator role: ${PG_MIGRATOR_USER}';
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '${PG_APP_USER}') THEN
    CREATE ROLE ${PG_APP_USER} LOGIN PASSWORD '${PG_APP_PASSWORD}';
    RAISE NOTICE 'Created application role: ${PG_APP_USER}';
  END IF;
END
\$\$;

-- Grant CREATE privilege on database to migrator (needed for extensions)
GRANT CREATE ON DATABASE ${POSTGRES_DB} TO ${PG_MIGRATOR_USER};

-- Security: prevent arbitrary CREATE in public schema
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- Allow both roles to connect to the database
GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO ${PG_MIGRATOR_USER}, ${PG_APP_USER};

-- Create dedicated app schema owned by migrator
CREATE SCHEMA IF NOT EXISTS ${PG_APP_SCHEMA} AUTHORIZATION ${PG_MIGRATOR_USER};
CREATE SCHEMA IF NOT EXISTS ${PG_APP_CODE_SCHEMA} AUTHORIZATION ${PG_MIGRATOR_USER};

-- Set helpful search_path defaults (so users don't need to qualify schema names)
ALTER ROLE ${PG_MIGRATOR_USER} IN DATABASE ${POSTGRES_DB} SET search_path = ${PG_APP_SCHEMA}, ${PG_APP_CODE_SCHEMA}, public;
ALTER ROLE ${PG_APP_USER}      IN DATABASE ${POSTGRES_DB} SET search_path = ${PG_APP_SCHEMA}, ${PG_APP_CODE_SCHEMA}, public;
-- Ensure superuser has the same convenient search_path
ALTER ROLE ${POSTGRES_USER}    IN DATABASE ${POSTGRES_DB} SET search_path = ${PG_APP_SCHEMA}, ${PG_APP_CODE_SCHEMA}, public;

-- Schema-level privileges
GRANT ALL   ON SCHEMA ${PG_APP_SCHEMA} TO ${PG_MIGRATOR_USER};  -- Full control
GRANT ALL   ON SCHEMA ${PG_APP_CODE_SCHEMA} TO ${PG_MIGRATOR_USER};  -- Full control
GRANT USAGE ON SCHEMA ${PG_APP_SCHEMA} TO ${PG_APP_USER};       -- Read-only access to schema
GRANT USAGE ON SCHEMA ${PG_APP_CODE_SCHEMA} TO ${PG_APP_USER};  -- Read-only access to schema
-- Superuser: explicit full access for clarity (even though superuser bypasses checks)
GRANT ALL   ON SCHEMA ${PG_APP_SCHEMA} TO ${POSTGRES_USER};
GRANT ALL   ON SCHEMA ${PG_APP_CODE_SCHEMA} TO ${POSTGRES_USER};

-- Default privileges: objects created by migrator are automatically accessible to app user
ALTER DEFAULT PRIVILEGES FOR ROLE ${PG_MIGRATOR_USER} IN SCHEMA ${PG_APP_SCHEMA}
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO ${PG_APP_USER};

ALTER DEFAULT PRIVILEGES FOR ROLE ${PG_MIGRATOR_USER} IN SCHEMA ${PG_APP_SCHEMA}
  GRANT USAGE, SELECT ON SEQUENCES TO ${PG_APP_USER};

ALTER DEFAULT PRIVILEGES FOR ROLE ${PG_MIGRATOR_USER} IN SCHEMA ${PG_APP_SCHEMA}
  GRANT EXECUTE ON FUNCTIONS TO ${PG_APP_USER};

-- Default privileges for app_code schema (functions typically live here)
ALTER DEFAULT PRIVILEGES FOR ROLE ${PG_MIGRATOR_USER} IN SCHEMA ${PG_APP_CODE_SCHEMA}
  GRANT EXECUTE ON FUNCTIONS TO ${PG_APP_USER};


-- Verification: show what we created
\echo 'Database setup completed!'
\echo 'Roles created:'
SELECT 
    rolname as role_name,
    rolcanlogin as can_login,
    rolcreatedb as can_create_db
FROM pg_roles 
WHERE rolname IN ('${PG_MIGRATOR_USER}', '${PG_APP_USER}')
ORDER BY rolname;

\echo 'Schemas created:'
SELECT schema_name, schema_owner
FROM information_schema.schemata
WHERE schema_name IN ('${PG_APP_SCHEMA}', '${PG_APP_CODE_SCHEMA}')
ORDER BY schema_name;

SQL

echo "ApexMemory database initialization complete!"
