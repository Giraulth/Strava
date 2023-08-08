CREATE TABLE kudos (
    activity_id BIGINT NOT NULL,
    username VARCHAR NOT NULL,
    PRIMARY KEY (activity_id, username)
);