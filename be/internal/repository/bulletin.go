package repository

import (
	"context"


	"github.com/lolwierd/weatherboy/be/internal/model"
)

const insertBulletinRaw = `
INSERT INTO bulletin_raw (path, fetched_at)
VALUES ($1, $2)
RETURNING id
`

// InsertBulletinRaw inserts a new bulletin raw record into the database.
func InsertBulletinRaw(ctx context.Context, br *model.BulletinRaw) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, insertBulletinRaw, br.Path, br.FetchedAt)
	if err := row.Scan(&br.ID); err != nil {
		return err
	}
	return nil
}

const insertParsedBulletin = `
INSERT INTO bulletin_parsed (bulletin_raw_id, location, forecast)
VALUES ($1, $2, $3)
RETURNING id, fetched_at
`

// InsertParsedBulletin inserts a new parsed bulletin record into the database.
func InsertParsedBulletin(ctx context.Context, b *model.BulletinParsed) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, insertParsedBulletin, b.BulletinRawID, b.Location, b.Forecast)
	if err := row.Scan(&b.ID, &b.FetchedAt); err != nil {
		return err
	}
	return nil
}