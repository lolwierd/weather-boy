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

const imdRiverBasinBaseURL = "https://mausam.imd.gov.in/api/basin_qpf_api.php"

type riverBasinResp struct {
	ObjID    string `json:"Obj_Id"`
	Date     string `json:"Date"`
	FMO      string `json:"FMO"`
	Basin    string `json:"Basin"`
	SubBasin string `json:"SubBasin"`
	Area     string `json:"Area"`
	Day1     string `json:"Day1"`
	Day2     string `json:"Day2"`
	Day3     string `json:"Day3"`
	Day4     string `json:"Day4"`
	Day5     string `json:"Day5"`
	AAP      string `json:"AAP"`
}

// FetchRiverBasinOnce fetches river basin data from the IMD API and stores it.
func FetchRiverBasinOnce(ctx context.Context, basinID int) error {
	url := fmt.Sprintf("%s?id=%d", imdRiverBasinBaseURL, basinID)
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
		return fmt.Errorf("imd river basin status %s: %s", resp.Status, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var arr []riverBasinResp
	if err := json.Unmarshal(body, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return fmt.Errorf("empty river basin response")
	}

	for _, r := range arr {
		basinID, err := strconv.Atoi(r.ObjID)
		if err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			return err
		}
		qpf := model.RiverBasinQPF{
			BasinID:  basinID,
			Date:     date,
			FMO:      r.FMO,
			Basin:    r.Basin,
			SubBasin: r.SubBasin,
			Area:     r.Area,
			Day1:     r.Day1,
			Day2:     r.Day2,
			Day3:     r.Day3,
			Day4:     r.Day4,
			Day5:     r.Day5,
			AAP:      r.AAP,
		}
		if err := repository.InsertRiverBasinQPF(ctx, &qpf); err != nil {
			return err
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
