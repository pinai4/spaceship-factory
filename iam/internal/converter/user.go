package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	commonV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/common/v1"
	userV1 "github.com/pinai4/spaceship-factory/shared/pkg/proto/user/v1"
)

func UserToProto(user model.User) *commonV1.User {
	protoNMs := make([]*commonV1.NotificationMethod, len(user.Info.NotificationMethods))
	for i, nm := range user.Info.NotificationMethods {
		protoNMs[i] = &commonV1.NotificationMethod{
			ProviderName: nm.Provider,
			Target:       nm.Target,
		}
	}

	var protoUpdatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		protoUpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &commonV1.User{
		Uuid: user.ID.String(),
		Info: &commonV1.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: protoNMs,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: protoUpdatedAt,
	}
}

func UserRegistrationInfoToModel(regInfo *userV1.UserRegistrationInfo) model.UserRegistrationInfo {
	if regInfo == nil {
		return model.UserRegistrationInfo{}
	}

	var info model.UserInfo
	if regInfo.GetInfo() != nil {
		notificationMethods := make([]model.NotificationMethod, len(regInfo.GetInfo().NotificationMethods))
		for i, notificationMethod := range regInfo.GetInfo().NotificationMethods {
			notificationMethods[i] = model.NotificationMethod{
				Provider: notificationMethod.GetProviderName(),
				Target:   notificationMethod.GetTarget(),
			}
		}

		info = model.UserInfo{
			Login:               regInfo.GetInfo().GetLogin(),
			Email:               regInfo.GetInfo().GetEmail(),
			NotificationMethods: notificationMethods,
		}
	}

	return model.UserRegistrationInfo{
		Info:     info,
		Password: regInfo.GetPassword(),
	}
}
