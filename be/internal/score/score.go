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
<<<<<<< Updated upstream
	NowcastPOP1H(ctx context.Context, loc string) (float64, error)
	LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error)
=======
	LatestNowcast(ctx context.Context, loc string) (*model.Nowcast, error)
	LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error)
>>>>>>> Stashed changes
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

func (dbRepo) LatestNowcast(ctx context.Context, loc string) (*model.Nowcast, error) {
	return repository.LatestNowcast(ctx, loc)
}
func (dbRepo) LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error) {
	return repository.LatestDistrictWarning(ctx, loc)
}
func (dbRepo) LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error) {
	return repository.LatestNowcastCategories(ctx, loc)
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

	if n, err := r.LatestNowcast(ctx, loc); err == nil {
		if n.POP >= 0.7 {
			res.Score += 0.5
			res.Breakdown["nowcast_pop"] = 0.5
		}
		if n.MMPerHr >= 4.0 { // Threshold for heavy rain, adjust as needed
			res.Score += 0.8 // Significant weight for heavy rain
			res.Breakdown["nowcast_mm_per_hr"] = 0.8
		}
	}

	if dw, err := r.LatestDistrictWarning(ctx, loc); err == nil {
		// Check for heavy rain or flood warnings for Day 1 (next 24 hours)
		if strings.Contains(strings.ToLower(dw.Day1Warning), "heavy rain") || strings.Contains(strings.ToLower(dw.Day1Warning), "very heavy rain") {
			res.Score += 0.8 // Significant weight for immediate flood risk
			res.Breakdown["district_warning_day1"] = 0.8
		}
		// Check for heavy rain or flood warnings for Day 2 (next 24-48 hours, still relevant for 12-24 hr flood potential)
		if strings.Contains(strings.ToLower(dw.Day2Warning), "heavy rain") || strings.Contains(strings.ToLower(dw.Day2Warning), "very heavy rain") {
			res.Score += 0.5 // Moderate weight for potential future flood risk
			res.Breakdown["district_warning_day2"] = 0.5
		}
	}

	if cats, err := r.LatestNowcastCategories(ctx, loc); err == nil {
		catScore := 0.0
		severeMap := map[int]float64{2: 0.1, 3: 0.1}
		for k, val := range cats {
			if val > 0 {
				if sc, ok := severeMap[k]; ok {
					catScore += sc
				}
			}
		}
		if catScore > 0 {
			res.Score += catScore
			res.Breakdown["categories"] = catScore
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
