-- +goose Up
CREATE TABLE IF NOT EXISTS ping (
    id SERIAL PRIMARY KEY,
    message VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ping;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
