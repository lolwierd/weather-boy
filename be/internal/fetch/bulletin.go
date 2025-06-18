package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/lolwierd/weatherboy/be/internal/logger"
	"github.com/lolwierd/weatherboy/be/internal/model"
	"github.com/lolwierd/weatherboy/be/internal/repository"
)

// FetchBulletinOnce downloads today's Gujarat bulletin PDF and stores it.
func FetchBulletinOnce(ctx context.Context) error {
	const url = "https://mausam.imd.gov.in/Backend/uploads/REVAMP_GMFC/mcdata/gujarat.pdf"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}

	config.LoadEnv()
	dir := filepath.Join(config.DataDir, "pdf")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	dateStr := time.Now().Format("2006-01-02")
	path := filepath.Join(dir, fmt.Sprintf("%s-gujarat.pdf", dateStr))

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	br := model.BulletinRaw{Path: path, FetchedAt: time.Now()}
	if err := repository.InsertBulletinRaw(ctx, &br); err != nil {
		logger.Error.Println("repository insert bulletin raw:", err)
	}
	return nil
}
