package cron

import (
	"time"

	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils"
)

// StartCronLoop starts an hourly cron job that runs domain updates at midnight.
// It uses dependency injection to access the domain store and configuration.
func StartCronLoop(store *state.DomainStore, config *state.Config) {
	println("[INFO] Starting cron job...")

	ticker := time.NewTicker(time.Hour)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if checkMidnight(t) {
					println("[INFO] Daily domains refresh cron job triggered at", t.Format("2006-01-02 15:04:05"))

					utils.UpdateDomains(store)

					utils.Notify(store, config)

					utils.ReportStatusInConsole(store)
				}
			}
		}
	}()
}

// Check if the current hour is 0 (midnight)
func checkMidnight(t time.Time) bool {
	return t.Hour() == 0
}
