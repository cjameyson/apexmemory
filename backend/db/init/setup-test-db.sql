-- Setup script for test database
-- Run as superuser against the _test database after creating it

-- Grant CREATE privilege on database to migrator (needed for extensions)
GRANT CREATE ON DATABASE :dbname TO :migrator;

-- Security: prevent arbitrary CREATE in public schema
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- Allow both roles to connect to the test database
GRANT CONNECT ON DATABASE :dbname TO :migrator, :appuser;

-- Create dedicated app schema owned by migrator
CREATE SCHEMA IF NOT EXISTS :appschema AUTHORIZATION :migrator;
CREATE SCHEMA IF NOT EXISTS :codeschema AUTHORIZATION :migrator;

-- Set helpful search_path defaults
ALTER ROLE :migrator IN DATABASE :dbname SET search_path = :appschema, :codeschema, public;
ALTER ROLE :appuser IN DATABASE :dbname SET search_path = :appschema, :codeschema, public;

-- Schema-level privileges
GRANT ALL ON SCHEMA :appschema TO :migrator;
GRANT ALL ON SCHEMA :codeschema TO :migrator;
GRANT USAGE ON SCHEMA :appschema TO :appuser;
GRANT USAGE ON SCHEMA :codeschema TO :appuser;

-- Default privileges for objects created by migrator
-- Note: TRUNCATE is needed for test isolation
ALTER DEFAULT PRIVILEGES FOR ROLE :migrator IN SCHEMA :appschema
  GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO :appuser;

ALTER DEFAULT PRIVILEGES FOR ROLE :migrator IN SCHEMA :appschema
  GRANT USAGE, SELECT ON SEQUENCES TO :appuser;

ALTER DEFAULT PRIVILEGES FOR ROLE :migrator IN SCHEMA :appschema
  GRANT EXECUTE ON FUNCTIONS TO :appuser;

ALTER DEFAULT PRIVILEGES FOR ROLE :migrator IN SCHEMA :codeschema
  GRANT EXECUTE ON FUNCTIONS TO :appuser;

-- Grant TRUNCATE on any existing tables (for re-runs after migrations)
GRANT TRUNCATE ON ALL TABLES IN SCHEMA :appschema TO :appuser;
