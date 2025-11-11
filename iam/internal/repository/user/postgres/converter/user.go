package converter

import (
	"database/sql"
	"time"

	"github.com/pinai4/spaceship-factory/iam/internal/model"
	repoModel "github.com/pinai4/spaceship-factory/iam/internal/repository/user/postgres/model"
)

func UserToModel(user repoModel.User) model.User {
	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		updatedAt = &user.UpdatedAt.Time
	}

	modelNPs := make([]model.NotificationMethod, len(user.NotificationMethods))
	for i, nm := range user.NotificationMethods {
		modelNPs[i] = model.NotificationMethod(nm)
	}

	return model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Login:               user.Login,
			Email:               user.Email,
			NotificationMethods: modelNPs,
		},
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    updatedAt,
	}
}

func UserToRepoModel(user model.User) repoModel.User {
	var updatedAt sql.NullTime
	if user.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: *user.UpdatedAt, Valid: true}
	}

	repoNMs := make([]repoModel.NotificationMethod, len(user.Info.NotificationMethods))
	for i, notificationMethod := range user.Info.NotificationMethods {
		repoNMs[i] = repoModel.NotificationMethod(notificationMethod)
	}

	return repoModel.User{
		ID:                  user.ID,
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		NotificationMethods: repoNMs,
		PasswordHash:        user.PasswordHash,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           updatedAt,
	}
}
