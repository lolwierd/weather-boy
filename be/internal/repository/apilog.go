package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertIMDAPICall stores an IMD API usage record.
func InsertIMDAPICall(ctx context.Context, l *model.IMDAPICall) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO imd_api_log (endpoint, bytes, requested_at)
         VALUES ($1,$2,$3)
         RETURNING id`,
		l.Endpoint, l.Bytes, l.RequestedAt,
	)
	if err := row.Scan(&l.ID); err != nil {
		
		return err
	}
	return nil
}
