CREATE TABLE IF NOT EXISTS nowcast_raw (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    data JSONB NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
