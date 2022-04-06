CREATE TABLE IF NOT EXISTS users (
--     id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4() UNIQUE NOT NULL,
--     full_name VARCHAR(255),
--     email VARCHAR(255) NOT NULL,
--     password VARCHAR(255) NOT NULL,
--     phone VARCHAR(255),
--     pin VARCHAR(255),
--     referral_link VARCHAR(255),
--     created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
--     updated_at TIMESTAMPTZ

    id VARCHAR PRIMARY KEY UNIQUE NOT NULL,
    email VARCHAR NOT NULL DEFAULT NULL,
    phone VARCHAR,
    user_name VARCHAR,
    full_name VARCHAR,
    password VARCHAR,
    hash VARCHAR,
    referee VARCHAR,
    role VARCHAR NOT NULL DEFAULT 'user',
    is_active bool NOT NULL DEFAULT FALSE,
    is_email_verified bool NOT NULL DEFAULT FALSE,
    is_phone_verified bool NOT NULL DEFAULT FALSE,
    activated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ
);