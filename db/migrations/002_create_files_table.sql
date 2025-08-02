-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    uploaded_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_files_uploaded_by ON files(uploaded_by);
CREATE INDEX idx_files_category ON files(category);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd