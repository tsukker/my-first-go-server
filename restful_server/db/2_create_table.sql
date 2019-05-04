\c app_db;

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    created_at text NOT NULL,
    updated_at text NOT NULL
);
