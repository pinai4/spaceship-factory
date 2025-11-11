package v1

import (
	"github.com/pinai4/spaceship-factory/iam/internal/service"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

var _ userV1.UserServiceServer = (*api)(nil)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{userService: userService}
}
