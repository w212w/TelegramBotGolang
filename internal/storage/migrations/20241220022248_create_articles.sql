-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    source_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sources;
-- +goose StatementEnd
