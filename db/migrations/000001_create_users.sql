-- +goose Up
CREATE TABLE users (
    id       BIGINT PRIMARY KEY,
    language TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
