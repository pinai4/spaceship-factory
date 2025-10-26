package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type GRPCServerConfig interface {
	Address() string
}
