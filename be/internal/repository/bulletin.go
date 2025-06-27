package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertBulletin inserts a bulletin into the database.
func InsertBulletin(ctx context.Context, b *model.Bulletin) error {
	conn, tx, err := GetConnTransaction(ctx)
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Release()
	}

	row := tx.QueryRow(ctx,
		`INSERT INTO bulletin (location, issued_at, text)
         VALUES ($1,$2,$3)
         RETURNING id, created_at`,
		b.Location, b.IssuedAt, b.Text,
	)
	if err := row.Scan(&b.ID, &b.CreatedAt); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// InsertBulletinRaw records a fetched bulletin file path.
func InsertBulletinRaw(ctx context.Context, br *model.BulletinRaw) error {
	conn, tx, err := GetConnTransaction(ctx)
	if err != nil {
		return err
	}
	if conn != nil {
		defer conn.Release()
	}

	row := tx.QueryRow(ctx,
		`INSERT INTO bulletin_raw (path, fetched_at)
         VALUES ($1,$2)
         RETURNING id`,
		br.Path, br.FetchedAt,
	)
	var id int
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
