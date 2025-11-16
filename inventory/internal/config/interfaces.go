package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type GRPCServerConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type AppEnvironmentConfig interface {
	IsDev() bool
	IsTest() bool
	IsProd() bool
}

type IAMGRPCClientConfig interface {
	Address() string
}
