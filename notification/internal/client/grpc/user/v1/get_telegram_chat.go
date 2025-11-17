package v1

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

func (c *client) GetTelegramChat(ctx context.Context, userUUID string) (int64, error) {
	resp, err := c.generatedClient.GetUser(ctx, &userV1.GetUserRequest{
		UserUuid: userUUID,
	})
	if err != nil {
		return 0, fmt.Errorf("client API call error: %w", err)
	}

	for _, n := range resp.GetUser().GetInfo().GetNotificationMethods() {
		if n.ProviderName == telegramNotificationMethodKey {
			telegramChat, err := c.parseInt64(n.Target)
			if err != nil {
				return 0, model.ErrUserTelegramChatInvalid
			}

			return telegramChat, nil
		}
	}

	return 0, model.ErrUserTelegramChatNotSpecified
}

func (c *client) parseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	return strconv.ParseInt(s, 10, 64)
}
