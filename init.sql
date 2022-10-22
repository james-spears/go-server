DO $do$ BEGIN IF EXISTS (
   SELECT
   FROM pg_catalog.pg_roles
   WHERE rolname = 'api'
) THEN RAISE NOTICE 'Role "api" already exists. Skipping.';
ELSE CREATE ROLE api LOGIN PASSWORD 'password';
END IF;
END $do$;
CREATE DATABASE app OWNER api;
\ connect app;
CREATE SCHEMA IF NOT EXISTS public AUTHORIZATION api;
ALTER SCHEMA public OWNER TO api;
CREATE TABLE IF NOT EXISTS public.users(
   id SERIAL PRIMARY KEY,
   email VARCHAR(100) NOT NULL,
   password VARCHAR(100) NOT NULL
);
/* Too broad, but OK for dev. */
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to api;
GRANT USAGE,
   SELECT ON ALL SEQUENCES IN SCHEMA public TO api;