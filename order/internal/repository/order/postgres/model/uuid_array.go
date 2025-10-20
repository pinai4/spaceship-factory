package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var _ sql.Scanner = (*UUIDArray)(nil)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return a.parse(v)
	case []byte:
		return a.parse(string(v))
	default:
		return fmt.Errorf("unsupported type %T for UUIDArray", v)
	}
}

func (a *UUIDArray) parse(s string) error {
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = []uuid.UUID{}
		return nil
	}
	parts := strings.Split(s, ",")
	res := make([]uuid.UUID, len(parts))
	for i, p := range parts {
		id, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			return err
		}
		res[i] = id
	}
	*a = res
	return nil
}
