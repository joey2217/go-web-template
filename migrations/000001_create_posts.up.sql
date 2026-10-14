CREATE TABLE IF NOT EXISTS posts (
    id         UUID         PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    status     SMALLINT     NOT NULL DEFAULT 0,
    author     VARCHAR(64)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_posts_status     ON posts (status);
CREATE INDEX IF NOT EXISTS idx_posts_author     ON posts (author);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);
