package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationMethod struct {
	Provider string
	Target   string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type User struct {
	ID           uuid.UUID
	Info         UserInfo
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}
