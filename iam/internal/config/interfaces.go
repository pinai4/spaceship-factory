package config

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
