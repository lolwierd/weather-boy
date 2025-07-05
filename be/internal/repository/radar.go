package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/model"
)

const insertRadar = `
INSERT INTO radar (location, max_dbz, captured_at)
VALUES ($1, $2, $3)
RETURNING id, fetched_at
`

// InsertRadar inserts a new radar record into the database.
func InsertRadar(ctx context.Context, r *model.Radar) error {
	return db.GetDBDriver().ConnPool.QueryRow(ctx, insertRadar, r.Location, r.MaxDBZ, r.CapturedAt).Scan(&r.ID, &r.FetchedAt)
}

const insertRadarSnapshot = `
INSERT INTO radar_snapshot (location, captured_at, max_dbz, bearing, range_km)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at
`

// InsertRadarSnapshot inserts a new radar snapshot record into the database.
func InsertRadarSnapshot(ctx context.Context, rs *model.RadarSnapshot) error {
	return db.GetDBDriver().ConnPool.QueryRow(ctx, insertRadarSnapshot, rs.Location, rs.CapturedAt, rs.MaxDBZ, rs.Bearing, rs.RangeKM).Scan(&rs.ID, &rs.CreatedAt)
}
