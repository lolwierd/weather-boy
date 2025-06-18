package model

import "time"

// Bulletin mirrors the `bulletin` table.
type Bulletin struct {
	ID        int       `db:"id"`
	Location  string    `db:"location"`
	IssuedAt  time.Time `db:"issued_at"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

// RadarSnapshot mirrors the `radar_snapshot` table.
type RadarSnapshot struct {
	ID         int       `db:"id"`
	Location   string    `db:"location"`
	CapturedAt time.Time `db:"captured_at"`
	MaxDBZ     float64   `db:"max_dbz"`
	Bearing    *float64  `db:"bearing"`
	RangeKM    *float64  `db:"range_km"`
	CreatedAt  time.Time `db:"created_at"`
}

// Nowcast mirrors the `nowcast` table.
type Nowcast struct {
	ID         int       `db:"id"`
	Location   string    `db:"location"`
	CapturedAt time.Time `db:"captured_at"`
	LeadMin    int       `db:"lead_min"`
	POP        float64   `db:"pop"`
	MMPerHr    float64   `db:"mm_per_hr"`
	CreatedAt  time.Time `db:"created_at"`
}

// BulletinRaw records a fetched bulletin PDF path and time.
type BulletinRaw struct {
	Path      string    `db:"path"`
	FetchedAt time.Time `db:"fetched_at"`
}
