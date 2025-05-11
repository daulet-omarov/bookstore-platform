CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW()
);