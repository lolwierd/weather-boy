package fetch

import (
    "context"

    "github.com/lolwierd/weatherboy/be/internal/config"
    "github.com/lolwierd/weatherboy/be/internal/logger"
)

// FetchRadarOnce is currently not implemented. It is stubbed so that the
// application can run without radar data. When radar fetching is implemented
// (see `radar_full` build tag implementation), this stub will be replaced.
func FetchRadarOnce(ctx context.Context, loc config.Location) error {
    logger.Info.Println("radar fetch stub: skipping actual download and processing for", loc.Name)
    return nil
}
