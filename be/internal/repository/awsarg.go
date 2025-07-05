package repository

import (
	"context"
	"fmt"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

const insertAWSARG = `
INSERT INTO aws_arg (
	station_id, call_sign, district, state, station_name, date, time, current_temp, dew_point_temp, rh,
	wind_direction, wind_speed, mslp, min_temp, max_temp, latitude, longitude, weather_code, nebulosity,
	feel_like, rainfall_sel, rainfall
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
RETURNING id, fetched_at
`

// InsertAWSARG inserts a new AWS/ARG record into the database.
func InsertAWSARG(ctx context.Context, a *model.AWSARG) error {
	conn, err := getConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		insertAWSARG,
		a.StationID, a.CallSign, a.District, a.State, a.StationName, a.Date, a.Time, a.CurrentTemp, a.DewPointTemp, a.RH,
		a.WindDirection, a.WindSpeed, a.MSLP, a.MinTemp, a.MaxTemp, a.Latitude, a.Longitude, a.WeatherCode, a.Nebulosity,
		a.FeelLike, a.RainfallSel, a.Rainfall,
	)
	if err := row.Scan(&a.ID, &a.FetchedAt); err != nil {
		return err
	}
	return nil
}

const getLatestAWSARG = `
SELECT id, station_id, call_sign, district, state, station_name, date, time, current_temp, dew_point_temp, rh,
       wind_direction, wind_speed, mslp, min_temp, max_temp, latitude, longitude, weather_code, nebulosity,
       feel_like, rainfall_sel, rainfall, fetched_at
FROM aws_arg
WHERE station_id = $1
ORDER BY fetched_at DESC
LIMIT 1
`

// LatestAWSARG retrieves the latest AWS/ARG record for a given station.
func LatestAWSARG(ctx context.Context, stationID string) (*model.AWSARG, error) {
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	a := &model.AWSARG{StationID: stationID}
	row := conn.QueryRow(ctx, getLatestAWSARG, stationID)
	err = row.Scan(
		&a.ID, &a.StationID, &a.CallSign, &a.District, &a.State, &a.StationName, &a.Date, &a.Time, &a.CurrentTemp, &a.DewPointTemp, &a.RH,
		&a.WindDirection, &a.WindSpeed, &a.MSLP, &a.MinTemp, &a.MaxTemp, &a.Latitude, &a.Longitude, &a.WeatherCode, &a.Nebulosity,
		&a.FeelLike, &a.RainfallSel, &a.Rainfall, &a.FetchedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get latest aws/arg: %w", err)
	}
	return a, nil
}
