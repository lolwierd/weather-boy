package score

import (
	"context"
	"strings"

	"github.com/lolwierd/weatherboy/be/internal/model"
	"github.com/lolwierd/weatherboy/be/internal/repository"
)

// Repo defines the data access used for scoring.
type Repo interface {
	LatestBulletin(ctx context.Context, loc string) (*model.Bulletin, error)
	LatestRadarSnapshot(ctx context.Context, loc string) (*model.RadarSnapshot, error)
	NowcastPOP1H(ctx context.Context, loc string) (float64, error)
}

// repo is the default backing repo used in production.
var repo Repo = dbRepo{}

// SetRepo allows tests to swap the repository implementation.
func SetRepo(r Repo) { repo = r }

type dbRepo struct{}

func (dbRepo) LatestBulletin(ctx context.Context, loc string) (*model.Bulletin, error) {
	return repository.LatestBulletin(ctx, loc)
}
func (dbRepo) LatestRadarSnapshot(ctx context.Context, loc string) (*model.RadarSnapshot, error) {
	return repository.LatestRadarSnapshot(ctx, loc)
}
func (dbRepo) NowcastPOP1H(ctx context.Context, loc string) (float64, error) {
	return repository.NowcastPOP1H(ctx, loc)
}

// Result is the risk score output.
type Result struct {
	Level     string             `json:"level"`
	Score     float64            `json:"score"`
	Breakdown map[string]float64 `json:"breakdown"`
}

// RiskLevel computes the risk level for a location.
func RiskLevel(ctx context.Context, loc string) (Result, error) {
	return riskLevel(ctx, repo, loc)
}

func riskLevel(ctx context.Context, r Repo, loc string) (Result, error) {
	res := Result{Breakdown: map[string]float64{}}

	if b, err := r.LatestBulletin(ctx, loc); err == nil {
		txt := strings.ToLower(b.Text)
		if strings.Contains(txt, "heavy") {
			res.Score += 0.4
			res.Breakdown["bulletin"] = 0.4
		}
	}

	if rad, err := r.LatestRadarSnapshot(ctx, loc); err == nil {
		if rad.MaxDBZ >= 45 {
			if rad.RangeKM == nil || *rad.RangeKM <= 40 {
				res.Score += 0.4
				res.Breakdown["radar"] = 0.4
			}
		}
	}

	if pop, err := r.NowcastPOP1H(ctx, loc); err == nil {
		if pop >= 0.7 {
			res.Score += 0.2
			res.Breakdown["nowcast"] = 0.2
		}
	}

	switch {
	case res.Score >= 0.8:
		res.Level = "RED"
	case res.Score >= 0.5:
		res.Level = "ORANGE"
	case res.Score >= 0.3:
		res.Level = "YELLOW"
	default:
		res.Level = "GREEN"
	}
	return res, nil
}
