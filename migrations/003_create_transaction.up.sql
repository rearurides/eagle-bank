CREATE TABLE IF NOT EXISTS transactions (
    id              TEXT PRIMARY KEY COLLATE BINARY, -- ^tan-[A-Za-z0-9]+$
    account_id      INTEGER NOT NULL,
    ammount         INTEGER NOT NULL DEFAULT 0,
    type            TEXT NOT NULL CHECK(type IN ('deposit', 'withdrawal')),
    reference       VARCHAR(255),
    currency        TEXT NOT NULL DEFAULT 'GBP',
    minor_unit      INTEGER NOT NULL DEFAULT 100,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
    ON DELETE RESTRICT
)
