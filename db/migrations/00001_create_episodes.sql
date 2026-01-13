-- +goose Up
-- +goose StatementBegin
CREATE TABLE episodes (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    feed_id UUID NOT NULL,
    feed_guid VARCHAR(255) NOT NULL,
    current_position INTEGER,
    played BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_episodes_feed_id ON episodes(feed_id);
CREATE INDEX idx_episodes_feed_guid ON episodes(feed_guid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_episodes_feed_guid;
DROP INDEX IF EXISTS idx_episodes_feed_id;
DROP TABLE IF EXISTS episodes;
-- +goose StatementEnd
