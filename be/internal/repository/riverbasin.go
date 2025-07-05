package repository

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/db"
	"github.com/lolwierd/weatherboy/be/internal/model"
)

const insertRiverBasinQPF = `
INSERT INTO river_basin_qpf (basin_id, date, fmo, basin, sub_basin, area, day1, day2, day3, day4, day5, aap)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, fetched_at
`

// InsertRiverBasinQPF inserts a new river basin QPF record into the database.
func InsertRiverBasinQPF(ctx context.Context, r *model.RiverBasinQPF) error {
	return db.GetDBDriver().ConnPool.QueryRow(ctx, insertRiverBasinQPF, 
		r.BasinID, r.Date, r.FMO, r.Basin, r.SubBasin, r.Area, r.Day1, r.Day2, r.Day3, r.Day4, r.Day5, r.AAP,
	).Scan(&r.ID, &r.FetchedAt)
}
