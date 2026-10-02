
package scheduler

import (
	"context"
	"log"
	"time"
	 _ "time/tzdata"
	"github.com/go-co-op/gocron/v2"
)

func Start(rsiScheduler *RSIScheduler) error {
	// Use Indian Standard Time, regardless of server timezone.
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}

	s, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
	)
	if err != nil {
		return err
	}

	// Run every weekday at 4:00 PM IST.
	_, err = s.NewJob(
		gocron.CronJob("0 16 * * 1-5", false),
		gocron.NewTask(func() {
			log.Println("Scheduled RSI update triggered")

			if err := rsiScheduler.Run(); err != nil {
				log.Printf("Scheduled RSI update failed: %v", err)
			}
		}),
	)
	if err != nil {
		return err
	}

	// Start the cron scheduler.
	s.Start()

	log.Println("RSI scheduler started: weekdays at 4:00 PM Asia/Kolkata")

	// Run once immediately after application startup.
	go func() {
		log.Println("Running initial RSI update")

		if err := rsiScheduler.Run(); err != nil {
			log.Printf("Initial RSI update failed: %v", err)
		}
	}()

	// Keep the function signature independent of a blocking wait.
	_ = context.Background()

	return nil
}