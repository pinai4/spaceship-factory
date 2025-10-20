package postgres

import (
	"github.com/jmoiron/sqlx"

	def "github.com/pinai4/spaceship-factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}
