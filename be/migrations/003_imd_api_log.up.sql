CREATE TABLE IF NOT EXISTS imd_api_log (
    id SERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL,
    bytes BIGINT NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
