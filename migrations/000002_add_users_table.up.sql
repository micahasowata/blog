CREATE TABLE IF NOT EXISTS users (
    id citext PRIMARY KEY NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    name citext NOT NULL,
    username citext NOT NULL,
    email citext UNIQUE NOT NULL
);