package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/pinai4/spaceship-factory/order/internal/config/env"
)

type Config struct {
	Logger                     LoggerConfig
	HTTPServer                 HTTPServerConfig
	Postgres                   PostgresConfig
	DBMigrations               DBMigrationsConfig
	InventoryGRPCClient        InventoryGRPCClientConfig
	PaymentGRPCClient          PaymentGRPCClientConfig
	IAMGRPCClient              IAMGRPCClientConfig
	Kafka                      KafkaConfig
	OrderPaidEventProducer     OrderPaidEventProducerConfig
	ShipAssembledEventConsumer ShipAssembledEventConsumerConfig
}

func Load(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	httpServerCfg, err := env.NewHTTPServerConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	dbMigrationsCfg, err := env.NewDBMigrationsConfig()
	if err != nil {
		return nil, err
	}

	inventoryGRPCClientCfg, err := env.NewInventoryGRPCClientConfig()
	if err != nil {
		return nil, err
	}

	paymentGRPCClientCfg, err := env.NewPaymentGRPCClientConfig()
	if err != nil {
		return nil, err
	}

	iamGRPCClientCfg, err := env.NewIAMGRPCClientConfig()
	if err != nil {
		return nil, err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return nil, err
	}

	orderPaidEventProducerCfg, err := env.NewOrderPaidEventProducerConfig()
	if err != nil {
		return nil, err
	}

	shipAssembledEventConsumerCfg, err := env.NewShipAssembledEventConsumerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:                     loggerCfg,
		HTTPServer:                 httpServerCfg,
		Postgres:                   postgresCfg,
		DBMigrations:               dbMigrationsCfg,
		InventoryGRPCClient:        inventoryGRPCClientCfg,
		PaymentGRPCClient:          paymentGRPCClientCfg,
		IAMGRPCClient:              iamGRPCClientCfg,
		Kafka:                      kafkaCfg,
		OrderPaidEventProducer:     orderPaidEventProducerCfg,
		ShipAssembledEventConsumer: shipAssembledEventConsumerCfg,
	}, nil
}
