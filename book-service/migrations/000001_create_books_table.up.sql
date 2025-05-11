CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    author text NOT NULL,
    price integer NOT NULL,
    stock integer NOT NULL,
    isbn text NOT NULL,
    image text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);