-- +migrate Up

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    password TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- faster lookups by email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert the admin user with a hashed password
INSERT INTO users (first_name, last_name, email, role, password, status)
VALUES (
    'Admin',
    'User',
    'admin@admin.com',
    'admin',
    '$2a$12$q3t4nJVrzQeU8XzaeJyVxOyFcVqhjUKapyl234VYjk7rJfvc8sENq', -- Hashed password for 'P@ssw0rd'
    'active'
);
