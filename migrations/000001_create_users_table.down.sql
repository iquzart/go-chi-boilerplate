-- +migrate Down

-- Drop table in auth schema
DROP TABLE IF EXISTS auth.users;

-- Optionally drop the schema if empty
DROP SCHEMA IF EXISTS auth CASCADE;

