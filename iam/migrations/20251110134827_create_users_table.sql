-- +goose Up
CREATE TABLE users
(
    id                   UUID PRIMARY KEY,
    login                TEXT      NOT NULL UNIQUE,
    password_hash        TEXT      NOT NULL,
    email                TEXT      NOT NULL UNIQUE,
    notification_methods JSONB     NOT NULL,
    created_at           TIMESTAMP NOT NULL,
    updated_at           TIMESTAMP
);

-- +goose Down
DROP TABLE users