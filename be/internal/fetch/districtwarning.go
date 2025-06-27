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

const imdDistrictWarningBaseURL = "https://mausam.imd.gov.in/api/warnings_district_api.php"

type districtWarningResp struct {
	ObjID    string `json:"Obj_id"`
	Date     string `json:"Date"`
	UTC      string `json:"UTC"`
	District string `json:"District"`
	Day1     string `json:"Day_1"`
	Day2     string `json:"Day_2"`
	Day3     string `json:"Day_3"`
	Day4     string `json:"Day_4"`
	Day5     string `json:"Day_5"`
	Day1Color string `json:"Day1_Color"`
	Day2Color string `json:"Day2_Color"`
	Day3Color string `json:"Day3_Color"`
	Day4Color string `json:"Day4_Color"`
	Day5Color string `json:"Day5_Color"`
}

// FetchDistrictWarnings fetches district-wise warning data from the IMD API for Vadodara and stores it.
func FetchDistrictWarnings(ctx context.Context) error {
	loc, ok := config.LocationByName("vadodara")
	if !ok || loc.DistrictID == 0 {
		return fmt.Errorf("district id for vadodara not set")
	}
	url := fmt.Sprintf("%s?id=%d", imdDistrictWarningBaseURL, loc.DistrictID)
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
		return fmt.Errorf("imd district warning status %s: %s", resp.Status, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var arr []districtWarningResp
	if err := json.Unmarshal(body, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return fmt.Errorf("empty district warning response")
	}

	raw := model.DistrictWarningRaw{
		Location:  "vadodara",
		Data:      body,
		FetchedAt: time.Now(),
	}
	if err := repository.InsertDistrictWarningRaw(ctx, &raw); err != nil {
		logger.Error.Println("insert district warning raw:", err)
	}

	// Assuming the API returns a single object in the array for the district
	dwResp := arr[0]

	issuedAt := time.Now()
	if dwResp.Date != "" && dwResp.UTC != "" {
		t, err := time.Parse("2006-01-02 15:04:05", dwResp.Date+" "+dwResp.UTC)
		if err == nil {
			issuedAt = t
		}
	}

	dw := model.DistrictWarning{
		Location:    "vadodara",
		IssuedAt:    issuedAt,
		Day1Warning: dwResp.Day1,
		Day2Warning: dwResp.Day2,
		Day3Warning: dwResp.Day3,
		Day4Warning: dwResp.Day4,
		Day5Warning: dwResp.Day5,
		Day1Color:   dwResp.Day1Color,
		Day2Color:   dwResp.Day2Color,
		Day3Color:   dwResp.Day3Color,
		Day4Color:   dwResp.Day4Color,
		Day5Color:   dwResp.Day5Color,
	}
	if err := repository.InsertDistrictWarning(ctx, &dw); err != nil {
		logger.Error.Println("insert district warning:", err)
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
