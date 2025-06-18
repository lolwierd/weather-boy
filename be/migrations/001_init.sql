-- +migrate Up
CREATE TABLE IF NOT EXISTS bulletin (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS radar_snapshot (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL,
    max_dbz REAL NOT NULL,
    bearing REAL,
    range_km REAL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS nowcast (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL,
    lead_min INT NOT NULL,
    pop NUMERIC NOT NULL,
    mm_per_hr NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS nowcast;
DROP TABLE IF EXISTS radar_snapshot;
DROP TABLE IF EXISTS bulletin;
