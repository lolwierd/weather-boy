package scheduler

import (
	"context"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/lolwierd/weatherboy/be/internal/fetch"
	"github.com/lolwierd/weatherboy/be/internal/logger"
)

var c *cron.Cron

// Start initializes and starts all cron jobs.
func Start() {
	if c != nil {
		return
	}
	// Cron in IST
	loc, _ := time.LoadLocation("Asia/Kolkata")
	c = cron.New(cron.WithLocation(loc))

	// Bulletin every day 18:30 IST
	_, err := c.AddFunc("CRON_TZ=Asia/Kolkata 30 18 * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: bulletin fetch")
		if err := fetch.FetchBulletinOnce(context.Background()); err != nil {
			logger.Error.Println("fetch bulletin:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add bulletin:", err)
	}

	// Nowcast every 15 minutes with jitter
	_, err = c.AddFunc("CRON_TZ=Asia/Kolkata */15 * * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: nowcast fetch")
		if err := fetch.FetchIMDNowcast(context.Background()); err != nil {
			logger.Error.Println("fetch nowcast:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add nowcast:", err)
	}

	// District warnings every day 18:00 IST
	_, err = c.AddFunc("CRON_TZ=Asia/Kolkata 0 18 * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: district warning fetch")
		if err := fetch.FetchDistrictWarnings(context.Background()); err != nil {
			logger.Error.Println("fetch district warning:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add district warning:", err)
	}

	// Radar every 5 minutes with jitter
	_, err = c.AddFunc("CRON_TZ=Asia/Kolkata */5 * * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: radar fetch")
		loc, ok := config.LocationByName("vadodara")
		if !ok {
			logger.Error.Println("location vadodara not found for radar fetch")
			return
		}
		if err := fetch.FetchRadarOnce(context.Background(), loc); err != nil {
			logger.Error.Println("fetch radar:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add radar:", err)
	}

	// River basin every day 19:00 IST
	_, err = c.AddFunc("CRON_TZ=Asia/Kolkata 0 19 * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: river basin fetch")
		if err := fetch.FetchRiverBasinOnce(context.Background(), 1); err != nil {
			logger.Error.Println("fetch river basin:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add river basin:", err)
	}

	// AWS/ARG every 30 minutes with jitter
	_, err = c.AddFunc("CRON_TZ=Asia/Kolkata */30 * * * *", func() {
		// jitter +/-30s
		jitter := time.Duration(rand.Intn(60)-30) * time.Second
		time.Sleep(jitter)
		logger.Info.Println("cron: aws/arg fetch")
		if err := fetch.FetchAWSARGOnce(context.Background(), "NDL"); err != nil {
			logger.Error.Println("fetch aws/arg:", err)
		}
	})
	if err != nil {
		logger.Error.Println("cron add aws/arg:", err)
	}

	c.Start()
	for _, e := range c.Entries() {
		logger.Info.Printf("cron next run %s\n", e.Next.Format(time.RFC3339))
	}

	// run all jobs once at startup
	go func() {
		logger.Info.Println("initial bulletin fetch")
		if err := fetch.FetchBulletinOnce(context.Background()); err != nil {
			logger.Error.Println("fetch bulletin:", err)
		}
	}()

	go func() {
		logger.Info.Println("initial nowcast fetch")
		if err := fetch.FetchIMDNowcast(context.Background()); err != nil {
			logger.Error.Println("fetch nowcast:", err)
		}
	}()

	go func() {
		logger.Info.Println("initial district warning fetch")
		if err := fetch.FetchDistrictWarnings(context.Background()); err != nil {
			logger.Error.Println("fetch district warning:", err)
		}
	}()

	go func() {
		logger.Info.Println("initial radar fetch")
		loc, ok := config.LocationByName("vadodara")
		if !ok {
			logger.Error.Println("location vadodara not found for initial radar fetch")
			return
		}
		if err := fetch.FetchRadarOnce(context.Background(), loc); err != nil {
			logger.Error.Println("fetch radar:", err)
		}
	}()

	go func() {
		logger.Info.Println("initial river basin fetch")
		if err := fetch.FetchRiverBasinOnce(context.Background(), 1); err != nil {
			logger.Error.Println("fetch river basin:", err)
		}
	}()

	go func() {
		logger.Info.Println("initial aws/arg fetch")
		if err := fetch.FetchAWSARGOnce(context.Background(), "NDL"); err != nil {
			logger.Error.Println("fetch aws/arg:", err)
		}
	}()
}
