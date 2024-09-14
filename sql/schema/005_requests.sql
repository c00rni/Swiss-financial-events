-- +goose Up
CREATE TABLE IF NOT EXISTS requests(
    id VARCHAR PRIMARY KEY NOT NULL,
    user_id VARCHAR NOT NULL,
    date DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE requests;
