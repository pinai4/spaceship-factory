package converter

import (
	"github.com/google/uuid"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoModel "github.com/pinai4/spaceship-factory/iam/internal/repository/session/redis/model"
)

func SessionToModel(session repoModel.Session) model.Session {
	modelNPs := make([]model.NotificationMethod, len(session.UserNotificationMethods))
	for i, nm := range session.UserNotificationMethods {
		modelNPs[i] = model.NotificationMethod(nm)
	}

	return model.Session{
		ID: uuid.MustParse(session.ID),
		User: model.User{
			ID: uuid.MustParse(session.UserID),
			Info: model.UserInfo{
				Login:               session.UserLogin,
				Email:               session.UserEmail,
				NotificationMethods: modelNPs,
			},
			CreatedAt: session.UserCreatedAt,
			UpdatedAt: session.UserUpdatedAt,
		},
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

func SessionToRepoModel(session model.Session) repoModel.Session {
	repoNMs := make([]repoModel.UserNotificationMethod, len(session.User.Info.NotificationMethods))
	for i, notificationMethod := range session.User.Info.NotificationMethods {
		repoNMs[i] = repoModel.UserNotificationMethod(notificationMethod)
	}

	return repoModel.Session{
		ID:                      session.ID.String(),
		UserID:                  session.User.ID.String(),
		UserLogin:               session.User.Info.Login,
		UserEmail:               session.User.Info.Email,
		UserNotificationMethods: repoNMs,
		UserCreatedAt:           session.User.CreatedAt,
		UserUpdatedAt:           session.User.UpdatedAt,
		CreatedAt:               session.CreatedAt,
		UpdatedAt:               session.UpdatedAt,
		ExpiresAt:               session.ExpiresAt,
	}
}
