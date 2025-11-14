package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	User      User
	CreatedAt time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}
