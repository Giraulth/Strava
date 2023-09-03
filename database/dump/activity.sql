CREATE TABLE activity (
    start_date TIMESTAMP NOT NULL,
    id BIGINT NOT NULL,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    distance NUMERIC NOT NULL,
    moving_time INT NOT NULL,
    elapsed_time INT NOT NULL,
    average_heartrate NUMERIC DEFAULT 0,
    max_heartrate NUMERIC DEFAULT 0,
    PRIMARY KEY (id)
);