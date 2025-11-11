package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres/converter"
)

func (r *repository) Create(ctx context.Context, user model.User) error {
	repoUser := repoConverter.UserToRepoModel(user)

	const q = `
	INSERT INTO users
		(id, login, email, notification_methods, password_hash, created_at, updated_at)
	VALUES
		(:id, :login, :email, :notification_methods, :password_hash, :created_at, :updated_at)`

	_, err := sqlx.NamedExecContext(ctx, r.db, q, repoUser)
	if err != nil {
		if r.isUniqueViolation(err) {
			return model.ErrUserAlreadyExists
		}
		return fmt.Errorf("UserRepository.Create create user error: %w", err)
	}

	return nil
}

func (r *repository) isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
