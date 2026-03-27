package service

import (
	"log"
	"time"
)

const dailyAIReportHour = 21

func StartAIReportDailyScheduler() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("load Asia/Shanghai location failed, fallback to local: %v", err)
		location = time.Local
	}

	adminService := &AdminService{}

	// Catch up once if service starts after 21:00 and today's system report is still missing.
	now := time.Now().In(location)
	if now.Hour() >= dailyAIReportHour {
		if err := adminService.GenerateDailyAIReportIfNeeded(location); err != nil {
			log.Printf("daily AI report catch-up failed: %v", err)
		}
	}

	go func() {
		for {
			now := time.Now().In(location)
			nextRun := nextDailyAIReportRun(now, location)
			waitDuration := time.Until(nextRun)
			if waitDuration < 0 {
				waitDuration = 0
			}

			log.Printf("daily AI report scheduler armed, next run at %s", nextRun.Format("2006-01-02 15:04:05"))
			timer := time.NewTimer(waitDuration)
			<-timer.C

			if err := adminService.GenerateDailyAIReportIfNeeded(location); err != nil {
				log.Printf("daily AI report generate failed: %v", err)
			}
		}
	}()
}

func nextDailyAIReportRun(now time.Time, location *time.Location) time.Time {
	runAt := time.Date(now.Year(), now.Month(), now.Day(), dailyAIReportHour, 0, 0, 0, location)
	if !now.Before(runAt) {
		runAt = runAt.AddDate(0, 0, 1)
	}
	return runAt
}
