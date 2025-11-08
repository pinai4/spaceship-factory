package notification

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pinai4/spaceship-factory/notification/internal/model"
)

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(event)
	if err != nil {
		return fmt.Errorf("NotificationService.SendOrderPaidNotification failed to build message: %w", err)
	}

	if err := s.telegramClient.SendMessage(ctx, s.chatID, message); err != nil {
		return fmt.Errorf("NotificationService.SendOrderPaidNotification failed to send message to telegram: %w", err)
	}

	return nil
}

func (s *service) buildOrderPaidMessage(event model.OrderPaidEvent) (string, error) {
	data := struct {
		OrderUUID       string
		UserUUID        string
		PaymentMethod   string
		TransactionUUID string
	}{
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}

	var buf bytes.Buffer
	err := s.orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
