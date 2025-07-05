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
	ID        int       `db:"id"`
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

// NowcastCategory stores category flags for a nowcast row.
type NowcastCategory struct {
	ID        int   `db:"id"`
	NowcastID int   `db:"nowcast_id"`
	Category  int   `db:"category"`
	Value     int16 `db:"value"`
}

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
}

// BulletinParsed mirrors the `bulletin_parsed` table.
type BulletinParsed struct {
	ID            int       `db:"id"`
	BulletinRawID int       `db:"bulletin_raw_id"`
	Location      string    `db:"location"`
	Forecast      string    `db:"forecast"`
	FetchedAt     time.Time `db:"fetched_at"`
}

// Radar mirrors the `radar` table.
type Radar struct {
	ID         int       `db:"id"`
	Location   string    `db:"location"`
	MaxDBZ     int       `db:"max_dbz"`
	CapturedAt time.Time `db:"captured_at"`
	FetchedAt  time.Time `db:"fetched_at"`
}

// RiverBasinQPF mirrors the `river_basin_qpf` table.
type RiverBasinQPF struct {
	ID        int       `db:"id"`
	BasinID   int       `db:"basin_id"`
	Date      time.Time `db:"date"`
	FMO       string    `db:"fmo"`
	Basin     string    `db:"basin"`
	SubBasin  string    `db:"sub_basin"`
	Area      string    `db:"area"`
	Day1      string    `db:"day1"`
	Day2      string    `db:"day2"`
	Day3      string    `db:"day3"`
	Day4      string    `db:"day4"`
	Day5      string    `db:"day5"`
	AAP       string    `db:"aap"`
	FetchedAt time.Time `db:"fetched_at"`
}

// AWSARG mirrors the `aws_arg` table.
type AWSARG struct {
	ID            int       `db:"id"`
	StationID     string    `db:"station_id"`
	CallSign      string    `db:"call_sign"`
	District      string    `db:"district"`
	State         string    `db:"state"`
	StationName   string    `db:"station_name"`
	Date          time.Time `db:"date"`
	Time          time.Time `db:"time"`
	CurrentTemp   float64   `db:"current_temp"`
	DewPointTemp  float64   `db:"dew_point_temp"`
	RH            float64   `db:"rh"`
	WindDirection float64   `db:"wind_direction"`
	WindSpeed     float64   `db:"wind_speed"`
	MSLP          float64   `db:"mslp"`
	MinTemp       float64   `db:"min_temp"`
	MaxTemp       float64   `db:"max_temp"`
	Latitude      float64   `db:"latitude"`
	Longitude     float64   `db:"longitude"`
	WeatherCode   string    `db:"weather_code"`
	Nebulosity    float64   `db:"nebulosity"`
	FeelLike      float64   `db:"feel_like"`
	RainfallSel   string    `db:"rainfall_sel"`
	Rainfall      float64   `db:"rainfall"`
	FetchedAt     time.Time `db:"fetched_at"`
}

