package redis

import (
	"context"

	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/converter"
	repoModel "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/model"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Session, error) {
	res := r.db.HGetAll(ctx, r.getSessionKey(id.String()))

	if len(res.Val()) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var session repoModel.Session
	if err := res.Scan(&session); err != nil {
		return model.Session{}, err
	}

	return repoConverter.SessionToModel(session), nil
}
