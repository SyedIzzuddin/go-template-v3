-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Create index for role-based queries
CREATE INDEX idx_users_role ON users(role);

-- Add check constraint to ensure valid roles
ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('admin', 'moderator', 'user'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove check constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role;

-- Remove index
DROP INDEX IF EXISTS idx_users_role;

-- Remove role column
ALTER TABLE users DROP COLUMN IF EXISTS role;
-- +goose StatementEnd