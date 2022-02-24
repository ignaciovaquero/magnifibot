package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
)

const (
	telegramParseMode = "MarkdownV2"
)

func (m *Magnifibot) SendTelegram(ctx context.Context, chatID string, message string) (int, error) {
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting chat ID from string to integer: %w", err)
	}
	telegramMessage, err := m.TelegramAPI.SendMessage(&telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: id},
		ParseMode: telegramParseMode,
		Text:      message,
	})
	if err != nil {
		return 0, fmt.Errorf("error sending telegram message: %w", err)
	}
	return telegramMessage.MessageID, nil
}
