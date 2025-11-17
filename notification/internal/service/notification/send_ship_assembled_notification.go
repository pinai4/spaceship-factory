package notification

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	message, err := s.buildShipAssembledMessage(event)
	if err != nil {
		return fmt.Errorf("NotificationService.SendShipAssembledNotification failed to build message: %w", err)
	}

	chatID, err := s.userClient.GetTelegramChat(ctx, event.UserUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserTelegramChatInvalid) || errors.Is(err, model.ErrUserTelegramChatNotSpecified) {
			return nil
		}
		return fmt.Errorf("NotificationService.SendShipAssembledNotification failed to get user telegram chat: %w", err)
	}

	if err := s.telegramClient.SendMessage(ctx, chatID, message); err != nil {
		return fmt.Errorf("NotificationService.SendShipAssembledNotification failed to send message to telegram: %w", err)
	}

	return nil
}

func (s *service) buildShipAssembledMessage(event model.ShipAssembledEvent) (string, error) {
	data := struct {
		OrderUUID    string
		UserUUID     string
		BuildTimeSec int64
	}{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := s.shipAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
