CREATE TABLE IF NOT EXISTS users (
    id         UUID         PRIMARY KEY,
    username   VARCHAR(64)  NOT NULL,
    email      VARCHAR(128) NOT NULL,
    nickname   VARCHAR(64),
    status     SMALLINT     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username   ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email      ON users (email);
CREATE INDEX        IF NOT EXISTS idx_users_status     ON users (status);
CREATE INDEX        IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX        IF NOT EXISTS idx_users_created_at ON users (created_at DESC);
