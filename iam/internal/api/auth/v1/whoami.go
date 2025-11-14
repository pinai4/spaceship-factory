package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pinai4/spaceship-factory/iam/internal/converter"
	"github.com/pinai4/spaceship-factory/iam/internal/model"
	authV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	sessID, err := uuid.Parse(req.GetSessionUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session uuid")
	}

	session, err := a.authService.Whoami(ctx, sessID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Errorf(codes.NotFound, "session not found")
		}
		return nil, status.Error(codes.Internal, "failed to check session")
	}

	return &authV1.WhoamiResponse{
		User:    converter.UserToProto(session.User),
		Session: converter.SessionToProto(session),
	}, nil
}
