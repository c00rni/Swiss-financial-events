-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    verified_email BOOLEAN NOT NULL,
    name VARCHAR NOT NULL,
    given_name VARCHAR NOT NULL,
    family_name VARCHAR NOT NULL,
    picture VARCHAR NOT NULL,
    token VARCHAR NOT NULL,
    api_key VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE users;
