package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationMethod struct {
	Provider string `json:"provider"`
	Target   string `json:"target"`
}

type NotificationMethods []NotificationMethod

func (nm NotificationMethods) Value() (driver.Value, error) {
	return json.Marshal(nm)
}

func (nm *NotificationMethods) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NotificationMethods: %v", value)
	}
	return json.Unmarshal(bytes, nm)
}

type User struct {
	ID                  uuid.UUID           `db:"id"`
	Login               string              `db:"login"`
	Email               string              `db:"email"`
	NotificationMethods NotificationMethods `db:"notification_methods"`
	PasswordHash        string              `db:"password_hash"`
	CreatedAt           time.Time           `db:"created_at"`
	UpdatedAt           sql.NullTime        `db:"updated_at"`
}
