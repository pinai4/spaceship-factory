package part

import (
	"github.com/pinai4/spaceship-factory/inventory/internal/repository"
	def "github.com/pinai4/spaceship-factory/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{partRepository: partRepository}
}
