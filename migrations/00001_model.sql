-- +goose Up
-- +goose StatementBegin
CREATE TABLE model (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    owned_by TEXT NOT NULL DEFAULT 'openai',
    object TEXT NOT NULL DEFAULT 'model'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE model;
-- +goose StatementEnd
