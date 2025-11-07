package notification

import def "github.com/pinai4/spaceship-factory/notification/internal/service"

var _ def.NotificationService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
