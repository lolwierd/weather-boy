package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertRadarSnapshot stores a radar snapshot.
func InsertRadarSnapshot(ctx context.Context, r *model.RadarSnapshot) error {
	conn, tx, err := GetConnTransaction(ctx)
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Release()
	}

	row := tx.QueryRow(ctx,
		`INSERT INTO radar_snapshot (location, captured_at, max_dbz, bearing, range_km)
         VALUES ($1,$2,$3,$4,$5)
         RETURNING id, created_at`,
		r.Location, r.CapturedAt, r.MaxDBZ, r.Bearing, r.RangeKM,
	)
	if err := row.Scan(&r.ID, &r.CreatedAt); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
