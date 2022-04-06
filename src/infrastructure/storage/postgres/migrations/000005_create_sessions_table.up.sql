CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR PRIMARY KEY UNIQUE NOT NULL,
    user_id VARCHAR NOT NULL,
    user_agent VARCHAR,
    client_ip VARCHAR,
    is_logout BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);