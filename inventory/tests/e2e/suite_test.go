package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const (
	testsTimeout     = 5 * time.Minute
	loggerLevelValue = "debug"

	// project name for Docker network and containers
	projectName = "inventory-e2e"

	// App environment variable keys
	envMongoImageNameKey    = "MONGO_IMAGE_NAME"
	envMongoExternalHostKey = "MONGO_EXTERNAL_HOST"
	envMongoExternalPortKey = "MONGO_EXTERNAL_PORT"
	envMongoDatabaseKey     = "MONGO_DATABASE"
	envMongoUsernameKey     = "MONGO_INITDB_ROOT_USERNAME"
	envMongoPasswordKey     = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	envMongoAuthDBKey       = "MONGO_AUTH_DB"

	appDockerfile = "deploy/docker/inventory/Dockerfile"

	mongoPartsCollectionName = "parts"

	// App environment variable keys
	envGRPCHostKey = "GRPC_HOST"
	envGRPCPortKey = "GRPC_PORT"
)

type E2ESuite struct {
	suite.Suite
	log               logger.Logger
	ctx               context.Context
	ctxCancel         context.CancelFunc
	suiteCtx          context.Context
	suiteCtxCancel    context.CancelFunc
	env               *TestEnvironment
	grpcClientConn    *grpc.ClientConn
	inventoryV1Client inventoryV1.InventoryServiceClient
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}
