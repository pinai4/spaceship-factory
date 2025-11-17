package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pinai4/spaceship-factory/notification/internal/client"
	userV1Client "github.com/pinai4/spaceship-factory/notification/internal/client/grpc/user/v1"
	"github.com/pinai4/spaceship-factory/notification/internal/client/telegram"
	"github.com/pinai4/spaceship-factory/notification/internal/config"
	"github.com/pinai4/spaceship-factory/notification/internal/service"
	notificationService "github.com/pinai4/spaceship-factory/notification/internal/service/notification"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	notificationService service.NotificationService

	telegramClient client.TelegramClient
	telegramBot    *bot.Bot

	userV1ClientGRPC userV1.UserServiceClient
	userV1Client     client.UserClient
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) NotificationService() service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = notificationService.NewService(d.TelegramClient(), d.UserV1Client())
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

func (d *diContainer) UserV1ClientGRPC() userV1.UserServiceClient {
	if d.userV1ClientGRPC == nil {
		conn, err := grpc.NewClient(
			d.Config().IAMGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to IAM Service (GRPC): %s\n", err.Error()))
		}
		d.closer.AddNamed("IAMClient connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.userV1ClientGRPC = userV1.NewUserServiceClient(conn)
	}

	return d.userV1ClientGRPC
}

func (d *diContainer) UserV1Client() client.UserClient {
	if d.userV1Client == nil {
		d.userV1Client = userV1Client.NewClient(d.UserV1ClientGRPC())
	}

	return d.userV1Client
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
