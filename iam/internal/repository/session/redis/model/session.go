package model

import (
	"encoding/json"
	"time"
)

type UserNotificationMethod struct {
	Provider string `json:"provider"`
	Target   string `json:"target"`
}

type UserNotificationMethods []UserNotificationMethod

func (nm UserNotificationMethods) MarshalBinary() ([]byte, error) {
	return json.Marshal(nm)
}

func (nm *UserNotificationMethods) UnmarshalText(data []byte) error {
	// alternative solution with alias
	// type alias UserNotificationMethods
	// return json.Unmarshal(data, (*alias)(nm))
	var tempNM []UserNotificationMethod
	if err := json.Unmarshal(data, &tempNM); err != nil {
		return err
	}
	*nm = tempNM

	return nil
}

type Session struct {
	ID                      string                  `redis:"id"`
	UserID                  string                  `redis:"user_id"`
	UserLogin               string                  `redis:"user_login"`
	UserEmail               string                  `redis:"user_email"`
	UserNotificationMethods UserNotificationMethods `redis:"user_notification_methods"`
	UserCreatedAt           time.Time               `redis:"user_created_at"`
	UserUpdatedAt           *time.Time              `redis:"user_updated_at,omitempty"`
	CreatedAt               time.Time               `redis:"created_at"`
	UpdatedAt               *time.Time              `redis:"updated_at,omitempty"`
	ExpiresAt               time.Time               `redis:"expires_at"`
}
