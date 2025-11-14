package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoConverter "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/converter"
)

func (r *repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	repoSession := repoConverter.SessionToRepoModel(session)

	_, err := r.db.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, r.getSessionKey(repoSession.ID), repoSession)
		pipe.Expire(ctx, r.getSessionKey(repoSession.ID), ttl)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
