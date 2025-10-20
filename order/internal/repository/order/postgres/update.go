package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres/converter"
)

func (r *repository) Update(ctx context.Context, orderUUID uuid.UUID, order model.Order) error {
	repoOrder := repoConverter.OrderToRepoModel(order)
	repoOrder.ID = orderUUID
	repoOrder.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	const q = `
	UPDATE
		orders
	SET
		"user_id" = :user_id,
		"part_ids" = :part_ids,
		"total_price" = :total_price,
		"transaction_uuid" = :transaction_uuid,
		"payment_method" = :payment_method,
		"status" = :status,
		"updated_at" = :updated_at
	WHERE
		id = :id`

	res, err := sqlx.NamedExecContext(ctx, r.db, q, repoOrder)
	if err != nil {
		return fmt.Errorf("OrderRepository.Update update order error: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("OrderRepository.Update update order rows affected error: %w", err)
	}
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
