CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_email on users (email);
