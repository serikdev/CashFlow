-- +goose Up
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'TMT',
    is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_accounts_created_at ON accounts(created_at);
CREATE INDEX idx_accounts_active ON accounts(id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE accounts;