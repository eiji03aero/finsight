-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create unique index on email (for case-insensitive uniqueness)
CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email);

-- Add comment to table
COMMENT ON TABLE users IS 'User accounts for the application';
COMMENT ON COLUMN users.email IS 'User email address (unique, lowercase)';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt-hashed password';
