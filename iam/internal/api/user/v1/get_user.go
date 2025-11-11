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

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	id, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user uuid")
	}

	user, err := a.userService.Get(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	return &userV1.GetUserResponse{User: converter.UserToProto(user)}, nil
}
