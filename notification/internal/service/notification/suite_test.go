//go:build unit || !integration

package notification_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/notification/internal/service"
	"github.com/pinai4/spaceship-factory/notification/internal/service/notification"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service service.NotificationService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.service = notification.NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
