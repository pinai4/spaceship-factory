package assembly

import def "github.com/pinai4/spaceship-factory/assembly/internal/service"

var _ def.AssemblyService = (*service)(nil)

type service struct {
	assemblyProducer def.AssemblyProducer
}

func NewService(assemblyProducer def.AssemblyProducer) *service {
	return &service{
		assemblyProducer: assemblyProducer,
	}
}
