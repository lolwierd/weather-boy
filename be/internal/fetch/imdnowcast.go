package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/model"
	"github.com/lolwierd/weatherboy/be/internal/repository"
)

const imdNowcastBaseURL = "https://mausam.imd.gov.in/api/nowcast_district_api.php"

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

// districtNowcastResp mirrors the IMD district nowcast response. The API
// returns an array with a single object containing categorical fields and a
// `color` code indicating rainfall likelihood.
type districtNowcastResp struct {
	ObjID           string `json:"Obj_id"`
	Date            string `json:"Date"`
	TOI             string `json:"toi"`
	VUpto           string `json:"vupto"`
	Color           string `json:"color"`
	Message         string `json:"message"`
	PrecipIntensity string `json:"precip_intensity"`
	Cat1            string `json:"cat1"`
	Cat2            string `json:"cat2"`
	Cat3            string `json:"cat3"`
	Cat4            string `json:"cat4"`
	Cat5            string `json:"cat5"`
	Cat6            string `json:"cat6"`
	Cat7            string `json:"cat7"`
	Cat8            string `json:"cat8"`
	Cat9            string `json:"cat9"`
	Cat10           string `json:"cat10"`
	Cat11           string `json:"cat11"`
	Cat12           string `json:"cat12"`
	Cat13           string `json:"cat13"`
	Cat14           string `json:"cat14"`
	Cat15           string `json:"cat15"`
	Cat16           string `json:"cat16"`
	Cat17           string `json:"cat17"`
	Cat18           string `json:"cat18"`
	Cat19           string `json:"cat19"`
}

// colorToPOP converts the IMD color code (1-4) to an approximate probability of
// precipitation value.
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
	loc, ok := config.LocationByName("vadodara")
	if !ok || loc.DistrictID == 0 {
		return fmt.Errorf("district id for vadodara not set")
	}
	url := fmt.Sprintf("%s?id=%d", imdNowcastBaseURL, loc.DistrictID)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var arr []districtNowcastResp
	if err := json.Unmarshal(body, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return fmt.Errorf("empty nowcast response")
	}

	raw := model.NowcastRaw{
		Location:  "vadodara",
		Data:      body,
		FetchedAt: time.Now(),
	}
	if err := repository.InsertNowcastRaw(ctx, &raw); err != nil {
		logger.Error.Println("insert nowcast raw:", err)
	}

	col, err := strconv.Atoi(arr[0].Color)
	if err != nil {
		col = 0
	}

	pi, err := strconv.Atoi(arr[0].PrecipIntensity)
	if err != nil {
		pi = 0
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
		MMPerHr:    bucketToMMPerHr(pi),
	}
	if err := repository.InsertNowcast(ctx, &n); err != nil {
		logger.Error.Println("insert nowcast:", err)
	}

	// store category flags
	cats := []string{
		arr[0].Cat1, arr[0].Cat2, arr[0].Cat3, arr[0].Cat4, arr[0].Cat5,
		arr[0].Cat6, arr[0].Cat7, arr[0].Cat8, arr[0].Cat9, arr[0].Cat10,
		arr[0].Cat11, arr[0].Cat12, arr[0].Cat13, arr[0].Cat14,
		arr[0].Cat15, arr[0].Cat16, arr[0].Cat17, arr[0].Cat18, arr[0].Cat19,
	}
	for i, v := range cats {
		if v == "" {
			continue
		}
		val, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		cat := model.NowcastCategory{
			NowcastID: n.ID,
			Category:  i + 1,
			Value:     int16(val),
		}
		if err := repository.InsertNowcastCategory(ctx, &cat); err != nil {
			logger.Error.Println("insert nowcast category:", err)
		}
	}

	call := model.IMDAPICall{
		Endpoint:    url,
		Bytes:       int64(len(body)),
		RequestedAt: time.Now(),
	}
	if err := repository.InsertIMDAPICall(ctx, &call); err != nil {
		logger.Error.Println("repository insert api log:", err)
	} else {
		logger.Info.Printf("IMD API call %s bytes=%d", url, len(body))
	}
	return nil
}
