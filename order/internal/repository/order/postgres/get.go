package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres/converter"
	repoModel "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres/model"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	const q = `
	SELECT
	    id, user_id, part_ids, total_price, transaction_uuid, payment_method, status, created_at, updated_at
	FROM
		orders
	WHERE
		id = $1`

	var repoOrder repoModel.Order
	if err := r.db.GetContext(ctx, &repoOrder, q, orderUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("OrderRepository.Get get order error: %w", err)
	}

	return repoConverter.OrderToModel(repoOrder), nil
}
