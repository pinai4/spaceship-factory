package app

import (
	"fmt"

	"github.com/go-telegram/bot"

	"github.com/pinai4/spaceship-factory/notification/internal/client"
	"github.com/pinai4/spaceship-factory/notification/internal/client/telegram"
	"github.com/pinai4/spaceship-factory/notification/internal/config"
	"github.com/pinai4/spaceship-factory/notification/internal/service"
	notificationService "github.com/pinai4/spaceship-factory/notification/internal/service/notification"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	notificationService service.NotificationService

	telegramClient client.TelegramClient
	telegramBot    *bot.Bot
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) NotificationService() service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = notificationService.NewService(d.TelegramClient(), d.Config().TelegramBot.ChatID())
	}

	return d.notificationService
}

func (d *diContainer) TelegramClient() client.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegram.NewClient(d.TelegramBot())
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(d.Config().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
