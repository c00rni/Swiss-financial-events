-- +goose Up
CREATE TABLE IF NOT EXISTS topics (
    id VARCHAR PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS eventTopics (
    event_id VARCHAR NOT NULL,
    topic_id VARCHAR NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, topic_id)
);

-- +goose Down
DROP TABLE eventTopics;
DROP TABLE topics;
