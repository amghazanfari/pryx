-- +goose Up
-- +goose StatementBegin
CREATE TABLE endpoint (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    model_id INT REFERENCES model(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    object TEXT NOT NULL DEFAULT 'endpoint',
    api_key TEXT,
    url_address TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE endpoint;
-- +goose StatementEnd
