-- Create user_workspaces join table for many-to-many relationship
CREATE TABLE IF NOT EXISTS user_workspaces (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, workspace_id)
);

-- Create indexes for efficient lookups
CREATE INDEX IF NOT EXISTS user_workspaces_user_id_idx ON user_workspaces (user_id);
CREATE INDEX IF NOT EXISTS user_workspaces_workspace_id_idx ON user_workspaces (workspace_id);

-- Add comment to table
COMMENT ON TABLE user_workspaces IS 'Many-to-many relationship between users and workspaces';
