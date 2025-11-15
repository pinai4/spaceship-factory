package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type GRPCServerConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}

type DBMigrationsConfig interface {
	Path() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}

type AppEnvironmentConfig interface {
	IsDev() bool
	IsTest() bool
	IsProd() bool
}
