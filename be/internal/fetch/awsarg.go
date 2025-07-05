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

const imdAWSARGBaseURL = "https://city.imd.gov.in/api/aws_data_api.php"

type awsArgResp struct {
	ID            string `json:"ID"`
	CallSign      string `json:"CALL_SIGN"`
	District      string `json:"DISTRICT"`
	State         string `json:"STATE"`
	Station       string `json:"STATION"`
	Date          string `json:"DATE"`
	Time          string `json:"TIME"`
	CurrTemp      string `json:"CURR_TEMP"`
	DewPointTemp  string `json:"DEW_POINT_TEMP"`
	RH            string `json:"RH"`
	WindDirection string `json:"WIND_DIRECTION"`
	WindSpeed     string `json:"WIND_SPEED"`
	MSLP          string `json:"MSLP"`
	MinTemp       string `json:"MIN_TEMP"`
	MaxTemp       string `json:"MAX_TEMP"`
	Latitude      string `json:"Latitude"`
	Longitude     string `json:"Longitude"`
	WeatherCode   string `json:"WEATHER_CODE"`
	Nebulosity    string `json:"NEBULOSITY"`
	FeelLike      string `json:"Feel Like"`
	RainfallSel   string `json:"RAINFALL_SEL"`
	Rainfall      string `json:"RAINFALL"`
}

// FetchAWSARGOnce fetches AWS/ARG data from the IMD API and stores it.
func FetchAWSARGOnce(ctx context.Context, stationID string) error {
	url := fmt.Sprintf("%s?id=%s", imdAWSARGBaseURL, stationID)
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
		return fmt.Errorf("imd aws/arg status %s: %s", resp.Status, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var arr []awsArgResp
	if err := json.Unmarshal(body, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return fmt.Errorf("empty aws/arg response")
	}

	// Assuming the API returns a single object in the array for the station
	awsArgData := arr[0]

	date, err := time.Parse("2006-01-02", awsArgData.Date)
	if err != nil {
		return fmt.Errorf("invalid date in AWS/ARG response: %w", err)
	}
	timeOfDay, err := time.Parse("15:04:05", awsArgData.Time)
	if err != nil {
		return fmt.Errorf("invalid time in AWS/ARG response: %w", err)
	}

	currentTemp, err := strconv.ParseFloat(awsArgData.CurrTemp, 64)
	if err != nil {
		return fmt.Errorf("invalid current_temp in AWS/ARG response: %w", err)
	}
	dewPointTemp, err := strconv.ParseFloat(awsArgData.DewPointTemp, 64)
	if err != nil {
		return fmt.Errorf("invalid dew_point_temp in AWS/ARG response: %w", err)
	}
	rh, err := strconv.ParseFloat(awsArgData.RH, 64)
	if err != nil {
		return fmt.Errorf("invalid RH in AWS/ARG response: %w", err)
	}
	windDirection, err := strconv.ParseFloat(awsArgData.WindDirection, 64)
	if err != nil {
		return fmt.Errorf("invalid wind_direction in AWS/ARG response: %w", err)
	}
	windSpeed, err := strconv.ParseFloat(awsArgData.WindSpeed, 64)
	if err != nil {
		return fmt.Errorf("invalid wind_speed in AWS/ARG response: %w", err)
	}
	mslp, err := strconv.ParseFloat(awsArgData.MSLP, 64)
	if err != nil {
		return fmt.Errorf("invalid MSLP in AWS/ARG response: %w", err)
	}
	minTemp, err := strconv.ParseFloat(awsArgData.MinTemp, 64)
	if err != nil {
		return fmt.Errorf("invalid min_temp in AWS/ARG response: %w", err)
	}
	maxTemp, err := strconv.ParseFloat(awsArgData.MaxTemp, 64)
	if err != nil {
		return fmt.Errorf("invalid max_temp in AWS/ARG response: %w", err)
	}
	latitude, err := strconv.ParseFloat(awsArgData.Latitude, 64)
	if err != nil {
		return fmt.Errorf("invalid latitude in AWS/ARG response: %w", err)
	}
	longitude, err := strconv.ParseFloat(awsArgData.Longitude, 64)
	if err != nil {
		return fmt.Errorf("invalid longitude in AWS/ARG response: %w", err)
	}
	nebulosity, err := strconv.ParseFloat(awsArgData.Nebulosity, 64)
	if err != nil {
		return fmt.Errorf("invalid nebulosity in AWS/ARG response: %w", err)
	}
	rainfall, err := strconv.ParseFloat(awsArgData.Rainfall, 64)
	if err != nil {
		return fmt.Errorf("invalid rainfall in AWS/ARG response: %w", err)
	}

	awsArg := model.AWSARG{
		StationID:     awsArgData.ID,
		CallSign:      awsArgData.CallSign,
		District:      awsArgData.District,
		State:         awsArgData.State,
		StationName:   awsArgData.Station,
		Date:          date,
		Time:          timeOfDay,
		CurrentTemp:   currentTemp,
		DewPointTemp:  dewPointTemp,
		RH:            rh,
		WindDirection: windDirection,
		WindSpeed:     windSpeed,
		MSLP:          mslp,
		MinTemp:       minTemp,
		MaxTemp:       maxTemp,
		Latitude:      latitude,
		Longitude:     longitude,
		WeatherCode:   awsArgData.WeatherCode,
		Nebulosity:    nebulosity,
		FeelLike:      0, // This field is not directly available in the API response, setting to 0 for now
		RainfallSel:   awsArgData.RainfallSel,
		Rainfall:      rainfall,
	}

	if err := repository.InsertAWSARG(ctx, &awsArg); err != nil {
		return err
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
