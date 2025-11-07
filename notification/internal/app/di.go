package app

import (
	"github.com/pinai4/spaceship-factory/notification/internal/config"
	"github.com/pinai4/spaceship-factory/notification/internal/service"
	notificationService "github.com/pinai4/spaceship-factory/notification/internal/service/notification"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	notificationService service.NotificationService
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) NotificationService() service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = notificationService.NewService()
	}

	return d.notificationService
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
