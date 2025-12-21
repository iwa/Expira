package state

// Config holds configuration loaded from environment variables at startup.
// This configuration should never change during the application lifecycle.
type Config struct {
	NotificationDays []int

	TelegramNotification bool
	TelegramChatID       string
	TelegramToken        string

	DiscordNotification bool
	DiscordWebhookURL   string

	NtfyNotification bool
	NtfyURL          string
}

// NewConfig creates a new Config instance with default values
func NewConfig() *Config {
	return &Config{
		NotificationDays: []int{},
	}
}
