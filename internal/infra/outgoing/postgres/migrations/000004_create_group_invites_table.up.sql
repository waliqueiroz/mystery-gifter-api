CREATE TABLE IF NOT EXISTS group_invites (
    id         UUID         NOT NULL PRIMARY KEY,
    group_id   UUID         NOT NULL REFERENCES groups(id),
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
