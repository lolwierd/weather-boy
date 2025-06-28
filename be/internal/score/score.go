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
	LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error)
	LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error)
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
func (dbRepo) LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error) {
	return repository.LatestNowcastCategories(ctx, loc)
}
func (dbRepo) LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error) {
	return repository.LatestDistrictWarning(ctx, loc)
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

	if cats, err := r.LatestNowcastCategories(ctx, loc); err == nil {
		catScore := 0.0
		severeMap := map[int]float64{2: 0.1, 3: 0.1}
		alertCats := []int{13, 14, 19}
		triggeredAlert := false
		for k, val := range cats {
			if val > 0 {
				if sc, ok := severeMap[k]; ok {
					catScore += sc
				}
				for _, ac := range alertCats {
					if k == ac {
						triggeredAlert = true
					}
				}
			}
		}
		if triggeredAlert {
			res.Score = 0.9
			res.Breakdown["nowcast_alert"] = 0.9
		} else if catScore > 0 {
			res.Score += catScore
			res.Breakdown["categories"] = catScore
		}
	}

	if dw, err := r.LatestDistrictWarning(ctx, loc); err == nil {
		color := strings.ToLower(dw.Day1Color)
		switch color {
		case "red":
			if res.Score < 0.8 {
				res.Score = 0.8
			}
			res.Breakdown["district_warning"] = 0.8
		case "orange":
			if res.Score < 0.5 {
				res.Score = 0.5
			}
			res.Breakdown["district_warning"] = 0.5
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
