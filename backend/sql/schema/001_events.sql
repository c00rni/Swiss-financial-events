-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id VARCHAR PRIMARY KEY NOT NULL,
    title VARCHAR NOT NULL,
    link VARCHAR NOT NULL UNIQUE,
    address VARCHAR NOT NULL,
    location VARCHAR NOT NULL,
    start_at DATETIME NOT NULL,
    end_at DATETIME NOT NULL,
    description TEXT NOT NULL
);

-- +goose Down
DROP TABLE events;
