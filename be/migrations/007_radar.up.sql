CREATE TABLE radar (
    id SERIAL PRIMARY KEY,
    location VARCHAR(255) NOT NULL,
    max_dbz INT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
