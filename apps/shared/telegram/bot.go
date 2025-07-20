package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/silasstoffel/invest-tracker/config"
)

type TelegramBot struct {
	Token  string
	ChatID string
}

func NewTelegramBot(config *config.Config) *TelegramBot {
	return &TelegramBot{
		Token:  config.TelegramConfig.Token,
		ChatID: config.TelegramConfig.ChatId,
	}
}

func (T *TelegramBot) SendMessage(message string) error {
	chatID := T.ChatID
	botToken := T.Token

	if chatID == "" {
		return fmt.Errorf("telegram chat ID is not set")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	body := map[string]interface{}{
		"chat_id":              chatID,
		"text":                 message,
		"parse_mode":           "Markdown",
		"disable_notification": false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to encode JSON body: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to send message to Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Telegram, status code: %d", resp.StatusCode)
	}

	return nil
}
