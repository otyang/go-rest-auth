CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR NOT NULL,
    user_id VARCHAR NOT NULL,
    user_agent VARCHAR,
    client_ip VARCHAR,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    PRIMARY KEY (session_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);