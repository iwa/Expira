package utils

import (
	"fmt"
	"time"

	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils/providers"
)

// Notify checks all domains and sends notifications if they are approaching expiry.
// It uses the provided store for domain data and config for notification settings.
func Notify(store *state.DomainStore, config *state.Config) {
	println("[INFO] Sending notifications...")

	domains := store.GetAllDomains()

	for domain, domainData := range domains {
		daysUntil, shouldNotify := checkDaysForNotification(domainData.ExpiryDate, config.NotificationDays)

		if shouldNotify {
			if config.TelegramNotification && (config.TelegramChatID != "" && config.TelegramToken != "") {
				message := fmt.Sprintf("<b>⚠️ Domain %s will expire in %d days </b>\nExpiry date: <code>%s</code>", domain, daysUntil, domainData.ExpiryDate.Format("2006-01-02 15:04:05"))

				err := providers.SendTelegramMessage(config, message)
				if err != nil {
					println("[ERROR] Failed to send notification for domain", domain, ":", err)
				}
			}

			if config.DiscordNotification && config.DiscordWebhookURL != "" {
				message := fmt.Sprintf("**⚠️ Domain %s will expire in %d days**\nExpiry date: `%s`", domain, daysUntil, domainData.ExpiryDate.Format("2006-01-02 15:04:05"))

				err := providers.SendDiscordMessage(config, message)
				if err != nil {
					println("[ERROR] Failed to send notification for domain", domain, ":", err)
				}
			}

			if config.NtfyNotification && config.NtfyURL != "" {
				message := fmt.Sprintf("Domain %s will expire in %d days \nExpiry date: %s", domain, daysUntil, domainData.ExpiryDate.Format("2006-01-02 15:04:05"))

				err := providers.SendNtfyMessage(config, message)
				if err != nil {
					println("[ERROR] Failed to send notification for domain", domain, ":", err)
				}
			}
		}
	}
}

func checkDaysForNotification(expriyDate time.Time, notificationDays []int) (int, bool) {
	currentTime := time.Now()
	daysLeft := int(expriyDate.Sub(currentTime).Hours()/24) + 1 // Add 1 to include the current day

	if daysLeft < 0 {
		println("[WARN] Domain", expriyDate.Format("2006-01-02 15:04:05"), "has already expired.")
		return 0, false
	}

	for _, days := range notificationDays {
		if daysLeft == days {
			println("[INFO] Domain expiry is exactly", days, "days away:", expriyDate.Format("2006-01-02 15:04:05"))
			return days, true
		}
	}

	return 0, false
}
