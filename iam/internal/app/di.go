package app

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	userV1API "github.com/pinai4/spaceship-factory/iam/internal/api/user/v1"
	"github.com/pinai4/spaceship-factory/iam/internal/config"
	"github.com/pinai4/spaceship-factory/iam/internal/repository"
	userRepository "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres"
	"github.com/pinai4/spaceship-factory/iam/internal/service"
	"github.com/pinai4/spaceship-factory/iam/internal/service/passwordhasher"
	userService "github.com/pinai4/spaceship-factory/iam/internal/service/user"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	userV1API userV1.UserServiceServer

	userService    service.UserService
	userRepository repository.UserRepository
	passwordHasher service.PasswordHasher

	postgresDB *sqlx.DB
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

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx), d.PasswordHasher())
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.userRepository
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

func (d *diContainer) Config() *config.Config {
	return d.config
}
