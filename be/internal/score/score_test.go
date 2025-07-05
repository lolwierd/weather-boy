package score

import (
	"context"
	"testing"
	"fmt"

	"github.com/lolwierd/weatherboy/be/internal/model"
)

type stubRepo struct {
	bulletin string
	dbz      float64
	rng      float64
	pop      float64
	cats     map[int]int16
	warn     string
	qpf      float64
	rainfall float64
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
func (s stubRepo) LatestDistrictWarning(ctx context.Context, loc string) (*model.DistrictWarning, error) {
	if s.warn == "" {
		return nil, context.Canceled
	}
	return &model.DistrictWarning{Day1Color: s.warn}, nil
}
func (s stubRepo) LatestRiverBasinQPF(ctx context.Context, loc string) (*model.RiverBasinQPF, error) {
	if s.qpf == 0 {
		return nil, context.Canceled
	}
	return &model.RiverBasinQPF{Day1: fmt.Sprintf("%.2f", s.qpf)}, nil
}
func (s stubRepo) LatestAWSARG(ctx context.Context, loc string) (*model.AWSARG, error) {
	if s.rainfall == 0 {
		return nil, context.Canceled
	}
	return &model.AWSARG{Rainfall: s.rainfall}, nil
}

func TestRiskLevels(t *testing.T) {
	cases := []struct {
		name  string
		repo  stubRepo
		level string
	}{
		{"red", stubRepo{"heavy rain", 50, 30, 0.8, map[int]int16{2: 1}, "", 0, 0}, "RED"},
		{"orange", stubRepo{"heavy rain", 0, 0, 0.8, map[int]int16{2: 1}, "", 0, 0}, "ORANGE"},
		{"orange2", stubRepo{"heavy rain", 0, 0, 0, map[int]int16{2: 1}, "", 0, 0}, "ORANGE"},
		{"catalert", stubRepo{"", 0, 0, 0, map[int]int16{14: 1}, "", 0, 0}, "RED"},
		{"warnorange", stubRepo{"", 0, 0, 0, nil, "Orange", 0, 0}, "ORANGE"},
		{"green", stubRepo{"", 0, 0, 0, nil, "", 0, 0}, "GREEN"},
		{"qpf_test", stubRepo{"", 0, 0, 0, nil, "", 10.0, 0}, "GREEN"},
		{"aws_arg_test", stubRepo{"", 0, 0, 0, nil, "", 0, 6.0}, "GREEN"},
	}
	for _, tc := range cases {
		SetRepo(tc.repo)
		got, _ := RiskLevel(context.Background(), "vadodara")
		if got.Level != tc.level {
			t.Errorf("%s: want %s got %s", tc.name, tc.level, got.Level)
		}
	}
}
