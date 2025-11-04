package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/pinai4/spaceship-factory/order/internal/api/order/v1"
	"github.com/pinai4/spaceship-factory/order/internal/client"
	inventoryV1Client "github.com/pinai4/spaceship-factory/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/pinai4/spaceship-factory/order/internal/client/grpc/payment/v1"
	"github.com/pinai4/spaceship-factory/order/internal/config"
	"github.com/pinai4/spaceship-factory/order/internal/repository"
	orderRepository "github.com/pinai4/spaceship-factory/order/internal/repository/order/postgres"
	"github.com/pinai4/spaceship-factory/order/internal/service"
	orderService "github.com/pinai4/spaceship-factory/order/internal/service/order"
	"github.com/pinai4/spaceship-factory/order/internal/service/order/producer"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	orderV1 "github.com/pinai4/spaceship-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	orderV1API orderV1.Handler

	orderService  service.OrderService
	orderProducer service.OrderProducer

	orderRepository repository.OrderRepository

	postgresDB *sqlx.DB

	inventoryV1ClientGRPC inventoryV1.InventoryServiceClient
	inventoryV1Client     client.InventoryClient

	paymentV1ClientGRPC paymentV1.PaymentServiceClient
	paymentV1Client     client.PaymentClient

	syncProducer sarama.SyncProducer
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.PaymentV1Client(), d.InventoryV1Client(), d.OrderProducer())
	}

	return d.orderService
}

func (d *diContainer) OrderProducer() service.OrderProducer {
	if d.orderProducer == nil {
		topics := map[string]string{
			producer.OrderPaidEventTopicKey: d.Config().OrderPaidEventProducer.Topic(),
		}
		d.orderProducer = producer.NewProducer(d.SyncProducer(), topics)
	}

	return d.orderProducer
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.orderRepository
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

func (d *diContainer) InventoryV1ClientGRPC() inventoryV1.InventoryServiceClient {
	if d.inventoryV1ClientGRPC == nil {
		conn, err := grpc.NewClient(
			d.Config().InventoryGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Inventory Service (GRPC): %s\n", err.Error()))
		}
		d.closer.AddNamed("InventoryClient connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryV1ClientGRPC = inventoryV1.NewInventoryServiceClient(conn)
	}

	return d.inventoryV1ClientGRPC
}

func (d *diContainer) InventoryV1Client() client.InventoryClient {
	if d.inventoryV1Client == nil {
		d.inventoryV1Client = inventoryV1Client.NewClient(d.InventoryV1ClientGRPC())
	}

	return d.inventoryV1Client
}

func (d *diContainer) PaymentV1ClientGRPC() paymentV1.PaymentServiceClient {
	if d.paymentV1ClientGRPC == nil {
		conn, err := grpc.NewClient(
			d.Config().PaymentGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Payment Service (GRPC): %s\n", err.Error()))
		}
		d.closer.AddNamed("PaymentClient connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentV1ClientGRPC = paymentV1.NewPaymentServiceClient(conn)
	}

	return d.paymentV1ClientGRPC
}

func (d *diContainer) PaymentV1Client() client.PaymentClient {
	if d.paymentV1Client == nil {
		d.paymentV1Client = paymentV1Client.NewClient(d.PaymentV1ClientGRPC())
	}

	return d.paymentV1Client
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			d.Config().Kafka.Brokers(),
			d.Config().OrderPaidEventProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		d.closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
