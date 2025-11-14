package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pinai4/spaceship-factory/iam/internal/converter"
	"github.com/pinai4/spaceship-factory/iam/internal/model"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	id := uuid.New()

	if err := a.userService.Register(ctx, id, converter.UserRegistrationInfoToModel(req.GetInfo())); err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &userV1.RegisterResponse{UserUuid: id.String()}, nil
}
