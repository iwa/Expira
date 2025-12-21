package cron

import (
	"fmt"
	"time"

	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils"
)

// runDailyUpdate performs the scheduled domain update
func RunDailyUpdate(store *state.DomainStore, config *state.Config) {
	fmt.Println("[INFO] Daily domains refresh cron job triggered at", time.Now().Format("2006-12-20 15:04:05"))

	utils.UpdateDomains(store)
	utils.Notify(store, config)
	utils.ReportStatusInConsole(store)

	fmt.Println("[INFO] Daily domain update completed")
}
