-- +goose Up
-- +goose StatementBegin
-- Add email verification fields to users table
ALTER TABLE users ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN email_verification_token VARCHAR(255);
ALTER TABLE users ADD COLUMN email_verification_expires_at TIMESTAMP;

-- Create index for token lookup performance
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove email verification fields
DROP INDEX IF EXISTS idx_users_email_verification_token;
ALTER TABLE users DROP COLUMN IF EXISTS email_verification_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS email_verification_token;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
-- +goose StatementEnd
