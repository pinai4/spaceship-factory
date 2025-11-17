package client

import "context"

type UserClient interface {
	GetTelegramChat(ctx context.Context, userUUID string) (int64, error)
}
