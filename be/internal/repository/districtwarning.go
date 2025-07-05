package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

// InsertDistrictWarning inserts a district warning record.
func InsertDistrictWarning(ctx context.Context, dw *model.DistrictWarning) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO district_warning (location, issued_at, day1_warning, day2_warning, day3_warning, day4_warning, day5_warning, day1_color, day2_color, day3_color, day4_color, day5_color)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
         RETURNING id, created_at`,
		dw.Location, dw.IssuedAt, dw.Day1Warning, dw.Day2Warning, dw.Day3Warning, dw.Day4Warning, dw.Day5Warning, dw.Day1Color, dw.Day2Color, dw.Day3Color, dw.Day4Color, dw.Day5Color,
	)
	if err := row.Scan(&dw.ID, &dw.CreatedAt); err != nil {
		return err
	}
	return nil
}

// InsertDistrictWarningRaw stores the raw district warning JSON.
func InsertDistrictWarningRaw(ctx context.Context, dwr *model.DistrictWarningRaw) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`INSERT INTO district_warning_raw (location, data, fetched_at)
         VALUES ($1,$2,$3)
         RETURNING id`,
		dwr.Location, dwr.Data, dwr.FetchedAt,
	)
	if err := row.Scan(&dwr.ID); err != nil {
		return err
	}
	return nil
}

// LatestDistrictWarning returns the latest district warning record for a location.
func LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	dw := &model.DistrictWarning{}
	row := conn.QueryRow(ctx,
		`SELECT id, location, issued_at, day1_warning, day2_warning, day3_warning, day4_warning, day5_warning, day1_color, day2_color, day3_color, day4_color, day5_color, created_at
         FROM district_warning
         WHERE location = $1
         ORDER BY issued_at DESC, created_at DESC
         LIMIT 1`,
		loc,
	)
	if err := row.Scan(&dw.ID, &dw.Location, &dw.IssuedAt, &dw.Day1Warning, &dw.Day2Warning, &dw.Day3Warning, &dw.Day4Warning, &dw.Day5Warning, &dw.Day1Color, &dw.Day2Color, &dw.Day3Color, &dw.Day4Color, &dw.Day5Color, &dw.CreatedAt); err != nil {
		return nil, err
	}
	return dw, nil
}
