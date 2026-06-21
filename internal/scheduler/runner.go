package scheduler

import (
	"log"
	"time"
)

func Start(
	rsiScheduler *RSIScheduler,
) {

	ticker := time.NewTicker(
		time.Minute,
	)

	go func() {

		defer ticker.Stop()

		for range ticker.C {

			log.Println(
				"Starting RSI update...",
			)

			err := rsiScheduler.Run()

			if err != nil {

				log.Printf(
					"Scheduler failed: %v",
					err,
				)

				continue
			}

			log.Println(
				"Scheduler completed",
			)
		}
	}()
}
