CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY UNIQUE NOT NULL,
    email VARCHAR UNIQUE,
    phone VARCHAR UNIQUE,
    user_name VARCHAR,
    full_name VARCHAR,
    password VARCHAR,
    hash VARCHAR,
    referral_link VARCHAR UNIQUE NOT NULL,
    referral VARCHAR,
    role VARCHAR NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_google_verified BOOLEAN NOT NULL DEFAULT FALSE,
    google_secret VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ
);