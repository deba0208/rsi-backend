package scheduler

import (
	"log"

	"github.com/go-co-op/gocron/v2"
)

func Start(
	rsiScheduler *RSIScheduler,
) {

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Printf("Failed to create gocron scheduler: %v", err)
		return
	}

	_, err = s.NewJob(
		gocron.CronJob(
			"0 16 * * 1-5",
			false,
		),
		gocron.NewTask(
			func() {
				log.Println("Starting RSI update...")

				err := rsiScheduler.Run()
				if err != nil {
					log.Printf("Scheduler failed: %v", err)
					return
				}

				log.Println("Scheduler completed")
			},
		),
	)

	if err != nil {
		log.Printf("Failed to schedule RSI job: %v", err)
		return
	}

	// Start the scheduler in the background
	s.Start()
}
