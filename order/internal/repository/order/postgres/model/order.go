package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID      `db:"id"`
	UserID        uuid.UUID      `db:"user_id"`
	PartIDs       UUIDArray      `db:"part_ids"`
	TotalPrice    float64        `db:"total_price"`
	TransactionID uuid.NullUUID  `db:"transaction_uuid"`
	PaymentMethod sql.NullString `db:"payment_method"`
	Status        string         `db:"status"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
}
