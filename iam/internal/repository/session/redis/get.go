package redis

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/converter"
	repoModel "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/model"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Session, error) {
	var session repoModel.Session
	if err := r.db.HGetAll(ctx, r.getSessionKey(id.String())).Scan(&session); err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Session{}, model.ErrSessionNotFound
		}
		return model.Session{}, err
	}

	return repoConverter.SessionToModel(session), nil
}
