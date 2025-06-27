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

// IMDAPICall records each call made to IMD endpoints.
type IMDAPICall struct {
	ID          int       `db:"id"`
	Endpoint    string    `db:"endpoint"`
	Bytes       int64     `db:"bytes"`
	RequestedAt time.Time `db:"requested_at"`
}

// NowcastRaw stores the unparsed nowcast JSON for historical reference.
type NowcastRaw struct {
	ID        int       `db:"id"`
	Location  string    `db:"location"`
	Data      []byte    `db:"data"`
	FetchedAt time.Time `db:"fetched_at"`
}

<<<<<<< Updated upstream
// NowcastCategory stores category flags for a nowcast row.
type NowcastCategory struct {
	ID        int   `db:"id"`
	NowcastID int   `db:"nowcast_id"`
	Category  int   `db:"category"`
	Value     int16 `db:"value"`
=======
// DistrictWarning mirrors the `district_warning` table.
type DistrictWarning struct {
	ID          int       `db:"id"`
	Location    string    `db:"location"`
	IssuedAt    time.Time `db:"issued_at"`
	Day1Warning string    `db:"day1_warning"`
	Day2Warning string    `db:"day2_warning"`
	Day3Warning string    `db:"day3_warning"`
	Day4Warning string    `db:"day4_warning"`
	Day5Warning string    `db:"day5_warning"`
	Day1Color   string    `db:"day1_color"`
	Day2Color   string    `db:"day2_color"`
	Day3Color   string    `db:"day3_color"`
	Day4Color   string    `db:"day4_color"`
	Day5Color   string    `db:"day5_color"`
	CreatedAt   time.Time `db:"created_at"`
}

// DistrictWarningRaw stores the unparsed district warning JSON for historical reference.
type DistrictWarningRaw struct {
	ID        int       `db:"id"`
	Location  string    `db:"location"`
	Data      []byte    `db:"data"`
	FetchedAt time.Time `db:"fetched_at"`
>>>>>>> Stashed changes
}
