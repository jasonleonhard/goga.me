CREATE TABLE matches (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    black_user_id bigint REFERENCES users (id),
    white_user_id bigint REFERENCES users (id)
);

