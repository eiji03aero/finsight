-- Create workspaces table
CREATE TABLE IF NOT EXISTS workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add comment to table
COMMENT ON TABLE workspaces IS 'Workspaces for organizing users and data';
COMMENT ON COLUMN workspaces.name IS 'Workspace name (user-provided)';
