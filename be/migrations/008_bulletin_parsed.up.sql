CREATE TABLE bulletin_parsed (
    id SERIAL PRIMARY KEY,
    bulletin_raw_id INT NOT NULL REFERENCES bulletin_raw(id),
    location VARCHAR(255) NOT NULL,
    forecast TEXT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
