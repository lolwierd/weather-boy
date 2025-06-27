CREATE TABLE IF NOT EXISTS district_warning (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL,
    day1_warning TEXT,
    day2_warning TEXT,
    day3_warning TEXT,
    day4_warning TEXT,
    day5_warning TEXT,
    day1_color TEXT,
    day2_color TEXT,
    day3_color TEXT,
    day4_color TEXT,
    day5_color TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS district_warning_raw (
    id SERIAL PRIMARY KEY,
    location TEXT NOT NULL,
    data JSONB NOT NULL,
    fetched_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);