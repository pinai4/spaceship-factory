//go:build unit || !integration

package assembly_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pinai4/spaceship-factory/assembly/internal/service"
	"github.com/pinai4/spaceship-factory/assembly/internal/service/assembly"
	serviceMocks "github.com/pinai4/spaceship-factory/assembly/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	assemblyProducer *serviceMocks.AssemblyProducer

	service service.AssemblyService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.assemblyProducer = serviceMocks.NewAssemblyProducer(s.T())

	s.service = assembly.NewService(s.assemblyProducer)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
