CREATE TABLE users (
    id BIGSERIAL  PRIMARY KEY,

    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,

    name TEXT NOT NULL,
    photo_url TEXT,
    bio TEXT,

    is_email_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE transactions (
    id BIGSERIAL  PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    amount NUMERIC(18,2) NOT NULL,
    currency TEXT,
    category SMALLINT NOT NULL,
    type SMALLINT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,

    description TEXT,

    date TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE auth_sessions (
    id BIGSERIAL  PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    refresh_token_hash TEXT NOT NULL,

    device_name TEXT,
    ip_address TEXT,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE password_reset_otps (
    id BIGSERIAL  PRIMARY KEY,

    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    otp_hash TEXT NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,

    used BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT NOW()
);