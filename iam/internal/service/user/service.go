package user

import (
	"github.com/pinai4/spaceship-factory/iam/internal/repository"
	def "github.com/pinai4/spaceship-factory/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	hasher         def.PasswordHasher
}

func NewService(userRepository repository.UserRepository, hasher def.PasswordHasher) *service {
	return &service{
		userRepository: userRepository,
		hasher:         hasher,
	}
}
