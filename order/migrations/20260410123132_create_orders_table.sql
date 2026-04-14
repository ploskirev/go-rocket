-- +goose Up
create table orders (
    order_uuid TEXT NOT NULL PRIMARY KEY,
    user_uuid TEXT NOT NULL,
    part_uuids TEXT NOT NULL,
    total_price REAL NOT NULL,
    transaction_uuid TEXT,
    payment_method INTEGER,
    order_status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down
drop table orders;

