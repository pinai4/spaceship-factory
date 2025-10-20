package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pinai4/spaceship-factory/order/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	repoOrder := repoConverter.OrderToRepoModel(order)
	repoOrder.CreatedAt = time.Now()

	const q = `
	INSERT INTO orders
		(id, user_id, part_ids, total_price, transaction_uuid, payment_method, status, created_at, updated_at)
	VALUES
		(:id, :user_id, :part_ids, :total_price, :transaction_uuid, :payment_method, :status, :created_at, :updated_at)`

	_, err := sqlx.NamedExecContext(ctx, r.db, q, repoOrder)
	if err != nil {
		return fmt.Errorf("OrderRepository.Create create order error: %w", err)
	}

	return nil
}
