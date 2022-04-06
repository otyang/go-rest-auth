CREATE TABLE IF NOT EXISTS user_configs (
    user_id VARCHAR NOT NULL,
    kyc_level NUMERIC,
    id_card VARCHAR,
    acc_flagged BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);