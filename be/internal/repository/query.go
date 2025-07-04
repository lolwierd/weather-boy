package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// LatestBulletin returns the most recent bulletin for a location.
func LatestBulletin(ctx context.Context, loc string) (*model.Bulletin, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx, `SELECT id, location, issued_at, text, created_at
        FROM bulletin WHERE location=$1 ORDER BY issued_at DESC LIMIT 1`, loc)
	var b model.Bulletin
	if err := row.Scan(&b.ID, &b.Location, &b.IssuedAt, &b.Text, &b.CreatedAt); err != nil {
		return nil, err
	}
	return &b, nil
}

// LatestRadarSnapshot returns the latest radar snapshot for a location.
func LatestRadarSnapshot(ctx context.Context, loc string) (*model.RadarSnapshot, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx, `SELECT id, location, captured_at, max_dbz, bearing, range_km, created_at
        FROM radar_snapshot WHERE location=$1 ORDER BY captured_at DESC LIMIT 1`, loc)
	var r model.RadarSnapshot
	if err := row.Scan(&r.ID, &r.Location, &r.CapturedAt, &r.MaxDBZ, &r.Bearing, &r.RangeKM, &r.CreatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}

// NowcastPOP1H returns the probability of precipitation for the first hour of the latest nowcast.
func NowcastPOP1H(ctx context.Context, loc string) (float64, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx, `SELECT pop FROM nowcast WHERE location=$1 AND captured_at=(SELECT MAX(captured_at) FROM nowcast WHERE location=$1) AND lead_min <= 60 ORDER BY lead_min DESC LIMIT 1`, loc)
	var pop float64
	if err := row.Scan(&pop); err != nil {
		return 0, err
	}
	return pop, nil
}

// NowcastSlice returns the latest nowcast rows up to lead_min 240 minutes.
func NowcastSlice(ctx context.Context, loc string) ([]model.Nowcast, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx, `SELECT id, location, captured_at, lead_min, pop, mm_per_hr, created_at
        FROM nowcast WHERE location=$1 AND captured_at=(SELECT MAX(captured_at) FROM nowcast WHERE location=$1) AND lead_min <= 240 ORDER BY lead_min`, loc)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Nowcast
	for rows.Next() {
		var n model.Nowcast
		if err := rows.Scan(&n.ID, &n.Location, &n.CapturedAt, &n.LeadMin, &n.POP, &n.MMPerHr, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

// LatestNowcastCategories returns category values for the latest nowcast row.
func LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `SELECT id FROM nowcast WHERE location=$1 ORDER BY captured_at DESC LIMIT 1`, loc)
	var nid int
	if err := row.Scan(&nid); err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, `SELECT category, value FROM nowcast_category WHERE nowcast_id=$1`, nid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[int]int16)
	for rows.Next() {
		var cat int
		var val int16
		if err := rows.Scan(&cat, &val); err != nil {
			return nil, err
		}
		m[cat] = val
	}
	return m, rows.Err()
}

// LatestRiverBasinQPF returns the latest river basin QPF for a location.
func LatestRiverBasinQPF(ctx context.Context, loc string) (*model.RiverBasinQPF, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx, `SELECT id, basin_id, date, fmo, basin, sub_basin, area, day1, day2, day3, day4, day5, aap, fetched_at
        FROM river_basin_qpf WHERE basin=$1 ORDER BY fetched_at DESC LIMIT 1`, loc)
	var r model.RiverBasinQPF
	if err := row.Scan(&r.ID, &r.BasinID, &r.Date, &r.FMO, &r.Basin, &r.SubBasin, &r.Area, &r.Day1, &r.Day2, &r.Day3, &r.Day4, &r.Day5, &r.AAP, &r.FetchedAt); err != nil {
		return nil, err
	}
	return &r, nil
}
