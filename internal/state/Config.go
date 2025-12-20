package state

// Config holds immutable configuration loaded from environment variables at startup.
// This configuration does not change during the application lifecycle.
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
