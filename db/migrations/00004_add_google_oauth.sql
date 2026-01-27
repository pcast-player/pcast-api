-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE;
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
CREATE INDEX idx_users_google_id ON users(google_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_google_id;
-- Delete OAuth-only users (users without a password) before restoring NOT NULL constraint
DELETE FROM users WHERE password IS NULL;
ALTER TABLE users DROP COLUMN google_id;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
-- +goose StatementEnd
