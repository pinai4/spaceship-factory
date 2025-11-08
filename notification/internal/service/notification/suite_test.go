//go:build unit || !integration

package notification_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/pinai4/spaceship-factory/notification/internal/client/mocks"
	"github.com/pinai4/spaceship-factory/notification/internal/service"
	"github.com/pinai4/spaceship-factory/notification/internal/service/notification"
)

const chatID int64 = 1

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	telegramClient *clientMocks.TelegramClient

	service service.NotificationService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.telegramClient = clientMocks.NewTelegramClient(s.T())

	s.service = notification.NewService(s.telegramClient, chatID)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
