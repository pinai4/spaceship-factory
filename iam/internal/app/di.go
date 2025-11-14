package app

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	authV1API "github.com/pinai4/spaceship-factory/iam/internal/api/auth/v1"
	userV1API "github.com/pinai4/spaceship-factory/iam/internal/api/user/v1"
	"github.com/pinai4/spaceship-factory/iam/internal/config"
	"github.com/pinai4/spaceship-factory/iam/internal/repository"
	sessionRepository "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis"
	userRepository "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres"
	"github.com/pinai4/spaceship-factory/iam/internal/service"
	authService "github.com/pinai4/spaceship-factory/iam/internal/service/auth"
	"github.com/pinai4/spaceship-factory/iam/internal/service/passwordhasher"
	userService "github.com/pinai4/spaceship-factory/iam/internal/service/user"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	userV1API userV1.UserServiceServer
	authV1API authV1.AuthServiceServer

	userService service.UserService
	authService service.AuthService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordHasher    service.PasswordHasher

	postgresDB *sqlx.DB
	redisDB    *redis.Client
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx), d.PasswordHasher())
	}

	return d.userService
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
			d.PasswordHasher(),
			d.Config().Session.TTL(),
		)
	}

	return d.authService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.userRepository
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisDB(ctx))
	}

	return d.sessionRepository
}

func (d *diContainer) PasswordHasher() service.PasswordHasher {
	if d.passwordHasher == nil {
		d.passwordHasher = passwordhasher.NewPasswordHasher(4)
	}

	return d.passwordHasher
}

func (d *diContainer) PostgresDB(ctx context.Context) *sqlx.DB {
	if d.postgresDB == nil {
		db, err := sqlx.Open("pgx", d.Config().Postgres.DSN())
		if err != nil {
			panic(fmt.Sprintf("failed to open PostgresDB: %s\n", err.Error()))
		}
		d.closer.AddNamed("PostgresDB handler", func(ctx context.Context) error {
			return db.Close()
		})

		if err := db.PingContext(ctx); err != nil {
			panic(fmt.Sprintf("failed to ping PostgresDB: %s\n", err.Error()))
		}

		d.postgresDB = db
	}

	return d.postgresDB
}

func (d *diContainer) RedisDB(ctx context.Context) *redis.Client {
	if d.redisDB == nil {
		db := redis.NewClient(&redis.Options{
			Addr:            d.Config().Redis.Address(),
			PoolSize:        d.Config().Redis.MaxIdle(),
			ConnMaxIdleTime: d.Config().Redis.IdleTimeout(),
			DialTimeout:     d.Config().Redis.ConnectionTimeout(),
		})
		d.closer.AddNamed("RedisDB client", func(ctx context.Context) error {
			return db.Close()
		})

		if err := db.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("failed to ping RedisDB: %s\n", err.Error()))
		}

		d.redisDB = db
	}

	return d.redisDB
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
