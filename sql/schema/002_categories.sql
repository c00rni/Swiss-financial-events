-- +goose Up
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS eventCategories (
    event_id VARCHAR NOT NULL,
    category_id VARCHAR NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, category_id)
);

-- +goose Down
DROP TABLE eventCategories;
DROP TABLE categories;
