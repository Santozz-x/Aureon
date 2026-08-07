CREATE TABLE wallet_keys (
    address                 TEXT PRIMARY KEY,
    encrypted_private_key   BYTEA NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE api_keys (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL,
    secret_hash     TEXT NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ
);
