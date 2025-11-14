package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
)

func SessionToProto(session model.Session) *commonV1.Session {
	var protoUpdatedAt *timestamppb.Timestamp
	if session.UpdatedAt != nil {
		protoUpdatedAt = timestamppb.New(*session.UpdatedAt)
	}

	return &commonV1.Session{
		Uuid:      session.ID.String(),
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: protoUpdatedAt,
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
