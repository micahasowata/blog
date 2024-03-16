CREATE TABLE IF NOT EXISTS users (
    id citext PRIMARY KEY NOT NULL,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    name citext NOT NULL,
    username citext UNIQUE NOT NULL,
    avatar citext NOT NULL,
    bio citext NOT NULL,
    links citext[] NOT NULL,
    email citext UNIQUE NOT NULL,
    verified boolean NOT NULL DEFAULT false
);