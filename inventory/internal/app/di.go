package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/pinai4/spaceship-factory/inventory/internal/api/inventory/v1"
	"github.com/pinai4/spaceship-factory/inventory/internal/config"
	"github.com/pinai4/spaceship-factory/inventory/internal/repository"
	partRepository "github.com/pinai4/spaceship-factory/inventory/internal/repository/part/mongodb"
	"github.com/pinai4/spaceship-factory/inventory/internal/service"
	partService "github.com/pinai4/spaceship-factory/inventory/internal/service/part"
	"github.com/pinai4/spaceship-factory/platform/pkg/closer"
	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	config *config.Config
	closer *closer.Closer

	inventoryV1API inventoryV1.InventoryServiceServer

	partService service.PartService

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	authV1ClientGRPC authV1.AuthServiceClient
}

func NewDiContainer(config *config.Config, closer *closer.Closer) *diContainer {
	return &diContainer{config: config, closer: closer}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.Config().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		d.closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(d.Config().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) AuthV1ClientGRPC() authV1.AuthServiceClient {
	if d.authV1ClientGRPC == nil {
		conn, err := grpc.NewClient(
			d.Config().IAMGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to IAM Service (GRPC): %s\n", err.Error()))
		}
		d.closer.AddNamed("IAMClient connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.authV1ClientGRPC = authV1.NewAuthServiceClient(conn)
	}

	return d.authV1ClientGRPC
}

func (d *diContainer) Config() *config.Config {
	return d.config
}
