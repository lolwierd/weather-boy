//go:build ignore
// +build ignore

package fetch

// Legacy radar implementation is disabled for now.
// The application builds with stub implementation in `radar_stub.go`.


import (
    "context"

    "github.com/lolwierd/weatherboy/be/internal/config"
    "github.com/lolwierd/weatherboy/be/internal/logger"
)

// FetchRadarOnce is currently stubbed out as obtaining radar data is not yet implemented.
// It logs the skip and returns nil so that callers treat it as non-fatal.
func FetchRadarOnce(ctx context.Context, loc config.Location) error {
    logger.Info.Println("radar fetch stub: skipping actual download and processing for", loc.Name)
    return nil
}

import (
	"context"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/lolwierd/weatherboy/be/internal/logger"
)



// FetchRadarOnce downloads the latest radar image for a given location and stores it.
func FetchRadarOnce(ctx context.Context, loc config.Location) error {
    // Stubbed implementation: radar download not yet implemented; skip gracefully
    // We simply log and return nil. Once the radar functionality is implemented,
    // this stub can be replaced with the real downloader + parser.
    logger.Info.Println("radar fetch stubbed: skipping actual download and processing for", loc.Name)
    return nil
}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}

	config.LoadEnv()
	dir := filepath.Join(config.DataDir, "radar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	timestampStr := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.png", timestampStr, loc.RadarCodes[0])
	path := filepath.Join(dir, fileName)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	maxDBZ, err := parse.ParseRadarImage(path, 40)
	if err != nil {
		return err
	}

	capturedAt := time.Now()
	// TODO: derive the actual capture time from metadata once available

	radar := model.RadarSnapshot{
		Location:   loc.Name,
		MaxDBZ:     float64(maxDBZ),
		CapturedAt: capturedAt,
	}
	if err := repository.InsertRadarSnapshot(ctx, &radar); err != nil {
		return err
	}

	call := model.IMDAPICall{
		Endpoint:    url,
		Bytes:       n,
		RequestedAt: time.Now(),
	}
	if err := repository.InsertIMDAPICall(ctx, &call); err != nil {
		logger.Error.Println("repository insert api log:", err)
	} else {
		logger.Info.Printf("IMD API call %s bytes=%d location=%s", url, n, loc.RadarCodes[0])
	}
	return nil
}
