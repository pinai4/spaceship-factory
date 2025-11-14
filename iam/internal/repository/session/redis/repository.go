package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	def "github.com/pinai4/spaceship-factory/iam/internal/repository"
)

const sessionKeyPrefix = "session:"

var _ def.SessionRepository = (*repository)(nil)

type repository struct {
	db *redis.Client
}

func NewRepository(db *redis.Client) *repository {
	return &repository{db: db}
}

func (r *repository) getSessionKey(id string) string {
	return fmt.Sprintf("%s%s", sessionKeyPrefix, id)
}
