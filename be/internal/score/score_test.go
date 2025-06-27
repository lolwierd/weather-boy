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
	cats     map[int]int16
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

func TestRiskLevels(t *testing.T) {
	cases := []struct {
		name  string
		repo  stubRepo
		level string
	}{
		{"red", stubRepo{"heavy rain", 50, 30, 0.8, map[int]int16{2: 1}}, "RED"},
		{"orange", stubRepo{"heavy rain", 0, 0, 0.8, map[int]int16{2: 1}}, "ORANGE"},
		{"orange2", stubRepo{"heavy rain", 0, 0, 0, map[int]int16{2: 1}}, "ORANGE"},
		{"green", stubRepo{"", 0, 0, 0, nil}, "GREEN"},
	}
	for _, tc := range cases {
		SetRepo(tc.repo)
		got, _ := RiskLevel(context.Background(), "vadodara")
		if got.Level != tc.level {
			t.Errorf("%s: want %s got %s", tc.name, tc.level, got.Level)
		}
	}
}
