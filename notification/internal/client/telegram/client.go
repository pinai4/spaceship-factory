package telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

type client struct {
	bot *bot.Bot
}

// NewClient creates a new client for the Telegram Bot API
func NewClient(bot *bot.Bot) *client {
	return &client{
		bot: bot,
	}
}

// SendMessage sends a message to the specified chat
func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	})
	if err != nil {
		return err
	}

	return nil
}
