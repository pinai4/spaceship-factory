package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres/converter"
	repoModel "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres/model"
)

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	const q = `
	SELECT
		id, login, email, notification_methods, password_hash, created_at, updated_at
	FROM
		users
	WHERE
		login = $1`

	var repoUser repoModel.User
	if err := r.db.GetContext(ctx, &repoUser, q, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("UserRepository.GetByLogin get user error: %w", err)
	}

	return repoConverter.UserToModel(repoUser), nil
}
