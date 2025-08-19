package cronjob

import (
	"context"
	"log"
	"time"

	"murim-helper/internal/usecase"

	"github.com/robfig/cron/v3"
)

func StartCronJobs(uc usecase.ScheduleUsecase) {
	c := cron.New()
	// Run every day at midnight
	c.AddFunc("0 0 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("[CRON] Processing repeating schedules...")
		if err := uc.ProcessRepeatingSchedules(ctx); err != nil {
			log.Printf("[CRON] Error: %v", err)
		} else {
			log.Println("[CRON] Repeating schedules processed successfully")
		}
	})
	c.Start()
}
