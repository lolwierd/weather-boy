package score

import (
	"context"
	"testing"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

type stubRepo struct {
	bulletin string
	dbz      float64
	rng      float64
	pop      float64
<<<<<<< Updated upstream
	cats     map[int]int16
=======
	mmPerHr     float64
	day1Warning string
	day2Warning string
>>>>>>> Stashed changes
}

func (s stubRepo) LatestBulletin(ctx context.Context, loc string) (*model.Bulletin, error) {
	if s.bulletin == "" {
		return nil, context.Canceled
	}
	return &model.Bulletin{Text: s.bulletin}, nil
}
func (s stubRepo) LatestRadarSnapshot(ctx context.Context, loc string) (*model.RadarSnapshot, error) {
	if s.dbz == 0 {
		return nil, context.Canceled
	}
	return &model.RadarSnapshot{MaxDBZ: s.dbz, RangeKM: &s.rng}, nil
}
func (s stubRepo) NowcastPOP1H(ctx context.Context, loc string) (float64, error) {
	if s.pop == 0 {
		return 0, context.Canceled
	}
	return s.pop, nil
}
func (s stubRepo) LatestNowcastCategories(ctx context.Context, loc string) (map[int]int16, error) {
	if s.cats == nil {
		return nil, context.Canceled
	}
	return s.cats, nil
}

func (s stubRepo) LatestNowcast(ctx context.Context, loc string) (*model.Nowcast, error) {
	if s.pop == 0 && s.mmPerHr == 0 {
		return nil, context.Canceled
	}
	return &model.Nowcast{POP: s.pop, MMPerHr: s.mmPerHr}, nil
}

func (s stubRepo) LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error) {
	if s.day1Warning == "" && s.day2Warning == "" {
		return nil, context.Canceled
	}
	return &model.DistrictWarning{Day1Warning: s.day1Warning, Day2Warning: s.day2Warning}, nil
}

func TestRiskLevels(t *testing.T) {
	cases := []struct {
		name  string
		repo  stubRepo
		level string
	}{
<<<<<<< Updated upstream
		{"red", stubRepo{"heavy rain", 50, 30, 0.8, map[int]int16{2: 1}}, "RED"},
		{"orange", stubRepo{"heavy rain", 0, 0, 0.8, map[int]int16{2: 1}}, "ORANGE"},
		{"orange2", stubRepo{"heavy rain", 0, 0, 0, map[int]int16{2: 1}}, "ORANGE"},
		{"green", stubRepo{"", 0, 0, 0, nil}, "GREEN"},
=======
		{"red_all_factors", stubRepo{bulletin: "heavy rain", dbz: 50, rng: 30, pop: 0.8, mmPerHr: 5.0, day1Warning: "Heavy Rain"}, "RED"},
		{"red_heavy_rain_only", stubRepo{mmPerHr: 4.5}, "RED"},
		{"orange_pop_only", stubRepo{pop: 0.8}, "ORANGE"},
		{"yellow_bulletin_only", stubRepo{bulletin: "heavy rain"}, "YELLOW"},
		{"red_day1_warning", stubRepo{day1Warning: "Heavy Rain"}, "RED"},
		{"orange_day2_warning", stubRepo{day2Warning: "Very Heavy Rain"}, "ORANGE"},
		{"green_no_factors", stubRepo{}, "GREEN"},
>>>>>>> Stashed changes
	}
	for _, tc := range cases {
		SetRepo(tc.repo)
		got, _ := RiskLevel(context.Background(), "vadodara")
		if got.Level != tc.level {
			t.Errorf("%s: want %s got %s", tc.name, tc.level, got.Level)
		}
	}
}
