CREATE TABLE IF NOT EXISTS orders (
    id bigserial PRIMARY KEY,
    user_id integer NOT NULL,
    book_id integer NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);