package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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

// districtNowcastResp represents the district level nowcast response.
// The IMD endpoint returns an array with a single object containing
// multiple categorical fields and a color code. We only care about the
// overall "color" as an indicator of rainfall likelihood.
type districtNowcastResp struct {
	ObjID string `json:"Obj_id"`
	Date  string `json:"Date"`
	TOI   string `json:"toi"`
	VUpto string `json:"vupto"`
	Color string `json:"color"`
}

// colorToPOP maps the IMD color code (0-4) to an approximate POP value.
func colorToPOP(c int) float64 {
	switch c {
	case 1:
		return 0.3
	case 2:
		return 0.6
	case 3:
		return 0.8
	case 4:
		return 1
	default:
		return 0
	}
}

// FetchIMDNowcast fetches nowcast data from the IMD API for Vadodara and stores it.
func FetchIMDNowcast(ctx context.Context) error {
	const url = "https://mausam.imd.gov.in/api/nowcast_district_api.php?id=244"
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
		return fmt.Errorf("imd nowcast status %s: %s", resp.Status, string(b))
	}

	var arr []districtNowcastResp
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return fmt.Errorf("empty nowcast response")
	}

	// Parse color as integer
	col, err := strconv.Atoi(arr[0].Color)
	if err != nil {
		col = 0
	}

	captured := time.Now()
	if arr[0].Date != "" && arr[0].TOI != "" {
		t, err := time.Parse("2006-01-02 1504", arr[0].Date+" "+arr[0].TOI)
		if err == nil {
			captured = t
		}
	}

	n := model.Nowcast{
		Location:   "vadodara",
		CapturedAt: captured,
		LeadMin:    0,
		POP:        colorToPOP(col),
		MMPerHr:    0,
	}
	if err := repository.InsertNowcast(ctx, &n); err != nil {
		logger.Error.Println("insert nowcast:", err)
	}

	// record API call metrics
	call := model.IMDAPICall{
		Endpoint:    url,
		Bytes:       resp.ContentLength,
		RequestedAt: time.Now(),
	}
	if err := repository.InsertIMDAPICall(ctx, &call); err != nil {
		logger.Error.Println("repository insert api log:", err)
	} else {
		logger.Info.Printf("IMD API call %s bytes=%d", url, resp.ContentLength)
	}
	return nil
}
