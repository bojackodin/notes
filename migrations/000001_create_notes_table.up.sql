CREATE TABLE IF NOT EXISTS notes (
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL,
    title text NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username varchar(255) unique NOT NULL,
    password bytea NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);


