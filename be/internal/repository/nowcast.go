package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertNowcast inserts a nowcast record.
func InsertNowcast(ctx context.Context, n *model.Nowcast) error {
	conn, tx, err := getConnTransaction(ctx)
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Release()
	}

	row := tx.QueryRow(ctx,
		`INSERT INTO nowcast (location, captured_at, lead_min, pop, mm_per_hr)
         VALUES ($1,$2,$3,$4,$5)
         RETURNING id, created_at`,
		n.Location, n.CapturedAt, n.LeadMin, n.POP, n.MMPerHr,
	)
	if err := row.Scan(&n.ID, &n.CreatedAt); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
