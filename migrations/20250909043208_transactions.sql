-- +goose Up
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount != 0),
    transaction_type VARCHAR(20) NOT NULL CHECK (
        transaction_type IN ('deposit', 'withdrawal', 'transfer')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_active ON transactions(id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE transactions;