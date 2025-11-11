package postgres

import (
	"github.com/jmoiron/sqlx"

	def "github.com/pinai4/spaceship-factory/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}
