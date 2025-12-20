package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iwa/Expira/internal/state"
)

type TelegramMessage struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
	ProtectContent      bool   `json:"protect_content"`
}

// SendTelegramMessage sends a notification message via Telegram API.
// It uses configuration from the provided Config instance.
func SendTelegramMessage(config *state.Config, message string) error {
	payload := TelegramMessage{
		ChatID:              config.TelegramChatID,
		Text:                message,
		ParseMode:           "HTML",
		DisableNotification: true,
		ProtectContent:      false,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("[INFO] Response from Telegram API:", string(responseBody))

	return nil
}
