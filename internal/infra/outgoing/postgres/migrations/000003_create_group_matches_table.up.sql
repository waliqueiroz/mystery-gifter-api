CREATE TABLE IF NOT EXISTS group_matches (
    group_id UUID NOT NULL REFERENCES groups(id),
    giver_id UUID NOT NULL REFERENCES users(id),
    receiver_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (group_id, giver_id),
    UNIQUE (group_id, receiver_id)
);