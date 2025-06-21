CREATE TABLE IF NOT EXISTS users (
    username    VARCHAR(30) UNIQUE NOT NULL,
    salt        CHAR(64)    NOT NULL,   -- 64 char hex string
    hash        CHAR(64)    NOT NULL,   -- 64 char hex string
    enc_symkey  CHAR(120)   NOT NULL,   -- 120 char hex string, encrypted SYMKEY
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,

    CHECK (username ~   '^[a-zA-Z0-9_]+$' AND length(username) >= 3),
    CHECK (salt ~       '^[0-9a-fA-F]{64}$'),
    CHECK (hash ~       '^[0-9a-fA-F]{64}$'),
    CHECK (enc_symkey ~ '^[0-9a-fA-F]{120}$')
);

