-- +goose Up
-- +goose StatementBegin
-- Add password reset fields to users table
ALTER TABLE users ADD COLUMN password_reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN password_reset_expires_at TIMESTAMP;

-- Create index for token lookup performance
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove password reset fields
DROP INDEX IF EXISTS idx_users_password_reset_token;
ALTER TABLE users DROP COLUMN IF EXISTS password_reset_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS password_reset_token;
-- +goose StatementEnd