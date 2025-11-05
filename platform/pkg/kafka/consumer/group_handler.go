package consumer

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

// MessageHandler — message handler.
type MessageHandler func(ctx context.Context, msg Message) error

// Middleware — middleware function for additional processing.
type Middleware func(next MessageHandler) MessageHandler

// groupHandler — wrapper for sarama.ConsumerGroupHandler.
type groupHandler struct {
	handler MessageHandler
	logger  Logger
}

// NewGroupHandler creates a new groupHandler with a middleware chain.
func NewGroupHandler(handler MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	// Apply the middleware chain
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (g *groupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Consumer message channel closed")
				return nil
			}

			msg := Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Error(session.Context(), "Consumer handler error", logger.Error(err))
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Consumer session context done")
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte)
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
