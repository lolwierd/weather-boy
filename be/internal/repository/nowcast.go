package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertNowcast inserts a nowcast record.
func InsertNowcast(ctx context.Context, n *model.Nowcast) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO nowcast (location, captured_at, lead_min, pop, mm_per_hr)
         VALUES ($1,$2,$3,$4,$5)
         RETURNING id, created_at`,
		n.Location, n.CapturedAt, n.LeadMin, n.POP, n.MMPerHr,
	)
	if err := row.Scan(&n.ID, &n.CreatedAt); err != nil {
		return err
	}
	return nil
}

// InsertNowcastRaw stores the raw nowcast JSON.
func InsertNowcastRaw(ctx context.Context, nr *model.NowcastRaw) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO nowcast_raw (location, data, fetched_at)
         VALUES ($1,$2,$3)
         RETURNING id`,
		nr.Location, nr.Data, nr.FetchedAt,
	)
	if err := row.Scan(&nr.ID); err != nil {
		return err
	}
	return nil
}

// InsertNowcastCategory stores a category value for a nowcast row.
func InsertNowcastCategory(ctx context.Context, c *model.NowcastCategory) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO nowcast_category (nowcast_id, category, value)
         VALUES ($1,$2,$3)
         RETURNING id`,
		c.NowcastID, c.Category, c.Value,
	)
	if err := row.Scan(&c.ID); err != nil {
		return err
	}
	return nil
}

// LatestNowcast returns the latest nowcast record for a location.
func LatestNowcast(ctx context.Context, loc string) (*model.Nowcast, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	n := &model.Nowcast{}
	row := conn.QueryRow(ctx,
		`SELECT id, location, captured_at, lead_min, pop, mm_per_hr, created_at
         FROM nowcast
         WHERE location = $1
         ORDER BY captured_at DESC, created_at DESC
         LIMIT 1`,
		loc,
	)
	if err := row.Scan(&n.ID, &n.Location, &n.CapturedAt, &n.LeadMin, &n.POP, &n.MMPerHr, &n.CreatedAt); err != nil {
		return nil, err
	}
	return n, nil
}
