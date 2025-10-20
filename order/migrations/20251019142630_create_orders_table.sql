-- +goose Up
CREATE TABLE orders
(
    id               UUID PRIMARY KEY,
    user_id          UUID           NOT NULL,
    part_ids         UUID[]         NOT NULL,
    total_price      NUMERIC(12, 2) NOT NULL,
    transaction_uuid UUID,
    payment_method   TEXT,
    status           TEXT           NOT NULL,
    created_at       TIMESTAMP      NOT NULL,
    updated_at       TIMESTAMP
);

-- +goose Down
DROP TABLE orders
