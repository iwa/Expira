package utils

import (
	"fmt"

	"github.com/iwa/Expira/internal/state"
	"github.com/iwa/Expira/internal/utils/providers"
)

// Notify checks all domains and sends notifications if they are approaching expiry.
// It uses the provided store for domain data and config for notification settings.
func Notify(store *state.DomainStore, config *state.Config) {
	println("[INFO] Sending notifications...")

	domains := store.GetAllDomains()

	for domain, domainData := range domains {
		daysUntil, shouldNotify := checkDaysForNotification(domainData, config.NotificationDays)

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

func checkDaysForNotification(domain state.Domain, notificationDays []int) (int, bool) {
	daysLeft := domain.GetDaysUntilExpiry()

	if daysLeft < 0 {
		println("[WARN] Domain", domain.Name, "has already expired.", domain.ExpiryDate.Format("2006-01-02 15:04:05"))
		return 0, false
	}

	for _, days := range notificationDays {
		if daysLeft == days {
			println("[INFO] Domain expiry is exactly", days, "days away:", domain.ExpiryDate.Format("2006-01-02 15:04:05"))
			return days, true
		}
	}

	return 0, false
}
