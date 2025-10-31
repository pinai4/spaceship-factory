//go:build integration

package e2e

import (
	"context"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/app"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/mongo"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/network"
	"github.com/pinai4/spaceship-factory/platform/pkg/testcontainers/path"
	inventoryV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (s *E2ESuite) SetupSuite() {
	var setupOK bool
	defer func() {
		if setupOK == false {
			s.cleanup()
		}
	}()

	s.log = logger.NewZapLogger(loggerLevelValue, true)

	s.suiteCtx, s.suiteCtxCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Load the .env file and set variables in the environment
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	s.Require().NoError(err, "Failed to load .env file")

	// Set variables in the process environment
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	s.initEnvironment()
	s.initInventoryV1Client()

	setupOK = true
}

func (s *E2ESuite) initEnvironment() {
	s.log.Info(s.suiteCtx, "🚀 Preparing test environment...")

	s.env = &TestEnvironment{}

	// Step 1: Create a shared Docker network
	generatedNetwork, err := network.NewNetwork(s.suiteCtx, projectName)
	s.Require().NoError(err, "Failed to create shared network")
	s.env.Network = generatedNetwork
	s.log.Info(s.suiteCtx, "✅ Network successfully created")

	// Step 2: Start the MongoDB container
	generatedMongo, err := mongo.NewContainer(s.suiteCtx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithImageName(s.getEnvVarWithCheck(envMongoImageNameKey)),
		mongo.WithDatabase(s.getEnvVarWithCheck(envMongoDatabaseKey)),
		mongo.WithAuth(s.getEnvVarWithCheck(envMongoUsernameKey), s.getEnvVarWithCheck(envMongoPasswordKey)),
		mongo.WithAuthDB(s.getEnvVarWithCheck(envMongoAuthDBKey)),
		mongo.WithLogger(s.log),
	)
	s.Require().NoError(err, "Failed to start MongoDB container")
	s.env.Mongo = generatedMongo
	s.log.Info(s.suiteCtx, "✅ MongoDB container successfully started")

	// Step 3: Start the application container
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		envAppEnvKey: appEnvTest,
		// Override MongoDB host for connecting to the container from testcontainers
		envMongoHostKey: generatedMongo.Config().ContainerName,
		envMongoPortKey: generatedMongo.Config().Port,
	}

	appContainer, err := app.NewContainer(s.suiteCtx,
		app.WithDockerfile(projectRoot, appDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithLogger(s.log),
	)
	s.Require().NoError(err, "Failed to start application container")
	s.env.App = appContainer
	s.log.Info(s.suiteCtx, "✅ Application container successfully started")

	s.log.Info(s.suiteCtx, "🎉 Test environment is ready")
}

func (s *E2ESuite) initInventoryV1Client() {
	var err error
	s.grpcClientConn, err = grpc.NewClient(
		s.env.App.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err, "failed to open gRPC client")
	s.log.Info(s.suiteCtx, "✅ grpcClientConn has been opened")

	s.inventoryV1Client = inventoryV1.NewInventoryServiceClient(s.grpcClientConn)
	s.log.Info(s.suiteCtx, "✅ inventoryV1Client has been created")
}

func (s *E2ESuite) getEnvVarWithCheck(key string) string {
	value := os.Getenv(key)
	if value == "" {
		s.log.Warn(s.suiteCtx, "Environment variable not set", logger.String("key", key))
	}

	return value
}
