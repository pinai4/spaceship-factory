package payment

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context) (string, error) {
	return uuid.NewString(), nil
}
