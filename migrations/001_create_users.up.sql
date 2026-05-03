CREATE TABLE IF NOT EXISTS users (
    id           TEXT PRIMARY KEY COLLATE BINARY, -- ^usr-[A-Za-z0-9]+$
    name         TEXT NOT NULL,
    email        TEXT NOT NULL UNIQUE,
    phone_number VARCHAR(16) NOT NULL, -- e.164 formated phone number ^\+[1-9]\d{1,14}$
    line_1       TEXT NOT NULL,
    line_2       TEXT NULL,
    line_3       TEXT NULL,
    town         TEXT NOT NULL,
    county       TEXT NOT NULL,
    postcode     TEXT NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);