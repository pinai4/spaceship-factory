package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type HTTPServerConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	DSN() string
}

type DBMigrationsConfig interface {
	Path() string
}

type InventoryGRPCClientConfig interface {
	Address() string
}

type PaymentGRPCClientConfig interface {
	Address() string
}
