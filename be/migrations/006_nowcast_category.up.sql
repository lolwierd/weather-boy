CREATE TABLE IF NOT EXISTS nowcast_category (
    id SERIAL PRIMARY KEY,
    nowcast_id INT REFERENCES nowcast(id),
    category INT NOT NULL,
    value SMALLINT NOT NULL
);
