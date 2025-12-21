package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iwa/Expira/internal/state"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

// SendDiscordMessage sends a notification message via Discord webhook.
// It uses configuration from the provided Config instance.
func SendDiscordMessage(config *state.Config, message string) error {
	payload := DiscordMessage{
		Content: message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.DiscordWebhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("[INFO] Response from Discord API:", string(responseBody))

	return nil
}
