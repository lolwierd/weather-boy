package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/model"
	"github.com/lolwierd/weatherboy/be/internal/repository"
)

// bucketToMMPerHr converts a precip_intensity bucket to mm_per_hr.
func bucketToMMPerHr(bucket int) float64 {
	mapping := []float64{0, 0.25, 0.5, 1, 2, 4, 8, 16, 32, 64}
	if bucket < 0 {
		return 0
	}
	if bucket >= len(mapping) {
		return mapping[len(mapping)-1]
	}
	return mapping[bucket]
}

// metNetResp represents the minimal subset of the MetNet response we care about.
type metNetResp struct {
	CapturedAt  time.Time `json:"start_time"`
	StepMinutes int       `json:"step_minutes"`
	POP         []float64 `json:"pop"`
	Intensity   []int     `json:"precip_intensity"`
}

// FetchMetNetNowcast fetches nowcast data from MetNet for Vadodara and stores it.
func FetchMetNetNowcast(ctx context.Context) error {
	config.LoadEnv()
	url := fmt.Sprintf("%s/%.2f,%.2f", config.MetNetBaseURL, 22.30, 73.20)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("metnet status %s: %s", resp.Status, string(b))
	}

	var data metNetResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if data.StepMinutes == 0 {
		data.StepMinutes = 5
	}
	for i := 0; i < len(data.POP) && i < len(data.Intensity); i++ {
		n := model.Nowcast{
			Location:   "vadodara",
			CapturedAt: data.CapturedAt,
			LeadMin:    i * data.StepMinutes,
			POP:        data.POP[i],
			MMPerHr:    bucketToMMPerHr(data.Intensity[i]),
		}
		if err := repository.InsertNowcast(ctx, &n); err != nil {
			logger.Error.Println("insert nowcast:", err)
		}
	}
	return nil
}
