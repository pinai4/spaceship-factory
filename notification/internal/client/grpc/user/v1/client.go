package v1

import (
	def "github.com/pinai4/spaceship-factory/notification/internal/client"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

const telegramNotificationMethodKey = "telegram"

var _ def.UserClient = (*client)(nil)

type client struct {
	generatedClient userV1.UserServiceClient
}

func NewClient(generatedClient userV1.UserServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
